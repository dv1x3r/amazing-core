# Data and databases

The `data/` directory contains `.sql` database migrations and a `cache.json` file with currently known asset files metadata.

## SQLite Databases

For simplicity and portability, the project uses SQLite.
The `data/sql/` directory is embedded into the server binary at compile time, and migrations are run automatically at startup.

### core.db

**Schema:** `data/sql/core_db/`

The main SQLite database. Contains the classified list of available assets and random names used in the sign-up scene.

Managed by `goose` using migrations in `data/sql/core_db/updates/`. Shortcuts for creating, applying, and resetting migrations are available in the `Makefile`.

An initial squashed migration (`data/sql/core_db/base.sql`) includes both schema and seed data.

Use `data/sql/core_db/squash.sh` to squash update migrations into `base.sql`.

Files under `data/sql/queries/` are named example SQL queries available to the admin dashboard`s SQL Explorer.

### blob.db

**Schema:** `data/sql/blob_db/`

A separate SQLite database storing the **raw binary content** of game asset files.

Downloaded on startup if `blob.download = true` and the database is missing.

When a client requests `/cdn/{cdnid}`, the server fetches the `blob` column for that `cdnid` (enabled when `settings.assetDeliveryAPI = true`).

The database can be populated using the [**blob-tool**](https://github.com/dv1x3r/amazing-core/tree/master/tools/blob-tool) or the **admin dashboard**.

## Cache Metadata

### cache.json

A summary of known cache files, including name (CDN ID), size, type, hash, GSF OID and basic bundle info (version, object counts, scene roots).

Generated with the [**cache-tool**](https://github.com/dv1x3r/amazing-core/tree/master/tools/cache-tool) Python script:

```sh
python tools/cache-tool/cache.py /path/to/cache/folder \
  --summary-file data/cache.json
```

> This file is already included in the `core.db` base migration.
> For future updates (if new unique cache assets are discovered), add a new migration under `updates/` with the additional data.
