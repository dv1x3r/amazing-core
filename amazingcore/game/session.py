from amazingcore.logger import log, LogLevel
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.message_header import MessageHeader
from amazingcore.messages.message_factory import MessageFactory


class Session:

    def __init__(self):
        self.message_factory = MessageFactory()

    async def process_message(self, peer_name: str, data: bytearray) -> BitStream:
        request_bs = BitStream(data)
        message_header = MessageHeader()
        message_header.deserialize(request_bs)
        message = self.message_factory.build_message(message_header)
        if message:
            message.request.deserialize(request_bs)
            log(f'{peer_name} Request <{message_header.message_type}> {message.request}', LogLevel.TRACE)
            await message.process()
            log(f'{peer_name} Response <{message_header.message_type}> {message.response}', LogLevel.TRACE)
            response_bs = BitStream()
            message_header.serialize(response_bs)
            message.response.serialize(response_bs)
            return response_bs
        else:
            log(f'{peer_name} Request <{message_header.message_type}> is not supported yet: {bytes(request_bs.data)}', LogLevel.WARN)
