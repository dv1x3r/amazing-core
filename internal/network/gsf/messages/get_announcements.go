package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAnnouncementsRequest requests login announcements for the player.
type GetAnnouncementsRequest struct {
	UnMarked bool
}

func (req *GetAnnouncementsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.UnMarked = reader.ReadBool()
}

// GetAnnouncementsResponse contains announcements displayed during login.
type GetAnnouncementsResponse struct {
	Announcements []types.Announcement
}

func (res *GetAnnouncementsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Announcements, func(value types.Announcement) {
		writer.WriteObject(&value)
	})
}
