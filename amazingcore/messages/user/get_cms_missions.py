from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.quest import Quest

import datetime as dt


class GetCmsMissionsMessage(Message):
    def __init__(self):
        self.request: GetCmsMissionsRequest = GetCmsMissionsRequest()
        self.response: GetCmsMissionsResponse = GetCmsMissionsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.missions = []  # nested classes are not implemented


class GetCmsMissionsRequest(SerializableMessage):
    def __init__(self):
        self.quest_type_id: ObjectID = None
        self.hierarchies: list[str] = None
        self.return_locked: bool = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.quest_type_id = ObjectID()
        self.quest_type_id.deserialize(bit_stream)
        self.hierarchies = []
        for _ in range(bit_stream.read_int()):
            self.hierarchies.append(bit_stream.read_str())
        self.return_locked = bit_stream.read_bool()

    def to_dict(self):
        return {
            'quest_type_id': self.quest_type_id.to_dict(),
            'hierarchies': self.hierarchies,
            'return_locked': self.return_locked,
        }


class GetCmsMissionsResponse(SerializableMessage):
    def __init__(self):
        self.missions: list[Quest] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.missions))
        for item in self.missions:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'missions': [item.to_dict() for item in self.missions]}
