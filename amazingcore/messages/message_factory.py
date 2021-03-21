from amazingcore.messages.user.selected_player_name import SelectedPlayerNameMessage
from amazingcore.messages.user.random_names import RandomNamesMessage
from amazingcore.messages.user.validate_name import ValidateNameMessage
from amazingcore.messages.user.client_version import ClientVersionMessage
from amazingcore.messages.message_header import MessageHeader
from amazingcore.messages.message_codes import ServiceClass, UserMessageTypes


class MessageFactory:

    def build_message(self, message_header: MessageHeader):
        if(message_header.service_class == ServiceClass.USER_SERVER):
            return self.__user__(message_header)
        if(message_header.service_class == ServiceClass.SYNC_SERVER):
            return self.__sync__(message_header)
        if(message_header.service_class == ServiceClass.LOCATION):
            return self.__location__(message_header)
        if(message_header.service_class == ServiceClass.CLIENT):
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

    def __sync__(self, message_header: MessageHeader):
        pass

    def __location__(self, message_header: MessageHeader):
        pass

    def __client__(self, message_header: MessageHeader):
        pass
