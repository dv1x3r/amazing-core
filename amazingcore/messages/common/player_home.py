from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_maze import PlayerMaze
from amazingcore.messages.common.home_theme import HomeTheme

import datetime as dt


class PlayerHome(SerializableMessage):
    def __init__(self,
                 player_maze: PlayerMaze = None,
                 player_name: str = None,
                 findable: bool = None,
                 findable_date: dt.datetime = None,
                 home_theme: HomeTheme = None,
                 player_id: ObjectID = None,
                 player_mazes: list[PlayerMaze] = None):
        self.player_maze = player_maze
        self.player_name = player_name
        self.findable = findable
        self.findable_date = findable_date
        self.home_theme = home_theme
        self.player_id = player_id
        self.player_mazes = player_mazes

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        if self.player_maze:
            self.player_maze.serialize(bit_stream)
        else:
            bit_stream.write_none()
        bit_stream.write_str(self.player_name)
        bit_stream.write_bool(self.findable)
        bit_stream.write_dt(self.findable_date)
        if self.home_theme:
            self.home_theme.serialize(bit_stream)
        else:
            bit_stream.write_none()
        if self.player_id:
            self.player_id.serialize(bit_stream)
        else:
            bit_stream.write_none()
        bit_stream.write_int(len(self.player_mazes))
        for item in self.player_mazes:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'player_maze': self.player_maze.to_dict() if self.player_maze else None,
            'player_name': self.player_name,
            'findable': self.findable,
            'findable_date': self.findable_date,
            'home_theme': self.home_theme.to_dict() if self.home_theme else None,
            'player_id': self.player_id.to_dict() if self.player_id else None,
            'player_mazes': [item.to_dict() for item in self.player_mazes],
        }
