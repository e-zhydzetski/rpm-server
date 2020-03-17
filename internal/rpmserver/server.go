package rpmserver

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)
import "github.com/go-chi/chi"

type Middleware func(handler http.Handler) http.Handler

func NewHandler(cfg Config) http.Handler {
	r := chi.NewRouter()
	r.Use(NewHTTPAuthInterceptor(cfg.AccessToken))
	r.Post("/packages", func(w http.ResponseWriter, r *http.Request) {
		pr, h, err := r.FormFile("package")
		if err != nil {
			log.Println("can't get package binary:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fName := filepath.Join(cfg.RepositoryPath, h.Filename)
		f, err := os.OpenFile(
			fName,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL, // will fail if file already exists
			os.ModePerm,
		)
		if err != nil {
			if os.IsExist(err) {
				log.Println("package already exists:", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Println("can't create file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(f, pr)
		_ = f.Close()
		if err != nil {
			_ = os.Remove(fName)
			log.Println("can't save file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cmd := exec.Command("createrepo", "-v", cfg.RepositoryPath)
		err = cmd.Run()
		if err != nil {
			_ = os.Remove(fName)
			log.Println("update repo index error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return r
}
