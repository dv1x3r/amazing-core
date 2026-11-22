package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetCurrenciesRequest is empty.
type GetCurrenciesRequest struct{}

func (req *GetCurrenciesRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetCurrenciesResponse contains the currencies.
type GetCurrenciesResponse struct {
	Currencies []types.Currency
}

func (res *GetCurrenciesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Currencies, func(value types.Currency) {
		writer.WriteObject(&value)
	})
}
