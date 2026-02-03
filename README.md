# <img src="web/images/logo.png" height="75"> Amazing Core

Amazing Core is an open-source server emulator for **Amazing World**, an MMO originally developed by Ganz and shut down in 2018.
This project provides a modular, configurable framework with tools for server management, asset handling, and game services, accessible via a web-based dashboard.

👉 Feel free to join our [Discord server!](https://discord.com/invite/TWfTBbfdA9)

## ⚠️ Warning!

Amazing Core is still in development and **not yet in a playable state** - _many_ message handlers currently return placeholder responses.

- **No multiplayer, NPCs, or quests yet**;
- **Only the intro level and the _empty_ Spring Bay map are accessible**;
- **Do not use your real username or password**;
- **Use any dummy username and password to log in**;

## 🎯 Features

- Open-source Amazing World server;
- Web-based control panel for server management;
- Lightweight and portable, no complex setup needed;

## 🎮 Install the game

1. Visit the game SteamDB page: https://steamdb.info/app/293500/;
2. Click the **Install** button in the top right corner to open Steam;
3. A Steam popup will appear, allowing you to install the game;
4. After installation, navigate to the game folder and open the `ServerConfig.xml` file in a text editor;
5. Modify the server address value as shown below:

If you **do not want to run a server yourself**, use the public demo server:

```xml
ServerIP = 'springbay.amazingcore.org'
```

To use your **own local server**, use the value below and proceed to the **Getting Started** section:

```xml
ServerIP = 'localhost'
```

## 🕹 In the game

- To play the intro level, click the `I'm new!` button in the main menu;
- To explore the Spring Bay, click the `Log in` button and enter any username and password;

## 🧪 Getting Started

Use this section if you want to download and run the prebuilt server - no setup or compilation required.

1. **Download** the latest [server release](https://github.com/dv1x3r/amazing-core/releases):
   - For regular Windows just use `amazing-core_Windows_x86_64.zip`
2. **Extract** the archive to a folder of your choice;
3. **Download** the game assets database [blob.db](https://drive.google.com/drive/folders/1K7k7ZHrL5KZTdsa5_BblgafPgeGWwKRc?usp=share_link) and place it inside the `data_db` folder;
4. **Run** the server binary;

Once started:

- The API server will be available at http://localhost:3000
- The Game server will listen on `localhost:8182`
- You can customize server settings using the `config.json`.

## 🧱 Build

To build the server from source, you will need **Go 1.25** or newer:

```sh
go build -o ./build/server ./cmd/server/main.go
```

To build and run with a single command:

```sh
go run ./cmd/server/main.go
```

You can choose between SQLite drivers by setting the `CGO_ENABLED` environment variable:

- Build with `CGO_ENABLED=0` to use `modernc.org/sqlite` driver (release version);
- Build with `CGO_ENABLED=1` to use `github.com/mattn/go-sqlite3` driver (default);

## 📁 Structure

```
cmd/           - entry point
data/          - sql migrations
internal/      
├── api/       - http server for admin dashboard and asset streaming
├── game/      - game server and message handling
├── network/   - tcp server protocol implementation
├── services/  - business logic and database interaction
├── config/    - configuration variables
├── lib/       - shared libraries (e.g. logging, helpers)
tools/         - development tools (e.g. asset importers)
web/           - embedded frontend for admin dashboard
```

## 📄 License

This project is licensed under the [GNU AGPL v3](LICENSE).

Amazing World™ is a registered trademark of Ganz. Amazing Core is an unofficial, fan-made project intended for personal and educational use only. It is not affiliated with or endorsed by Ganz or Amazing World™ in any way.
