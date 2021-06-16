from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class RuleProperty(SerializableMessage):
    def __init__(self,
                 id: ObjectID = None,
                 parent_id: ObjectID = None,
                 components: list[str] = None,
                 parent_components: list[str] = None,
                 create_time: dt.datetime = None,
                 modified_time: dt.datetime = None,
                 properties: dict[str, str] = None,
                 children_group: dict[str, list] = None):  # list of RuleProperty
        self.id = id
        self.parent_id = parent_id
        self.components = components
        self.parent_components = parent_components
        self.create_time = create_time
        self.modified_time = modified_time
        self.properties = properties
        self.children_group = children_group

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.id = ObjectID()
        self.id.serialize(bit_stream)
        self.parent_id = ObjectID()
        self.parent_id.serialize(bit_stream)
        bit_stream.write_int(len(self.components))
        for item in self.components:
            bit_stream.write_str(item)
        bit_stream.write_int(len(self.parent_components))
        for item in self.parent_components:
            bit_stream.write_str(item)
        bit_stream.write_dt(self.create_time)
        bit_stream.write_dt(self.modified_time)
        bit_stream.write_int(len(self.properties))
        for key in self.properties:
            bit_stream.write_str(key)
            bit_stream.write_int(len(self.properties[key]))
            for item in self.properties[key]:
                bit_stream.write_str(item)
        bit_stream.write_int(len(self.children_group))
        for key in self.children_group:
            bit_stream.write_str(key)
            bit_stream.write_int(len(self.children_group[key]))
            for item in self.children_group[key]:
                bit_stream.write_int(len(item))
                for list_item in item:
                    list_item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        children_group_dict = {}
        for key in self.children_group:
            for item in self.children_group[key]:
                children_group_dict[key] = item.to_dict()

        return {
            'id': self.id.to_dict(),
            'parent_id': self.parent_id.to_dict(),
            'components': [i.to_dict() for i in self.components],
            'parent_components': [i.to_dict() for i in self.parent_components],
            'create_time': self.create_time,
            'modified_time': self.modified_time,
            'properties': self.properties,
            'children_group_dict': children_group_dict,
        }
