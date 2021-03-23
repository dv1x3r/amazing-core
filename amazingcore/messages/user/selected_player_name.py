from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class SelectedPlayerNameMessage(Message):
    def __init__(self):
        self.request: SelectedPlayerNameRequest = SelectedPlayerNameRequest()
        self.response: SelectedPlayerNameResponse = SelectedPlayerNameResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class SelectedPlayerNameRequest(SerializableMessage):
    def __init__(self):
        self.family_name: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        bit_stream.read_start()
        self.family_name = bit_stream.read_str()

    def to_dict(self):
        return {'family_name': self.family_name}


class SelectedPlayerNameResponse(SerializableMessage):

    def serialize(self, bit_stream: BitStream):
        pass  # only header needed

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        pass
