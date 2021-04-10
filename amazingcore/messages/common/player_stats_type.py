from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackage

from enum import Enum


class PlayerStatsTypeValue(Enum):
    ENERGY = 0
    EXPERIENCE = 1
    SATISFACTION = 2
    MONEY = 3
    POINTS = 4
    LEVEL = 5
    ESTORE = 6
    TENDER = 7
    TOKEN = 8


class PlayerStatsType(SerializableMessage):
    def __init__(self,
                 container_aw_object_id: ObjectID = None,
                 container_asset_map: dict[str, list[Asset]] = None,
                 container_asset_pkg: list[AssetPackage] = None,
                 player_stats_type_value: PlayerStatsTypeValue = None,
                 is_avatar: bool = None,
                 name: str = None):
        self.container_aw_object_id = container_aw_object_id
        self.container_asset_map = container_asset_map
        self.container_asset_pkg = container_asset_pkg
        self.player_stats_type_value = player_stats_type_value
        self.is_avatar = is_avatar
        self.name = name

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.container_aw_object_id.serialize(bit_stream)
        bit_stream.write_int(len(self.container_asset_map))
        for key in self.container_asset_map:
            bit_stream.write_str(key)
            bit_stream.write_int(len(self.container_asset_map[key]))
            for item in self.container_asset_map[key]:
                item.serialize(bit_stream)
        bit_stream.write_int(len(self.container_asset_pkg))
        for item in self.container_asset_pkg:
            item.serialize(bit_stream)
        bit_stream.write_str(self.player_stats_type_value.name)
        bit_stream.write_bool(self.is_avatar)
        bit_stream.write_str(self.name)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        asset_map_dict = {}
        for key in self.container_asset_map:
            for asset in self.container_asset_map[key]:
                asset_map_dict[key] = asset.to_dict()

        return {
            'container_aw_object_id': self.container_aw_object_id.to_dict(),
            'container_asset_map': asset_map_dict,
            'container_asset_pkg': [item.to_dict() for item in self.container_asset_pkg],
            'player_stats_type_value': self.player_stats_type_value,
            'is_avatar': self.is_avatar,
            'name': self.name,
        }
