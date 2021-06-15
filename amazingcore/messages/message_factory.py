from amazingcore.messages.user.list_limits import ListLimitsMessage
from amazingcore.messages.user.get_player_npcs import GetPlayerNpcsMessage
from amazingcore.messages.user.get_npc_relationships import GetNpcRelationshipsMessage
from amazingcore.messages.user.get_online_statuses import GetOnlineStatusesMessage
from amazingcore.messages.user.enter_building import EnterBuildingMessage
from amazingcore.messages.user.get_announcements import GetAnnouncementsMessage
from amazingcore.messages.user.get_chat_channel_types import GetChatChannelTypesMessage
from amazingcore.messages.user.get_inventory_objects import GetInventoryObjectsMessage
from amazingcore.messages.user.get_build_objects import GetBuildObjectsMessage
from amazingcore.messages.user.get_maze_items import GetMazeItemsMessage
from amazingcore.messages.sync.sync_login import SyncLoginMessage
from amazingcore.messages.user.init_location import InitLocationMessage
from amazingcore.messages.user.get_zones import GetZonesMessage
from amazingcore.messages.user.get_outfits import GetOutfitsMessage
from amazingcore.messages.user.get_avatars import GetAvatarsMessage
from amazingcore.messages.user.get_outfit_items import GetOutfitItemsMessage
from amazingcore.messages.user.get_item_categories import GetItemCategoriesMessage
from amazingcore.messages.user.get_required_experience import GetRequiredExperienceMessage
from amazingcore.messages.user.get_site_frame import GetSiteFrameMessage
from amazingcore.messages.user.get_player_stats import GetPlayerStatsMessage
from amazingcore.messages.user.manage_home_invitations import ManageHomeInvitationsMessage
from amazingcore.messages.user.get_crisp_actions import GetCrispActionsMessage
from amazingcore.messages.user.get_cms_missions import GetCmsMissionsMessage
from amazingcore.messages.user.get_currencies import GetCurrenciesMessage
from amazingcore.messages.user.get_player_missions import GetPlayerMissionsMessage
from amazingcore.messages.user.get_npc_relationship_levels import GetNpcRelationshipLevelsMessage
from amazingcore.messages.user.get_home_invitations import GetHomeInvitationsMessage
from amazingcore.messages.user.get_notification_by_player_id import GetNotificationByPlayerIdMessage
from amazingcore.messages.user.get_cms_notifications import GetCmsNotificationsMessage
from amazingcore.messages.user.get_stats_type import GetStatsTypeMessage
from amazingcore.messages.user.get_tiers import GetTiersMessage
from amazingcore.messages.user.register_avatar_for_registration import RegisterAvatarForRegistrationMessage
from amazingcore.messages.user.register_player import RegisterPlayerMessage
from amazingcore.messages.user.check_username import CheckUsernameMessage
from amazingcore.messages.user.login import LoginMessage
from amazingcore.messages.user.selected_player_name import SelectedPlayerNameMessage
from amazingcore.messages.user.random_names import RandomNamesMessage
from amazingcore.messages.user.validate_name import ValidateNameMessage
from amazingcore.messages.user.client_version import ClientVersionMessage
from amazingcore.messages.message_header import MessageHeader
from amazingcore.messages.message_codes import ServiceClass, SyncMessageTypes, UserMessageTypes


class MessageFactory:

    def build_message(self, message_header: MessageHeader):
        if message_header.service_class == ServiceClass.USER_SERVER:
            return self.__user__(message_header)
        if message_header.service_class == ServiceClass.SYNC_SERVER:
            return self.__sync__(message_header)
        if message_header.service_class == ServiceClass.LOCATION:
            return self.__location__(message_header)
        if message_header.service_class == ServiceClass.CLIENT:
            return self.__client__(message_header)

    def __user__(self, message_header: MessageHeader):
        if message_header.message_type == UserMessageTypes.GET_CLIENT_VERSION_INFO:
            return ClientVersionMessage()
        if message_header.message_type == UserMessageTypes.VALIDATE_NAME:
            return ValidateNameMessage()
        if message_header.message_type == UserMessageTypes.GET_RANDOM_NAMES:
            return RandomNamesMessage()
        if message_header.message_type == UserMessageTypes.SELECT_PLAYER_NAME:
            return SelectedPlayerNameMessage()
        if message_header.message_type == UserMessageTypes.LOGIN:
            return LoginMessage()
        if message_header.message_type == UserMessageTypes.CHECK_USERNAME:
            return CheckUsernameMessage()
        if message_header.message_type == UserMessageTypes.REGISTER_PLAYER:
            return RegisterPlayerMessage()
        if message_header.message_type == UserMessageTypes.REGISTER_AVATAR_FOR_REGISTRATION:
            return RegisterAvatarForRegistrationMessage()
        if message_header.message_type == UserMessageTypes.GET_TIERS:
            return GetTiersMessage()
        if message_header.message_type == UserMessageTypes.GET_STATS_TYPE:
            return GetStatsTypeMessage()
        if message_header.message_type == UserMessageTypes.GET_CMS_NOTIFICATIONS:
            return GetCmsNotificationsMessage()
        if message_header.message_type == UserMessageTypes.GET_NOTIFICATION_BY_PLAYER_ID:
            return GetNotificationByPlayerIdMessage()
        if message_header.message_type == UserMessageTypes.GET_HOME_INVITATIONS:
            return GetHomeInvitationsMessage()
        if message_header.message_type == UserMessageTypes.GET_NPC_RELATIONSHIP_LEVELS:
            return GetNpcRelationshipLevelsMessage()
        if message_header.message_type == UserMessageTypes.GET_PLAYER_MISSIONS:
            return GetPlayerMissionsMessage()
        if message_header.message_type == UserMessageTypes.GET_CURRENCIES:
            return GetCurrenciesMessage()
        if message_header.message_type == UserMessageTypes.GET_CMS_MISSIONS:
            return GetCmsMissionsMessage()
        if message_header.message_type == UserMessageTypes.GET_CRISP_ACTIONS:
            return GetCrispActionsMessage()
        if message_header.message_type == UserMessageTypes.MANAGE_HOME_INVITATIONS:
            return ManageHomeInvitationsMessage()
        if message_header.message_type == UserMessageTypes.GET_PLAYER_STATS:
            return GetPlayerStatsMessage()
        if message_header.message_type == UserMessageTypes.GET_SITE_FRAME:
            return GetSiteFrameMessage()
        if message_header.message_type == UserMessageTypes.GET_REQUIRED_EXPERIENCE:
            return GetRequiredExperienceMessage()
        if message_header.message_type == UserMessageTypes.GET_CMS_ITEMCATEGORIES:
            return GetItemCategoriesMessage()
        if message_header.message_type == UserMessageTypes.GET_OUTFIT_ITEMS:
            return GetOutfitItemsMessage()
        if message_header.message_type == UserMessageTypes.GET_AVATARS:
            return GetAvatarsMessage()
        if message_header.message_type == UserMessageTypes.GET_OUTFITS:
            return GetOutfitsMessage()
        if message_header.message_type == UserMessageTypes.GET_ZONES:
            return GetZonesMessage()
        if message_header.message_type == UserMessageTypes.INIT_LOCATION:
            return InitLocationMessage()
        if message_header.message_type == UserMessageTypes.GET_MAZE_ITEMS:
            return GetMazeItemsMessage()
        if message_header.message_type == UserMessageTypes.GET_BUILD_OBJECTS:
            return GetBuildObjectsMessage()
        if message_header.message_type == UserMessageTypes.GET_INVENTORY_OBJECTS:
            return GetInventoryObjectsMessage()
        if message_header.message_type == UserMessageTypes.GET_CHAT_CHANNEL_TYPES:
            return GetChatChannelTypesMessage()
        if message_header.message_type == UserMessageTypes.GET_ANNOUNCEMENTS:
            return GetAnnouncementsMessage()
        if message_header.message_type == UserMessageTypes.ENTER_BUILDING:
            return EnterBuildingMessage()
        if message_header.message_type == UserMessageTypes.GET_ONLINE_STATUSES:
            return GetOnlineStatusesMessage()
        if message_header.message_type == UserMessageTypes.GET_NPC_RELATIONSHIPS:
            return GetNpcRelationshipsMessage()
        if message_header.message_type == UserMessageTypes.GET_PLAYER_NPCS:
            return GetPlayerNpcsMessage()
        if message_header.message_type == UserMessageTypes.LIST_LIMITS:
            return ListLimitsMessage()

    def __sync__(self, message_header: MessageHeader):
        if message_header.message_type == SyncMessageTypes.LOGIN:
            return SyncLoginMessage()

    def __location__(self, message_header: MessageHeader):
        pass

    def __client__(self, message_header: MessageHeader):
        pass
