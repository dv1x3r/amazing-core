from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.notification import Notification
from amazingcore.messages.common.notification_category import NotificationCategory
from amazingcore.messages.common.player_notification import PlayerNotification

import datetime as dt


class GetNotificationByPlayerIdMessage(Message):
    def __init__(self):
        self.request: GetNotificationByPlayerIdRequest = GetNotificationByPlayerIdRequest()
        self.response: GetNotificationByPlayerIdResponse = GetNotificationByPlayerIdResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        notification = Notification(
            aw_object_id=ObjectID(1, 2, 3, 4),
            notification_category=NotificationCategory(
                ObjectID(1, 2, 3, 4), 1, 2),
            notification_type=3,
            requires_email=False
        )

        self.response.player_notifications = [PlayerNotification(
            aw_object_id=ObjectID(1, 2, 3, 4),
            from_player_id=ObjectID(1, 2, 3, 4),
            from_player_name='from_player_name',
            to_player_id=ObjectID(1, 2, 3, 4),
            village_name='from_player_name',
            village_id=ObjectID(1, 2, 3, 4),
            object_id=ObjectID(1, 2, 3, 4),
            notification_text='notification_text',
            notification=notification,
            is_read=False,
            create_date=dt.datetime.now(),
            expiry_date=(dt.datetime.now() + dt.timedelta(days=1))
        )]


class GetNotificationByPlayerIdRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetNotificationByPlayerIdResponse(SerializableMessage):
    def __init__(self):
        self.player_notifications: list[PlayerNotification] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_notifications))
        for item in self.player_notifications:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'player_notifications': [item.to_dict() for item in self.player_notifications]}
