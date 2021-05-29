from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_stats_type import PlayerStatsType, PlayerStatsTypeValue


class GetStatsTypeMessage(Message):
    def __init__(self):
        self.request: GetStatsTypeRequest = GetStatsTypeRequest()
        self.response: GetStatsTypeResponse = GetStatsTypeResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        # self.response.player_stats_types = [PlayerStatsType(
        #     container_aw_object_id=ObjectID(1, 2, 3, 4),
        #     container_asset_map={},
        #     container_asset_pkg=[],
        #     player_stats_type_value=PlayerStatsTypeValue.LEVEL,
        #     is_avatar=False,
        #     name='name'
        # )]

        self.response.player_stats_types = []


class GetStatsTypeRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetStatsTypeResponse(SerializableMessage):
    def __init__(self):
        self.player_stats_types: list[PlayerStatsType] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_stats_types))
        for item in self.player_stats_types:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'player_stats_types': [item.to_dict() for item in self.player_stats_types]}
