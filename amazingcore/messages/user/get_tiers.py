from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.tier import Tier


class GetTiersMessage(Message):
    def __init__(self):
        self.request: GetTiersRequest = GetTiersRequest()
        self.response: GetTiersResponse = GetTiersResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.tiers = [Tier(
            container_aw_object_id=ObjectID(1, 2, 3, 4),
            container_asset_map={},
            container_asset_pkg=[],
            rotation_days=1,
            rotation_rate=2,
            reporting_level_id=ObjectID(1, 2, 3, 4),
            paid=True,
            premium=True,
            closed=False,
            pricing_info='pricing_info',
            expiry_period=3,
            ordinal=4,
            expiry_tier_id=ObjectID(1, 2, 3, 4)
        )]


class GetTiersRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetTiersResponse(SerializableMessage):
    def __init__(self):
        self.tiers: list[Tier] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.tiers))
        for item in self.tiers:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'tiers': [item.to_dict() for item in self.tiers]}
