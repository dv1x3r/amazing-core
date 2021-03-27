from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class ObjectID(SerializableMessage):
    def __init__(self):
        self.object_class: int = None
        self.object_type: int = None
        self.server: int = None
        self.object_number: int = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(self.object_class)
        bit_stream.write_int(self.object_type)
        bit_stream.write_int(self.server)
        bit_stream.write_long(self.object_number)

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.object_class = bit_stream.read_int()
        self.object_type = bit_stream.read_int()
        self.server = bit_stream.read_int()
        self.object_number = bit_stream.read_long()

    def to_dict(self):
        return {
            'object_class': self.object_class,
            'object_type': self.object_type,
            'server': self.server,
            'object_number': self.object_number,
        }
