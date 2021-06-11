from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.player_avatar import PlayerAvatar
from amazingcore.messages.common.avatar import Avatar
from amazingcore.messages.common.asset import Asset

import datetime as dt


class GetAvatarsMessage(Message):
    def __init__(self):
        self.request: GetAvatarsRequest = GetAvatarsRequest()
        self.response: GetAvatarsResponse = GetAvatarsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        object_id = ObjectID(0, 0, 0, 0)

        avatar = Avatar(
            aw_object_id=object_id,
            asset_map={
                "Prefab_Unity3D": [
                    Asset(ObjectID(0, 0, 0, 0), 'asset_type', 'assets/Player_Avatar.unity3d',
                          'Player_Avatar.unity3d', 'asset_group', 59109),  # !(item.resName == "PF__Avatar.unity3d")
                ],
            },
            asset_packages=[],
            dimensions='1',
            weight=1,
            height=1,
            max_outfits=1,
            name='avatar_name')

        # activePlayerAvatar should be in AvatarManager.Instance.GSFPlayerAvatars: LoadAvatarsCommand.cs -> Step2()
        self.response.avatars = [
            PlayerAvatar(
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
        ]


class GetAvatarsRequest(SerializableMessage):
    def __init__(self):
        self.start: int = None
        self.max: int = None
        self.filter_ids: list[ObjectID] = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.start = bit_stream.read_int()
        self.max = bit_stream.read_int()
        self.filter_ids = []
        for _ in range(bit_stream.read_int()):
            item = ObjectID()
            item.deserialize(bit_stream)
            self.filter_ids.append(item)

    def to_dict(self):
        return {
            'start': self.start,
            'max': self.max,
            'filter_ids': [item.to_dict() for item in self.filter_ids],
        }


class GetAvatarsResponse(SerializableMessage):
    def __init__(self):
        self.avatars: list[PlayerAvatar] = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.avatars))
        for item in self.avatars:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'avatars': [item.to_dict() for item in self.avatars],
        }
