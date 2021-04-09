from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class CheckUsernameMessage(Message):
    def __init__(self):
        self.request: CheckUsernameRequest = CheckUsernameRequest()
        self.response: CheckUsernameResponse = CheckUsernameResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class CheckUsernameRequest(SerializableMessage):
    def __init__(self):
        self.username: str = None
        self.password: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.username = bit_stream.read_str()
        self.password = bit_stream.read_str()

    def to_dict(self):
        return {'username': self.username, 'password': self.password}


class CheckUsernameResponse(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        pass

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {}
