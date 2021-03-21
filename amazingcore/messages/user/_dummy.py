from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class DummyMessage(Message):
    def __init__(self):
        self.request: DummyRequest = DummyRequest()
        self.response: DummyResponse = DummyResponse()

    async def process(self):
        pass


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
