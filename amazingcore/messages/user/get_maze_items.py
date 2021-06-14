from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class GetMazeItemsMessage(Message):
    def __init__(self):
        self.request: GetMazeItemsRequest = GetMazeItemsRequest()
        self.response: GetMazeItemsResponse = GetMazeItemsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.maze_items = []


class GetMazeItemsRequest(SerializableMessage):
    def __init__(self):
        self.player_maze_id: ObjectID
        self.player_id: ObjectID

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_maze_id = ObjectID()
        self.player_maze_id.deserialize(bit_stream)
        self.player_id = ObjectID()
        self.player_id.deserialize(bit_stream)

    def to_dict(self):
        return {
            'player_maze_id': self.player_maze_id.to_dict(),
            'player_id': self.player_id.to_dict(),
        }


class GetMazeItemsResponse(SerializableMessage):
    def __init__(self):
        self.maze_items: list = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(0)  # maze_items size

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'maze_items': self.maze_items}
