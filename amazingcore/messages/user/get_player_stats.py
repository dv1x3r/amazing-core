from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_stats import PlayerStats


class GetPlayerStatsMessage(Message):
    def __init__(self):
        self.request: GetPlayerStatsRequest = GetPlayerStatsRequest()
        self.response: GetPlayerStatsResponse = GetPlayerStatsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        # self.response.player_stats = [PlayerStats(
        #     aw_object_id=ObjectID(1, 2, 3, 4),
        #     player_avatar_id=ObjectID(1, 2, 3, 4),
        #     stats_type_id=ObjectID(1, 2, 3, 4),
        #     level=12,
        #     object_id=ObjectID(1, 2, 3, 4)
        # )]

        self.response.player_stats = []


class GetPlayerStatsRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetPlayerStatsResponse(SerializableMessage):
    def __init__(self):
        self.player_stats: list[PlayerStats] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_stats))
        for item in self.player_stats:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'player_stats': [item.to_dict() for item in self.player_stats]}
