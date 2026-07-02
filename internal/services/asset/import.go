package asset

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

func (s *Service) ImportAssets(ctx context.Context) error {
	const op = "asset.Service.ImportAssets"
	const query = `
		insert or ignore into asset (cdnid, gsfoid, hash, size, file_type_id, res_name, bundle_version)
		select
			af.cdnid,
			af.metadata ->> '$.file.oid' as gsfoid,
			af.metadata ->> '$.file.hash' as hash,
			af.metadata ->> '$.file.size' as size,
			ft.id as file_type_id,
			replace(replace(replace(replace(
				af.metadata ->> '$.bundle.assets[0].name',
				'BuildPlayer-', ''),
				'.sharedAssets', ''),
				'CustomAssetBundle-', ''),
				'CAB-', '') || '.unity3d' as res_name,
			nullif(concat_ws(' ',
				af.metadata ->> '$.bundle.assets[0].target_platform',
				af.metadata ->> '$.bundle.info.version_engine')
				,'') as bundle_version
			from blob.asset_file as af
			left join file_type as ft on ft.name = af.metadata ->> '$.file.type';
		`
	_, err := s.store.DB().ExecContext(ctx, query)
	return wrap.IfErr(op, err)
}
