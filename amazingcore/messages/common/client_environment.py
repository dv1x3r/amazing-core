from amazingcore.messages.message_interfaces import SerializableMessage
from amazingcore.codec.bit_stream import BitStream
import datetime as dt


class ClientEnvironment(SerializableMessage):

    def __init__(self):
        self.unity_version: str = None
        self.user_agent: str = None
        self.screen_resolution: str = None
        self.machine_os: str = None
        self.user_time: dt.datetime = None
        self.utc_offset_in_minutes: int = None
        self.ip_address: str = None

    def serialize(self, bit_stream: BitStream):
        raise NotImplementedError

    def deserialize(self, bit_stream: BitStream):
        if not bit_stream.read_start():
            return
        self.unity_version = bit_stream.read_str()
        self.user_agent = bit_stream.read_str()
        self.screen_resolution = bit_stream.read_str()
        self.machine_os = bit_stream.read_str()
        self.user_time = bit_stream.read_dt()
        self.utc_offset_in_minutes = bit_stream.read_int()
        self.ip_address = bit_stream.read_str()

    def to_dict(self):
        return {
            'unity_version': self.unity_version,
            'user_agent': self.user_agent,
            'screen_resolution': self.screen_resolution,
            'machine_os': self.machine_os,
            'user_time': self.user_time,
            'utc_offset_in_minutes': self.utc_offset_in_minutes,
            'ip_address': self.ip_address,
        }
