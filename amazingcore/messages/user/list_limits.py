from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class ListLimitsMessage(Message):
    def __init__(self):
        self.request: ListLimitsRequest = ListLimitsRequest()
        self.response: ListLimitsResponse = ListLimitsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.limit_values = []


class ListLimitsRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class ListLimitsResponse(SerializableMessage):
    def __init__(self):
        self.limit_values = None  # GSFLimit

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.limit_values))
        for item in self.limit_values:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'limit_values': [item.to_dict() for item in self.limit_values]}
