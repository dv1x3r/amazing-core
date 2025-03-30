//go:generate go tool stringer -type=OnlineStatusValue

package onlinestatusvalue

type OnlineStatusValue int32

const (
	OFFLINE OnlineStatusValue = iota
	APPEAR_OFFLINE
	ONLINE
	BUSY
)

func Parse(s string) OnlineStatusValue {
	switch s {
	case "OFFLINE":
		return OFFLINE
	case "APPEAR_OFFLINE":
		return APPEAR_OFFLINE
	case "ONLINE":
		return ONLINE
	case "BUSY":
		return BUSY
	default:
		return -1
	}
}
