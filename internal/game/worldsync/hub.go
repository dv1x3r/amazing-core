package worldsync

import (
	"sync"
	"time"

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

	// live players by sync connection
	playersByConn map[*gsf.Conn]*Player

	// live players grouped by location OID, then player OID
	playersByLocation map[int64]map[int64]*Player
}

// NewHub creates an empty world sync hub.
func NewHub() *Hub {
	return &Hub{
		playersByOID:      map[int64]*Player{},
		playersByConn:     map[*gsf.Conn]*Player{},
		playersByLocation: map[int64]map[int64]*Player{},
	}
}

// Join registers a sync connection and exchanges ADD_PLAYER notifications with players in the same location.
func (h *Hub) Join(conn *gsf.Conn, playerOID types.OID) error {
	return SendAll(h.join(conn, playerOID))
}

func (h *Hub) join(conn *gsf.Conn, playerOID types.OID) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()

	outbound := []OutboundNotify{}
	current := h.playersByConn[conn]
	player := h.ensurePlayer(playerOID)

	// This connection is already joined as this player.
	if current == player && player.Conn == conn {
		return nil
	}

	// This connection is already linked to another player.
	if current != nil && current != player {
		outbound = append(outbound, h.removeLivePlayer(current)...)
	}

	// A player can only have one live sync connection.
	if player.Conn != nil {
		outbound = append(outbound, h.removeLivePlayer(player)...)
	}

	existing := h.playersInLocation(player.Location, nil)
	h.addLivePlayer(player, conn)

	newPlayerNotify := player.AddPlayerNotify()
	for _, other := range existing {
		outbound = append(outbound,
			NewOutboundNotify(conn, clientmessagetype.ADD_PLAYER, other.AddPlayerNotify()),
			NewOutboundNotify(other.Conn, clientmessagetype.ADD_PLAYER, newPlayerNotify),
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

	if player.Conn == nil {
		return nil
	}

	notify := player.ChangeObjectNotify()
	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, player) {
		outbound = append(outbound, NewOutboundNotify(recipient.Conn, clientmessagetype.CHANGE_OBJECT, notify))
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

	// Remember the target location even before the sync connection joins.
	if player.Conn == nil {
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
		outbound = append(outbound, NewOutboundNotify(recipient.Conn, clientmessagetype.REMOVE_PLAYER, removeNotify))
	}
	for _, other := range newRecipients {
		outbound = append(outbound,
			NewOutboundNotify(player.Conn, clientmessagetype.ADD_PLAYER, other.AddPlayerNotify()),
			NewOutboundNotify(other.Conn, clientmessagetype.ADD_PLAYER, addNotify),
		)
	}
	return outbound
}

// Move updates the latest movement state and relays it to players in the same location.
func (h *Hub) Move(conn *gsf.Conn, move *notify.Move) error {
	return SendAll(h.move(conn, move))
}

func (h *Hub) move(conn *gsf.Conn, move *notify.Move) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()

	player := h.playersByConn[conn]
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
		outbound = append(outbound, NewOutboundNotify(recipient.Conn, clientmessagetype.MOVE_PLAYER, notify))
	}
	return outbound
}

// Action relays a player action to players in the same location.
// It reports throttled when the action repeats a recent one and is not relayed.
func (h *Hub) Action(conn *gsf.Conn, action *notify.Action) (throttled bool, err error) {
	outbound, throttled := h.action(conn, action)
	return throttled, SendAll(outbound)
}

func (h *Hub) action(conn *gsf.Conn, action *notify.Action) ([]OutboundNotify, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	player := h.playersByConn[conn]
	if player == nil {
		return nil, false
	}
	if player.throttlesItemAnimate(action, time.Now()) {
		return nil, true
	}

	// TODO: keep isState actions for the ADD_PLAYER actionState of players joining later.
	notify := &types.PlayerActionNotify{
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
	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, player) {
		outbound = append(outbound, NewOutboundNotify(recipient.Conn, clientmessagetype.ACTION, notify))
	}
	return outbound, false
}

// Chat relays a chat notify to connected players in the sender's current location.
func (h *Hub) Chat(senderOID types.OID, chat *notify.Chat) error {
	return SendAll(h.chat(senderOID, chat))
}

func (h *Hub) chat(senderOID types.OID, chat *notify.Chat) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()

	player := h.playersByOID[senderOID.Int64()]
	if player == nil || player.Conn == nil {
		return nil
	}

	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, nil) {
		outbound = append(outbound, NewOutboundNotify(recipient.Conn, clientmessagetype.CHAT, chat))
	}
	return outbound
}

// Leave removes a sync connection and tells nearby players to remove that player.
func (h *Hub) Leave(conn *gsf.Conn) error {
	return SendAll(h.leave(conn))
}

func (h *Hub) leave(conn *gsf.Conn) []OutboundNotify {
	h.mu.Lock()
	defer h.mu.Unlock()
	player := h.playersByConn[conn]
	if player == nil {
		return nil
	}
	return h.removeLivePlayer(player)
}

// addLivePlayer indexes a connected sync player by connection and location.
func (h *Hub) addLivePlayer(player *Player, conn *gsf.Conn) {
	player.Conn = conn
	h.playersByConn[conn] = player
	h.addPlayerToLocation(player)
}

// removeLivePlayer removes a connected player from live indexes.
func (h *Hub) removeLivePlayer(player *Player) []OutboundNotify {
	if player.Conn == nil {
		return nil
	}

	notify := player.RemovePlayerNotify()
	outbound := make([]OutboundNotify, 0)
	for _, recipient := range h.playersInLocation(player.Location, player) {
		outbound = append(outbound, NewOutboundNotify(recipient.Conn, clientmessagetype.REMOVE_PLAYER, notify))
	}

	h.removePlayerFromLocation(player)
	delete(h.playersByConn, player.Conn)

	// Keep the player in playersByOID so location/version survive reconnect.
	player.Conn = nil
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
	if player.Conn == nil {
		return
	}
	if h.playersByLocation[player.Location] == nil {
		h.playersByLocation[player.Location] = map[int64]*Player{}
	}
	h.playersByLocation[player.Location][player.OID.Int64()] = player
}

// removePlayerFromLocation removes a connected player from its current location bucket.
func (h *Hub) removePlayerFromLocation(player *Player) {
	if player.Conn == nil {
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
