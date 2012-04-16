package cosmofs

type chunk struct {
	Name string
	RemPath string
	//peer peer.Peer
}

type File struct {
	Path string
	Filename string
	Size int64
	Chunks []chunk
	NumChunks bool
	MaintainChunks bool
	Online bool
	KeepCopy bool
	Public bool
	Backup bool
}

type Directory struct {
	File
	Recursive bool
}
