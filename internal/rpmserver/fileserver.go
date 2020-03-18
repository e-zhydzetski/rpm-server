package rpmserver

import (
	"net/http"
	"os"
)

func NewFileServer(baseFolder string) http.Handler {
	return http.FileServer(safeFileSystem{http.Dir(baseFolder)})
}

type safeFileSystem struct {
	fs http.FileSystem
}

func (sfs safeFileSystem) Open(path string) (http.File, error) {
	f, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}
