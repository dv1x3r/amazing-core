import os
import asyncio
from amazingcore.core import Core

try:
    asyncio.run(Core().main('127.0.0.1', 8182))
except KeyboardInterrupt:
    os._exit(0)  # No Traceback when we still have connected client
