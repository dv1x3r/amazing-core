from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackage
from amazingcore.messages.common.rule_property import RuleProperty
from amazingcore.messages.common.dimensions import Dimensions


class Zone(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 asset_map: dict[str, list[Asset]] = None,  # GSFAssetContainer
                 asset_packages: list[AssetPackage] = None,
                 rule_property: RuleProperty = None,  # GSFRuleContainer
                 locked: bool = None,
                 is_multiplayer: bool = None,
                 is_player_hosted: bool = None,
                 is_played_offline: bool = None,
                 dimensions: Dimensions = None,  # Zone
                 buildings: list = None,  # GSFBuilding
                 ptag: str = None,
                 capacity: int = None):
        self.aw_object_id = aw_object_id
        self.asset_map = asset_map
        self.asset_packages = asset_packages
        self.rule_property = rule_property
        self.locked = locked
        self.is_multiplayer = is_multiplayer
        self.is_player_hosted = is_player_hosted
        self.is_played_offline = is_played_offline
        self.dimensions = dimensions
        self.buildings = buildings
        self.ptag = ptag
        self.capacity = capacity

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
        self.rule_property.serialize(bit_stream)
        bit_stream.write_bool(self.locked)
        bit_stream.write_bool(self.is_multiplayer)
        bit_stream.write_bool(self.is_player_hosted)
        bit_stream.write_bool(self.is_played_offline)
        self.dimensions.serialize(bit_stream)
        bit_stream.write_int(len(self.buildings))
        for item in self.buildings:
            item.serialize(bit_stream)
        bit_stream.write_str(self.ptag)
        bit_stream.write_int(self.capacity)

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
            'rule_property': self.rule_property.to_dict(),
            'locked': self.locked,
            'is_multiplayer': self.is_multiplayer,
            'is_player_hosted': self.is_player_hosted,
            'is_played_offline': self.is_played_offline,
            'dimensions': self.dimensions.to_dict(),
            'buildings': [item.to_dict() for item in self.buildings],
            'ptag': self.ptag,
            'capacity': self.capacity,
        }
