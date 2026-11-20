package game

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
)

// ── Sync ─────────────────────────────────────────────────────────────────────

// SyncLogin authenticates the player on the SYNC server.
// The SYNC server handles real-time positional and social updates.
// It also handles RELOGIN, the client sends it with the same payload after the connection is restored.
func (h *Handler) SyncLogin(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SyncLoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session := h.sessions.AttachSync(r.Conn(), req.UOID, req.Token)
	if session == nil {
		return errNoSession
	}
	h.world.Join(session.player, r.Conn())
	res := &messages.SyncLoginResponse{}
	return w.Write(res)
}

// MovePlayer receives local player movement and relays it to other sync peers.
func (h *Handler) MovePlayer(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &notify.Move{}
	if err := r.Read(req); err != nil {
		return err
	}
	session := h.sessions.Get(r.Conn())
	if session == nil {
		return nil
	}
	h.world.Move(session.player, req)
	return nil
}

// SyncAction receives a local player action and relays it to other sync peers.
func (h *Handler) SyncAction(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &notify.Action{}
	if err := r.Read(req); err != nil {
		return err
	}
	session := h.sessions.Get(r.Conn())
	if session == nil {
		return nil
	}
	if h.world.Action(session.player, req) {
		r.MarkDropped()
	}
	return nil
}

// HeartbeatNotify keeps the SYNC connection alive, only SYNC_LOGIN binds it to a session.
func (h *Handler) HeartbeatNotify(w gsf.ResponseWriter, r *gsf.Request) error {
	return r.Read(&notify.Heartbeat{})
}
