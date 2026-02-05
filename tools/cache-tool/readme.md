# 🧰 Cache Tool

**Part of the [amazing-core](https://github.com/dv1x3r/amazing-core) project.**

This CLI tool inspects and unpacks game cache files and asset bundles.

## 🔧 Requirements

- [UnityPy](https://github.com/K0lb3/UnityPy): `pip install unitypy`
- [ffmpeg](https://ffmpeg.org/download.html): required for audio conversion to mp3

## 📦 Usage

### Write a single JSON summary file

```bash
python cache.py /path/to/cache/file-or-folder --summary-file output.json
```

### Write JSON summary files per cache file

```bash
python cache.py /path/to/cache/file-or-folder --summaries-dir /path/to/json/output
```

### Unpack audio, meshes, and textures from asset bundles

```bash
python cache.py /path/to/cache/file-or-folder --unpack-dir /path/to/assets/output
```

Flags can be combined, e.g. `--summaries-dir` and `--unpack-dir` together.

## ⚙️ Parameters

| Parameter         | Description                                                                                    |
| ----------------- | ---------------------------------------------------------------------------------------------- |
| `path`            | Path to a single cache file or a folder containing multiple files.                             |
| `--parse-scene`   | Parse and include scene hierarchy in the summary. Can produce large output for complex scenes. |
| `--summary-file`  | File to write a single JSON summary file.                                                      |
| `--summaries-dir` | Directory to write JSON summary files (one per input file).                                    |
| `--unpack-dir`    | Directory to unpack assets (audio, images, models) into.                                       |
