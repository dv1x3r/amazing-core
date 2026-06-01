package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AcceptQuestRequest requests a quest acceptance.
type AcceptQuestRequest struct {
	PlayerQuestOID types.OID
	QuestOID       types.OID
	LocationOID    types.OID
}

func (req *AcceptQuestRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerQuestOID)
	reader.ReadObject(&req.QuestOID)
	reader.ReadObject(&req.LocationOID)
}

// AcceptQuestResponse contains accepted quest status.
type AcceptQuestResponse struct {
	IsAccepted    bool
	CanStart      bool
	CanPlay       bool
	JoinedPlayers []types.OID
}

func (res *AcceptQuestResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsAccepted)
	writer.WriteBool(res.CanStart)
	writer.WriteBool(res.CanPlay)
	gsf.WriteSlice(writer, res.JoinedPlayers, func(value types.OID) {
		writer.WriteObject(&value)
	})
}
