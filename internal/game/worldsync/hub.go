package worldsync

import (
	"sync"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
)

// Hub tracks live sync players and builds visibility notifications between them.
type Hub struct {
	mu sync.Mutex

	// remembered players, connected or not
	playersByOID map[int64]*Player

	// live players by sync session
	playersBySession map[*gsf.Session]*Player

	// live players grouped by location OID, then player OID
	playersByLocation map[int64]map[int64]*Player
}

// NewHub creates an empty world sync hub.
func NewHub() *Hub {
	return &Hub{
		playersByOID:      map[int64]*Player{},
		playersBySession:  map[*gsf.Session]*Player{},
		playersByLocation: map[int64]map[int64]*Player{},
	}
}

// Join registers a sync session and exchanges ADD_PLAYER notifications with players in the same location.
func (h *Hub) Join(session *gsf.Session, playerOID types.OID) error {
	return SendAll(h.join(session, playerOID))
}

func (h *Hub) join(session *gsf.Session, playerOID types.OID) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()

	outbound := []OutboundNotify{}
	current := h.playersBySession[session]
	player := h.ensurePlayer(playerOID)

	// This session is already joined as this player.
	if current == player && player.Session == session {
		return nil
	}

	// This session is already linked to another player.
	if current != nil && current != player {
		outbound = append(outbound, h.removeLivePlayer(current)...)
	}

	// A player can only have one live sync session.
	if player.Session != nil {
		outbound = append(outbound, h.removeLivePlayer(player)...)
	}

	existing := h.playersInLocation(player.Location, nil)
	h.addLivePlayer(player, session)

	newPlayerNotify := player.AddPlayerNotify()
	for _, other := range existing {
		outbound = append(outbound,
			NewOutboundNotify(session, clientmessagetype.ADD_PLAYER, other.AddPlayerNotify()),
			NewOutboundNotify(other.Session, clientmessagetype.ADD_PLAYER, newPlayerNotify),
		)
	}
	return outbound
}

// AppearanceChanged bumps the cached player-details version and asks nearby players to reload this player.
func (h *Hub) AppearanceChanged(playerOID types.OID) error {
	return SendAll(h.appearanceChanged(playerOID))
}

func (h *Hub) appearanceChanged(playerOID types.OID) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()

	player := h.ensurePlayer(playerOID)
	player.PlayerDetailVersion++

	if player.Session == nil {
		return nil
	}

	notify := player.ChangeObjectNotify()
	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, player) {
		outbound = append(outbound, NewOutboundNotify(recipient.Session, clientmessagetype.CHANGE_OBJECT, notify))
	}
	return outbound
}

// SetLocation moves a player between locations and notifies players that lose or gain visibility.
func (h *Hub) SetLocation(playerOID types.OID, locationOID types.OID) error {
	return SendAll(h.setLocation(playerOID, locationOID))
}

func (h *Hub) setLocation(playerOID types.OID, locationOID types.OID) []OutboundNotify {
	location := locationOID.Int64()

	h.mu.Lock()
	defer h.mu.Unlock()

	player := h.ensurePlayer(playerOID)
	oldLocation := player.Location
	if oldLocation == location {
		return nil
	}

	// Remember the target location even before the sync session joins.
	if player.Session == nil {
		player.Location = location
		return nil
	}

	oldRecipients := h.playersInLocation(oldLocation, player)
	newRecipients := h.playersInLocation(location, player)

	removeNotify := player.RemovePlayerNotify()
	h.movePlayer(player, location)
	addNotify := player.AddPlayerNotify()

	outbound := make([]OutboundNotify, 0, len(oldRecipients)+len(newRecipients)*2)
	for _, recipient := range oldRecipients {
		outbound = append(outbound, NewOutboundNotify(recipient.Session, clientmessagetype.REMOVE_PLAYER, removeNotify))
	}
	for _, other := range newRecipients {
		outbound = append(outbound,
			NewOutboundNotify(player.Session, clientmessagetype.ADD_PLAYER, other.AddPlayerNotify()),
			NewOutboundNotify(other.Session, clientmessagetype.ADD_PLAYER, addNotify),
		)
	}
	return outbound
}

// Move updates the latest movement state and relays it to players in the same location.
func (h *Hub) Move(session *gsf.Session, move *notify.Move) error {
	return SendAll(h.move(session, move))
}

func (h *Hub) move(session *gsf.Session, move *notify.Move) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()

	player := h.playersBySession[session]
	if player == nil {
		return nil
	}

	player.Pos = append(player.Pos[:0], move.Pos...)
	player.QTH = move.QTH
	player.SecondQTH = move.SecondQTH
	player.Seq = move.Seq

	notify := player.PlayerMoveNotify()
	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, player) {
		outbound = append(outbound, NewOutboundNotify(recipient.Session, clientmessagetype.MOVE_PLAYER, notify))
	}
	return outbound
}

// Leave removes a sync session and tells nearby players to remove that player.
func (h *Hub) Leave(session *gsf.Session) error {
	return SendAll(h.leave(session))
}

func (h *Hub) leave(session *gsf.Session) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()
	player := h.playersBySession[session]
	if player == nil {
		return nil
	}
	return h.removeLivePlayer(player)
}

// addLivePlayer indexes a connected sync player by session and location.
func (h *Hub) addLivePlayer(player *Player, session *gsf.Session) {
	player.Session = session
	h.playersBySession[session] = player
	h.addPlayerToLocation(player)
}

// removeLivePlayer removes a connected player from live indexes.
func (h *Hub) removeLivePlayer(player *Player) []OutboundNotify {
	if player.Session == nil {
		return nil
	}

	notify := player.RemovePlayerNotify()
	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, player) {
		outbound = append(outbound, NewOutboundNotify(recipient.Session, clientmessagetype.REMOVE_PLAYER, notify))
	}

	h.removePlayerFromLocation(player)
	delete(h.playersBySession, player.Session)

	// Keep the player in playersByOID so location/version survive reconnect.
	player.Session = nil
	return outbound
}

// movePlayer updates a connected player's location and location index.
func (h *Hub) movePlayer(player *Player, location int64) {
	if player.Location == location {
		return
	}
	h.removePlayerFromLocation(player)
	player.Location = location
	h.addPlayerToLocation(player)
}

// addPlayerToLocation inserts a connected player into its current location bucket.
func (h *Hub) addPlayerToLocation(player *Player) {
	if player.Session == nil {
		return
	}
	if h.playersByLocation[player.Location] == nil {
		h.playersByLocation[player.Location] = map[int64]*Player{}
	}
	h.playersByLocation[player.Location][player.OID.Int64()] = player
}

// removePlayerFromLocation removes a connected player from its current location bucket.
func (h *Hub) removePlayerFromLocation(player *Player) {
	if player.Session == nil {
		return
	}
	players := h.playersByLocation[player.Location]
	if players == nil {
		return
	}
	delete(players, player.OID.Int64())
	if len(players) == 0 {
		delete(h.playersByLocation, player.Location)
	}
}

// playersInLocation returns the connected players currently visible in a location.
func (h *Hub) playersInLocation(location int64, except *Player) []*Player {
	players := h.playersByLocation[location]
	if len(players) == 0 {
		return nil
	}
	result := make([]*Player, 0, len(players))
	for _, player := range players {
		if player != except {
			result = append(result, player)
		}
	}
	return result
}

// ensurePlayer returns the remembered sync player for an OID.
func (h *Hub) ensurePlayer(playerOID types.OID) *Player {
	playerKey := playerOID.Int64()
	player := h.playersByOID[playerKey]
	if player == nil {
		player = &Player{}
		h.playersByOID[playerKey] = player
	}
	player.OID = playerOID
	return player
}
