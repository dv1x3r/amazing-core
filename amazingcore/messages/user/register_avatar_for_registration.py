from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset
from amazingcore.messages.common.asset_package import AssetPackageContainer

import datetime as dt


class RegisterAvatarForRegistrationMessage(Message):
    def __init__(self):
        self.request: RegisterAvatarForRegistrationRequest = RegisterAvatarForRegistrationRequest()
        self.response: RegisterAvatarForRegistrationResponse = RegisterAvatarForRegistrationResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.invalid_code_count = 0
        self.response.invalid_code_threshold = 0

        object_id = ObjectID(1, 2, 3, 4)

        asset = Asset(
            aw_object_id=object_id,
            asset_type_name='asset_type_name',
            cdn_id='cdn_id',
            res_name='res_name',
            group_name='group_name',
            file_size=1024)

        apc = AssetPackageContainer(
            p_tag='p_tag', create_date=dt.datetime.now())
        apc.container_aw_object_id = object_id
        apc.container_asset_map = {}
        apc.container_asset_pkg = []

        avatar = Avatar(
            aw_object_id=object_id,
            asset_map={'asset_map_key': [asset]},
            asset_packages=[apc],
            dimensions='6',
            weight=7,
            height=8,
            max_outfits=9,
            name='avatar_name')

        self.response.player_avatar = PlayerAvatar(
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


class RegisterAvatarForRegistrationRequest(SerializableMessage):
    def __init__(self):
        self.player_id: ObjectID = None
        self.secret_code: str = None
        self.name: str = None
        self.bio: str = None
        self.avatar_id: ObjectID = None
        self.given_inventory_ids: list[ObjectID] = None
        self.given_item_slot_ids: list[ObjectID] = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.player_id = ObjectID()
        self.player_id.deserialize(bit_stream)
        self.secret_code = bit_stream.read_str()
        self.name = bit_stream.read_str()
        self.bio = bit_stream.read_str()
        self.avatar_id = ObjectID()
        self.avatar_id.deserialize(bit_stream)
        self.given_inventory_ids = []
        for _ in range(bit_stream.read_int()):
            object_id = ObjectID()
            object_id.deserialize(bit_stream)
            self.given_inventory_ids.append(object_id)
        self.given_item_slot_ids = []
        for _ in range(bit_stream.read_int()):
            object_id = ObjectID()
            object_id.deserialize(bit_stream)
            self.given_item_slot_ids.append(object_id)

    def to_dict(self):
        return {
            'player_id': self.player_id.to_dict(),
            'secret_code': self.secret_code,
            'name': self.name,
            'bio': self.bio,
            'avatar_id': self.avatar_id.to_dict(),
            'given_inventory_ids': [item.to_dict() for item in self.given_inventory_ids],
            'given_item_slot_ids': [item.to_dict() for item in self.given_item_slot_ids],
        }


class RegisterAvatarForRegistrationResponse(SerializableMessage):
    def __init__(self):
        self.player_avatar: PlayerAvatar = None,
        self.invalid_code_count: int = None,
        self.invalid_code_threshold: int = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.player_avatar.serialize(bit_stream)
        bit_stream.write_int(self.invalid_code_count)
        bit_stream.write_int(self.invalid_code_threshold)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'player_avatar': self.player_avatar.to_dict(),
            'invalid_code_count': self.invalid_code_count,
            'invalid_code_threshold': self.invalid_code_threshold,
        }
