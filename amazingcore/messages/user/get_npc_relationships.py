from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset

import datetime as dt


class GetNpcRelationshipsMessage(Message):
    def __init__(self):
        self.request: GetNpcRelationshipsRequest = GetNpcRelationshipsRequest()
        self.response: GetNpcRelationshipsResponse = GetNpcRelationshipsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.relationships = []


class GetNpcRelationshipsRequest(SerializableMessage):
    def __init__(self):
        self.zone_id: ObjectID = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.zone_id = ObjectID()
        self.zone_id.deserialize(bit_stream)

    def to_dict(self):
        return {'zone_id': self.zone_id.to_dict()}


class GetNpcRelationshipsResponse(SerializableMessage):
    def __init__(self):
        self.relationships: list = None  # GSFNpcRelationship

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.relationships))
        for item in self.relationships:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'relationships': [item.to_dict() for item in self.relationships],
        }
