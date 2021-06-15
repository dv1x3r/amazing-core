from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class Orientation(SerializableMessage):
    def __init__(self,
                 sw: bool = None,
                 sx: bool = None,
                 sy: bool = None,
                 sz: bool = None,
                 cx: int = None,
                 cy: int = None,
                 cz: int = None):
        self.sw = sw
        self.sx = sx
        self.sy = sy
        self.sz = sz
        self.cx = cx
        self.cy = cy
        self.cz = cz

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_bool(self.sw)
        bit_stream.write_bool(self.sx)
        bit_stream.write_bool(self.sy)
        bit_stream.write_bool(self.sz)
        bit_stream.write_short(self.cx)
        bit_stream.write_short(self.cy)
        bit_stream.write_short(self.cz)

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.sw = bit_stream.read_bool()
        self.sx = bit_stream.read_bool()
        self.sy = bit_stream.read_bool()
        self.sz = bit_stream.read_bool()
        self.cx = bit_stream.read_short()
        self.cy = bit_stream.read_short()
        self.cz = bit_stream.read_short()

    def to_dict(self):
        return {
            'sw': self.sw,
            'sx': self.sx,
            'sy': self.sy,
            'sz': self.sz,
            'cx': self.cx,
            'cy': self.cy,
            'cz': self.cz,
        }
