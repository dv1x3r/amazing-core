from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.npc_relationship_level import NpcRelationshipLevel


class GetNpcRelationshipLevelsMessage(Message):
    def __init__(self):
        self.request: GetNpcRelationshipLevelsRequest = GetNpcRelationshipLevelsRequest()
        self.response: GetNpcRelationshipLevelsResponse = GetNpcRelationshipLevelsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        # self.response.npc_relationship_levels = [NpcRelationshipLevel(
        #     aw_object_id=ObjectID(1, 2, 3, 4),
        #     level=0,
        #     relationship_value=1,
        #     npc_id=ObjectID(1, 2, 3, 4)
        # )]

        self.response.npc_relationship_levels = []


class GetNpcRelationshipLevelsRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetNpcRelationshipLevelsResponse(SerializableMessage):
    def __init__(self):
        self.npc_relationship_levels: list[NpcRelationshipLevel] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.npc_relationship_levels))
        for item in self.npc_relationship_levels:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'npc_relationship_levels': [item.to_dict() for item in self.npc_relationship_levels]}
