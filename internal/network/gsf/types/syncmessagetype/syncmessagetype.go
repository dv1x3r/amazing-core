//go:generate go tool stringer -type=SyncMessageType

package syncmessagetype

type SyncMessageType int32

const (
	ADD_OBJECT             SyncMessageType = 1
	MOVE_OBJECT            SyncMessageType = 2
	REMOVE_OBJECT          SyncMessageType = 3
	CHANGE_OBJECT          SyncMessageType = 4
	VILLAGE_HANDOFF_QUERY  SyncMessageType = 5
	SERVER_CHANGE_OBJECT   SyncMessageType = 6
	MOVE_PLAYER            SyncMessageType = 7
	REMOVE_PLAYER          SyncMessageType = 8
	CHAT                   SyncMessageType = 9
	START_EVENT            SyncMessageType = 10
	STOP_EVENT             SyncMessageType = 11
	EMOTE                  SyncMessageType = 12
	BIND_USER_NOTIFY       SyncMessageType = 13
	PAUSE_NPC              SyncMessageType = 14
	RESUME_NPC             SyncMessageType = 15
	UPDATE_NPC_SCRIPT      SyncMessageType = 16
	NOTIFICATION           SyncMessageType = 17
	CLIENT_TEST            SyncMessageType = 18
	ADD_ITEM               SyncMessageType = 19
	EVICT                  SyncMessageType = 20
	ENTER_LOC              SyncMessageType = 21
	EXIT_LOC               SyncMessageType = 22
	ECHO                   SyncMessageType = 23
	GET_VILLAGE            SyncMessageType = 24
	REFRESH                SyncMessageType = 25
	BIND                   SyncMessageType = 26
	BIND_QUERY             SyncMessageType = 27
	BIND_VILLAGE_NOTIFY    SyncMessageType = 28
	VILLAGE_HANDOFF        SyncMessageType = 29
	FIND_SERVER            SyncMessageType = 30
	UPDATE_FILTER          SyncMessageType = 31
	SEND_NOTIFY            SyncMessageType = 32
	LOGIN                  SyncMessageType = 33
	UPDATE_USER_INFO       SyncMessageType = 34
	UPDATE_USER_FRIENDS    SyncMessageType = 35
	UPDATE_USER_GROUPS     SyncMessageType = 36
	MANAGE_USER_FRIENDS    SyncMessageType = 37
	MANAGE_USER_GROUPS     SyncMessageType = 38
	LIST_SERVERS           SyncMessageType = 39
	LIST_VILLAGES          SyncMessageType = 40
	LIST_USERS             SyncMessageType = 41
	MOVE_VILLAGE           SyncMessageType = 42
	START_NPCS             SyncMessageType = 43
	STOP_NPCS              SyncMessageType = 44
	USER_SESSION_HANDOFF   SyncMessageType = 45
	EMOTE_SVC              SyncMessageType = 46
	ACTION                 SyncMessageType = 47
	LOGOUT                 SyncMessageType = 48
	CLOSE_ZONE             SyncMessageType = 49
	TERMINATE_USER_SESSION SyncMessageType = 50
	GET_PLAYER_COUNT       SyncMessageType = 51
	UPDATE_LOCATION        SyncMessageType = 52
	HEARTBEAT_NOTIFY       SyncMessageType = 54
	RELOGIN                SyncMessageType = 55
)
