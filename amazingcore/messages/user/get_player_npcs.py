from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset

import datetime as dt


class GetPlayerNpcsMessage(Message):
    def __init__(self):
        self.request: GetPlayerNpcsRequest = GetPlayerNpcsRequest()
        self.response: GetPlayerNpcsResponse = GetPlayerNpcsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.npcs = []


class GetPlayerNpcsRequest(SerializableMessage):
    def __init__(self):
        self.player_id: ObjectID = None
        self.zone_id: ObjectID = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_id = ObjectID()
        self.player_id.deserialize(bit_stream)
        self.zone_id = ObjectID()
        self.zone_id.deserialize(bit_stream)

    def to_dict(self):
        return {
            'player_id': self.player_id.to_dict(),
            'zone_id': self.zone_id.to_dict(),
        }


class GetPlayerNpcsResponse(SerializableMessage):
    def __init__(self):
        self.npcs: list = None  # GSFNPC

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.npcs))
        for item in self.npcs:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'npcs': [item.to_dict() for item in self.npcs],
        }
