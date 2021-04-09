from amazingcore.messages.common.player_setting import PlayerSetting
from amazingcore.messages.common.crisp_data import CrispData
from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class PlayerInfo(SerializableMessage):
    def __init__(self,
                 tier_id: ObjectID = None,
                 player_name: str = None,
                 world_name: str = None,
                 crisp_data: CrispData = None,
                 verified: bool = None,
                 verification_expiry: dt.datetime = None,
                 scs_block_expiry: dt.datetime = None,
                 eula_id: ObjectID = None,
                 current_eula_id: ObjectID = None,
                 u_13: bool = None,
                 chat_blocked_parent: bool = None,
                 chat_allowed: bool = None,
                 chat_blocked_expiry: dt.datetime = None,
                 findable: bool = None,
                 findable_date: dt.datetime = None,
                 user_subscription_expiry_date: dt.datetime = None,
                 qa: bool = None,
                 player_settings: list[PlayerSetting] = None):
        self.tier_id = tier_id
        self.player_name = player_name
        self.world_name = world_name
        self.crisp_data = crisp_data
        self.verified = verified
        self.verification_expiry = verification_expiry
        self.scs_block_expiry = scs_block_expiry
        self.eula_id = eula_id
        self.current_eula_id = current_eula_id
        self.u_13 = u_13
        self.chat_blocked_parent = chat_blocked_parent
        self.chat_allowed = chat_allowed
        self.chat_blocked_expiry = chat_blocked_expiry
        self.findable = findable
        self.findable_date = findable_date
        self.user_subscription_expiry_date = user_subscription_expiry_date
        self.qa = qa
        self.player_settings = player_settings

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.tier_id.serialize(bit_stream)
        bit_stream.write_str(self.player_name)
        bit_stream.write_str(self.world_name)
        self.crisp_data.serialize(bit_stream)
        bit_stream.write_bool(self.verified)
        bit_stream.write_dt(self.verification_expiry)
        bit_stream.write_dt(self.scs_block_expiry)
        self.eula_id.serialize(bit_stream)
        self.current_eula_id.serialize(bit_stream)
        bit_stream.write_bool(self.u_13)
        bit_stream.write_bool(self.chat_blocked_parent)
        bit_stream.write_bool(self.chat_allowed)
        bit_stream.write_dt(self.chat_blocked_expiry)
        bit_stream.write_bool(self.findable)
        bit_stream.write_dt(self.findable_date)
        bit_stream.write_dt(self.user_subscription_expiry_date)
        bit_stream.write_bool(self.qa)
        bit_stream.write_int(len(self.player_settings))
        for ps in self.player_settings:
            ps.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'tier_id': self.tier_id.to_dict(),
            'player_name': self.player_name,
            'world_name': self.world_name,
            'crisp_data': self.crisp_data.to_dict(),
            'verified': self.verified,
            'verification_expiry': self.verification_expiry,
            'scs_block_expiry': self.scs_block_expiry,
            'eula_id': self.eula_id.to_dict(),
            'current_eula_id': self.current_eula_id.to_dict(),
            'u_13': self.u_13,
            'chat_blocked_parent': self.chat_blocked_parent,
            'chat_allowed': self.chat_allowed,
            'chat_blocked_expiry': self.chat_blocked_expiry,
            'findable': self.findable,
            'findable_date': self.findable_date,
            'user_subscription_expiry_date': self.user_subscription_expiry_date,
            'qa': self.qa,
            'player_settings': [i.to_dict() for i in self.player_settings],
        }
