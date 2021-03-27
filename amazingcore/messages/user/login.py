from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.client_environment import ClientEnvironment


class LoginMessage(Message):
    def __init__(self):
        self.request: LoginRequest = LoginRequest()
        self.response: LoginResponse = LoginResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class LoginRequest(SerializableMessage):
    def __init__(self):
        self.login_id: str = None
        self.password: str = None
        self.site_pin: int = None
        self.language_local_pair_id: ObjectID = None
        self.user_queueing_token: str = None
        self.client_environment: ClientEnvironment = None
        self.token: str = None
        self.login_type: int = None
        self.cnl: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        bit_stream.read_start()
        self.login_id = bit_stream.read_str()
        self.password = bit_stream.read_str()
        self.site_pin = bit_stream.read_int()
        self.language_local_pair_id = ObjectID()
        self.language_local_pair_id.deserialize(bit_stream)
        self.user_queueing_token = bit_stream.read_str()
        # self.client_environment = ClientEnvironment()
        # self.client_environment.deserialize(bit_stream)
        # self.token = bit_stream.read_str()
        # self.login_type = bit_stream.read_int()
        # self.cnl = bit_stream.read_str()

    def to_dict(self):
        return {
            'login_id': self.login_id,
            'password': self.password,
            'site_pin': self.site_pin,
            'language_local_pair_id': self.language_local_pair_id.to_dict(),
            'user_queueing_token': self.user_queueing_token,
        }


class LoginResponse(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        pass

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {}
