from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class SelectedPlayerNameMessage(Message):
    def __init__(self):
        self.request: SelectedPlayerNameRequest = SelectedPlayerNameRequest()
        self.response = None

    async def process(self):
        pass


class SelectedPlayerNameRequest(SerializableMessage):
    def __init__(self):
        self.family_name: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        bit_stream.read_start()
        self.family_name = bit_stream.read_str()

    def __str__(self):
        return str({'family_name': self.family_name})
