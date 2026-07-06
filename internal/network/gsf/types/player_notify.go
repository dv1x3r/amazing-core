package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerNotify is a queued player notification returned by HEARTBEAT.
type PlayerNotify struct {
	FromPlayerOID   OID
	ToPlayerOID     OID
	NotificationOID OID
	ObjectOID       OID
	ObjectOID2      OID
	EndDate         gsf.UnixTime
	Source          string
	SendingObjectID OID
	Quantity        int32
}

func (pn *PlayerNotify) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pn.FromPlayerOID)
	writer.WriteObject(&pn.ToPlayerOID)
	writer.WriteObject(&pn.NotificationOID)
	writer.WriteObject(&pn.ObjectOID)
	writer.WriteObject(&pn.ObjectOID2)
	writer.WriteUtcDate(pn.EndDate)
	writer.WriteString(pn.Source)
	writer.WriteObject(&pn.SendingObjectID)
	writer.WriteInt32(pn.Quantity)
}

func (pn *PlayerNotify) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pn.FromPlayerOID)
	reader.ReadObject(&pn.ToPlayerOID)
	reader.ReadObject(&pn.NotificationOID)
	reader.ReadObject(&pn.ObjectOID)
	reader.ReadObject(&pn.ObjectOID2)
	pn.EndDate = reader.ReadUtcDate()
	pn.Source = reader.ReadString()
	reader.ReadObject(&pn.SendingObjectID)
	pn.Quantity = reader.ReadInt32()
}
