from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.avatar import Avatar

import datetime as dt


class PlayerAvatar(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 avatar: Avatar = None,
                 player_id: ObjectID = None,
                 name: str = None,
                 bio: str = None,
                 secret_code: str = None,
                 create_ts: dt.datetime = None,
                 player_avatar_outfit_id: ObjectID = None,
                 outfit_no: int = None,
                 play_time: int = None,
                 last_play: dt.datetime = None):
        self.aw_object_id = aw_object_id
        self.avatar = avatar
        self.player_id = player_id
        self.name = name
        self.bio = bio
        self.secret_code = secret_code
        self.create_ts = create_ts
        self.player_avatar_outfit_id = player_avatar_outfit_id
        self.outfit_no = outfit_no
        self.play_time = play_time
        self.last_play = last_play

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        self.avatar.serialize(bit_stream)
        self.player_id.serialize(bit_stream)
        bit_stream.write_str(self.name)
        bit_stream.write_str(self.bio)
        bit_stream.write_str(self.secret_code)
        bit_stream.write_dt(self.create_ts)
        self.player_avatar_outfit_id.serialize(bit_stream)
        bit_stream.write_short(self.outfit_no)
        bit_stream.write_long(self.play_time, nullable=True)
        bit_stream.write_dt(self.last_play)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'avatar': self.avatar.to_dict(),
            'player_id': self.player_id.to_dict(),
            'name': self.name,
            'bio': self.bio,
            'secret_code': self.secret_code,
            'create_ts': self.create_ts,
            'player_avatar_outfit_id': self.player_avatar_outfit_id.to_dict(),
            'outfit_no': self.outfit_no,
            'play_time': self.play_time,
            'last_play': self.last_play,
        }
