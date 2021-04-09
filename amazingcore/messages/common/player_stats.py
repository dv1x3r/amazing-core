from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class PlayerStats(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 player_avatar_id: ObjectID = None,
                 stats_type_id: ObjectID = None,
                 level: int = None,
                 object_id: ObjectID = None):
        self.aw_object_id = aw_object_id
        self.player_avatar_id = player_avatar_id
        self.stats_type_id = stats_type_id
        self.level = level
        self.object_id = object_id

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        self.player_avatar_id.serialize(bit_stream)
        self.stats_type_id.serialize(bit_stream)
        bit_stream.write_int(self.level)
        self.object_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'player_avatar_id': self.player_avatar_id.to_dict(),
            'stats_type_id': self.stats_type_id.to_dict(),
            'level': self.level,
            'object_id': self.object_id.to_dict(),
        }
