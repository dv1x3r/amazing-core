from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.notification_category import NotificationCategory


class Notification(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 notification_category: NotificationCategory = None,
                 notification_type: int = None,
                 requires_email: bool = None):
        self.aw_object_id = aw_object_id
        self.notification_category = notification_category
        self.notification_type = notification_type
        self.requires_email = requires_email

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        self.notification_category.serialize(bit_stream)
        bit_stream.write_int(self.notification_type)
        bit_stream.write_bool(self.requires_email)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'notification_category': self.notification_category.to_dict(),
            'notification_type': self.notification_type,
            'requires_email': self.requires_email,
        }
