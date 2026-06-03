package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerBuildObject is a player-owned maze piece.
type PlayerBuildObject struct {
	OID             OID
	MazePiece       MazePiece
	PlayerMazeOID   OID
	IsTaskCompleted bool
	IsLocked        bool
	FeatureCode     string
	PlayerOID       OID
}

func (pbo *PlayerBuildObject) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pbo.OID)
	writer.WriteObject(&pbo.MazePiece)
	writer.WriteObject(&pbo.PlayerMazeOID)
	writer.WriteBool(pbo.IsTaskCompleted)
	writer.WriteBool(pbo.IsLocked)
	writer.WriteString(pbo.FeatureCode)
	writer.WriteObject(&pbo.PlayerOID)
}

func (pbo *PlayerBuildObject) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pbo.OID)
	reader.ReadObject(&pbo.MazePiece)
	reader.ReadObject(&pbo.PlayerMazeOID)
	pbo.IsTaskCompleted = reader.ReadBool()
	pbo.IsLocked = reader.ReadBool()
	pbo.FeatureCode = reader.ReadString()
	reader.ReadObject(&pbo.PlayerOID)
}
