//go:generate go tool stringer -type=OIDType

package oidtype

type OIDType int

const (
	UNTYPED      OIDType = 1
	PLAYER       OIDType = 2
	AVATAR_SLOT  OIDType = 3
	ITEM         OIDType = 4
	AVATAR       OIDType = 5
	NPC          OIDType = 6
	ZONE         OIDType = 7
	PLAYER_QUEST OIDType = 8
)
