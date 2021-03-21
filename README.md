# Amazing Core

Here we will try to create a custom server for the Amazing World game, which was closed in 2018.

- Game page on Steam: https://steamcommunity.com/app/293500
- Download game with Steam: steam://install/293500
- Discord Server: https://discord.gg/TWfTBbfdA9

**P.S.** I am not accepting / related to any donations as long as we do not have ready to go product and hosted server.  

## Reverse Engineering the Game

You can check some description in the /docs/intro.ipynb, or dig into the game files manually using ILSpy.

## How To Run

1. Install the game using the link above
2. Remove (or rename) : ```steam_api.dll```, ```steam_appid.txt``` and ```SteamworksNative.dll```, otherwise Steam will count hours in game
3. Edit ServerConfig.xml : ```ServerIP = '127.0.0.1'```
4. Run Amazing Core using command line and python 3.9+ : ```python app.py```
5. Start the Game
6. Cry a lot
