-- +goose Up
create table avatar (
    [id] integer primary key,
    [name] text not null unique collate nocase,
    [max_outfits] integer not null,
    [container_id] integer not null references asset_container(id) on delete restrict
) strict;

