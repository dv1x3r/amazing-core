-- +goose Up
create table zone (
    [id] integer primary key,
    [container_id] integer not null unique references asset_container(id) on delete restrict
) strict;

