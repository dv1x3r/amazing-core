# Amazing Core mdBook

Docs are built with [mdBook](https://github.com/rust-lang/mdBook).

Markdown formatting is handled using [dprint](https://dprint.dev/plugins/markdown/).

To run the docs locally in watch mode:

```sh
make docs-serve
# or
mdbook serve docs/
```

To self-host cache archive:

1. Use `cache.py` from `cache-tool` to unpack cache bundles:

```sh
# make index cache.json
python cache.py /path/to/cache/folder --summary-file /path/to/output/cache.json

# make summary.json per cache file, and unpack assets
python cache.py /path/to/cache/folder \
  --parse-scene --ffmpeg-mp3 --zip \
  --summaries-dir /path/to/output/unpacked \
  --unpack-dir /path/to/output/unpacked
```

2. Serve the output directory using npm or bun:

```sh
make docs-serve-data
# or
npmx serve -l 8080 --cors
# or
bunx serve -l 8080 --cors
```

3. Update `src/vars/cache-url.md` to point to the served address.
