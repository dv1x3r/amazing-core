-- +goose Up
create table avatar (
    [id] integer primary key,
    [container_id] integer not null unique references asset_container(id) on delete restrict,
    [name] text not null unique collate nocase,
    [max_outfits] integer not null,
) strict;

