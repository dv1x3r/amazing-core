-- +goose Up
-- +goose StatementBegin

/*

Database workflow until the first stable release
https://github.com/dv1x3r/amazing-core/issues/4

base.sql holds the full schema and data, a new core.db is created from it.

During development add migrations to this folder as usual, pull requests should include them for an easier review.
From time to time the migrations are squashed into base.sql (squash.sh) and removed from this folder,
data changes can touch thousands of rows and keeping every migration is not worth it at this stage.

base.sql is applied only to a new database, so the squashed changes never reach an existing core.db.
After updating the server to a new version or pulling a squashed base.sql, delete your core.db
and the server will recreate it. Local changes such as players and dashboard edits are lost.

*/

-- +goose StatementEnd

