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

`cache.py` inspects one cache file or all files in a folder. It detects file type, calculates SHA-1 hashes, adds duration metadata for mp3/ogg files, adds dimensions metadata for png files, extracts Unity bundle metadata, parses scene hierarchy for Unity files, and can unpack meshes, textures, and audio.

Requirements:

- `UnityPy`
- `mutagen`
- `Pillow`
- `ffmpeg`, only when using `--mp3`

Print metadata to stdout:

```sh
python tools/cache.py /path/to/cache/file-or-folder --stdout
```

Write one metadata file per cache file:

```sh
python tools/cache.py /path/to/cache/file-or-folder --json
```

JSON metadata files are written beside each input file as `<filename>.meta.json`.

Unpack supported Unity assets:

```sh
python tools/cache.py /path/to/cache/file-or-folder --unpack
```

Unpacked assets are written beside each input file as `<filename>_assets`.

Zip supported Unity assets:

```sh
python tools/cache.py /path/to/cache/file-or-folder --unpack --zip
```

Zipped assets are written beside each input file as `<filename>.zip`, and the temporary `<filename>_assets` folder is removed after the zip is created.

Flags can be combined, for example `--stdout --json --unpack`.

### ⚙️ Parameters

| Parameter    | Description                                                                                  |
| ------------ | -------------------------------------------------------------------------------------------- |
| `path`       | Path to a single cache file or a folder containing multiple files.                           |
| `--manifest` | JSON manifest with `file_path` and `json_path` entries.                                      |
| `--stdout`   | Write JSON metadata to stdout. Single file outputs one object; folders output an array.      |
| `--json`     | Write metadata beside each input file as `<filename>.meta.json`.                             |
| `--unpack`   | Unpack assets beside each input file into `<filename>_assets` folders.                       |
| `--mp3`      | Convert unpacked audio to mp3 using ffmpeg. Use with `--unpack`.                             |
| `--zip`      | Zip unpacked assets to `<filename>.zip` and remove `<filename>_assets`. Use with `--unpack`. |

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
