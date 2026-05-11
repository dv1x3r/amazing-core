select
	it.id,
    it.name,
    it.container_id,
    (ac.gsfoid || ' - ' || ac.name) as container,
    cat.categories,
    slt.slots
from item as it
join asset_container as ac on ac.id = it.container_id
left join (
	select
		icm.item_id,
		json_group_array(
			json_object('id', ic.id, 'text', ic.name)
			order by ic.name
		) as categories
	from item_category_map as icm
	join item_category as ic on ic.id = icm.category_id
	group by icm.item_id
) as cat on cat.item_id = it.id
left join (
	select
		ias.item_id,
		json_group_array(
			json_object('id', avs.id, 'text', avs.slot)
			order by avs.slot
		) as slots
	from item_acceptable_slot as ias
	join avatar_slot as avs on avs.id = ias.avatar_slot_id
	group by ias.item_id
) as slt on slt.item_id = it.id

