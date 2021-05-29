from amazingcore.messages.message_interfaces import Message, SerializableMessage
from amazingcore.messages.message_codes import ResultCode, AppCode
from amazingcore.messages.message_header import MessageHeader
from amazingcore.codec.bit_stream import BitStream


class GetRequiredExperienceMessage(Message):
    def __init__(self):
        self.request: GetRequiredExperienceRequest = GetRequiredExperienceRequest()
        self.response: GetRequiredExperienceResponse = GetRequiredExperienceResponse()

    async def process(self, message_header: MessageHeader):
        message_header.result_code = ResultCode.OK
        message_header.app_code = AppCode.OK


class GetRequiredExperienceRequest(SerializableMessage):
    def __init__(self):
        self.level: int = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.level = bit_stream.read_short()

    def to_dict(self):
        return {'level': self.level}


class GetRequiredExperienceResponse(SerializableMessage):
    def __init__(self):
        self.min_experience: int = None
        self.max_experience: int = None
        self.min_time_played: int = None
        self.max_time_played: int = None
        self.player_experience_level = None

    def serialize(self, bit_stream: BitStream):
        bit_stream.write_start()
        bit_stream.write_int(self.min_experience)
        bit_stream.write_int(self.max_experience)
        bit_stream.write_int(self.min_time_played)
        bit_stream.write_int(self.max_time_played)
        bit_stream.write_none()  # GSFPlayerExperienceLevel

    def deserialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def to_dict(self):
        return {
            'min_experience': self.min_experience,
            'max_experience': self.max_experience,
            'min_time_played': self.min_time_played,
            'max_time_played': self.max_time_played,
            'player_experience_level': self.player_experience_level.to_dict() if self.player_experience_level else None,
        }
