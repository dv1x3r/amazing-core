<div style="display:flex; align-items:center; gap:12px;">
  <img src="/images/logo.png" height="60" />
  <h1 class="text-red">Amazing Core</h1>
</div>

Amazing Core is an **open-source server emulator** for Amazing World, an MMO originally developed by Ganz and shut down in 2018.
This project provides a modular, configurable framework with tools for server management, asset handling, and game services, accessible via a web-based dashboard.

> ⚠️ **Still in development** - not yet in a playable state.
>
> - No multiplayer, NPCs, or quests yet;
> - Only the intro level and the _empty_ Spring Bay map are accessible;
> - Do not use your real username or password;
> - Use any dummy username and password to log in;

But you can check out the work-in-progress prototype and join our community!

<div style="display:flex; flex-wrap:wrap; justify-content:space-between; gap:16px;">
  <div style="border-radius:16px; overflow:hidden; margin:0 auto;">
    <iframe src="https://discord.com/widget?id=822788246973972510&theme=dark" width="350" height="340" allowtransparency="true" frameborder="0" sandbox="allow-popups allow-popups-to-escape-sandbox allow-same-origin allow-scripts"></iframe>
  </div>
  <div style="border-radius:16px; overflow:hidden; background:#1e1f22; box-sizing:border-box; width:350px; display:flex; flex-direction:column; margin:0 auto;">
    <div style="background:#08872b; color:#ffffff; display:flex; justify-content:space-between; align-items:center; padding:15px 20px;">
      <div style="display:flex; align-items:center; gap:10px;">
        <i class="fa-brands fa-github" style="font-size:32px; color:#e6edf3;"></i>
        <span style="font-size:18px; font-weight:600;">GitHub</span>
      </div>
      <a href="https://github.com/dv1x3r/amazing-core/">
        <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/dv1x3r/amazing-core" style="display:block;"/>
      </a>
    </div>
    <div style="display:flex; flex-direction:column; gap:15px; padding:20px;">
      <div style="background:#151b23; color:#8b949e; border:1px solid #30363d; border-radius:16px; overflow:hidden; font-size:13px;">
        <div style="display:flex; justify-content:space-between; padding:8px 12px; font-weight:600; border-bottom:1px solid #30363d;">
          <span>Latest release binaries</span>
          <span style="background:#0a241b; border:1px solid #08872b; color:#0fbf3e; font-size:10px; font-weight:700; font-family:monospace; padding:2px 7px; border-radius:20px;">{{#include vars/version.md}}</span>
        </div>
        <a href="https://github.com/dv1x3r/amazing-core/releases/download/{{#include vars/version.md}}/amazing-core_Windows_x86_64.zip" style="display:flex; align-items:center; gap:8px; padding:8px 12px; color:#58a6ff; text-decoration:none; border-bottom:1px solid #30363d;">
          <i class="fa-brands fa-microsoft" style="color:#e6edf3;"></i>
          <span style="flex:1;">Windows x86_64</span>
          <span style="font-family:monospace; font-size:10px; color:#ffffff; background:#303133; padding:2px 6px; border-radius:4px;">.zip</span>
        </a>
        <a href="https://github.com/dv1x3r/amazing-core/releases/download/{{#include vars/version.md}}/amazing-core_Linux_x86_64.tar.gz" style="display:flex; align-items:center; gap:8px; padding:8px 12px; color:#58a6ff; text-decoration:none; border-bottom:1px solid #30363d;">
          <i class="fa-brands fa-linux" style="color:#e6edf3;"></i>
          <span style="flex:1;">Linux x86_64</span>
          <span style="font-family:monospace; font-size:10px; color:#ffffff; background:#303133; padding:2px 6px; border-radius:4px;">.tar.gz</span>
        </a>
        <a href="https://github.com/dv1x3r/amazing-core/releases/download/{{#include vars/version.md}}/amazing-core_Darwin_x86_64.tar.gz" style="display:flex; align-items:center; gap:8px; padding:8px 12px; color:#58a6ff; text-decoration:none; border-bottom:1px solid #30363d;">
          <i class="fa-brands fa-apple" style="color:#e6edf3;"></i>
          <span style="flex:1;">macOS x86_64</span>
          <span style="font-family:monospace; font-size:10px; color:#ffffff; background:#303133; padding:2px 6px; border-radius:4px;">.tar.gz</span>
        </a>
        <a href="https://github.com/dv1x3r/amazing-core/releases/download/{{#include vars/version.md}}/amazing-core_Darwin_arm64.tar.gz" style="display:flex; align-items:center; gap:8px; padding:8px 12px; color:#58a6ff; text-decoration:none">
          <i class="fa-brands fa-apple" style="color:#e6edf3;"></i>
          <span style="flex:1;">macOS arm64</span>
          <span style="font-family:monospace; font-size:10px; color:#ffffff; background:#303133; padding:2px 6px; border-radius:4px;">.tar.gz</span>
        </a>
      </div>
      <div style="display:flex; gap:8px;">
        <a href="https://github.com/dv1x3r/amazing-core" style="flex:1; text-align:center; background:#08872b; color:#ffffff; text-decoration:none; padding:10px; border-radius:8px; font-size:14px;">View on GitHub</a>
        <a href="https://github.com/dv1x3r/amazing-core/releases" style="flex:1; text-align:center; background:#303133; border:1px solid #30363d; color:#ffffff; text-decoration:none; padding:10px; border-radius:8px; font-size:14px;">All Releases</a>
      </div>
    </div>
  </div>
</div>

## Download the game

You can install the latest published version with the following Steam link:
<a
  href="steam://launch/293500"
  aria-label="Launch this game on Steam You must own it steam://launch/293500"
  style="margin-left:10px; border:1px solid #30363d; background:#303133; color:#e6edf3; text-decoration:none; padding:7px 10px; border-radius:4px;">
<i class="fa-brands fa-steam"></i>
<span>Install</span>
</a>

The game has its page on [SteamDB](https://steamdb.info/app/293500/), where you can also see additional information.

## Connect to the demo server

After installation, navigate to the game folder and open the `ServerConfig.xml` file in a text editor.

Modify the server address value as shown below:

```xml
ServerIP = 'springbay.amazingcore.org'
```

Now you can start the game.

- To play the intro level, click the `I'm new!` button in the main menu;
- To explore the Spring Bay, click the `Log in` button and enter any username and password;

## Host your own server

With your own server, you can access the configuration dashboard to configure skins, maps, NPCs, and other features (work in progress).

### Pre-compiled binaries

1. **Download** the latest [server release](#) from GitHub.
2. **Extract** the archive to a folder of your choice;
3. **Run** the server binary;

Once started:

- The API server will be available at [http://localhost:3000](http://localhost:3000)
  - Use `admin / admin` to log in to the configuration dashboard
- The Game server will listen on `localhost:8182`
- You can customize server settings using the `config.json`

### Build from source using Go

To build the server from source, you will need **Go 1.25** or newer:

```sh
make
# or
go build -o ./build/server ./cmd/server/main.go
```

To build and run with a single command:

```sh
make run
# or
go run ./cmd/server/main.go
```

You can choose between SQLite drivers by setting the `CGO_ENABLED` environment variable:

- Build with `CGO_ENABLED=0` to use `modernc.org/sqlite` driver (default);
- Build with `CGO_ENABLED=1` to use `github.com/mattn/go-sqlite3` driver;

## License

This project is licensed under the [GNU AGPL v3](https://www.gnu.org/licenses/agpl-3.0.html).

Amazing World™ is a registered trademark of Ganz. Amazing Core is an unofficial, fan-made project intended for personal and educational use only. It is not affiliated with or endorsed by Ganz or Amazing World™ in any way.
