package game

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
)

// ── Sync ─────────────────────────────────────────────────────────────────────

// SyncLogin binds the SYNC connection to the session, so the player sees the other players and is seen.
// The client also sends it as RELOGIN after the SYNC connection is restored.
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

// MovePlayer relays the movement of the player to the players nearby.
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

// SyncAction relays the actions of the player, e.g. animations, to the players nearby.
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

// HeartbeatNotify keeps the server idle timeout from closing the SYNC connection.
func (h *Handler) HeartbeatNotify(w gsf.ResponseWriter, r *gsf.Request) error {
	return r.Read(&notify.Heartbeat{})
}
