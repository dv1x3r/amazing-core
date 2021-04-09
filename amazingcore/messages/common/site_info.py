from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class SiteInfo(SerializableMessage):
    def __init__(self,
                 site_id: ObjectID = None,
                 nickname_first: str = None,
                 nickname_last: str = None,
                 site_user_id: ObjectID = None):
        self.site_id = site_id
        self.nickname_first = nickname_first
        self.nickname_last = nickname_last
        self.site_user_id = site_user_id

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.site_id.serialize(bit_stream)
        bit_stream.write_str(self.nickname_first)
        bit_stream.write_str(self.nickname_last)
        self.site_user_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'site_id': self.site_id.to_dict(),
            'nickname_first': self.nickname_first,
            'nickname_last': self.nickname_last,
            'site_user_id': self.site_user_id.to_dict(),
        }
