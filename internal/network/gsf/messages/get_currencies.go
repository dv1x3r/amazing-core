package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetCurrenciesRequest requests the currency definitions.
type GetCurrenciesRequest struct{}

func (req *GetCurrenciesRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetCurrenciesResponse contains the currency definitions, the client reads the currency name from its Text asset.
type GetCurrenciesResponse struct {
	Currencies []types.Currency
}

func (res *GetCurrenciesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Currencies, func(value types.Currency) {
		writer.WriteObject(&value)
	})
}
