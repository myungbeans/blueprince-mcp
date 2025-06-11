package storage

type Store interface {
	GetFiles(filename string) ([]string, error)
	ListFiles() ([]string, error)
	MoveFile(filename, destination string) error
}
