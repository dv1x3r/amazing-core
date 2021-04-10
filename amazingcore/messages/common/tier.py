from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackage


class Tier(SerializableMessage):
    def __init__(self,
                 container_aw_object_id: ObjectID = None,
                 container_asset_map: dict[str, list[Asset]] = None,
                 container_asset_pkg: list[AssetPackage] = None,
                 rotation_days: int = None,
                 rotation_rate: int = None,
                 reporting_level_id: ObjectID = None,
                 paid: bool = None,
                 premium: bool = None,
                 closed: bool = None,
                 pricing_info: str = None,
                 expiry_period: int = None,
                 ordinal: int = None,
                 expiry_tier_id: ObjectID = None):
        self.container_aw_object_id = container_aw_object_id
        self.container_asset_map = container_asset_map
        self.container_asset_pkg = container_asset_pkg
        self.rotation_days = rotation_days
        self.rotation_rate = rotation_rate
        self.reporting_level_id = reporting_level_id
        self.paid = paid
        self.premium = premium
        self.closed = closed
        self.pricing_info = pricing_info
        self.expiry_period = expiry_period
        self.ordinal = ordinal
        self.expiry_tier_id = expiry_tier_id

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
        bit_stream.write_short(self.rotation_days)
        bit_stream.write_short(self.rotation_rate)
        self.reporting_level_id.serialize(bit_stream)
        bit_stream.write_bool(self.paid)
        bit_stream.write_bool(self.premium)
        bit_stream.write_bool(self.closed)
        bit_stream.write_str(self.pricing_info)
        bit_stream.write_short(self.expiry_period)
        bit_stream.write_int(self.ordinal)
        self.expiry_tier_id.serialize(bit_stream)

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
            'rotation_days': self.rotation_days,
            'rotation_rate': self.rotation_rate,
            'reporting_level_id': self.reporting_level_id.to_dict(),
            'paid': self.paid,
            'premium': self.premium,
            'closed': self.closed,
            'pricing_info': self.pricing_info,
            'expiry_period': self.expiry_period,
            'ordinal': self.ordinal,
            'expiry_tier_id': self.expiry_tier_id.to_dict(),
        }
