from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class CrispData(SerializableMessage):
    def __init__(self,
                 crisp_action_id: ObjectID = None,
                 crisp_message: str = None,
                 crisp_expiry_date: dt.datetime = None,
                 crisp_confirmed: bool = None):
        self.crisp_action_id = crisp_action_id
        self.crisp_message = crisp_message
        self.crisp_expiry_date = crisp_expiry_date
        self.crisp_confirmed = crisp_confirmed

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.crisp_action_id.serialize(bit_stream)
        bit_stream.write_str(self.crisp_message)
        bit_stream.write_dt(self.crisp_expiry_date)
        bit_stream.write_bool(self.crisp_confirmed)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'crisp_action_id': self.crisp_action_id.to_dict(),
            'crisp_message': self.crisp_message,
            'crisp_expiry_date': self.crisp_expiry_date,
            'crisp_confirmed': self.crisp_confirmed,
        }
