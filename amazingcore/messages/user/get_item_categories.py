from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class GetItemCategoriesMessage(Message):
    def __init__(self):
        self.request: GetItemCategoriesRequest = GetItemCategoriesRequest()
        self.response: GetItemCategoriesResponse = GetItemCategoriesResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class GetItemCategoriesRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetItemCategoriesResponse(SerializableMessage):
    def __init__(self):
        self.item_categories = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(0)  # item_categories size

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'item_categories': self.item_categories}
