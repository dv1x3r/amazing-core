//go:generate go tool stringer -type=ResultCode

package resultcode

type ResultCode int32

const (
	INCOMPLETE   ResultCode = -1
	OK           ResultCode = 0
	APP          ResultCode = 5
	APP_DB       ResultCode = 6
	ERR          ResultCode = 10
	QUEUE        ResultCode = 11
	DB           ResultCode = 20
	DB_QUEUE     ResultCode = 21
	DB_NO_RETRY  ResultCode = 22
	NO_MEM       ResultCode = 30
	COMM         ResultCode = 40
	CONN_FAILED  ResultCode = 41
	DISCONNECT   ResultCode = 42
	SHUTDOWN     ResultCode = 43
	IO           ResultCode = 44
	TIMEOUT      ResultCode = 45
	BUSY         ResultCode = 46
	COMM_INIT    ResultCode = 47
	CANCEL       ResultCode = 48
	WOULD_BLOCK  ResultCode = 49
	PROTOCOL_VER ResultCode = 50
	SERIALIZE    ResultCode = 51
	PENDING_IO   ResultCode = 52
	ASYNC        ResultCode = 53
	CONN         ResultCode = 54
	CHNL_CLOSED  ResultCode = 55
	CONN_EXHST   ResultCode = 56
	NO_DEST      ResultCode = 57
	CHRONO       ResultCode = 58
	NOT_READY    ResultCode = 59
)
