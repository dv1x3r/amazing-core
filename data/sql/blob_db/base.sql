CREATE TABLE asset_file (
    [id] integer primary key,
    [cdnid] text not null unique collate nocase,
    [blob] blob not null,
    [metadata] blob,
    [hash] text not null unique collate nocase
) strict;
