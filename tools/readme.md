# Tools

Developer utilities for working with Amazing World client files, cache files, Unity assets, and repository tooling.

## Files

| File             | Purpose                                                                                |
| ---------------- | -------------------------------------------------------------------------------------- |
| `cache.py`       | Inspect cache files, generate metadata JSON, and optionally unpack Unity assets.       |
| `silicon.sh`     | Patch a macOS Amazing World client bundle with the Unity 5.3.6p8 macOS player/runtime. |
| `steamdepots.sh` | Download known Amazing World Steam depot manifests using DepotDownloader.              |
| `unitypy.ipynb`  | Scratch notebook for UnityPy bundle inspection and metadata extraction experiments.    |
| `textures.ipynb` | Scratch notebook for Unity texture/material experiments.                               |

## cache.py

`cache.py` inspects one cache file or all files in a folder. It detects file type, calculates SHA-1 hashes, extracts Unity bundle metadata, parses scene hierarchy for Unity files, and can unpack meshes, textures, and audio.

Requirements:

- `UnityPy`
- `ffmpeg`, only when using `--ffmpeg-mp3`

Print metadata to stdout:

```sh
python tools/cache.py /path/to/cache/file-or-folder --stdout
```

Write one metadata file per cache file:

```sh
python tools/cache.py /path/to/cache/file-or-folder --metadata-dir /path/to/json/output
```

Unpack supported Unity assets:

```sh
python tools/cache.py /path/to/cache/file-or-folder --unpack-dir /path/to/assets/output
```

Flags can be combined, for example `--stdout --metadata-dir out --unpack-dir out`.

### ⚙️ Parameters

| Parameter        | Description                                                                                     |
| ---------------- | ----------------------------------------------------------------------------------------------- |
| `path`           | Path to a single cache file or a folder containing multiple files.                              |
| `--stdout`       | Write JSON metadata to stdout. Single file outputs one object; folders output an array.         |
| `--metadata-dir` | Directory to write one metadata JSON file per input file.                                       |
| `--unpack-dir`   | Directory to unpack assets (audio, images, models) into.                                       |
| `--ffmpeg-mp3`   | Convert unpacked audio to mp3 using ffmpeg. Use with `--unpack-dir`.                            |
| `--zip`          | Zip unpacked assets after extraction. Use with `--unpack-dir`.                                  |

## silicon.sh

`silicon.sh` patches a local macOS Amazing World client bundle.

Run:

```sh
tools/silicon.sh
```

The script downloads Unity 5.3.6p8, extracts the Unity player/runtime, copies the required runtime files into `AmazingWorld.app`, replaces the executable, and patches `ServerConfig.xml` to point at `springbay.amazingcore.org`.

## steamdepots.sh

`steamdepots.sh` downloads known Amazing World Steam depot manifests.

Requirements:

- `depotdownloader` available in `PATH`
- Steam credentials for an account that can access the depots

Run:

```sh
tools/steamdepots.sh
```

The script prompts for Steam username and password, then downloads each configured depot manifest into `293500/<depot>/<manifest>/`.
