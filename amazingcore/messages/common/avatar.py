from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackageContainer


class Avatar(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 asset_map: dict[str, list[Asset]] = None,  # GSFAssetContainer
                 asset_packages: list[AssetPackageContainer] = None,
                 dimensions: str = None,
                 weight: int = None,
                 height: int = None,
                 max_outfits: int = None,
                 name: str = None):
        self.aw_object_id = aw_object_id
        self.asset_map = asset_map
        self.asset_packages = asset_packages
        self.dimensions = dimensions
        self.weight = weight
        self.height = height
        self.max_outfits = max_outfits
        self.name = name

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
        bit_stream.write_str(self.dimensions)
        bit_stream.write_double(self.weight)
        bit_stream.write_double(self.height)
        bit_stream.write_short(self.max_outfits)
        bit_stream.write_str(self.name)

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
            'dimensions': self.dimensions,
            'weight': self.weight,
            'height': self.height,
            'max_outfits': self.max_outfits,
            'name': self.name,
        }
