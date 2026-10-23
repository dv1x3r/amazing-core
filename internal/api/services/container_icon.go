package services

import "net/url"

// containerIconJoin selects the first png image of every container as icon.cdnid (icon.rn = 1).
const containerIconJoin = `(
		select
			am.container_id,
			a.cdnid,
			row_number() over (partition by am.container_id order by am.position) as rn
		from asset_container_assetmap as am
		join asset as a on a.id = am.win_asset_id
		join file_type as ft on ft.id = a.file_type_id
		where ft.name = 'image/png'
	) as icon`

// iconURL returns the CDN URL of the container icon, or an empty string when there is none.
func (s *Service) iconURL(cdnid *string) string {
	if cdnid == nil {
		return ""
	}
	u, _ := url.JoinPath(s.deliveryURL, *cdnid)
	return u
}
