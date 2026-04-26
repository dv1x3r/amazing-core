-- +goose Up
create table player (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [is_tutorial_completed] integer not null default 0 check ([is_tutorial_completed] in (0, 1)),
    [is_qa] integer not null default 0 check ([is_qa] in (0, 1)),
    [created_at] integer not null default (unixepoch())
) strict;

insert into player ([gsfoid], [is_tutorial_completed], [is_qa])
values (1, 1, 1);

create table player_avatar (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [player_id] integer not null references player(id) on delete cascade,
    [avatar_id] integer not null references avatar(id) on delete cascade,
    [name] text not null unique collate nocase,
    [outfit_no] integer not null default 1,
    [is_active] integer not null default 0 check ([is_active] in (0, 1)),
    unique([player_id], [avatar_id])
) strict;

insert into player_avatar ([gsfoid], [player_id], [avatar_id], [name], [is_active])
select id, 1, id, name, case when id = 1 then 1 else 0 end
from avatar;

create table player_avatar_outfit (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [player_avatar_id] integer not null references player_avatar(id) on delete cascade,
    [outfit_no] integer not null default 1,
    unique([player_avatar_id], [outfit_no])
) strict;

insert into player_avatar_outfit ([gsfoid], [player_avatar_id])
values (1, 1);

