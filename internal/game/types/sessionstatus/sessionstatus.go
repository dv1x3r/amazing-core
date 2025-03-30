//go:generate go tool stringer -type=SessionStatus

package sessionstatus

type SessionStatus int32

const (
	IN_PROGRESS SessionStatus = iota
	LOGGED_OUT
	TIMED_OUT
	FORCED_LOGOUT
	LOGOUT_IN_PROGRESS
	INCORRECT_PASSWORD
	INVALID_LOGIN
	SYSTEM_ERROR
	USER_NOT_ACTIVE
	USER_UNAPPROVED
	ACCOUNT_EXPIRED
	USER_EXISTS
)

func Parse(s string) SessionStatus {
	switch s {
	case "IN_PROGRESS":
		return IN_PROGRESS
	case "LOGGED_OUT":
		return LOGGED_OUT
	case "TIMED_OUT":
		return TIMED_OUT
	case "FORCED_LOGOUT":
		return FORCED_LOGOUT
	case "LOGOUT_IN_PROGRESS":
		return LOGOUT_IN_PROGRESS
	case "INCORRECT_PASSWORD":
		return INCORRECT_PASSWORD
	case "INVALID_LOGIN":
		return INVALID_LOGIN
	case "SYSTEM_ERROR":
		return SYSTEM_ERROR
	case "USER_NOT_ACTIVE":
		return USER_NOT_ACTIVE
	case "USER_UNAPPROVED":
		return USER_UNAPPROVED
	case "ACCOUNT_EXPIRED":
		return ACCOUNT_EXPIRED
	case "USER_EXISTS":
		return USER_EXISTS
	default:
		return -1
	}
}
