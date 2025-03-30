package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type VillagePlot struct {
	OID        OID
	CreateDate time.Time
	PlotNo     gsf.Null[int32]
	VillageID  OID
	PlotTypeID OID
	PlayerID   OID
	IsStore    bool
}

func (vp *VillagePlot) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&vp.OID)
	writer.WriteUtcDate(vp.CreateDate)
	gsf.WriteNullable(writer, vp.PlotNo, writer.WriteInt32)
	writer.WriteObject(&vp.VillageID)
	writer.WriteObject(&vp.PlotTypeID)
	writer.WriteObject(&vp.PlayerID)
	writer.WriteBool(vp.IsStore)
}

func (vp *VillagePlot) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&vp.OID)
	vp.CreateDate = reader.ReadUtcDate()
	vp.PlotNo = gsf.ReadNullable(reader, reader.ReadInt32)
	reader.ReadObject(&vp.VillageID)
	reader.ReadObject(&vp.PlotTypeID)
	reader.ReadObject(&vp.PlayerID)
	vp.IsStore = reader.ReadBool()
}
