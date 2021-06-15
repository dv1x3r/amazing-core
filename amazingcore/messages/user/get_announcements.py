from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class GetAnnouncementsMessage(Message):
    def __init__(self):
        self.request: GetAnnouncementsRequest = GetAnnouncementsRequest()
        self.response: GetAnnouncementsResponse = GetAnnouncementsResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK

        self.response.announcements = []


class GetAnnouncementsRequest(SerializableMessage):
    def __init__(self):
        self.un_marked: bool = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.un_marked = bit_stream.read_bool()

    def to_dict(self):
        return {'un_marked': self.un_marked}


class GetAnnouncementsResponse(SerializableMessage):
    def __init__(self):
        self.announcements: list = None  # GSFAnnouncement

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(len(self.announcements))
        for item in self.announcements:
            item.serialize(bit_stream)

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {'announcements': [item.to_dict() for item in self.announcements]}
