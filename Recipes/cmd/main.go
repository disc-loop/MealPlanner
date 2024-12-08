package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

type Specification struct {
	Paths map[string]map[string]interface{} `yaml:"paths"`
}

const (
	port     = "3333"
	specPath = "edamam-spec.yaml"
)

func main() {
	spec := readSpec(specPath)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	for path, methods := range spec.Paths {
		r.Route(path, func(r chi.Router) {
			for method := range methods {
				if method == "get" {
					msg := fmt.Sprintf("Called %s %s", strings.ToUpper(method), path)
					r.Get("/", func(w http.ResponseWriter, req *http.Request) {
						w.Write([]byte(msg))
					})
				}
			}
		})
	}

	slog.Info("Listening on http://localhost:" + port)
	http.ListenAndServe(":"+port, r)
}

func readSpec(path string) Specification {
	bs, err := os.ReadFile(path)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not read file %s. Error: %s", path, err))
	}

	s := Specification{}
	err = yaml.Unmarshal(bs, &s)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not read spec. Error: %s", err))
	}

	return s
}
