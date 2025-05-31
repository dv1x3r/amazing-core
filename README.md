# <img src="web/images/logo.png" height="75"> Amazing Core

Amazing Core is an open-source server emulator for **Amazing World**, an MMO originally developed by Ganz and shut down in 2018.
This project provides a modular, configurable framework with tools for server management, asset handling, and game services, accessible via a web-based dashboard.

Feel free to join our [Discord server!](https://discord.com/invite/TWfTBbfdA9)

## Features

Amazing Core is still in development and **not yet in a playable state** - _many_ message handlers currently return placeholder responses.
However, the project already includes:

- Modular architecture for game services;
- Implementation of the networking protocol;
- SQLite storage for simplicity and portability;
- Lightweight Web UI for server and database management;

## Install the game

1. Visit the game SteamDB page: https://steamdb.info/app/293500/;
2. Click the **Install** button in the top right corner;
3. Install the game using Steam;
4. Navigate to the game folder and open the `ServerConfig.xml` file in a text editor;
5. Modify the server address as shown below:

If you **do not want to run a server yourself**, you can use the public development server:

```xml
ServerIP = 'springbay.amazingcore.org'
```

**To use your your own local server**:

```xml
ServerIP = 'localhost'
```

- To play the intro level, click the `I'm new!` button in the main menu;
- To explore Spring Bay, click the `Log in` button and enter any username and password;

## Getting Started

1. **Download** the latest [server release](https://github.com/dv1x3r/amazing-core/releases):
   - For regular Windows use `amazing-core_Windows_x86_64.zip`
   - For Apple silicon Mac use `amazing-core_Darwin_arm64.tar.gz`
2. **Extract** the archive to a folder of your choice;
3. **Download** the [blob.db](https://drive.google.com/drive/folders/1K7k7ZHrL5KZTdsa5_BblgafPgeGWwKRc?usp=share_link) database file and place it inside the `data_db` folder;
4. **Run** the server binary;

Once started:

- The API server will be available at http://localhost:3000
- The Game server will listen on `localhost:8182`
- You can customize server settings using the `config.json`.

## Build

Make sure you have the following installed:

- **Go >= 1.24**;

To build the project:

```sh
go build -o ./build/server ./cmd/server/main.go
```

To build and run it with a single command:

```sh
go run ./cmd/server/main.go
```

## Structure

```
cmd/           - entry point
data/          - sql migrations
internal/      
├── api/       - http server for admin dashboard and asset streaming
├── game/      - tcp game server and message handling
├── config/    - configuration variables
├── lib/       - shared libraries (e.g. logging, helpers)
├── services/  - business logic and database interaction
tools/         - development tools (e.g. asset importers)
web/           - embedded frontend for admin dashboard
```

## License

This project is licensed under the [GNU AGPL v3](LICENSE).

Amazing World™ is a registered trademark of Ganz. Amazing Core is an unofficial, fan-made project intended for personal and educational use only. It is not affiliated with or endorsed by Ganz or Amazing World™ in any way.
