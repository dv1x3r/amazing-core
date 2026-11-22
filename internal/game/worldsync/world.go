package worldsync

import (
	"sync"
	"time"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
)

// World tracks the players in the world by location and notifies them about each other.
// The session owns its Player, World guards the Player fields.
type World struct {
	mu sync.Mutex

	// livePlayers are the players with a SYNC connection by OID.
	livePlayers map[int64]*Player

	// playersByLocation are the live players that entered a location, by location OID, then by player OID.
	playersByLocation map[int64]map[int64]*Player

	// lastDetailVersion makes every player details version unique,
	// so the clients reload the details after any change or a new session.
	lastDetailVersion int64
}

// NewWorld creates an empty world.
func NewWorld() *World {
	return &World{
		livePlayers:       map[int64]*Player{},
		playersByLocation: map[int64]map[int64]*Player{},
	}
}

// NewPlayer creates the player of a new session.
func (w *World) NewPlayer(playerOID types.OID) *Player {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastDetailVersion++
	return &Player{OID: playerOID, location: noLocation, detailVersion: w.lastDetailVersion}
}

// Join puts the player into the world on its SYNC connection.
// The player gets into a location on ENTER_BUILDING, a SYNC reconnect returns it to its location.
func (w *World) Join(player *Player, syncConn *gsf.Conn) {
	w.send(func() []outboundNotify {
		// Already joined on this connection.
		if player.syncConn == syncConn {
			return nil
		}

		// Remove the player that had this connection, the player from its old connection,
		// and the player of an older session with the same OID.
		outbound := w.removeLivePlayer(w.playerOn(syncConn))
		outbound = append(outbound, w.removeLivePlayer(player)...)
		if old := w.livePlayers[player.OID.Int64()]; old != nil {
			outbound = append(outbound, w.removeLivePlayer(old)...)
		}

		player.syncConn = syncConn
		w.livePlayers[player.OID.Int64()] = player
		if player.location != noLocation {
			outbound = append(outbound, w.enterLocation(player)...)
		}
		return outbound
	})
}

// Disconnect removes the player of the closed SYNC connection from the world.
// It looks the player up by the connection, the session can be gone after another LOGIN.
func (w *World) Disconnect(syncConn *gsf.Conn) {
	w.send(func() []outboundNotify {
		return w.removeLivePlayer(w.playerOn(syncConn))
	})
}

// AppearanceChanged bumps the details version and tells nearby players to reload them.
func (w *World) AppearanceChanged(player *Player) {
	w.send(func() []outboundNotify {
		w.lastDetailVersion++
		player.detailVersion = w.lastDetailVersion
		if player.syncConn == nil {
			return nil
		}
		return w.broadcast(player.location, player, clientmessagetype.CHANGE_OBJECT, player.changeObjectNotify())
	})
}

// Leave takes the player out of its location when it travels (INIT_LOCATION), nearby players remove it.
func (w *World) Leave(player *Player) {
	w.send(func() []outboundNotify {
		outbound := w.leaveLocation(player)
		player.location = noLocation
		return outbound
	})
}

// Enter puts the player into the loaded location at its spawn point (ENTER_BUILDING).
// The client clears its players right before ENTER_BUILDING, so it gets ADD_PLAYER of everyone there again.
func (w *World) Enter(player *Player, locationOID types.OID, pos types.Position, qth types.QTH) {
	w.send(func() []outboundNotify {
		outbound := w.leaveLocation(player)
		player.location = locationOID.Int64()
		player.pos = []types.Position{pos}
		player.qth = qth
		if player.syncConn == nil {
			return outbound
		}
		return append(outbound, w.enterLocation(player)...)
	})
}

// Move saves the player position and sends it to nearby players.
func (w *World) Move(player *Player, move *notify.Move) {
	w.send(func() []outboundNotify {
		if player.syncConn == nil {
			return nil
		}
		player.pos = move.Pos
		player.qth = move.QTH
		player.secondQTH = move.SecondQTH
		player.seq = move.Seq
		return w.broadcast(player.location, player, clientmessagetype.MOVE_PLAYER, player.playerMoveNotify())
	})
}

// Action sends a player action to nearby players.
// throttled is true when the action is dropped as a repeat.
func (w *World) Action(player *Player, action *notify.Action) (throttled bool) {
	w.send(func() []outboundNotify {
		if player.syncConn == nil {
			return nil
		}
		if player.throttlesItemAnimate(action, time.Now()) {
			throttled = true
			return nil
		}
		relayed := types.PlayerActionNotify{
			POID:       player.OID,
			ActionOID:  action.ActionOID,
			TargetOID:  action.TargetOID,
			SourceOID:  action.SourceOID,
			ActionType: action.ActionType,
			Parameters: action.Parameters,
			Pos:        action.Pos,
			QTH:        action.QTH,
			IsState:    action.IsState,
			ClearState: action.ClearState,
			StateIndex: action.StateIndex,
		}
		player.setState(relayed)
		return w.broadcast(player.location, player, clientmessagetype.ACTION, &relayed)
	})
	return throttled
}

// Chat sends a chat message to everyone in the sender's location, the sender too.
func (w *World) Chat(sender *Player, chat *notify.Chat) {
	w.send(func() []outboundNotify {
		if sender.syncConn == nil {
			return nil
		}
		return w.broadcast(sender.location, nil, clientmessagetype.CHAT, chat)
	})
}

// KeepAlive sends the empty interest list to every live player, so an idle SYNC connection gets data:
// the client closes it after 60 seconds without data (GSFSocketBase.SocketReceiveTimeout) and logs in again.
func (w *World) KeepAlive() {
	w.send(func() []outboundNotify {
		outbound := make([]outboundNotify, 0, len(w.livePlayers))
		for _, player := range w.livePlayers {
			outbound = append(outbound, newOutboundNotify(player.syncConn, clientmessagetype.INT_LIST, &notify.InterestList{}))
		}
		return outbound
	})
}

// send builds the notifies under the lock and sends them after unlocking it.
func (w *World) send(build func() []outboundNotify) {
	w.mu.Lock()
	notifies := build()
	w.mu.Unlock()
	sendAll(notifies)
}

// broadcast builds the notify for every player in the location except one.
func (w *World) broadcast(location int64, except *Player, msgType clientmessagetype.ClientMessageType, body gsf.Serializable) []outboundNotify {
	players := w.playersInLocation(location, except)
	outbound := make([]outboundNotify, 0, len(players))
	for _, recipient := range players {
		outbound = append(outbound, newOutboundNotify(recipient.syncConn, msgType, body))
	}
	return outbound
}

// enterLocation puts a live player into its location and sends ADD_PLAYER both ways.
func (w *World) enterLocation(player *Player) []outboundNotify {
	existing := w.playersInLocation(player.location, player)
	if w.playersByLocation[player.location] == nil {
		w.playersByLocation[player.location] = map[int64]*Player{}
	}
	w.playersByLocation[player.location][player.OID.Int64()] = player
	addNotify := player.addPlayerNotify()
	var outbound []outboundNotify
	for _, other := range existing {
		outbound = append(outbound,
			newOutboundNotify(player.syncConn, clientmessagetype.ADD_PLAYER, other.addPlayerNotify()),
			newOutboundNotify(other.syncConn, clientmessagetype.ADD_PLAYER, addNotify),
		)
	}
	return outbound
}

// leaveLocation takes a live player out of its location and tells nearby players to remove it.
func (w *World) leaveLocation(player *Player) []outboundNotify {
	if player.syncConn == nil || player.location == noLocation {
		return nil
	}
	outbound := w.broadcast(player.location, player, clientmessagetype.REMOVE_PLAYER, player.removePlayerNotify())
	players := w.playersByLocation[player.location]
	delete(players, player.OID.Int64())
	if len(players) == 0 {
		delete(w.playersByLocation, player.location)
	}
	return outbound
}

// removeLivePlayer removes the player from the world and tells nearby players.
// The player keeps its state and can join again.
func (w *World) removeLivePlayer(player *Player) []outboundNotify {
	if player == nil || player.syncConn == nil {
		return nil
	}
	outbound := w.leaveLocation(player)
	delete(w.livePlayers, player.OID.Int64())
	player.syncConn = nil
	return outbound
}

// playersInLocation returns the players in the location except one.
func (w *World) playersInLocation(location int64, except *Player) []*Player {
	players := w.playersByLocation[location]
	result := make([]*Player, 0, len(players))
	for _, player := range players {
		if player != except {
			result = append(result, player)
		}
	}
	return result
}

// playerOn returns the live player on the SYNC connection, or nil.
func (w *World) playerOn(syncConn *gsf.Conn) *Player {
	for _, player := range w.livePlayers {
		if player.syncConn == syncConn {
			return player
		}
	}
	return nil
}
