from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class RegisterPlayerMessage(Message):
    def __init__(self):
        self.request: RegisterPlayerRequest = RegisterPlayerRequest()
        self.response: RegisterPlayerResponse = RegisterPlayerResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK
        self.response.player_id = ObjectID(1, 2, 3, 4)


class RegisterPlayerRequest(SerializableMessage):
    def __init__(self):
        self.token: str = None
        self.password: str = None
        self.parent_email_address: str = None
        self.birth_date: dt.datetime = None
        self.gender: str = None
        self.location_id: ObjectID = None
        self.username: str = None
        self.worldname: str = None
        self.chat_allowed: bool = None
        self.cnl: str = None
        self.referred_by_worldname: str = None
        self.login_type: int = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.token = bit_stream.read_str()
        self.password = bit_stream.read_str()
        self.parent_email_address = bit_stream.read_str()
        self.birth_date = bit_stream.read_dt()
        self.gender = bit_stream.read_str()
        self.location_id = ObjectID()
        self.location_id.deserialize(bit_stream)
        self.username = bit_stream.read_str()
        self.worldname = bit_stream.read_str()
        self.chat_allowed = bit_stream.read_bool()
        self.cnl = bit_stream.read_str()
        self.referred_by_worldname = bit_stream.read_str()
        self.login_type = bit_stream.read_int()

    def to_dict(self):
        return {
            'token': self.token,
            'password': self.password,
            'parent_email_address': self.parent_email_address,
            'birth_date': self.birth_date,
            'gender': self.gender,
            'location_id': self.location_id.to_dict(),
            'username': self.username,
            'worldname': self.worldname,
            'chat_allowed': self.chat_allowed,
            'cnl': self.cnl,
            'referred_by_worldname': self.referred_by_worldname,
            'login_type': self.login_type,
        }


class RegisterPlayerResponse(SerializableMessage):
    def __init__(self, player_id: ObjectID = None):
        self.player_id = player_id

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.player_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'player_id': self.player_id.to_dict()}
