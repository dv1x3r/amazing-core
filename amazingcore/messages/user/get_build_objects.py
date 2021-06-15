from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset

import datetime as dt


class GetBuildObjectsMessage(Message):
    def __init__(self):
        self.request: GetBuildObjectsRequest = GetBuildObjectsRequest()
        self.response: GetBuildObjectsResponse = GetBuildObjectsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.player_build_objects = []


class GetBuildObjectsRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetBuildObjectsResponse(SerializableMessage):
    def __init__(self):
        self.player_build_objects: list = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_build_objects))
        for item in self.player_build_objects:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'player_build_objects': [item.to_dict() for item in self.player_build_objects],
        }
