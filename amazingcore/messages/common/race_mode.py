from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class RaceMode(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 name: str = None):
        self.aw_object_id = aw_object_id
        self.name = name

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_str(self.name)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'name': self.name,
        }
