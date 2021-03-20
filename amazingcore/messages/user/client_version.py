from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.codec.bit_stream import BitStream


class ClientVersionMessage(Message):

    def __init__(self):
        self.request: ClientVersionRequest = ClientVersionRequest()
        self.response: ClientVersionResponse = ClientVersionResponse()

    async def process(self):
        if self.request.client_name == 'AmazingWorld':
            self.response.client_version = '133852.true'
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
        bit_stream.read_start()
        self.client_name = bit_stream.read_str()

    def __str__(self):
        return str({'client_name': self.client_name})


class ClientVersionResponse(SerializableMessage):
    """
    Client Version format: "ClientVersion.ForceUpdate" \n
    For example <133852.true> stands for the latest version \n
    With <133853.true> game will require to update
    """

    def __init__(self,):
        self.client_version: str = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_str(self.client_version)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def __str__(self):
        return str({'client_version': self.client_version})
