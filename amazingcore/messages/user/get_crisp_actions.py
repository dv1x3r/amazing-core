from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.crisp_action import CrispAction


class GetCrispActionsMessage(Message):
    def __init__(self):
        self.request: GetCrispActionsRequest = GetCrispActionsRequest()
        self.response: GetCrispActionsResponse = GetCrispActionsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.crisp_actions = []


class GetCrispActionsRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetCrispActionsResponse(SerializableMessage):
    def __init__(self):
        self.crisp_actions: list[CrispAction] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.crisp_actions))
        for item in self.crisp_actions:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'crisp_actions': [item.to_dict() for item in self.crisp_actions]}
