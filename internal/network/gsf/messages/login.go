package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/sessionstatus"
)

// LoginRequest is the client login payload sent to the USER server.
type LoginRequest struct {
	// Player username.
	LoginID string

	// Player password.
	Password string

	// The client always sends 1234.
	SitePIN int32

	// The client always sends 293578400718237473.
	LanguageLocalePairID types.OID

	// The client always sends "Token".
	UserQueueingToken string

	// OS, resolution, Unity version, etc.
	ClientEnvInfo types.ClientEnvironmentData

	Token     string
	LoginType int32
	CNL       string
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

// LoginResponse is the server payload that initializes the client session.
type LoginResponse struct {
	SiteInfo       types.SiteInfo
	Status         sessionstatus.SessionStatus
	SessionID      types.OID
	ConversationID int64

	// Base URL used for asset downloads.
	AssetDeliveryURL string

	// Information about the player including his active avatar.
	Player types.Player

	// Number of maximum outfit presets available (subscription based).
	MaxOutfit int16

	PlayerStats             []types.PlayerStats
	PlayerInfoTO            types.PlayerInfoTO
	CurrentServerTime       gsf.UnixTime
	SystemLockoutTime       gsf.UnixTime
	SystemShutdownTime      gsf.UnixTime
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
