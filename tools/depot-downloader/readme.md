# 🧰 Depot Downloader Tool

**Part of the [amazing-core](https://github.com/dv1x3r/amazing-core) project.**

This CLI tool downloads all Amazing World depot files available on Steam.

### 🔧 Features

- Uses [**SteamRE/DepotDownloader**](https://github.com/SteamRE/DepotDownloader) CLI.

### 📦 Usage

#### Download Steam depots

```bash
go run main.go --username <steam username> --password <steam password>
```

### ⚙️ Parameters

| Parameter    | Description                                                     |
| ------------ | --------------------------------------------------------------- |
| `--username` | the username of the account to login to for restricted content. |
| `--password` | the password of the account to login to for restricted content. |
