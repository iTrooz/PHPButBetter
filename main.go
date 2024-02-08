package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"phpbutbetter/handlers"
	"strconv"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetSpecificHandler(ext string) func(w http.ResponseWriter, filepath string) error {
	switch ext {
	case ".c":
		return handlers.CHandler
	case ".cpp":
		return handlers.CppHandler
	case ".go":
		return handlers.GoHandler
	case ".java":
		return handlers.JavaHandler
	case ".b":
		return handlers.PipeWrapper("bf")
	case ".py":
		return handlers.PipeWrapper("python")
	case ".js":
		return handlers.PipeWrapper("node")
	case ".rs":
		return handlers.RustHandler
	case ".php":
		return handlers.PipeWrapper("php")
	case ".sh":
		return handlers.PipeWrapper("sh")
	case ".cmake":
		return handlers.CMakeHandler
	}
	return nil
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	err := MainHandlerInt(w, r)
	if err != nil {
		errorStr := fmt.Sprintf("Failed to run handler for path %v: %v", r.URL.Path, err)
		log.Error(errorStr)
		fmt.Fprint(w, errorStr)
	} else {
		log.Infof("Served %v successfully", r.URL.Path)
	}
}
func MainHandlerInt(w http.ResponseWriter, r *http.Request) error {
	ext := filepath.Ext(r.URL.Path)
	specificHandler := GetSpecificHandler(ext)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if specificHandler == nil {
		return fmt.Errorf("Could not find an handler for file extension '%s'", ext)
	}

	rootPath := getEnv("ROOTPATH", "./root")
	fullPath := path.Join(rootPath, r.URL.Path)

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File at '%s' does not exist", r.URL.Path)
		} else {
			return fmt.Errorf("Failed to check file at '%s'", r.URL.Path)
		}
	}

	err := specificHandler(w, fullPath)
	if err != nil {
		return fmt.Errorf("Failed to execute handler: %w", err)
	}

	return nil
}

func main() {
	http.HandleFunc("/", MainHandler)

	port := "8080"
	port_env := os.Getenv("PORT")
	if _, err := strconv.Atoi(port_env); err == nil {
		port = port_env
	}

	log.Infof("Serving on 0.0.0.0:%s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
