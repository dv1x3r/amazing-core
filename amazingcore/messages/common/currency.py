from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackage


class Currency(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 asset_map: dict[str, list[Asset]] = None,  # GSFAssetContainer
                 asset_packages: list[AssetPackage] = None,
                 stats_type_id: ObjectID = None,
                 is_default: bool = None):
        self.aw_object_id = aw_object_id
        self.asset_map = asset_map
        self.asset_packages = asset_packages
        self.stats_type_id = stats_type_id
        self.is_default = is_default

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_int(len(self.asset_map))
        for key in self.asset_map:
            bit_stream.write_str(key)
            bit_stream.write_int(len(self.asset_map[key]))
            for item in self.asset_map[key]:
                item.serialize(bit_stream)
        bit_stream.write_int(len(self.asset_packages))
        for item in self.asset_packages:
            item.serialize(bit_stream)
        self.stats_type_id.serialize(bit_stream)
        bit_stream.write_bool(self.is_default, nullable=True)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        asset_map_dict = {}
        for key in self.asset_map:
            for asset in self.asset_map[key]:
                asset_map_dict[key] = asset.to_dict()

        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'asset_map': asset_map_dict,
            'asset_packages': [item.to_dict() for item in self.asset_packages],
            'stats_type_id': self.stats_type_id.to_dict(),
            'is_default': self.is_default,
        }
