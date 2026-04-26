package siteframe

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetGSFSiteFrame(ctx context.Context, platform gsf.Platform, typeValue int32) (types.SiteFrame, error) {
	const op = "siteframe.Service.GetGSFSiteFrame"

	var sf types.SiteFrame
	sf.TypeValue = typeValue

	var containerID int
	row := s.store.DB().QueryRowContext(ctx, "select container_id from site_frame where type_value = ?;", typeValue)
	if err := row.Scan(&containerID); err != nil {
		return sf, wrap.IfErr(op, err)
	}

	container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
	if err != nil {
		return sf, wrap.IfErr(op, err)
	}
	sf.AssetContainer = container

	// TODO: should not be empty, but I do not know what kind of assets should be here
	if sf.AssetMap["Preload_PrefabUnity3D"] == nil {
		sf.AssetMap["Preload_PrefabUnity3D"] = []types.Asset{}
	}

	// TODO: should not be empty, but I do not know what kind of assets should be here
	if sf.AssetMap["Audio"] == nil {
		sf.AssetMap["Audio"] = []types.Asset{}
	}

	return sf, nil
}
