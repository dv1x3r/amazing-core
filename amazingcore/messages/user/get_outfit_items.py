from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.currency import Currency


class GetOutfitItemsMessage(Message):
    def __init__(self):
        self.request: GetOutfitItemsRequest = GetOutfitItemsRequest()
        self.response: GetOutfitItemsResponse = GetOutfitItemsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.outfit_items = []
        # DessAvatarManager.cs (new GSF)
        # OutfitsManager.cs (new GSF)


class GetOutfitItemsRequest(SerializableMessage):
    def __init__(self):
        self.player_avatar_outfit_id: ObjectID = None
        self.player_id: ObjectID = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_avatar_outfit_id = ObjectID()
        self.player_avatar_outfit_id.deserialize(bit_stream)
        self.player_id = ObjectID()
        self.player_id.deserialize(bit_stream)

    def to_dict(self):
        return {
            'player_avatar_outfit_id': self.player_avatar_outfit_id.to_dict(),
            'player_id': self.player_id.to_dict(),
        }


class GetOutfitItemsResponse(SerializableMessage):
    def __init__(self):
        self.outfit_items: list = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.outfit_items))
        for item in self.outfit_items:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'outfit_items': [item.to_dict() for item in self.outfit_items]}
