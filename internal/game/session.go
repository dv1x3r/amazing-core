package game

import (
	"crypto/rand"
	"encoding/binary"
	"strings"
	"sync"
	"time"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// detachedSessionTTL is how long a session without connections waits for RELOGIN.
const detachedSessionTTL = 5 * time.Minute

// Session is a logged in player.
type Session struct {
	playerOID      types.OID
	platform       gsf.Platform
	loginName      string
	sessionOID     types.OID
	conversationID int64

	// guarded by SessionManager.mu
	userConn   *gsf.Conn
	syncConn   *gsf.Conn
	detachedAt time.Time
}

// SessionManager keeps the sessions by the RELOGIN token, by the player and by the connection.
type SessionManager struct {
	mu       sync.Mutex
	byToken  map[types.OID]*Session
	byPlayer map[types.OID]*Session
	byConn   map[*gsf.Conn]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		byToken:  map[types.OID]*Session{},
		byPlayer: map[types.OID]*Session{},
		byConn:   map[*gsf.Conn]*Session{},
	}
}

// Login starts a new session of the player on the USER connection, an older session of the player is closed.
func (m *SessionManager) Login(conn *gsf.Conn, playerOID types.OID, platform gsf.Platform, loginName string) *Session {
	s := &Session{
		playerOID:      playerOID,
		platform:       platform,
		loginName:      loginName,
		sessionOID:     types.OID{Number: randomInt63() & (1<<40 - 1)},
		conversationID: randomInt63(),
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.removeExpiredLocked(time.Now())
	if old, ok := m.byPlayer[playerOID]; ok {
		m.removeLocked(old)
	}
	m.byToken[s.sessionOID] = s
	m.byPlayer[playerOID] = s
	m.bindLocked(conn, s, false)
	return s
}

// Relogin binds the USER connection to the session of the LOGIN response,
// the conversation ID and the login name must match the LOGIN ones.
func (m *SessionManager) Relogin(conn *gsf.Conn, sessionOID types.OID, conversationID int64, loginName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.byToken[sessionOID]
	if !ok || s.conversationID != conversationID || !strings.EqualFold(s.loginName, loginName) {
		return false
	}
	if s.expired(time.Now()) {
		m.removeLocked(s)
		return false
	}
	m.bindLocked(conn, s, false)
	return true
}

// AttachSync binds the SYNC connection to the session of the player, the client sends no token there.
func (m *SessionManager) AttachSync(conn *gsf.Conn, playerOID types.OID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.byPlayer[playerOID]
	if !ok || s.expired(time.Now()) {
		return false
	}
	m.bindLocked(conn, s, true)
	return true
}

// Logout closes the session of the connection.
func (m *SessionManager) Logout(conn *gsf.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.byConn[conn]; ok {
		m.removeLocked(s)
	}
}

// Detach unbinds a closed connection, the session waits detachedSessionTTL for RELOGIN after the last one.
func (m *SessionManager) Detach(conn *gsf.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.unbindLocked(conn)
}

// Get returns the session of the connection, or nil before LOGIN.
func (m *SessionManager) Get(conn *gsf.Conn) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.byConn[conn]
}

func (m *SessionManager) bindLocked(conn *gsf.Conn, s *Session, isSync bool) {
	m.unbindLocked(conn)
	old := &s.userConn
	if isSync {
		old = &s.syncConn
	}
	if *old != nil {
		delete(m.byConn, *old)
	}
	*old = conn
	s.detachedAt = time.Time{}
	m.byConn[conn] = s
}

func (m *SessionManager) unbindLocked(conn *gsf.Conn) {
	s, ok := m.byConn[conn]
	if !ok {
		return
	}
	delete(m.byConn, conn)
	if s.userConn == conn {
		s.userConn = nil
	}
	if s.syncConn == conn {
		s.syncConn = nil
	}
	if s.userConn == nil && s.syncConn == nil {
		s.detachedAt = time.Now()
	}
}

func (m *SessionManager) removeLocked(s *Session) {
	delete(m.byToken, s.sessionOID)
	if m.byPlayer[s.playerOID] == s {
		delete(m.byPlayer, s.playerOID)
	}
	for _, conn := range []*gsf.Conn{s.userConn, s.syncConn} {
		if conn != nil {
			delete(m.byConn, conn)
		}
	}
}

func (m *SessionManager) removeExpiredLocked(now time.Time) {
	for _, s := range m.byToken {
		if s.expired(now) {
			m.removeLocked(s)
		}
	}
}

func (s *Session) expired(now time.Time) bool {
	return !s.detachedAt.IsZero() && now.Sub(s.detachedAt) > detachedSessionTTL
}

func randomInt63() int64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return int64(binary.LittleEndian.Uint64(b[:]) >> 1)
}
