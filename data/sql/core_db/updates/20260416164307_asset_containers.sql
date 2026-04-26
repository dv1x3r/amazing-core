-- +goose Up
create table asset_container (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [name] text not null unique collate nocase
    [ptag] text collate nocase,
    [created_at] integer not null default (unixepoch())
) strict;

create table asset_container_assetmap (
    [id] integer primary key not null,
    [container_id] integer not null references asset_container(id) on delete cascade,
    [win_asset_id] integer not null references asset(id) on delete cascade,
    [osx_asset_id] integer references asset(id) on delete set null,
    [position] integer not null default 0,
    unique([container_id], [win_asset_id])
) strict;

create table asset_container_package (
    [id] integer primary key not null,
    [container_id] integer not null references asset_container(id) on delete cascade,
    [pkg_container_id] integer not null references asset_container(id) on delete cascade,
    [position] integer not null default 0,
    unique([container_id], [pkg_container_id])
) strict;

insert into asset_container ([id], [gsfoid], [name])
values (1, 0, 'Site Frame container');

create table site_frame (
    [id] integer primary key not null,
    [type_value] integer not null unique,
    [container_id] integer not null references asset_container(id) on delete restrict
) strict;

insert into site_frame ([id], [type_value], [container_id])
values (1, 1, 1);

