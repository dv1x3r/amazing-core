# 🧰 Blob Tool

**Part of the [amazing-core](https://github.com/dv1x3r/amazing-core) project.**

This CLI tool allows you to sync asset files with `blob.db`.

It supports importing cache files into `blob.db` as well as exporting database blobs back to disk.

## 📦 Usage

### Import cache files into `blob.db`

```bash
go run main.go --mode import --db /path/to/blob.db --dir /path/to/game/cache
```

### Export cache files from `blob.db`

```bash
go run main.go --mode export --db /path/to/blob.db --dir /path/to/output/cache --overwrite
```

## ⚙️ Parameters

| Flag          | Description                                                                         |
| ------------- | ----------------------------------------------------------------------------------- |
| `--mode`      | Operation mode: `import` or `export`.                                               |
| `--db`        | Path to the `blob.db` database.                                                     |
| `--dir`       | `cache` source directory for import mode and destination directory for export mode. |
| `--overwrite` | Overwrite files on disk in export mode. Default: `false`.                           |
