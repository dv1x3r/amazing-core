-- +goose Up
-- +goose StatementBegin
create table file_type (
    [id] integer primary key,
    [name] text not null unique collate nocase
) strict;

insert into file_type ([id], [name]) values
(1, 'audio/mp3'),
(2, 'audio/ogg'),
(3, 'image/png'),
(4, 'AssetBundle/UnityFS'),
(5, 'AssetBundle/UnityWeb'),
(6, 'TreeNode/Announcement'),
(7, 'TreeNode/Areas'),
(8, 'TreeNode/AvatarProperty'),
(9, 'TreeNode/BuildingCompletion'),
(10, 'TreeNode/BuildingUI'),
(11, 'TreeNode/DressAvatarSlots'),
(12, 'TreeNode/Fish'),
(13, 'TreeNode/Game'),
(14, 'TreeNode/Item'),
(15, 'TreeNode/LevelUp'),
(16, 'TreeNode/Mission'),
(17, 'TreeNode/NPCs'),
(18, 'TreeNode/NPCAnimations'),
(19, 'TreeNode/NPCRelationships'),
(20, 'TreeNode/Property'),
(21, 'TreeNode/cQuest'),
(22, 'TreeNode/Quest'),
(23, 'TreeNode/Root'),
(24, 'TreeNode/SpawnPoints'),
(25, 'TreeNode/UIWidget'),
(26, 'json');

create table asset_type (
    [id] integer primary key,
    [name] text not null unique collate nocase
) strict;

insert into asset_type ([id], [name]) values
(1, 'Animation_Unity3D'),
(2, 'Audio'),
(3, 'Conditional_Prefab_Unity3d'),
(4, 'Config_Text'),
(5, 'Font_Unity3D'),
(6, 'Free_Play_Ad'),
(7, 'Game_Layout'),
(8, 'Help_Guide_Image'),
(9, 'Icon'),
(10, 'Images'),
(11, 'Intro_Unity3D'),
(12, 'Level_Up_Image'),
(13, 'Localized_Data'),
(14, 'LocalizedText'),
(15, 'Location_Unity3D'),
(16, 'Movie_Unity3D'),
(17, 'Nix_Animation_Unity3D'),
(18, 'Prefab_Unity3D'),
(19, 'Preload_PrefabUnity3D'),
(20, 'Property Text'),
(21, 'Relationship_Text'),
(22, 'Scene_Unity3D'),
(23, 'Text'),
(24, 'Textures'),
(25, 'Treenode'),
(26, 'zone_scene');

create table asset_group (
    [id] integer primary key,
    [name] text not null unique collate nocase
) strict;

insert into asset_group ([id], [name]) values
(1, 'Main_Scene'),
(2, 'Locked'),
(3, '2D Sprite'),
(4, '3D Components'),
(5, '3D Prefab'),
(6, 'Exterior Icon'),
(7, 'Interior Icon'),
(8, 'Inventory Icon'),
(9, 'SpawnPoints'),
(10, 'Wheel Icon');

create table asset (
    [id] integer primary key,
    [cdnid] text not null unique collate nocase,
    [gsfoid] integer not null unique,
    [file_type_id] integer not null references file_type(id) on delete restrict,
    [asset_type_id] integer references asset_type(id) on delete restrict,
    [asset_group_id] integer references asset_group(id) on delete restrict,
    [res_name] text collate nocase,
    [description] text,
    [hash] text not null,
    [size] integer not null
) strict;

create table asset_metadata (
    [id] integer primary key,
    [asset_id] integer unique references asset(id) on delete cascade,
    [metadata] blob
) strict;
-- +goose StatementEnd
