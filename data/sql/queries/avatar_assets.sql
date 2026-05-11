select
	pa.id as player_avatar_id,
    pa.player_id as player_id,
    pa.name as player_avatar_name,
    pa.is_active,
    a.name as avatar_name,
    ac.id as container_id,
    ac.name as container_name,
    asset_win.gsfoid || ' - ' || asset_win.res_name as windows_asset,
    asset_osx.gsfoid || ' - ' || asset_osx.res_name as osx_asset
from player_avatar as pa
join avatar as a on a.id = pa.avatar_id
join asset_container as ac on ac.id = a.container_id
left join asset_container_assetmap as aca on aca.container_id = a.container_id
left join asset as asset_win on asset_win.id = aca.win_asset_id
left join asset as asset_osx on asset_osx.id = aca.osx_asset_id

