from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID

import datetime as dt


class PlayerQuest(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 create_date: dt.datetime = None,
                 accepted_date: dt.datetime = None,
                 started_date: dt.datetime = None,
                 expiry_date: dt.datetime = None,
                 completed_date: dt.datetime = None,
                 player_id: ObjectID = None,
                 quest_state_id: ObjectID = None,
                 quest_state_name: str = None,
                 player_avatar_id: ObjectID = None,
                 parent_player_quest_id: ObjectID = None,
                 quest=None,  # : Quest = None,
                 player_level: int = None,
                 npc_relationship_level: int = None,
                 npc_relationship_points: int = None,
                 player_money: int = None,
                 player_xp: int = None,
                 unlocked: bool = None,
                 rule_property=None):  # : RuleProperty = None):
        self.aw_object_id = aw_object_id
        self.create_date = create_date
        self.accepted_date = accepted_date
        self.started_date = started_date
        self.expiry_date = expiry_date
        self.completed_date = completed_date
        self.player_id = player_id
        self.quest_state_id = quest_state_id
        self.quest_state_name = quest_state_name
        self.player_avatar_id = player_avatar_id
        self.parent_player_quest_id = parent_player_quest_id
        self.quest = quest
        self.player_level = player_level
        self.npc_relationship_level = npc_relationship_level
        self.npc_relationship_points = npc_relationship_points
        self.player_money = player_money
        self.player_xp = player_xp
        self.unlocked = unlocked
        self.rule_property = rule_property

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_dt(self.create_date)
        bit_stream.write_dt(self.accepted_date)
        bit_stream.write_dt(self.started_date)
        bit_stream.write_dt(self.expiry_date)
        bit_stream.write_dt(self.completed_date)
        self.player_id.serialize(bit_stream)
        self.quest_state_id.serialize(bit_stream)
        bit_stream.write_str(self.quest_state_name)
        self.player_avatar_id.serialize(bit_stream)
        self.parent_player_quest_id.serialize(bit_stream)
        self.quest.serialize(bit_stream)
        bit_stream.write_int(self.player_level)
        bit_stream.write_int(self.npc_relationship_level)
        bit_stream.write_int(self.npc_relationship_points)
        bit_stream.write_int(self.player_money)
        bit_stream.write_int(self.player_xp)
        bit_stream.write_bool(self.unlocked)
        self.rule_property.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'create_date': self.create_date,
            'accepted_date': self.accepted_date,
            'started_date': self.started_date,
            'expiry_date': self.expiry_date,
            'completed_date': self.completed_date,
            'player_id': self.player_id.to_dict(),
            'quest_state_id': self.quest_state_id.to_dict(),
            'quest_state_name': self.quest_state_name,
            'player_avatar_id': self.player_avatar_id.to_dict(),
            'parent_player_quest_id': self.parent_player_quest_id.to_dict(),
            'quest': self.quest.to_dict(),
            'player_level': self.player_level,
            'npc_relationship_level': self.npc_relationship_level,
            'npc_relationship_points': self.npc_relationship_points,
            'player_money': self.player_money,
            'player_xp': self.player_xp,
            'unlocked': self.unlocked,
            'rule_property': self.rule_property.to_dict(),
        }
