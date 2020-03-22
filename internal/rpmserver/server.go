package rpmserver

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

		pkgFileName := h.Filename
		pkgFileName = filepath.Base(pkgFileName)
		if !strings.HasSuffix(pkgFileName, ".rpm") {
			log.Println("invalid package file extension, only *.rpm supported:", pkgFileName)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid package file extension, only *.rpm supported"))
			return
		}

		fName := filepath.Join(cfg.PushRepoPath, pkgFileName)
		f, err := os.OpenFile(
			fName,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL, // will fail if file already exists
			os.ModePerm,
		)
		if err != nil {
			if os.IsExist(err) {
				log.Println("package already exists:", err)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("package already exists"))
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

		//nolint:gosec // G204(shell injection) safe as command arg is a configuration parameter
		cmd := exec.Command("createrepo", "-v", cfg.PushRepoPath)

		err = cmd.Run()
		if err != nil {
			_ = os.Remove(fName)
			log.Println("update repo index error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("package", pkgFileName, "successfully added to repo", cfg.PushRepoPath)
	})
	return r
}
