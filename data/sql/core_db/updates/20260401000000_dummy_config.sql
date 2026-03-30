-- +goose Up
create table dummy_config (
    [key] text,
    [value] text
) strict;

insert into dummy_config ([key], [value])
values
('map', 'OTYwOTUyODk5OTk1MA'),
('avatar', 'OTQ0NDE2MzMyMTg3MA');
