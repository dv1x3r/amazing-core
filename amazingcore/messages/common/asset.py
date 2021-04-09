from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID


class Asset(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 asset_type_name: str = None,
                 cdn_id: str = None,
                 res_name: str = None,
                 group_name: str = None,
                 file_size: int = None):
        self.aw_object_id = aw_object_id
        self.asset_type_name = asset_type_name
        self.cdn_id = cdn_id
        self.res_name = res_name
        self.group_name = group_name
        self.file_size = file_size

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_str(self.asset_type_name)
        bit_stream.write_str(self.cdn_id)
        bit_stream.write_str(self.res_name)
        bit_stream.write_str(self.group_name)
        bit_stream.write_long(self.file_size)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'asset_type_name': self.asset_type_name,
            'cdn_id': self.cdn_id,
            'res_name': self.res_name,
            'group_name': self.group_name,
            'file_size': self.file_size,
        }
