//go:generate go tool stringer -type=NotificationType

package notificationtype

type NotificationType int32

const (
	FRIEND_REQUEST_PENDING           NotificationType = 1
	FRIEND_REQUEST_ACCEPTED          NotificationType = 2
	FRIEND_REQUEST_DECLINED          NotificationType = 3
	REFRESH_NOTIFICATION_LIST        NotificationType = 4
	FRIEND_DELETED                   NotificationType = 5
	JOIN_CHAT_GROUP                  NotificationType = 6
	PLAYER_LEFT_CHAT_GROUP           NotificationType = 7
	NEW_PLAYER_JOINED                NotificationType = 8
	PLAYER_DECLINE_REQUEST           NotificationType = 9
	PLAYER_REJOINED                  NotificationType = 10
	SEND_GIFT                        NotificationType = 11
	ACCEPT_GIFT                      NotificationType = 12
	DECLINE_GIFT                     NotificationType = 13
	AWARD_OFFERRED                   NotificationType = 14
	QUEST_FAILED                     NotificationType = 15
	QUEST_ITEM_COLLECTED             NotificationType = 16
	SEND_QUEST_INVITE                NotificationType = 17
	ACCEPT_QUEST_INVITE              NotificationType = 18
	DECLINE_QUEST_INVITE             NotificationType = 19
	QUEST_STARTED                    NotificationType = 20
	QUEST_COMPLETED                  NotificationType = 21
	STATUS_CHANGED                   NotificationType = 22
	HOME_INVITE                      NotificationType = 23
	FORCED_LOGOUT                    NotificationType = 24
	REQUEST_EVENT                    NotificationType = 25
	ACCEPT_EVENT                     NotificationType = 26
	REJECT_EVENT                     NotificationType = 27
	SEND_VILLAGE_INVITE              NotificationType = 28
	SYSTEM_SHUTDOWN                  NotificationType = 29
	NEW_AVATAR_REGISTERED            NotificationType = 30
	SEND_GAME                        NotificationType = 31
	VILLAGE_KICK_OUT                 NotificationType = 32
	MAX_PLAYER_LEVEL                 NotificationType = 33
	MAX_NPC_LEVEL                    NotificationType = 34
	DYNAMIC_SURPRISE_AWARDED         NotificationType = 35
	GIFT_PLAYER_IMMEDIATE            NotificationType = 36
	GIFT_OTHER_IMMEDIATE             NotificationType = 37
	ACCEPT_VILLAGE_INVITE            NotificationType = 38
	REJECT_VILLAGE_INVITE            NotificationType = 39
	ACCEPT_HOME_INVITE               NotificationType = 40
	REJECT_HOME_INVITE               NotificationType = 41
	EXPIRED_GIFT                     NotificationType = 42
	MIN_PLAYERS_TO_START_REACHED     NotificationType = 43
	NPC_ACCEPTED_REQUEST             NotificationType = 44
	MEMBERSHIP_SUBSCRIPTION_AWARDED  NotificationType = 45
	CRISP_MESSAGE                    NotificationType = 46
	EVENT_IN_PROGRESS                NotificationType = 47
	EVENT_ENDED                      NotificationType = 48
	SEND_MESSAGE                     NotificationType = 49
	PLAYER_LEFT_QUEST                NotificationType = 50
	MIN_PLAYERS_TO_START_NOT_REACHED NotificationType = 51
	PLAYER_CONFIRMED_QUEST           NotificationType = 52
	OBJECT_UNPLACED                  NotificationType = 53
	MIN_PLAYERS_TO_PLAY_NOT_REACHED  NotificationType = 54
	INVENTORY_ITEM_REMOVED           NotificationType = 55
	TRANSACTION_ACCEPTED             NotificationType = 56
	TRANSACTION_COMPLETED            NotificationType = 57
	ANNOUCEMENT                      NotificationType = 58
	QUEST_MORE_PLAYERS_NEEDED        NotificationType = 59
)
