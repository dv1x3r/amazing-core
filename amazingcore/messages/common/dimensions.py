from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class Dimensions(SerializableMessage):
    def __init__(self,
                 cx: int = None,
                 cy: int = None,
                 cz: int = None,
                 box: int = None,
                 boy: int = None,
                 boz: int = None):
        self.cx = cx
        self.cy = cy
        self.cz = cz
        self.box = box
        self.boy = boy
        self.boz = boz

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(self.cx)
        bit_stream.write_int(self.cy)
        bit_stream.write_int(self.cz)
        bit_stream.write_int(self.box)
        bit_stream.write_int(self.boy)
        bit_stream.write_int(self.boz)

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.cx = bit_stream.read_int()
        self.cy = bit_stream.read_int()
        self.cz = bit_stream.read_int()
        self.box = bit_stream.read_int()
        self.boy = bit_stream.read_int()
        self.boz = bit_stream.read_int()

    def to_dict(self):
        return {
            'xc': self.cx,
            'cy': self.cy,
            'cz': self.cz,
            'box': self.box,
            'boy': self.boy,
            'boz': self.boz,
        }
