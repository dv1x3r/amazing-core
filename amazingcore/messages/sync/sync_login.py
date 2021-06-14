from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class SyncLoginMessage(Message):

    def __init__(self):
        self.request: SyncLoginRequest = SyncLoginRequest()
        self.response: SyncLoginResponse = SyncLoginResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class SyncLoginRequest(SerializableMessage):

    def __init__(self):
        self.uid: ObjectID = None
        self.token: str = None
        self.max_vis_size: int = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.uid = ObjectID()
        self.uid.deserialize(bit_stream)
        self.token = bit_stream.read_str()
        self.max_vis_size = bit_stream.read_int()

    def to_dict(self):
        return {
            'uid': self.uid.to_dict(),
            'token': self.token,
            'max_vis_size': self.max_vis_size,
        }


class SyncLoginResponse(SerializableMessage):

    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        pass

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {}
