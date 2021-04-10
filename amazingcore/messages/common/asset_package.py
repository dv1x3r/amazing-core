from amazingcore.messages.common.asset import Asset
from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class AssetPackage(SerializableMessage):

    def __init__(self,
                 container_aw_object_id: ObjectID = None,
                 container_asset_map: dict[str, list[Asset]] = None,
                 container_asset_pkg: list[any] = None,  # list[AssetPackage]
                 p_tag: str = None,
                 create_date: dt.datetime = None):
        self.container_aw_object_id = container_aw_object_id
        self.container_asset_map = container_asset_map
        self.container_asset_pkg: list[AssetPackage] = container_asset_pkg
        self.p_tag = p_tag
        self.create_date = create_date

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
        bit_stream.write_str(self.p_tag)
        bit_stream.write_dt(self.create_date)

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
        }
