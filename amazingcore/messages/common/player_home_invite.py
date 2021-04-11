from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class PlayerHomeInvite(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,
                 invite_player_id: ObjectID = None,
                 player_id: ObjectID = None,
                 is_player_home: bool = None,
                 invite_status: int = None,
                 invite_date: dt.datetime = None,
                 blocked_date: dt.datetime = None):
        self.aw_object_id = aw_object_id
        self.invite_player_id = invite_player_id
        self.player_id = player_id
        self.is_player_home = is_player_home
        self.invite_status = invite_status
        self.invite_date = invite_date
        self.blocked_date = blocked_date

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        self.invite_player_id.serialize(bit_stream)
        self.player_id.serialize(bit_stream)
        bit_stream.write_bool(self.is_player_home)
        bit_stream.write_int(self.invite_status)
        bit_stream.write_dt(self.invite_date)
        bit_stream.write_dt(self.blocked_date)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'invite_player_id': self.invite_player_id.to_dict(),
            'player_id': self.player_id.to_dict(),
            'is_player_home': self.is_player_home,
            'invite_status': self.invite_status,
            'invite_date': self.invite_date,
            'blocked_date': self.blocked_date,
        }
