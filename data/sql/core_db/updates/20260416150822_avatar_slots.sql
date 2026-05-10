-- +goose Up
create table avatar_slot (
    [id] integer primary key,
    [gsfoid] integer not null unique,
    [slot] text not null unique collate nocase,
    [animation] text
) strict;

insert into avatar_slot ([id], [gsfoid], [slot], [animation])
values
(1, 289356276061314050, 'Back', 'clothingback'),
(2, 289356276061314056, 'Face', 'clothingHead'),
(3, 289356276061314062, 'Hands', 'clothingHands'),
(4, 289356276061314068, 'Hat', 'clothingHead'),
(5, 289356276061314075, 'Jacket', 'clothingUpper'),
(6, 289356276061314002, 'LowerBody', 'clothingLower'),
(7, 289356276061314174, 'Shoes', 'clothingLower'),
(8, 289356276061314180, 'UpperBody', 'clothingUpper'),
(9, 289356276061582038, 'RightHand', 'idle001.oneHanded'),
(10, 289356276067676408, 'Vehicle', null),
(11, 289356276067681519, 'Pet', null);

