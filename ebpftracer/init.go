package ebpftracer

type file struct {
	pid uint32
	fd  uint64
}

type sock struct {
	pid uint32
	fd  uint64
	// proc.Sock
}

func readFds(pids []uint32) (files []file, socks []sock) {
	
	return
}
