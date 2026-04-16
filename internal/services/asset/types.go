package asset

import (
	"errors"

	"github.com/dv1x3r/w2go/w2"
)

var (
	ErrContainerExists      = errors.New("asset container with the same GSF OID already exists")
	ErrContainerInUse       = errors.New("asset container is still referenced and cannot be removed")
	ErrContainerAssetExists = errors.New("the same primary asset already exists in the container")
)

type Asset struct {
	ID          int              `json:"id"`
	OID         int64            `json:"oid"`
	OIDStr      string           `json:"oid_str"`
	CDNID       string           `json:"cdnid"`
	URL         string           `json:"url"`
	FileType    w2.Dropdown      `json:"file_type"`
	AssetType   w2.Dropdown      `json:"asset_type"`
	AssetGroup  w2.Dropdown      `json:"asset_group"`
	ResName     w2.Field[string] `json:"res_name"`
	Description w2.Field[string] `json:"description"`
	Hash        string           `json:"hash"`
	Size        int              `json:"size"`
	SizeStr     string           `json:"size_str"`
	Metadata    w2.Field[string] `json:"metadata"`
	Version     w2.Field[string] `json:"version"`
}

type Container struct {
	ID     int              `json:"id"`
	OID    w2.Field[int64]  `json:"oid"`
	OIDStr string           `json:"oid_str"`
	Name   w2.Field[string] `json:"name"`
}

type ContainerAsset struct {
	ID       int         `json:"id"`
	Position int         `json:"position"`
	WINAsset w2.Dropdown `json:"win_asset"`
	OSXAsset w2.Dropdown `json:"osx_asset"`
}

type ContainerPackage struct {
	ID        int    `json:"id"`
	Position  int    `json:"position"`
	Container string `json:"container"`
	Name      string `json:"name"`
	PTag      string `json:"ptag"`
}
