from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class NotificationCategory(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 ordinal: int = None,
                 default_expiry_days: int = None):
        self.aw_object_id = aw_object_id
        self.ordinal = ordinal
        self.default_expiry_days = default_expiry_days

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_int(self.ordinal)
        bit_stream.write_int(self.default_expiry_days)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'ordinal': self.ordinal,
            'default_expiry_days': self.default_expiry_days,
        }
