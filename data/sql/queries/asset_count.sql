select
	ft.name as file_type
    ,at.name as asset_type
    ,ag.name as asset_group
    ,count(*) as asset_count
from asset as a
join file_type as ft on ft.id = a.file_type_id
left join asset_type as at on at.id = a.asset_type_id
left join asset_group as ag on ag.id = a.asset_group_id
group by ft.name, at.name, ag.name
order by 1,2,3

