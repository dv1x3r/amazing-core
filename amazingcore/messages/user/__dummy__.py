from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class DummyMessage(Message):
    def __init__(self):
        self.request: DummyRequest = DummyRequest()
        self.response: DummyResponse = DummyResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class DummyRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def __str__(self):
        return str({})


class DummyResponse(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        pass

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def __str__(self):
        return str({})
