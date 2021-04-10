from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.client_environment import ClientEnvironment
from amazingcore.messages.common.site_info import SiteInfo
from amazingcore.messages.common.session_status import SessionStatus
from amazingcore.messages.common.player import Player
from amazingcore.messages.common.player_stats import PlayerStats
from amazingcore.messages.common.player_info import PlayerInfo
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackage
from amazingcore.messages.common.race_mode import RaceMode
from amazingcore.messages.common.crisp_data import CrispData
from amazingcore.messages.common.player_setting import PlayerSetting

import datetime as dt


class LoginMessage(Message):
    def __init__(self):
        self.request: LoginRequest = LoginRequest()
        self.response: LoginResponse = LoginResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        object_id = ObjectID(1, 2, 3, 4)

        self.response.site_info = SiteInfo(
            site_id=object_id,
            nickname_first='nickname_first',
            nickname_last='nickname_last',
            site_user_id=object_id)

        self.response.session_status = SessionStatus.IN_PROGRESS
        self.response.session_id = object_id
        self.response.conversation_id = 5
        self.response.asset_delivery_url = 'localhost'

        asset = Asset(
            aw_object_id=object_id,
            asset_type_name='asset_type_name',
            cdn_id='cdn_id',
            res_name='res_name',
            group_name='group_name',
            file_size=1024)

        apc = AssetPackage(
            container_aw_object_id=object_id,
            container_asset_map={},
            container_asset_pkg=[],
            p_tag='p_tag',
            create_date=dt.datetime.now())

        avatar = Avatar(
            aw_object_id=object_id,
            asset_map={'asset_map_key': [asset]},
            asset_packages=[apc],
            dimensions='6',
            weight=7,
            height=8,
            max_outfits=9,
            name='avatar_name')

        player_avatar = PlayerAvatar(
            aw_object_id=object_id,
            avatar=avatar,
            player_id=object_id,
            name='player_name',
            bio='player_bio',
            secret_code='1234',
            create_ts=dt.datetime.now(),
            player_avatar_outfit_id=object_id,
            outfit_no=10,
            play_time=11,
            last_play=dt.datetime.now())

        race_mode = RaceMode(aw_object_id=object_id, name='race_name')

        self.response.player = Player(
            aw_object_id=object_id,
            create_date=dt.datetime.now(),
            active_player_avatar=player_avatar,
            home_theme_id=object_id,
            current_race_mode=race_mode,
            workshop_options='workshop_options',
            is_tutorial_completed=False,
            yard_building_id=object_id,
            last_login=dt.datetime.now(),
            play_time=12,
            is_qa=True,
            home_village_plot_id=object_id,
            store_village_plot_id=object_id,
            player_store_id=object_id,
            player_maze_id=object_id,
            village_id=object_id)

        self.response.max_outfit = 13

        player_stats_1 = PlayerStats(
            aw_object_id=object_id,
            player_avatar_id=object_id,
            stats_type_id=object_id,
            level=14,
            object_id=object_id)

        self.response.player_stats = [player_stats_1]

        crisp_data = CrispData(
            crisp_action_id=object_id,
            crisp_message='crisp_message',
            crisp_expiry_date=dt.datetime.now() + dt.timedelta(days=1),
            crisp_confirmed=True)

        player_setting_1 = PlayerSetting(
            aw_object_id=object_id,
            player_id=object_id,
            setting_name='setting_name',
            value='value')

        self.response.player_info = PlayerInfo(
            tier_id=object_id,
            player_name='omnio',
            world_name='coreworld',
            crisp_data=crisp_data,
            verified=True,
            verification_expiry=dt.datetime.now() + dt.timedelta(days=1),
            scs_block_expiry=dt.datetime.now() + dt.timedelta(days=1),
            eula_id=object_id,
            current_eula_id=object_id,
            u_13=True,
            chat_blocked_parent=False,
            chat_allowed=True,
            chat_blocked_expiry=dt.datetime.now() + dt.timedelta(days=1),
            findable=True,
            findable_date=dt.datetime.now() + dt.timedelta(days=1),
            user_subscription_expiry_date=dt.datetime.now() + dt.timedelta(days=1),
            qa=True,
            player_settings=[player_setting_1])

        self.response.current_server_time = dt.datetime.now()
        self.response.system_lockout_time = dt.datetime.now() + dt.timedelta(days=1)
        self.response.system_shutdown_time = dt.datetime.now() + dt.timedelta(days=1)
        self.response.client_inactivity_timeout = 365
        self.response.cnl = 'cnl'


class LoginRequest(SerializableMessage):
    def __init__(self):
        self.login_id: str = None
        self.password: str = None
        self.site_pin: int = None
        self.language_local_pair_id: ObjectID = None
        self.user_queueing_token: str = None
        self.client_environment: ClientEnvironment = None
        self.token: str = None
        self.login_type: int = None
        self.cnl: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.login_id = bit_stream.read_str()
        self.password = bit_stream.read_str()
        self.site_pin = bit_stream.read_int()
        self.language_local_pair_id = ObjectID()
        self.language_local_pair_id.deserialize(bit_stream)
        self.user_queueing_token = bit_stream.read_str()
        self.client_environment = ClientEnvironment()
        self.client_environment.deserialize(bit_stream)
        self.token = bit_stream.read_str()
        self.login_type = bit_stream.read_int()
        self.cnl = bit_stream.read_str()

    def to_dict(self):
        return {
            'login_id': self.login_id,
            'password': self.password,
            'site_pin': self.site_pin,
            'language_local_pair_id': self.language_local_pair_id.to_dict(),
            'user_queueing_token': self.user_queueing_token,
            'client_environment': self.client_environment.to_dict(),
            'token': self.token,
            'login_type': self.login_type,
            'cnl': self.cnl,
        }


class LoginResponse(SerializableMessage):
    def __init__(self):
        self.site_info: SiteInfo = None,
        self.session_status: SessionStatus = None,
        self.session_id: ObjectID = None,
        self.conversation_id: int = None,
        self.asset_delivery_url: str = None,
        self.player: Player = None,
        self.max_outfit: int = None,
        self.player_stats: list[PlayerStats] = None,
        self.player_info: PlayerInfo = None,
        self.current_server_time: dt.datetime = None,
        self.system_lockout_time: dt.datetime = None,
        self.system_shutdown_time: dt.datetime = None,
        self.client_inactivity_timeout: int = None,
        self.cnl: str = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.site_info.serialize(bit_stream)
        bit_stream.write_str(self.session_status.name)
        self.session_id.serialize(bit_stream)
        bit_stream.write_long(self.conversation_id)
        bit_stream.write_str(self.asset_delivery_url)
        self.player.serialize(bit_stream)
        bit_stream.write_short(self.max_outfit)
        bit_stream.write_int(len(self.player_stats))
        for ps in self.player_stats:
            ps.serialize(bit_stream)
        self.player_info.serialize(bit_stream)
        bit_stream.write_dt(self.current_server_time)
        bit_stream.write_dt(self.system_lockout_time)
        bit_stream.write_dt(self.system_shutdown_time)
        bit_stream.write_int(self.client_inactivity_timeout)
        bit_stream.write_str(self.cnl)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'site_info': self.site_info.to_dict(),
            'session_status': self.session_status,
            'session_id': self.session_id.to_dict(),
            'conversation_id': self.conversation_id,
            'asset_delivery_url': self.asset_delivery_url,
            'player': self.player.to_dict(),
            'max_outfit': self.max_outfit,
            'player_stats': [i.to_dict() for i in self.player_stats],
            'player_info': self.player_info.to_dict(),
            'current_server_time': self.current_server_time,
            'system_lockout_time': self.system_lockout_time,
            'system_shutdown_time': self.system_shutdown_time,
            'client_inactivity_timeout': self.client_inactivity_timeout,
            'cnl': self.cnl,
        }
