//go:generate go tool stringer -type=ChatChannelTypeValue

package chatchanneltypevalue

type ChatChannelTypeValue int32

const (
	PRIVATE ChatChannelTypeValue = iota
	PRIVATE_GROUP
	LOCAL
)

func Parse(s string) ChatChannelTypeValue {
	switch s {
	case "PRIVATE":
		return PRIVATE
	case "PRIVATE_GROUP":
		return PRIVATE_GROUP
	case "LOCAL":
		return LOCAL
	default:
		return -1
	}
}
