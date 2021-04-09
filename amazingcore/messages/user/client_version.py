from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class ClientVersionMessage(Message):

    def __init__(self):
        self.request: ClientVersionRequest = ClientVersionRequest()
        self.response: ClientVersionResponse = ClientVersionResponse()

    async def process(self, message_header: MessageHeader):
        if self.request.client_name == 'AmazingWorld':
            self.response.client_version = '133852.true'
            message_header.result_code = ResultCode.OK
            message_header.app_code = AppCode.OK
        else:
            raise ValueError('invalid client_name')


class ClientVersionRequest(SerializableMessage):
    """
    Client Name should always be "AmazingWorld"
    """

    def __init__(self):
        self.client_name: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.client_name = bit_stream.read_str()

    def to_dict(self):
        return {'client_name': self.client_name}


class ClientVersionResponse(SerializableMessage):
    """
    Client Version format: "ClientVersion.ForceUpdate" \n
    For example <133852.true> stands for the latest version \n
    With <133853.true> game will require to update
    """

    def __init__(self):
        self.client_version: str = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_str(self.client_version)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'client_version': self.client_version}
