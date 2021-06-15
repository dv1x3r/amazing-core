from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.position import Position
from amazingcore.messages.common.orientation import Orientation


class EnterBuildingMessage(Message):
    def __init__(self):
        self.request: EnterBuildingRequest = EnterBuildingRequest()
        self.response: EnterBuildingResponse = EnterBuildingResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.building_id = ObjectID(0, 0, 0, 0)


class EnterBuildingRequest(SerializableMessage):
    def __init__(self):
        self.loc_id: ObjectID = None
        self.building_id: ObjectID = None
        self.pos: Position = None
        self.orientation: Orientation = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.loc_id = ObjectID()
        self.loc_id.deserialize(bit_stream)
        self.building_id = ObjectID()
        self.building_id.deserialize(bit_stream)
        self.pos = Position()
        self.pos.deserialize(bit_stream)
        self.orientation = Orientation()
        self.orientation.deserialize(bit_stream)

    def to_dict(self):
        return {
            'loc_id': self.loc_id.to_dict() if self.loc_id else None,
            'building_id': self.building_id.to_dict() if self.building_id else None,
            'pos': self.pos.to_dict() if self.pos else None,
            'orientation': self.orientation.to_dict() if self.orientation else None,
        }


class EnterBuildingResponse(SerializableMessage):
    def __init__(self):
        self.building_id: ObjectID = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.building_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'building_id': self.building_id.to_dict(),
        }
