#!/bin/bash
set -e

echo squashing core.db...

sqlite3 squash.db < base.sql
go tool goose -dir updates sqlite3 squash.db up
sqlite3 squash.db .dump > base.sql
rm squash.db

echo "✅ squash complete!"

