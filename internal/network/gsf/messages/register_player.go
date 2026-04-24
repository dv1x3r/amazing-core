package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

type RegisterPlayerRequest struct {
	Token               string
	Password            string
	ParentEmailAddress  string
	BirthDate           gsf.UnixTime
	Gender              string
	LocationID          types.OID
	Username            string
	Worldname           string
	ChatAllowed         bool
	CNL                 string
	ReferredByWorldname string
	LoginType           int32
}

func (req *RegisterPlayerRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Token = reader.ReadString()
	req.Password = reader.ReadString()
	req.ParentEmailAddress = reader.ReadString()
	req.BirthDate = reader.ReadUtcDate()
	req.Gender = reader.ReadString()
	reader.ReadObject(&req.LocationID)
	req.Username = reader.ReadString()
	req.Worldname = reader.ReadString()
	req.ChatAllowed = reader.ReadBool()
	req.CNL = reader.ReadString()
	req.ReferredByWorldname = reader.ReadString()
	req.LoginType = reader.ReadInt32()
}

type RegisterPlayerResponse struct {
	PlayerID types.OID
}

func (res *RegisterPlayerResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.PlayerID)
}
