package worldsync

import (
	"maps"
	"slices"
	"time"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

const (
	// actionItemAnimate is SyncManager.ActionNotifyType.ItemAnimate.
	actionItemAnimate = 3

	// itemAnimateThrottle drops repeats of the same ItemAnimate, level objects can send it every frame.
	itemAnimateThrottle = 250 * time.Millisecond

	// noLocation is the location of a player that has not entered one, the OIDs are not negative.
	noLocation = -1
)

// Player is the player of a session, it keeps its state between SYNC reconnects.
type Player struct {
	OID types.OID

	// syncConn is nil when the player is not in the world.
	syncConn *gsf.Conn

	// location is noLocation until ENTER_BUILDING and while the player travels.
	location int64

	// detailVersion is the version of the player avatar and outfit, the clients reload them
	// with GET_OTHER_PLAYER_DETAILS when ADD_PLAYER or CHANGE_OBJECT has another version.
	detailVersion int64

	pos       []types.Position
	qth       types.QTH
	secondQTH types.QTH
	seq       byte

	// states are the last state actions by stateIndex, ADD_PLAYER sends them to players who join later.
	states map[int32]types.PlayerActionNotify

	// itemAnimates is the last ItemAnimate time by object, for the throttle.
	itemAnimates map[string]time.Time
}

// setState saves a state action: clearState clears its stateIndex slot, isState sets it.
func (p *Player) setState(action types.PlayerActionNotify) {
	if action.ClearState {
		delete(p.states, action.StateIndex)
	}
	if action.IsState {
		if p.states == nil {
			p.states = map[int32]types.PlayerActionNotify{}
		}
		p.states[action.StateIndex] = action
	}
}

// actionState returns the state actions for ADD_PLAYER.
func (p *Player) actionState() []types.PlayerActionNotify {
	return slices.Collect(maps.Values(p.states))
}

// throttlesItemAnimate reports whether the same ItemAnimate was sent less than itemAnimateThrottle ago.
func (p *Player) throttlesItemAnimate(action *notify.Action, now time.Time) bool {
	// State actions, like sitting, are never throttled.
	if action.ActionType != actionItemAnimate || action.IsState {
		return false
	}
	if last, ok := p.itemAnimates[action.Parameters]; ok && now.Sub(last) < itemAnimateThrottle {
		return true
	}
	if p.itemAnimates == nil {
		p.itemAnimates = map[string]time.Time{}
	}
	for object, last := range p.itemAnimates {
		if now.Sub(last) >= itemAnimateThrottle {
			delete(p.itemAnimates, object)
		}
	}
	p.itemAnimates[action.Parameters] = now
	return false
}

// addPlayerNotify is the ADD_PLAYER payload.
func (p *Player) addPlayerNotify() *notify.AddPlayer {
	return &notify.AddPlayer{
		POID:              p.OID,
		PlayerVillagerOID: types.OID{},
		Ver:               p.detailVersion,
		LOID:              types.OIDFromInt64(p.location),
		LCP:               false,
		TimeOffset:        0,
		Pos:               p.firstPosition(),
		WPos:              nil,
		QTH:               p.qth,
		SecondQTH:         p.secondQTH,
		Weight:            0,
		Seq:               p.seq,
		Type:              0,
		ActionState:       p.actionState(),
	}
}

// changeObjectNotify is the CHANGE_OBJECT payload, the clients reload the player details.
func (p *Player) changeObjectNotify() *notify.ChangeObject {
	return &notify.ChangeObject{
		InstanceOID: p.OID,
		TemplateOID: types.OID{},
		PlayerOID:   p.OID,
		VillageOID:  types.OID{},
		LocationOID: types.OIDFromInt64(p.location),
		Ver:         int32(p.detailVersion),
		Pos:         p.firstPosition(),
		QTH:         p.qth,
		State:       nil,
	}
}

// playerMoveNotify is the MOVE_PLAYER payload.
func (p *Player) playerMoveNotify() *notify.PlayerMove {
	return &notify.PlayerMove{
		POID:              p.OID,
		PlayerVillagerOID: types.OID{},
		LCP:               false,
		Pos:               p.pos,
		QTH:               p.qth,
		SecondQTH:         p.secondQTH,
		Seq:               p.seq,
	}
}

// removePlayerNotify is the REMOVE_PLAYER payload.
func (p *Player) removePlayerNotify() *notify.RemovePlayer {
	return &notify.RemovePlayer{
		POID:              p.OID,
		PlayerVillagerOID: types.OID{},
		LCP:               false,
		Type:              0,
	}
}

// firstPosition returns the first position, or zero before ENTER_BUILDING.
func (p *Player) firstPosition() types.Position {
	if len(p.pos) == 0 {
		return types.Position{}
	}
	return p.pos[0]
}
