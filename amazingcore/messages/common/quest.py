from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class Quest(SerializableMessage):
    def __init__(self,
                 # RuleContainer properties
                 ):
        pass

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {}
