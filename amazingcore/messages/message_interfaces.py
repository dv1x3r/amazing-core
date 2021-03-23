from abc import ABC, abstractmethod
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.message_header import MessageHeader


class SerializableMessage(ABC):
    @abstractmethod
    def serialize(self, bit_stream: BitStream):
        pass

    @abstractmethod
    def deserialize(self, bit_stream: BitStream):
        pass

    @abstractmethod
    def to_dict(self):
        pass


class Message(ABC):
    request: SerializableMessage
    response: SerializableMessage

    @abstractmethod
    async def process(self, message_header: MessageHeader):
        pass
