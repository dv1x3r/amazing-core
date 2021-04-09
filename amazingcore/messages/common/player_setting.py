from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class PlayerSetting(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 player_id: ObjectID = None,
                 setting_name: str = None,
                 value: str = None):
        self.aw_object_id = aw_object_id
        self.player_id = player_id
        self.setting_name = setting_name
        self.value = value

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        self.player_id.serialize(bit_stream)
        bit_stream.write_str(self.setting_name)
        bit_stream.write_str(self.value)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'player_id': self.player_id.to_dict(),
            'setting_name': self.setting_name,
            'value': self.value,
        }
