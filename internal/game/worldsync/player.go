package worldsync

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

const (
	// actionItemAnimate is SyncManager.ActionNotifyType.ItemAnimate.
	actionItemAnimate = 3

	// itemAnimateThrottle is the minimum interval between relays of a decorative ItemAnimate of the same object.
	// Level object triggers can fire every frame, while the receiver only restarts the same animation clip.
	itemAnimateThrottle = 250 * time.Millisecond
)

// Player stores world sync state that survives reconnects and zone joins.
type Player struct {
	OID                 types.OID
	Conn                *gsf.Conn
	Location            int64
	PlayerDetailVersion int64
	Pos                 []types.Position
	QTH                 types.QTH
	SecondQTH           types.QTH
	Seq                 byte

	// itemAnimates holds the last relay time of decorative ItemAnimate actions by the object path hash.
	itemAnimates map[string]time.Time
}

// throttlesItemAnimate reports whether a decorative ItemAnimate of the same object was relayed less than itemAnimateThrottle ago.
// Actions with isState (stopPlayerMovement objects such as seats) are never throttled.
// The position is not compared, the receiver finds the object by the path hash only.
func (p *Player) throttlesItemAnimate(action *notify.Action, now time.Time) bool {
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

// AddPlayerNotify builds the payload used when this player becomes visible.
func (p *Player) AddPlayerNotify() *notify.AddPlayer {
	return &notify.AddPlayer{
		POID:              p.OID,
		PlayerVillagerOID: types.OID{},
		Ver:               p.PlayerDetailVersion,
		LOID:              types.OIDFromInt64(p.Location),
		LCP:               false,
		TimeOffset:        0,
		Pos:               p.firstPosition(),
		WPos:              nil,
		QTH:               p.QTH,
		SecondQTH:         p.SecondQTH,
		Weight:            0,
		Seq:               p.Seq,
		Type:              0,
		ActionState:       nil,
	}
}

// ChangeObjectNotify builds the payload used to invalidate cached player details.
func (p *Player) ChangeObjectNotify() *notify.ChangeObject {
	return &notify.ChangeObject{
		InstanceOID: p.OID,
		TemplateOID: types.OID{},
		PlayerOID:   p.OID,
		VillageOID:  types.OID{},
		LocationOID: types.OIDFromInt64(p.Location),
		Ver:         int32(p.PlayerDetailVersion),
		Pos:         p.firstPosition(),
		QTH:         p.QTH,
		State:       nil,
	}
}

// PlayerMoveNotify builds the payload used to relay this player's latest movement.
func (p *Player) PlayerMoveNotify() *notify.PlayerMove {
	return &notify.PlayerMove{
		POID:              p.OID,
		PlayerVillagerOID: types.OID{},
		LCP:               false,
		Pos:               append([]types.Position(nil), p.Pos...),
		QTH:               p.QTH,
		SecondQTH:         p.SecondQTH,
		Seq:               p.Seq,
	}
}

// RemovePlayerNotify builds the payload used when this player is no longer visible.
func (p *Player) RemovePlayerNotify() *notify.RemovePlayer {
	return &notify.RemovePlayer{
		POID:              p.OID,
		PlayerVillagerOID: types.OID{},
		LCP:               false,
		Type:              0,
	}
}

// firstPosition returns the latest primary position, or zero when no movement arrived yet.
func (p *Player) firstPosition() types.Position {
	if len(p.Pos) == 0 {
		return types.Position{}
	}
	return p.Pos[0]
}
