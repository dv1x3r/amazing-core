from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset

import datetime as dt


class GetInventoryObjectsMessage(Message):
    def __init__(self):
        self.request: GetInventoryObjectsRequest = GetInventoryObjectsRequest()
        self.response: GetInventoryObjectsResponse = GetInventoryObjectsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.player_items = []


class GetInventoryObjectsRequest(SerializableMessage):
    def __init__(self):
        self.container_id: ObjectID = None
        self.item_category_ids: list[ObjectID] = None
        self.player_item_ids: list[ObjectID] = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.container_id = ObjectID()
        self.container_id.deserialize(bit_stream)
        self.item_category_ids = []
        for _ in range(bit_stream.read_int()):
            item = ObjectID()
            item.deserialize(bit_stream)
            self.item_category_ids.append(item)
        self.player_item_ids = []
        for _ in range(bit_stream.read_int()):
            item = ObjectID()
            item.deserialize(bit_stream)
            self.player_item_ids.append(item)

    def to_dict(self):
        return {
            'container_id': self.container_id.to_dict(),
            'item_category_ids': [item.to_dict() for item in self.item_category_ids],
            'player_item_ids': [item.to_dict() for item in self.player_item_ids],
        }


class GetInventoryObjectsResponse(SerializableMessage):
    def __init__(self):
        self.player_items: list = None  # GSFPlayerItem

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_items))
        for item in self.player_items:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'player_items': [item.to_dict() for item in self.player_items],
        }
