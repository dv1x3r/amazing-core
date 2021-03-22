import asyncio
from asyncio.streams import StreamReader, StreamWriter
import traceback
from amazingcore.logger import LogLevel, log
from amazingcore.codec.bit_protocol import BitProtocol
from amazingcore.session import Session


class Core:

    def __init__(self):
        self.bit_protocol = BitProtocol()

    async def main(self, host, port):
        tcp_server = await asyncio.start_server(self.client_connected, host, port)
        log('server is listening for connections...', LogLevel.INFO)
        await tcp_server.serve_forever()

    async def client_connected(self, reader: StreamReader, writer: StreamWriter):
        peer_name = writer.transport.get_extra_info('peername')
        log(f'{peer_name} connected', LogLevel.INFO)
        session = Session()
        try:
            while True:
                try:
                    data = await self.bit_protocol.read_data(reader)
                    response = await session.process_message(peer_name, data)
                    if response:
                        await self.bit_protocol.write_message(writer, response.data)
                except NotImplementedError as err:
                    log(f'{peer_name} {err}', LogLevel.ERROR)
                    traceback.print_exc()
        except ConnectionError as err:
            log(f'{peer_name} disconnected: {err}', LogLevel.INFO)
        except Exception as err:
            log(f'{peer_name} disconnected: {err}', LogLevel.FATAL)
            traceback.print_exc()
        finally:
            writer.close()
