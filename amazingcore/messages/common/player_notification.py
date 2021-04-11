from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.notification import Notification

import datetime as dt


class PlayerNotification(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 from_player_id: ObjectID = None,
                 from_player_name: str = None,
                 to_player_id: ObjectID = None,
                 village_name: str = None,
                 village_id: ObjectID = None,
                 object_id: ObjectID = None,
                 notification_text: str = None,
                 notification: Notification = None,
                 is_read: bool = None,
                 create_date: dt.datetime = None,
                 expiry_date: dt.datetime = None):
        self.aw_object_id = aw_object_id
        self.from_player_id = from_player_id
        self.from_player_name = from_player_name
        self.to_player_id = to_player_id
        self.village_name = village_name
        self.village_id = village_id
        self.object_id = object_id
        self.notification_text = notification_text
        self.notification = notification
        self.is_read = is_read
        self.create_date = create_date
        self.expiry_date = expiry_date

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        self.from_player_id.serialize(bit_stream)
        bit_stream.write_str(self.from_player_name)
        self.to_player_id.serialize(bit_stream)
        bit_stream.write_str(self.village_name)
        self.village_id.serialize(bit_stream)
        self.object_id.serialize(bit_stream)
        bit_stream.write_str(self.notification_text)
        self.notification.serialize(bit_stream)
        bit_stream.write_bool(self.is_read)
        bit_stream.write_dt(self.create_date)
        bit_stream.write_dt(self.expiry_date)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'from_player_id': self.from_player_id.to_dict(),
            'from_player_name': self.from_player_name,
            'to_player_id': self.to_player_id.to_dict(),
            'village_name': self.village_name,
            'village_id': self.village_id.to_dict(),
            'object_id': self.object_id.to_dict(),
            'notification_text': self.notification_text,
            'notification': self.notification.to_dict(),
            'is_read': self.is_read,
            'create_date': self.create_date,
            'expiry_date': self.expiry_date,
        }
