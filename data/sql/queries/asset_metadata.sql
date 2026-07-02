select
	a.id
    ,a.cdnid
    ,a.gsfoid
    ,ft.name as file_type
    ,at.name as asset_type
    ,ag.name as asset_group
    ,a.res_name
    ,replace(replace(replace(replace(
      af.metadata ->> '$.bundle.assets[0].name',
      'BuildPlayer-', ''),
      '.sharedAssets', ''),
      'CustomAssetBundle-', ''),
      'CAB-', '') || '.unity3d' as clean_asset_name
    ,concat_ws(' ', af.metadata ->> '$.bundle.assets[0].target_platform', af.metadata ->> '$.bundle.info.version_engine') as bundle_version
    , (
        SELECT group_concat(name, ', ')
        FROM (
            SELECT r.value ->> '$.name' AS name
            FROM json_each(af.metadata, '$.bundle.scene') AS r
            LIMIT 3
        )
      ) as top3_roots
from asset as a
join file_type as ft on ft.id = a.file_type_id
left join asset_type as at on at.id = a.asset_type_id
left join asset_group as ag on ag.id = a.asset_group_id
left join blob.asset_file as af on af.cdnid = a.cdnid

