select
	a.id
    ,a.cdnid
    ,ft.name as file_type
    ,am.metadata ->> '$.info.version_engine' as engine
    ,am.metadata ->> '$.assets[0].target_platform' as platform
    ,replace(replace(replace(replace(am.metadata ->> '$.assets[0].name', 'BuildPlayer-', ''), '.sharedAssets', ''), 'CustomAssetBundle-', ''), 'CAB-', '') as asset
    , (
        SELECT group_concat(name, ', ')
        FROM (
            SELECT r.value ->> '$.name' AS name
            FROM json_each(am.metadata, '$.roots') AS r
            LIMIT 3
        )
      ) as top3_roots
from asset as a
join file_type as ft on ft.id = a.file_type_id
join asset_metadata as am on am.asset_id = a.id
where ft.name like 'AssetBundle/%'

