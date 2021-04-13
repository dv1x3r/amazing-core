from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_quest import PlayerQuest

import datetime as dt


class GetPlayerMissionsMessage(Message):
    def __init__(self):
        self.request: GetPlayerMissionsRequest = GetPlayerMissionsRequest()
        self.response: GetPlayerMissionsResponse = GetPlayerMissionsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.player_missions = []  # nested classes are not implemented


class GetPlayerMissionsRequest(SerializableMessage):
    def __init__(self):
        self.limit: int = None
        self.offset: int = None
        self.action: int = None
        self.hierarchies: list[str] = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.limit = bit_stream.read_int()
        self.offset = bit_stream.read_int()
        self.action = bit_stream.read_int()
        self.hierarchies = []
        for _ in range(bit_stream.read_int()):
            self.hierarchies.append(bit_stream.read_str())

    def to_dict(self):
        return {
            'limit': self.limit,
            'offset': self.offset,
            'action': self.action,
            'hierarchies': self.hierarchies,
        }


class GetPlayerMissionsResponse(SerializableMessage):
    def __init__(self):
        self.player_missions: list[PlayerQuest] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_missions))
        for item in self.player_missions:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'player_missions': [item.to_dict() for item in self.player_missions]}
