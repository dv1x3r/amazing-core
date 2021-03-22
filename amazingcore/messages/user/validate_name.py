from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class ValidateNameMessage(Message):
    def __init__(self):
        self.request: ValidateNameRequest = ValidateNameRequest()
        self.response: ValidateNameResponse = ValidateNameResponse()

    async def process(self, message_header: MessageHeader):
        # restrict new username patterns
        if self.request.name.lower().find('fuck') != -1:
            self.response.filter_name = 'bad'
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class ValidateNameRequest(SerializableMessage):
    """
    Name check for new player registration
    """

    def __init__(self):
        self.name: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        bit_stream.read_start()
        self.name = bit_stream.read_str()

    def __str__(self):
        return str({'name': self.name})


class ValidateNameResponse(SerializableMessage):
    """
    If filter_name is not None then 'not appropriate'
    """

    def __init__(self, filter_name: str = None):
        self.filter_name = filter_name

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_str(self.filter_name)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def __str__(self):
        return str({'filter_name': self.filter_name})
