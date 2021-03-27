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

        if not message:
            log(LogLevel.WARN, '{} [bold yellow]<{}>[/] is not supported yet: {}'.format(
                peer_name, message_header.message_type, bytes(request_bs.data)))
            return

        if message_header.is_response:
            log(LogLevel.WARN, '{} [bold yellow]<{}>[/] responses from client are not supported yet: {}'.format(
                peer_name, message_header.message_type, bytes(request_bs.data)))
            return

        message.request.deserialize(request_bs)
        await message.process(message_header)
        message_header.is_response = True

        log(LogLevel.INFO, '{} Processed [bold blue]<{}>[/] with [bold blue]<{}> <{}>[/]'.format(
            peer_name, message_header.message_type, message_header.result_code, message_header.app_code))

        response_bs = BitStream()
        message_header.serialize(response_bs)
        message.response.serialize(response_bs)

        log(LogLevel.DEBUG, '=>', {
            'request_bs': bytes(request_bs.data),
            'response_bs': bytes(response_bs.data),
            'request': message.request.to_dict(),
            'response': message.response.to_dict()})

        return response_bs
