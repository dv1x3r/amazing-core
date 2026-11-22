package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Store is a store of an NPC or of a zone.
type Store struct {
	AssetContainer
	CreateDate gsf.UnixTime
	ZoneOID    OID
	NPCOID     OID
	SpawnPoint string
}

func (s *Store) Serialize(writer gsf.ProtocolWriter) {
	s.AssetContainer.Serialize(writer)
	writer.WriteUtcDate(s.CreateDate)
	writer.WriteObject(&s.ZoneOID)
	writer.WriteObject(&s.NPCOID)
	writer.WriteString(s.SpawnPoint)
}

func (s *Store) Deserialize(reader gsf.ProtocolReader) {
	s.AssetContainer.Deserialize(reader)
	s.CreateDate = reader.ReadUtcDate()
	reader.ReadObject(&s.ZoneOID)
	reader.ReadObject(&s.NPCOID)
	s.SpawnPoint = reader.ReadString()
}
