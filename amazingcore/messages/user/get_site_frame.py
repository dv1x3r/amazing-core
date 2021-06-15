from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.common.object_id import ObjectID
from amazingcore.messages.common.site_frame import SiteFrame
from amazingcore.messages.common.asset import Asset


import datetime as dt


class GetSiteFrameMessage(Message):
    def __init__(self):
        self.request: GetSiteFrameRequest = GetSiteFrameRequest()
        self.response: GetSiteFrameResponse = GetSiteFrameResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.site_frame = SiteFrame(
            aw_object_id=ObjectID(0, 0, 0, 0),
            asset_map={
                # this is used to load hardcoded assets (instead of using Resources.Load())
                # LoadLoginScene.cs -> LoadAvatar -> DownloadManager.LoadAsset("Player_Base.unity3d")
                # Should contain child object with animation (PlayerController GetAnimObject checks for child objects)
                'Amazing_Core': [Asset(ObjectID(0, 0, 0, 0), 'asset_type', 'assets/Player_Base.unity3d', 'Player_Base.unity3d', 'asset_group', 59051)],

                # LoadLoginScene.cs -> AvatarLoadedHandler() -> LoadSlotIds()
                # DressAvatarManager.cs -> LoadSlotIds -> ClientManager.Instance.configList
                'Config_Text': [],

                # preload objects are downloaded to the Cached folder ClientManager.cs -> LoadPreloadedAssets()
                # also used to pass AWLoadingScreen: OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> preloadList
                'Preload_PrefabUnity3D': [],

                # also used to pass AWLoadingScreen: OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> audioClipList
                'Audio': [],
            },
            asset_packages=[],
            type_value=1
        )

        self.response.asset_delivery_url = 'http://localhost:8080/'  # + asset.cdn_id


class GetSiteFrameRequest(SerializableMessage):
    def __init__(self):
        self.type_value: int = None
        self.lang_locale_pair_id: ObjectID = None
        self.tier_id: ObjectID = None
        self.birth_date: dt.datetime = None
        self.registration_date: dt.datetime = None
        self.preview_date: dt.datetime = None
        self.is_preview_enabled: bool = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.type_value = bit_stream.read_int()
        self.lang_locale_pair_id = ObjectID()
        self.lang_locale_pair_id.deserialize(bit_stream)
        self.tier_id = ObjectID()
        self.tier_id.deserialize(bit_stream)
        self.birth_date = bit_stream.read_dt()
        self.registration_date = bit_stream.read_dt()
        self.preview_date = bit_stream.read_dt()
        self.is_preview_enabled = bit_stream.read_bool()

    def to_dict(self):
        return {
            'type_value': self.type_value,
            'lang_locale_pair_id': self.lang_locale_pair_id.to_dict(),
            'tier_id': self.tier_id.to_dict(),
            'birth_date': self.birth_date,
            'registration_date': self.registration_date,
            'preview_date': self.preview_date,
            'is_preview_enabled': self.is_preview_enabled,
        }


class GetSiteFrameResponse(SerializableMessage):
    def __init__(self,
                 site_frame: SiteFrame = None,
                 asset_delivery_url: str = None):
        self.site_frame = site_frame
        self.asset_delivery_url = asset_delivery_url

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        self.site_frame.serialize(bit_stream)
        bit_stream.write_str(self.asset_delivery_url)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'site_frame': self.site_frame.to_dict(),
            'asset_delivery_url': self.asset_delivery_url,
        }
