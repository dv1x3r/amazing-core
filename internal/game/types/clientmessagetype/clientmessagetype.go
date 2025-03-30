//go:generate go tool stringer -type=ClientMessageType

package clientmessagetype

type ClientMessageType int32

const (
	ADD_OBJECT           ClientMessageType = 1
	MOVE_OBJECT          ClientMessageType = 2
	REMOVE_OBJECT        ClientMessageType = 3
	CHANGE_OBJECT        ClientMessageType = 4
	SERVER_CHANGE_OBJECT ClientMessageType = 5
	MOVE_PLAYER          ClientMessageType = 6
	REMOVE_PLAYER        ClientMessageType = 7
	CHAT                 ClientMessageType = 8
	START_EVENT          ClientMessageType = 9
	STOP_EVENT           ClientMessageType = 10
	EMOTE                ClientMessageType = 11
	CHANGE_WEIGHT        ClientMessageType = 12
	PAUSE                ClientMessageType = 13
	RESUME               ClientMessageType = 14
	NOTIFICATION         ClientMessageType = 16
	CHANGE_SERVER        ClientMessageType = 17
	SEND_NOTIFY          ClientMessageType = 18
	INT_LIST             ClientMessageType = 19
	ADD_PLAYER           ClientMessageType = 20
	UPDATE_NPCS          ClientMessageType = 21
	STOP_NPC             ClientMessageType = 22
	MINIMAP              ClientMessageType = 23
	POS_RECAP            ClientMessageType = 24
	ACTION               ClientMessageType = 25
	EVICT                ClientMessageType = 26
	ONLINE_STATUS        ClientMessageType = 27
	CHANGE_OBJECT_STATE  ClientMessageType = 28
)
