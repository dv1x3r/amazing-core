package filetype

type FileType int

const (
	Unknown FileType = iota + 1
	MP3
	OGG
	PNG
	JSON
	TreeNode
	UnityFS
	UnityWeb
)
