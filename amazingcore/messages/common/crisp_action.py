from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class CrispAction(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 code: int = None,
                 name: str = None,
                 chat_send_allowed: bool = None,
                 chat_receive_allowed: bool = None,
                 login_allowed: bool = None,
                 ban_length: int = None):
        self.aw_object_id = aw_object_id
        self.code = code
        self.name = name
        self.chat_send_allowed = chat_send_allowed
        self.chat_receive_allowed = chat_receive_allowed
        self.login_allowed = login_allowed
        self.ban_length = ban_length

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_int(self.code)
        bit_stream.write_str(self.name)
        bit_stream.write_bool(self.chat_send_allowed)
        bit_stream.write_bool(self.chat_receive_allowed)
        bit_stream.write_bool(self.login_allowed)
        bit_stream.write_int(self.ban_length)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'code': self.code,
            'name': self.name,
            'chat_send_allowed': self.chat_send_allowed,
            'chat_receive_allowed': self.chat_receive_allowed,
            'login_allowed': self.login_allowed,
            'ban_length': self.ban_length,
        }
