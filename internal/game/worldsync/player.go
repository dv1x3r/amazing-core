package worldsync

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// Player stores world sync state that survives reconnects and zone joins.
type Player struct {
	OID                 types.OID
	Session             *gsf.Session
	Location            int64
	PlayerDetailVersion int64
	Pos                 []types.Position
	QTH                 types.QTH
	SecondQTH           types.QTH
	Seq                 byte
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
