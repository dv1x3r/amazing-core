package notify

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// InterestList is the players to add, change and remove nearby, always empty.
type InterestList struct{}

func (n *InterestList) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(0) // adds
	writer.WriteInt32(0) // changes
	writer.WriteInt32(0) // removes
}
