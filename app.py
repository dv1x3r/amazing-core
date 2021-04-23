import os
import asyncio
from amazingcore.core import Core
from amazingcore.cdn import Cdn


async def main():
    tasks = [Core().main('127.0.0.1', 8182), Cdn().main('127.0.0.1', 8080)]
    await asyncio.gather(*tasks)

try:
    asyncio.run(main())
except KeyboardInterrupt:
    os._exit(0)  # No Traceback when we still have connected client
