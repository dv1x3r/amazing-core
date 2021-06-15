from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class GetChatChannelTypesMessage(Message):
    def __init__(self):
        self.request: GetChatChannelTypesRequest = GetChatChannelTypesRequest()
        self.response: GetChatChannelTypesResponse = GetChatChannelTypesResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.chat_channel_types = []


class GetChatChannelTypesRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetChatChannelTypesResponse(SerializableMessage):
    def __init__(self):
        self.chat_channel_types: list = None  # GSFChatChannelType

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.chat_channel_types))
        for item in self.chat_channel_types:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'chat_channel_types': [item.to_dict() for item in self.chat_channel_types]}
