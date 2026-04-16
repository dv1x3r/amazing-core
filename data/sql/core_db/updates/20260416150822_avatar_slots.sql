-- +goose Up
create table avatar_slot (
    [id] integer primary key,
    [slot] text not null unique collate nocase,
    [animation] text
) strict;

insert into avatar_slot ([id], [slot], [animation])
values
(289356276061314050, 'Back', 'clothingback'),
(289356276061314056, 'Face', 'clothingHead'),
(289356276061314062, 'Hands', 'clothingHands'),
(289356276061314068, 'Hat', 'clothingHead'),
(289356276061314075, 'Jacket', 'clothingUpper'),
(289356276061314002, 'LowerBody', 'clothingLower'),
(289356276061314174, 'Shoes', 'clothingLower'),
(289356276061314180, 'UpperBody', 'clothingUpper'),
(289356276061582038, 'RightHand', 'idle001.oneHanded'),
(289356276067676408, 'Vehicle', null),
(289356276067681519, 'Pet', null);

