package storage

type Store interface {
	GetFile(filename, dirname string) (any, error) //TODO: What's the right format to represent the img file?
	ListFiles(dirname string) ([]string, error)
	MoveFile(filename, destination string) error
}
