from amazingcore.messages.common.race_mode import RaceMode
from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar

import datetime as dt


class Player(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,  # GSFAwObject
                 create_date: dt.datetime = None,
                 active_player_avatar: PlayerAvatar = None,
                 home_theme_id: ObjectID = None,
                 current_race_mode: RaceMode = None,
                 workshop_options: str = None,
                 is_tutorial_completed: bool = None,
                 yard_building_id: ObjectID = None,
                 last_login: dt.datetime = None,
                 play_time: int = None,
                 is_qa: bool = None,
                 home_village_plot_id: ObjectID = None,
                 store_village_plot_id: ObjectID = None,
                 player_store_id: ObjectID = None,
                 player_maze_id: ObjectID = None,
                 village_id: ObjectID = None):
        self.aw_object_id = aw_object_id
        self.create_date = create_date
        self.active_player_avatar = active_player_avatar
        self.home_theme_id = home_theme_id
        self.current_race_mode = current_race_mode
        self.workshop_options = workshop_options
        self.is_tutorial_completed = is_tutorial_completed
        self.yard_building_id = yard_building_id
        self.last_login = last_login
        self.play_time = play_time
        self.is_qa = is_qa
        self.home_village_plot_id = home_village_plot_id
        self.store_village_plot_id = store_village_plot_id
        self.player_store_id = player_store_id
        self.player_maze_id = player_maze_id
        self.village_id = village_id

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_dt(self.create_date)
        self.active_player_avatar.serialize(bit_stream)
        self.home_theme_id.serialize(bit_stream)
        self.current_race_mode.serialize(bit_stream)
        bit_stream.write_str(self.workshop_options)
        bit_stream.write_bool(self.is_tutorial_completed)
        self.yard_building_id.serialize(bit_stream)
        bit_stream.write_dt(self.last_login)
        bit_stream.write_long(self.play_time, nullable=True)
        bit_stream.write_bool(self.is_qa)
        self.home_village_plot_id.serialize(bit_stream)
        self.store_village_plot_id.serialize(bit_stream)
        self.player_store_id.serialize(bit_stream)
        self.player_maze_id.serialize(bit_stream)
        self.village_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'create_date': self.create_date,
            'active_player_avatar': self.active_player_avatar.to_dict(),
            'home_theme_id': self.home_theme_id.to_dict(),
            'current_race_mode': self.current_race_mode.to_dict(),
            'workshop_options': self.workshop_options,
            'is_tutorial_completed': self.is_tutorial_completed,
            'yard_building_id': self.yard_building_id.to_dict(),
            'last_login': self.last_login,
            'play_time': self.play_time,
            'is_qa': self.is_qa,
            'home_village_plot_id': self.home_village_plot_id.to_dict(),
            'store_village_plot_id': self.store_village_plot_id.to_dict(),
            'player_store_id': self.player_store_id.to_dict(),
            'player_maze_id': self.player_maze_id.to_dict(),
            'village_id': self.village_id.to_dict(),
        }
