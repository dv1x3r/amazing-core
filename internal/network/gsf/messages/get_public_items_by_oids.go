package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

type GetPublicItemsByOIDsRequest struct {
	OIDs             []types.OID
	LangLocalePairID types.OID
	TierID           types.OID
	BirthDate        gsf.UnixTime
	RegistrationDate gsf.UnixTime
	PreviewDate      gsf.UnixTime
	IsPreviewEnabled bool
}

func (req *GetPublicItemsByOIDsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.OIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
	reader.ReadObject(&req.LangLocalePairID)
	reader.ReadObject(&req.TierID)
	req.BirthDate = reader.ReadUtcDate()
	req.RegistrationDate = reader.ReadUtcDate()
	req.PreviewDate = reader.ReadUtcDate()
	req.IsPreviewEnabled = reader.ReadBool()
}

type GetPublicItemsByOIDsResponse struct {
	Items []types.Item
}

func (res *GetPublicItemsByOIDsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Items, func(value types.Item) {
		writer.WriteObject(&value)
	})
}
