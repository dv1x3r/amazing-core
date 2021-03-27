from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.message_codes import ServiceClass, UserMessageTypes, SyncMessageTypes, ClientMessageTypes, ResultCode, AppCode


class MessageHeader:

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
        self.request_id: int = request_id
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

    @result_code.setter
    def result_code(self, value: ResultCode):
        self.__result_code__ = value.value

    @property
    def app_code(self):
        return AppCode(self.__app_code__)

    @app_code.setter
    def app_code(self, value: AppCode):
        self.__app_code__ = value.value

    @property
    def is_service(self):
        return self.__flags__ & 2 == 0

    @property
    def is_response(self):
        return self.is_service & (self.__flags__ & 1) != 0

    @is_response.setter
    def is_response(self, value: bool):
        if value:
            self.__flags__ |= 1
        else:
            self.__flags__ &= ~1

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
            bit_stream.write_int(self.request_id)
        if self.is_request:
            bit_stream.write_str(None)  # log_correlator (always empty)
        if self.is_response:
            bit_stream.write_int(self.__result_code__)
            bit_stream.write_int(self.__app_code__)
            if self.__app_code__ != 0:
                app_code_name = self.app_code.name if self.__app_code__ is not None else ''
                bit_stream.write_str(app_code_name)
            if self.__app_code__ == 17:
                bit_stream.write_int(0)  # app_codes size (always empty)

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():  # GSFMessage
            return
        if not bit_stream.read_start():  # MessageHeader
            return
        self.__flags__ = bit_stream.read_int()
        self.__service_class__ = bit_stream.read_int()
        self.__message_type__ = bit_stream.read_int()
        if self.is_service:
            self.request_id = bit_stream.read_int()
        if self.is_request:
            bit_stream.read_str()  # log_correlator (always empty)

    def to_dict(self):
        return {
            'is_service': self.is_service,
            'is_response': self.is_response,
            'is_request': self.is_request,
            'is_notify': self.is_notify,
            'is_discardable': self.is_discardable,
            'service_class': self.service_class if self.__service_class__ is not None else None,
            'message_type': self.message_type if self.__message_type__ is not None else None,
            'request_id': self.request_id,
            'result_code': self.result_code if self.__result_code__ is not None else None,
            'app_code': self.app_code if self.__app_code__ is not None else None,
        }
