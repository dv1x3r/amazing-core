from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.currency import Currency


class GetCurrenciesMessage(Message):
    def __init__(self):
        self.request: GetCurrenciesRequest = GetCurrenciesRequest()
        self.response: GetCurrenciesResponse = GetCurrenciesResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        # self.response.currencies = [Currency(
        #     aw_object_id=ObjectID(1, 2, 3, 4),
        #     asset_map={},
        #     asset_packages=[],
        #     stats_type_id=ObjectID(1, 2, 3, 4),
        #     is_default=True
        # ), Currency(
        #     aw_object_id=ObjectID(1, 2, 3, 4),
        #     asset_map={},
        #     asset_packages=[],
        #     stats_type_id=ObjectID(1, 2, 3, 4),
        #     is_default=False
        # ), Currency(
        #     aw_object_id=ObjectID(1, 2, 3, 4),
        #     asset_map={},
        #     asset_packages=[],
        #     stats_type_id=ObjectID(1, 2, 3, 4),
        #     is_default=None
        # )]

        self.response.currencies = []


class GetCurrenciesRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetCurrenciesResponse(SerializableMessage):
    def __init__(self):
        self.currencies: list[Currency] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.currencies))
        for item in self.currencies:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'currencies': [item.to_dict() for item in self.currencies]}
