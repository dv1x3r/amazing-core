from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class ManageHomeInvitationsMessage(Message):
    def __init__(self):
        self.request: ManageHomeInvitationsRequest = ManageHomeInvitationsRequest()
        self.response: ManageHomeInvitationsResponse = ManageHomeInvitationsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class ManageHomeInvitationsRequest(SerializableMessage):
    def __init__(self):
        self.player_ids: list[ObjectID] = None
        self.is_invite_or_accept: bool = None
        self.is_host: bool = None
        self.permanent: bool = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_ids = []
        for _ in range(bit_stream.read_int()):
            item = ObjectID()
            item.deserialize(bit_stream)
            self.player_ids.append(item)
        self.is_invite_or_accept = bit_stream.read_bool()
        self.is_host = bit_stream.read_bool()
        self.permanent = bit_stream.read_bool()

    def to_dict(self):
        return {
            'player_ids': [item.to_dict() for item in self.player_ids],
            'is_invite_or_accept': self.is_invite_or_accept,
            'is_host': self.is_host,
            'permanent': self.permanent,
        }


class ManageHomeInvitationsResponse(SerializableMessage):

    def serialize(self, bit_stream: BitStream):
        pass

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        pass
