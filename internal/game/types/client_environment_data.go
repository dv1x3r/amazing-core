package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type ClientEnvironmentData struct {
	UnityVersion       string
	UserAgent          string
	ScreenResolution   string
	MachineOS          string
	UserTime           time.Time
	UtcOffsetInMinutes int32
	IpAddress          string
}

func (ce *ClientEnvironmentData) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(ce.UnityVersion)
	writer.WriteString(ce.UserAgent)
	writer.WriteString(ce.ScreenResolution)
	writer.WriteString(ce.MachineOS)
	writer.WriteUtcDate(ce.UserTime)
	writer.WriteInt32(ce.UtcOffsetInMinutes)
	writer.WriteString(ce.IpAddress)
}

func (ce *ClientEnvironmentData) Deserialize(reader gsf.ProtocolReader) {
	ce.UnityVersion = reader.ReadString()
	ce.UserAgent = reader.ReadString()
	ce.ScreenResolution = reader.ReadString()
	ce.MachineOS = reader.ReadString()
	ce.UserTime = reader.ReadUtcDate()
	ce.UtcOffsetInMinutes = reader.ReadInt32()
	ce.IpAddress = reader.ReadString()
}
