# 🧰 Cache Importer Tool

**Part of the [amazing-core](https://github.com/dv1x3r/amazing-core) project.**

This CLI tool processes an asset cache folder and helps keep your server's asset database in sync.

### 🔧 Features

- **Imports assets** into `blob.db` from a local cache folder.

### 📦 Usage

#### Import cache file into `blob.db`

```bash
go run main.go --cache-dir /path/to/game/cache --db /path/to/blob.db
```

### ⚙️ Parameters

| Parameter      | Description                                             |
| -------------- | ------------------------------------------------------- |
| `--cache-dir`  | Path to the directory containing the game client cache. |
| `--db`         | Path to the target `blob.db` database.                  |
| `--everything` | Do not skip any files. Optional, default: false.        |
