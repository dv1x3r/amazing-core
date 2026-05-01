-- +goose Up
create table item_category (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [name] text not null unique collate nocase,
    [parent_id] integer references item_category(id) on delete restrict,
    [is_public] integer not null default 0 check ([is_public] in (0, 1)),
    [is_outdoor] integer not null default 0 check ([is_outdoor] in (0, 1)),
    [is_walkover] integer not null default 0 check ([is_walkover] in (0, 1)),
    [show_in_dock] integer not null default 0 check ([show_in_dock] in (0, 1))
) strict;

insert into item_category ([id], [gsfoid], [name])
values
(1, 1, 'Clothing'),
(2, 2, 'Decoration'),
(3, 3, 'Door'),
(4, 4, 'EmoteItem'),
(5, 5, 'Flooring'),
(6, 6, 'Hallpaper'),
(7, 7, 'Housepaper'),
(8, 8, 'MazePiece'),
(9, 9, 'Placable'),
(10, 10, 'RoofDecoration'),
(11, 11, 'Rug'),
(12, 12, 'WallDecoration'),
(13, 13, 'Wallpaper'),
(14, 14, 'Window'),
(15, 15, 'Yard'),
(16, 16, 'Fishing'),
(17, 17, 'House'),
(18, 18, 'Upgrade'),
(19, 19, 'Dances'),
(20, 20, 'Feelings'),
(21, 21, 'Gestures'),
(22, 22, 'ShowOff'),
(23, 23, 'Aquarium'),
(24, 24, 'Fish'),
(25, 25, 'AcceptsMountable'),
(26, 26, 'Abilities'),
(27, 27, 'ConsumableVanity'),
(28, 28, 'ConsumableAbility'),
(29, 29, 'VotingSigns'),
(30, 30, 'Crafting'),
(31, 31, 'Enhanceables'),
(32, 32, 'Colors'),
(33, 33, 'Particles'),
(34, 34, 'Vehicles'),
(35, 35, 'Pets'),
(36, 36, 'ConsumablePack');

create table item (
    [id] integer primary key,
    [container_id] integer not null unique references asset_container(id) on delete restrict,
    [name] text not null unique collate nocase
) strict;

create table item_category_map (
    [id] integer primary key,
    [item_id] integer not null references item(id) on delete cascade,
    [category_id] integer not null references item_category(id) on delete restrict,
    unique([item_id], [category_id])
) strict;

create table item_acceptable_slot (
    [id] integer primary key,
    [item_id] integer not null references item(id) on delete cascade,
    [slot_id] integer not null references avatar_slot(id) on delete restrict,
    unique([item_id], [slot_id])
) strict;

create table player_item (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [player_id] integer not null references player(id) on delete cascade,
    [item_id] integer not null references item(id) on delete cascade,
    [slot_id] integer references avatar_slot(id) on delete restrict,
    [player_avatar_id] integer references player_avatar(id) on delete set null,
    [avatar_outfit_id] integer references player_avatar_outfit(id) on delete set null,
    [quantity] integer not null default 1,
    unique([player_id], [item_id])
) strict;

