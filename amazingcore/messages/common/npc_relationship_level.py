from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class NpcRelationshipLevel(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 level: int = None,
                 relationship_value: int = None,
                 npc_id: ObjectID = None):
        self.aw_object_id = aw_object_id
        self.level = level
        self.relationship_value = relationship_value
        self.npc_id = npc_id

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_int(self.level)
        bit_stream.write_int(self.relationship_value)
        self.npc_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'level': self.level,
            'relationship_value': self.relationship_value,
            'npc_id': self.npc_id.to_dict(),
        }
