from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset

import datetime as dt


class GetZonesMessage(Message):
    def __init__(self):
        self.request: GetZonesRequest = GetZonesRequest()
        self.response: GetZonesResponse = GetZonesResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.zones = []


class GetZonesRequest(SerializableMessage):
    def __init__(self):
        pass

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        pass

    def to_dict(self):
        return {}


class GetZonesResponse(SerializableMessage):
    def __init__(self):
        self.zones: list = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.zones))
        for item in self.zones:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'zones': [item.to_dict() for item in self.zones],
        }
