//go:generate go tool stringer -type=QuestState

package queststate

type QuestState int

const (
	AVAILABLE        QuestState = 1
	ACCEPTED         QuestState = 2
	COMPLETED        QuestState = 3
	CLOSED           QuestState = 4
	CLOSED_FAIL      QuestState = 5
	FAILED           QuestState = 6
	STARTED          QuestState = 7
	REJECTED         QuestState = 8
	CLOSED_ABANDONED QuestState = 9
)
