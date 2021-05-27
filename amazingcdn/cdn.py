import asyncio
from aiohttp import web
from amazingcore.logger import LogLevel, log


class Cdn:

    async def handler(self, request: web.Request):
        log(LogLevel.INFO, 'HTTP CDN ' + str(request))
        # log(LogLevel.DEBUG, 'headers =>', dict(request.headers))

        try:
            val = open('amazingcdn\\' + str(request.rel_url), 'rb').read()
            return web.Response(body=val)
        except:
            log(LogLevel.WARN, 'HTTP CDN Not Found:   ' + str(request.rel_url))
            return web.HTTPNotFound()

    async def main(self, host, port):
        app = web.Application()
        app.add_routes([web.get('/{tail:.*}', self.handler)])
        runner = web.AppRunner(app)
        await runner.setup()
        site = web.TCPSite(runner)
        await site.start()
        log(LogLevel.INFO, 'HTTP CDN server is listening for connections...')
        await asyncio.Event().wait()
