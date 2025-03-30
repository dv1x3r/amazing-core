//go:generate go tool stringer -type=ServiceClass

package serviceclass

type ServiceClass int32

const (
	USER_SERVER ServiceClass = 18
	SYNC_SERVER ServiceClass = 19
	LOCATION    ServiceClass = 20
	CLIENT      ServiceClass = -1
)
