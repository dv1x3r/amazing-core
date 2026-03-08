-- +goose Up
-- +goose StatementBegin
create table file_type (
    [id] integer primary key,
    [name] text not null unique collate nocase
) strict;

insert into file_type ([id], [name]) values
(1, 'Unknown'),
(2, 'audio/mp3'),
(3, 'audio/ogg'),
(4, 'image/png'),
(5, 'JSON'),
(6, 'TreeNode'),
(7, 'UnityFS'),
(8, 'UnityWeb');

create table asset_type (
    [id] integer primary key,
    [name] text not null unique collate nocase
) strict;

insert into asset_type ([id], [name]) values
(1, 'Unknown'),
(2, 'Animation_Unity3D'),
(3, 'Audio'),
(4, 'Conditional_Prefab_Unity3d'),
(5, 'Config_Text'),
(6, 'Font_Unity3D'),
(7, 'Free_Play_Ad'),
(8, 'Game_Layout'),
(9, 'Help_Guide_Image'),
(10, 'Icon'),
(11, 'Images'),
(12, 'Intro_Unity3D'),
(13, 'Level_Up_Image'),
(14, 'Localized_Data'),
(15, 'LocalizedText'),
(16, 'Location_Unity3D'),
(17, 'Movie_Unity3D'),
(18, 'Nix_Animation_Unity3D'),
(19, 'Prefab_Unity3D'),
(20, 'Preload_PrefabUnity3D'),
(21, 'Property Text'),
(22, 'Relationship_Text'),
(23, 'Scene_Unity3D'),
(24, 'Text'),
(25, 'Textures'),
(26, 'Treenode'),
(27, 'zone_scene');

create table asset_group (
    [id] integer primary key,
    [name] text not null unique collate nocase
) strict;

insert into asset_group ([id], [name]) values
(1, 'Unknown'),
(2, 'Main_Scene'),
(3, 'Locked'),
(4, '2D Sprite'),
(5, '3D Components'),
(6, '3D Prefab'),
(7, 'Exterior Icon'),
(8, 'Interior Icon'),
(9, 'Inventory Icon'),
(10, 'SpawnPoints'),
(11, 'Wheel Icon');

create table asset (
    [id] integer primary key,
    [cdnid] text not null unique collate nocase,
    [oid] integer,
    [name] text,
    [size] integer not null,
    [file_type_id] integer not null references file_type(id) on delete cascade,
    [asset_type_id] integer not null references asset_type(id) on delete cascade,
    [asset_group_id] integer not null references asset_group(id) on delete cascade
) strict;
-- +goose StatementEnd
