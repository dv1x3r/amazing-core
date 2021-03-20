from abc import ABC, abstractmethod
from amazingcore.codec.bit_stream import BitStream


class SerializableMessage(ABC):
    @abstractmethod
    def serialize(self, bit_stream: BitStream):
        pass

    @abstractmethod
    def deserialize(self, bit_stream: BitStream):
        pass


class Message(ABC):
    request: SerializableMessage
    response: SerializableMessage

    @abstractmethod
    async def process(self):
        pass
