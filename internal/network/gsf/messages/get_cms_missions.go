package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetCMSMissionsRequest contains the quest type of the missions.
type GetCMSMissionsRequest struct {
	QuestTypeOID types.OID
	Hierarchies  []string
	ReturnLocked bool
}

func (req *GetCMSMissionsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.QuestTypeOID)
	req.Hierarchies = gsf.ReadSlice(reader, reader.ReadString)
	req.ReturnLocked = reader.ReadBool()
}

// GetCMSMissionsResponse contains the mission definitions.
type GetCMSMissionsResponse struct {
	Missions []types.Quest
}

func (res *GetCMSMissionsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Missions, func(value types.Quest) {
		writer.WriteObject(&value)
	})
}
