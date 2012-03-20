package file

type chunk struct {
	name string
	remPath string
	//peer peer.Peer
}

type File struct {
	Path string
	Filename string
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
