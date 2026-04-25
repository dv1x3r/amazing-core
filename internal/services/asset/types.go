package asset

import (
	"errors"

	"github.com/dv1x3r/w2go/w2"
)

var (
	ErrContainerExists         = errors.New("container with this GSF OID already exists")
	ErrContainerInUse          = errors.New("container is still in use and cannot be removed")
	ErrContainerAssetExists    = errors.New("container already contains this primary asset")
	ErrContainerPackageExists  = errors.New("container already contains this package")
	ErrPackageCyclicDependency = errors.New("circular dependency detected (A → B → A)")
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
	ID        int              `json:"id"`
	OID       w2.Field[int64]  `json:"oid"`
	OIDStr    string           `json:"oid_str"`
	Name      w2.Field[string] `json:"name"`
	PTag      w2.Field[string] `json:"ptag"`
	Assets    int              `json:"assets"`
	Packages  int              `json:"packages"`
	CreatedAt w2.UnixTime      `json:"created_at"`
}

type ContainerAsset struct {
	ID       int         `json:"id"`
	Position int         `json:"position"`
	WINAsset w2.Dropdown `json:"win_asset"`
	OSXAsset w2.Dropdown `json:"osx_asset"`
}

type ContainerPackage struct {
	ID           int         `json:"id"`
	Position     int         `json:"position"`
	PkgContainer w2.Dropdown `json:"pkg_container"`
	Assets       int         `json:"assets"`
	Packages     int         `json:"packages"`
}
