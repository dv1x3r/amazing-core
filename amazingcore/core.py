import asyncio
from asyncio.streams import StreamReader, StreamWriter
from amazingcore.logger import LogLevel, log
from amazingcore.codec.bit_protocol import BitProtocol
from amazingcore.session import Session


class Core:

    def __init__(self):
        self.bit_protocol = BitProtocol()

    async def main(self, host, port):
        tcp_server = await asyncio.start_server(self.client_connected, host, port)
        log(LogLevel.INFO, 'Server is listening for connections...')
        await tcp_server.serve_forever()

    async def client_connected(self, reader: StreamReader, writer: StreamWriter):
        peer_name = writer.transport.get_extra_info('peername')
        log(LogLevel.INFO, f'{peer_name} connected')
        session = Session()
        try:
            while True:
                try:
                    data = await self.bit_protocol.read_data(reader)
                    response = await session.process_message(peer_name, data)
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
