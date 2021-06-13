from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.home_theme import HomeTheme

import datetime as dt


class PlayerMaze(SerializableMessage):
    def __init__(self,
                 aw_object_id: ObjectID = None,
                 name: str = None,
                 size: int = None,
                 thumbnail: bytes = None,
                 publish_timestamp: dt.datetime = None,
                 num_rooms: int = None,
                 num_tubes: int = None,
                 rating: int = None,
                 is_locked: bool = None,
                 is_home_maze: bool = None,
                 is_published: bool = None,
                 is_publish_expired: bool = None,
                 player_id: ObjectID = None,
                 maze_pieces: list = None,  # GSFPlayerMazePiece
                 home_theme: HomeTheme=None,
                 parent_id: ObjectID = None,
                 source_id: ObjectID = None):
        self.aw_object_id = aw_object_id
        self.name = name
        self.size = size
        self.thumbnail = thumbnail
        self.publish_timestamp = publish_timestamp
        self.num_rooms = num_rooms
        self.num_tubes = num_tubes
        self.rating = rating
        self.is_locked = is_locked
        self.is_home_maze = is_home_maze
        self.is_published = is_published
        self.is_publish_expired = is_publish_expired
        self.player_id = player_id
        self.maze_pieces = maze_pieces
        self.home_theme = home_theme
        self.parent_id = parent_id
        self.source_id = source_id

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.aw_object_id.serialize(bit_stream)
        bit_stream.write_str(self.name)
        bit_stream.write_long(self.size)
        bit_stream.write_bytes(self.thumbnail)
        bit_stream.write_dt(self.publish_timestamp)
        bit_stream.write_short(self.num_rooms)
        bit_stream.write_short(self.num_tubes)
        bit_stream.write_short(self.rating, True)
        bit_stream.write_bool(self.is_locked)
        bit_stream.write_bool(self.is_home_maze)
        bit_stream.write_bool(self.is_published)
        bit_stream.write_bool(self.is_publish_expired)
        self.player_id.serialize(bit_stream)
        bit_stream.write_int(0)  # GSFPlayerMazePieces
        self.home_theme.serialize(bit_stream)
        self.parent_id.serialize(bit_stream)
        self.source_id.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'aw_object_id': self.aw_object_id.to_dict(),
            'name': self.name,
            'size': self.size,
            'thumbnail': self.thumbnail,
            'publish_timestamp': self.publish_timestamp,
            'num_rooms': self.num_rooms,
            'num_tubes': self.num_tubes,
            'rating': self.rating,
            'is_locked': self.is_locked,
            'is_home_maze': self.is_home_maze,
            'is_published': self.is_published,
            'is_publish_expired': self.is_publish_expired,
            'player_id': self.player_id.to_dict(),
            'maze_pieces': [i.to_dict() for i in self.maze_pieces],
            'home_theme': self.home_theme.to_dict(),
            'parent_id': self.parent_id.to_dict(),
            'source_id': self.source_id.to_dict(),
        }
