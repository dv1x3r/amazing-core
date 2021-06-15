from amazingcore.messages.common.home_theme import HomeTheme
from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.player_home import PlayerHome
from amazingcore.messages.common.player_maze import PlayerMaze

import datetime as dt


class InitLocationMessage(Message):
    def __init__(self):
        self.request: InitLocationRequest = InitLocationRequest()
        self.response: InitLocationResponse = InitLocationResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.sync_server_token = '1234'
        self.response.sync_server_ip = 'localhost'
        self.response.sync_server_port = 8182

        object_id = ObjectID(0, 0, 0, 0)

        home_theme = HomeTheme(object_id, {
            # LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()
            'Scene_Unity3D': [Asset(object_id, 'asset_type', 'non_existing_cdn_id', 'HomeLotSmall.unity3d', 'Main_Scene', 0)]
            # 'Scene_Unity3D': [Asset(object_id, 'asset_type', 'non_existing_cdn_id', 'Springtime003.unity3d', 'Main_Scene', 0)]
        }, [])

        player_maze = PlayerMaze(
            object_id, 'coremaze', 42, None, dt.datetime.now(
            ), 1, 1, 1, False, True, True, False, object_id, [], home_theme, object_id, object_id
        )

        self.response.home = PlayerHome(
            player_maze, 'kek', True, dt.datetime.now(), home_theme, object_id, [player_maze])


class InitLocationRequest(SerializableMessage):
    def __init__(self):
        self.loc_id: ObjectID = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.loc_id = ObjectID()
        self.loc_id.deserialize(bit_stream)

    def to_dict(self):
        return {
            'loc_id': self.loc_id.to_dict(),
        }


class InitLocationResponse(SerializableMessage):
    def __init__(self):
        self.zone_instance = None  # GSFZoneInstance
        self.village = None  # GSFVillage
        self.home: PlayerHome = None
        self.sync_server_token: str = None
        self.sync_server_ip: str = None
        self.sync_server_port: int = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_none()  # GSFZoneInstance
        bit_stream.write_none()  # GSFVillage
        self.home.serialize(bit_stream)
        bit_stream.write_str(self.sync_server_token)
        bit_stream.write_str(self.sync_server_ip)
        bit_stream.write_int(self.sync_server_port)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'zone_instance': self.zone_instance.to_dict() if self.zone_instance else None,
            'village': self.village.to_dict() if self.village else None,
            'home': self.home.to_dict() if self.home else None,
            'sync_server_token': self.sync_server_token,
            'sync_server_ip': self.sync_server_ip,
            'sync_server_port': self.sync_server_port,
        }
