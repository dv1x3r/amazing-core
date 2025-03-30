package messages

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
	"github.com/dv1x3r/amazing-core/internal/game/types/sessionstatus"
)

type LoginRequest struct {
	LoginID              string
	Password             string
	SitePIN              int32
	LanguageLocalePairID types.OID
	UserQueueingToken    string
	ClientEnvInfo        types.ClientEnvironmentData
	Token                string
	LoginType            int32
	CNL                  string
}

func (req *LoginRequest) Deserialize(reader gsf.ProtocolReader) {
	req.LoginID = reader.ReadString()
	req.Password = reader.ReadString()
	req.SitePIN = reader.ReadInt32()
	reader.ReadObject(&req.LanguageLocalePairID)
	req.UserQueueingToken = reader.ReadString()
	reader.ReadObject(&req.ClientEnvInfo)
	req.Token = reader.ReadString()
	req.LoginType = reader.ReadInt32()
	req.CNL = reader.ReadString()
}

type LoginResponse struct {
	SiteInfo                types.SiteInfo
	Status                  sessionstatus.SessionStatus
	SessionID               types.OID
	ConversationID          int64
	AssetDeliveryURL        string
	Player                  types.Player
	MaxOutfit               int16
	PlayerStats             []types.PlayerStats
	PlayerInfoTO            types.PlayerInfoTO
	CurrentServerTime       time.Time
	SystemLockoutTime       time.Time
	SystemShutdownTime      time.Time
	ClientInactivityTimeout int32
	CNL                     string
}

func (res *LoginResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.SiteInfo)
	writer.WriteString(res.Status.String())
	writer.WriteObject(&res.SessionID)
	writer.WriteInt64(res.ConversationID)
	writer.WriteString(res.AssetDeliveryURL)
	writer.WriteObject(&res.Player)
	writer.WriteInt16(res.MaxOutfit)
	gsf.WriteSlice(writer, res.PlayerStats, func(value types.PlayerStats) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&res.PlayerInfoTO)
	writer.WriteUtcDate(res.CurrentServerTime)
	writer.WriteUtcDate(res.SystemLockoutTime)
	writer.WriteUtcDate(res.SystemShutdownTime)
	writer.WriteInt32(res.ClientInactivityTimeout)
	writer.WriteString(res.CNL)
}
