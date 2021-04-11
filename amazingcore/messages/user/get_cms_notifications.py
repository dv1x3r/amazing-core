from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.notification import Notification
from amazingcore.messages.common.notification_category import NotificationCategory


class GetCmsNotificationsMessage(Message):
    def __init__(self):
        self.request: GetCmsNotificationsRequest = GetCmsNotificationsRequest()
        self.response: GetCmsNotificationsResponse = GetCmsNotificationsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.notifications = [Notification(
            aw_object_id=ObjectID(1, 2, 3, 4),
            notification_category=NotificationCategory(
                ObjectID(1, 2, 3, 4), 1, 2),
            notification_type=3,
            requires_email=False
        )]


class GetCmsNotificationsRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetCmsNotificationsResponse(SerializableMessage):
    def __init__(self):
        self.notifications: list[Notification] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.notifications))
        for item in self.notifications:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'notifications': [item.to_dict() for item in self.notifications]}
