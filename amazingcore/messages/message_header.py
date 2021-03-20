from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.messages.message_codes import ServiceClass, UserMessageTypes, SyncMessageTypes, ClientMessageTypes, ResultCode, AppCode


class MessageHeader(SerializableMessage):

    def __init__(self,
                 flags: int = None,
                 service_class: int = None,
                 message_type: int = None,
                 request_id: int = None,
                 result_code: int = None,
                 app_code: int = None):
        self.__flags__: int = flags
        self.__service_class__: int = service_class
        self.__message_type__: int = message_type
        self.__request_id__: int = request_id
        self.__result_code__: int = result_code
        self.__app_code__: int = app_code

    @property
    def service_class(self):
        return ServiceClass(self.__service_class__)

    @property
    def message_type(self):
        if self.service_class == ServiceClass.USER_SERVER:
            return UserMessageTypes(self.__message_type__)
        elif self.service_class == ServiceClass.SYNC_SERVER:
            return SyncMessageTypes(self.__message_type__)
        elif self.service_class == ServiceClass.CLIENT:
            return ClientMessageTypes(self.__message_type__)
        elif self.service_class == ServiceClass.LOCATION:
            return None  # 404 Type Not Found ;P

    @property
    def result_code(self):
        return ResultCode(self.__result_code__)

    @property
    def app_code(self):
        return AppCode(self.__app_code__)

    @property
    def is_service(self):
        return self.__flags__ & 2 == 0

    @property
    def is_response(self):
        return self.is_service & (self.__flags__ & 1) != 0

    @property
    def is_request(self):
        return self.is_service & (not self.is_response)

    @property
    def is_notify(self):
        return self.__flags__ & 2 != 0

    @property
    def is_discardable(self):
        return self.__flags__ & 16 != 0

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()  # GSFMessage
        bit_stream.write_start()  # MessageHeader
        bit_stream.write_int(self.__flags__)
        bit_stream.write_int(self.__service_class__)
        bit_stream.write_int(self.__message_type__)
        if self.is_service:
            bit_stream.write_int(self.__request_id__)
        if self.is_request:
            bit_stream.write_str(None)  # log_correlator (always empty)
        if self.is_response:
            bit_stream.write_int(self.__result_code__)
            bit_stream.write_int(self.__app_code__)
            if self.__app_code__ != 0:
                bit_stream.write_str(self.app_code.name)
            if self.__app_code__ == 17:
                bit_stream.write_int(0)  # app_codes size (always empty)

    def deserialize(self, bit_stream: BitStream):
        bit_stream.read_start()  # GSFMessage
        bit_stream.read_start()  # MessageHeader
        self.__flags__ = bit_stream.read_int()
        self.__service_class__ = bit_stream.read_int()
        self.__message_type__ = bit_stream.read_int()
        if self.is_service:
            self.__request_id__ = bit_stream.read_int()
        if self.is_request:
            bit_stream.read_str()  # log_correlator (always empty)

    def __str__(self):
        return str({
            'is_service': self.is_service,
            'is_response': self.is_response,
            'is_request': self.is_request,
            'is_notify': self.is_notify,
            'is_discardable': self.is_discardable,
            'service_class': self.service_class,
            'message_type': self.message_type,
            'request_id': self.__request_id__,
        })
