import os
import asyncio
from amazingcore.core import Core
from amazingcdn.cdn import Cdn

small_slant = '''
           ___   __  ______ ____  _____  _______
          / _ | /  |/  / _ /_  / /  _/ |/ / ___/
         / __ |/ /|_/ / __ |/ /__/ //    / (_ / 
        /_/ |_/_/  /_/_/ |_/___/___/_/|_/\___/  
                        / ___/ __ \/ _ \/ __/   
                       / /__/ /_/ / , _/ _/     
                       \___/\____/_/|_/___/     
  AmazingCore 0.0.1 - github.com/dv1x3r/amazing-core
'''
print(small_slant)


async def main():
    tasks = [Core().main('127.0.0.1', 8182), Cdn().main('127.0.0.1', 8080)]
    await asyncio.gather(*tasks)

try:
    asyncio.run(main())
except KeyboardInterrupt:
    os._exit(0)  # No Traceback when we still have connected client
