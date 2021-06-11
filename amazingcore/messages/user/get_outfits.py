from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class GetOutfitsMessage(Message):
    def __init__(self):
        self.request: GetOutfitsRequest = GetOutfitsRequest()
        self.response: GetOutfitsResponse = GetOutfitsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.player_avatar_outfits = []


class GetOutfitsRequest(SerializableMessage):
    def __init__(self):
        self.player_avatar_id: ObjectID = None
        self.player_id: ObjectID = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_avatar_id = ObjectID()
        self.player_avatar_id.deserialize(bit_stream)
        self.player_id = ObjectID()
        self.player_id.deserialize(bit_stream)

    def to_dict(self):
        return {
            'player_avatar_id': self.player_avatar_id.to_dict(),
            'player_id': self.player_id.to_dict(),
        }


class GetOutfitsResponse(SerializableMessage):
    def __init__(self):
        self.player_avatar_outfits: list = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_avatar_outfits))
        for item in self.player_avatar_outfits:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'player_avatar_outfits': [item.to_dict() for item in self.player_avatar_outfits],
        }
