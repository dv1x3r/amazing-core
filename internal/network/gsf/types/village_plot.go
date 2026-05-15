package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// VillagePlot is a plot assignment within a village.
type VillagePlot struct {
	OID         OID
	CreateDate  gsf.UnixTime
	PlotNo      gsf.Null[int32]
	VillageOID  OID
	PlotTypeOID OID
	PlayerOID   OID
	IsStore     bool
}

func (vp *VillagePlot) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&vp.OID)
	writer.WriteUtcDate(vp.CreateDate)
	gsf.WriteNullable(writer, vp.PlotNo, writer.WriteInt32)
	writer.WriteObject(&vp.VillageOID)
	writer.WriteObject(&vp.PlotTypeOID)
	writer.WriteObject(&vp.PlayerOID)
	writer.WriteBool(vp.IsStore)
}

func (vp *VillagePlot) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&vp.OID)
	vp.CreateDate = reader.ReadUtcDate()
	vp.PlotNo = gsf.ReadNullable(reader, reader.ReadInt32)
	reader.ReadObject(&vp.VillageOID)
	reader.ReadObject(&vp.PlotTypeOID)
	reader.ReadObject(&vp.PlayerOID)
	vp.IsStore = reader.ReadBool()
}
