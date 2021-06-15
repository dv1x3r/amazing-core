from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class GetOnlineStatusesMessage(Message):
    def __init__(self):
        self.request: GetOnlineStatusesRequest = GetOnlineStatusesRequest()
        self.response: GetOnlineStatusesResponse = GetOnlineStatusesResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.online_statuses = []


class GetOnlineStatusesRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetOnlineStatusesResponse(SerializableMessage):
    def __init__(self):
        self.online_statuses = None  # GSFOnlineStatus

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.online_statuses))
        for item in self.online_statuses:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'zones': [item.to_dict() for item in self.online_statuses],
        }
