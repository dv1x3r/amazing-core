from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_home_invite import PlayerHomeInvite

import datetime as dt


class GetHomeInvitationsMessage(Message):
    def __init__(self):
        self.request: GetHomeInvitationsRequest = GetHomeInvitationsRequest()
        self.response: GetHomeInvitationsResponse = GetHomeInvitationsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        # self.response.player_home_invites = [PlayerHomeInvite(
        #     aw_object_id=ObjectID(1, 2, 3, 4),
        #     invite_player_id=ObjectID(1, 2, 3, 4),
        #     player_id=ObjectID(1, 2, 3, 4),
        #     is_player_home=True,
        #     invite_status=0,
        #     invite_date=dt.datetime.now(),
        #     blocked_date=(dt.datetime.now() + dt.timedelta(days=1))
        # )]

        self.response.player_home_invites = []


class GetHomeInvitationsRequest(SerializableMessage):
    def __init__(self):
        self.player_id: ObjectID = None
        self.is_host: bool = None
        self.is_status_pending: bool = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_id = ObjectID()
        self.player_id.deserialize(bit_stream)
        self.is_host = bit_stream.read_bool()
        self.is_status_pending = bit_stream.read_bool()

    def to_dict(self):
        return {
            'player_id': self.player_id.to_dict(),
            'is_host': self.is_host,
            'is_status_pending': self.is_status_pending,
        }


class GetHomeInvitationsResponse(SerializableMessage):
    def __init__(self):
        self.player_home_invites: list[PlayerHomeInvite] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.player_home_invites))
        for item in self.player_home_invites:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'player_home_invites': [item.to_dict() for item in self.player_home_invites]}
