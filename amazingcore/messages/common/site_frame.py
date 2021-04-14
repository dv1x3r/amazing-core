from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackage


class SiteFrame(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 asset_map: dict[str, list[Asset]] = None,  # GSFAssetContainer
                 asset_packages: list[AssetPackage] = None,
                 type_value: int = None):
        self.aw_object_id = aw_object_id
        self.asset_map = asset_map
        self.asset_packages = asset_packages
        self.type_value = type_value

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
        bit_stream.write_int(self.type_value)

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
            'type_value': self.type_value,
        }
