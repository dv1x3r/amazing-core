package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type CrispDataTO struct {
	CrispActionID   OID
	CrispMessage    string
	CrispExpiryDate time.Time
	CrispConfirmed  bool
}

func (cd *CrispDataTO) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&cd.CrispActionID)
	writer.WriteString(cd.CrispMessage)
	writer.WriteUtcDate(cd.CrispExpiryDate)
	writer.WriteBool(cd.CrispConfirmed)
}

func (cd *CrispDataTO) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&cd.CrispActionID)
	cd.CrispMessage = reader.ReadString()
	cd.CrispExpiryDate = reader.ReadUtcDate()
	cd.CrispConfirmed = reader.ReadBool()
}
