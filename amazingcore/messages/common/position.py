from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class Position(SerializableMessage):
    def __init__(self,
                 x: int = None,
                 y: int = None,
                 z: int = None,
                 t: int = None):
        self.x = x
        self.y = y
        self.z = z
        self.t = t

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(self.x)
        bit_stream.write_int(self.y)
        bit_stream.write_int(self.z)
        bit_stream.write_int(self.t)

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.x = bit_stream.read_int()
        self.y = bit_stream.read_int()
        self.z = bit_stream.read_int()
        self.t = bit_stream.read_int()

    def to_dict(self):
        return {
            'x': self.x,
            'y': self.y,
            'z': self.z,
            't': self.t,
        }
