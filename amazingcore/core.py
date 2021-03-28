import asyncio
from asyncio.streams import StreamReader, StreamWriter
from amazingcore.logger import LogLevel, log
from amazingcore.codec.bit_protocol import BitProtocol
from amazingcore.codec.bit_stream import BitStream
from amazingcore.messages.message_header import MessageHeader
from amazingcore.messages.message_factory import MessageFactory


class Core:

    def __init__(self):
        self.bit_protocol = BitProtocol()
        self.message_factory = MessageFactory()

    async def main(self, host, port):
        tcp_server = await asyncio.start_server(self.client_connected, host, port)
        log(LogLevel.INFO, 'Server is listening for connections...')
        await tcp_server.serve_forever()

    async def client_connected(self, reader: StreamReader, writer: StreamWriter):
        peer_name = writer.transport.get_extra_info('peername')
        log(LogLevel.INFO, f'{peer_name} connected')
        try:
            while True:
                try:
                    data = await self.bit_protocol.read_data(reader)
                    response = await self.process_message(peer_name, data)
                    if response:
                        await self.bit_protocol.write_message(writer, response.data)
                except NotImplementedError as err:
                    log(LogLevel.ERROR, f'{peer_name}')
        except ConnectionError as err:
            log(LogLevel.INFO, f'{peer_name} disconnected: {err}')
        except Exception as err:
            log(LogLevel.FATAL, f'{peer_name} disconnected')
        finally:
            writer.close()

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
