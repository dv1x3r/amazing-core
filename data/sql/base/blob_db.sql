create table asset_file (
    [id] integer primary key,
    [cdnid] text not null unique collate nocase,
    [blob] blob not null,
    [hash] text not null
) strict;

