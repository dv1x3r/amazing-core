import asyncio
from aiohttp import web
from amazingcore.logger import LogLevel, log


class Cdn:

    async def handler(self, request: web.Request):
        log(LogLevel.INFO, 'CDN ' + str(request))
        log(LogLevel.DEBUG, 'headers =>', dict(request.headers))
        # return web.Response(text='hello world')

    async def main(self, host, port):
        app = web.Application()
        app.add_routes([web.get('/{tail:.*}', self.handler)])
        runner = web.AppRunner(app)
        await runner.setup()
        site = web.TCPSite(runner)
        await site.start()
        log(LogLevel.INFO, 'HTTP CDN server is listening for connections...')
        await asyncio.Event().wait()
