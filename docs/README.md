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

1. Use `cache.py` to unpack cache bundles:

```sh
# make meta.json per cache file, and unpack assets
python tools/cache.py /path/to/cache/folder \
  --json --unpack --mp3 --zip
```

2. Download `index.json` using the Amazing Core server.

3. Serve the output directory using npm or bun:

```sh
make docs-serve-data
# or
npmx serve -l 8080 --cors
# or
bunx serve -l 8080 --cors
```

4. Update `src/vars/cache-url.md` to point to the served address.
