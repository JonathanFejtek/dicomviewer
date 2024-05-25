package http

import (
	"dicomviewer/dicom"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const DefaultPort = "3000"

type Server struct {
	dicomFiles *dicomFiles
	router     chi.Router

	port string
}

// ServerOptions options when instantiating a server
type ServerOptions struct {
	port string
}

// UsePort option to specify a port for the server to listen on
func UsePort(port string) func(opts *ServerOptions) {
	return func(opts *ServerOptions) {
		opts.port = port
	}
}

// NewServer constructs a new application server
func NewServer(options ...func(opts *ServerOptions)) *Server {

	var opts = ServerOptions{
		port: DefaultPort,
	}

	for _, optFn := range options {
		optFn(&opts)
	}

	service := &Server{
		dicomFiles: &dicomFiles{
			fileRepository: dicom.NewLocalFileAdapter(),
		},
		router: chi.NewRouter(),
		port:   opts.port,
	}

	service.registerRoutes()

	return service
}

func (s *Server) registerRoutes() {
	s.router.Use(middleware.Logger)

	s.router.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(apiV1 chi.Router) {
			apiV1.Route("/files", func(files chi.Router) {

				// POST /api/v1/files
				files.Post("/", s.dicomFiles.Create)

				// Get /api/v1/files
				files.Get("/", s.dicomFiles.GetAll)

				files.Route("/{id}", func(filesByID chi.Router) {

					// GET /api/v1/files/{id}
					filesByID.Get("/", s.dicomFiles.GetByID)

					// GET /api/v1/files/{id}/png
					filesByID.Get("/png", s.dicomFiles.GetAsPNG)

					// GET /api/v1/files/{id}/attributes
					filesByID.Get("/attributes", s.dicomFiles.SearchAttributes)
				})

			})
		})
	})
}

func (s *Server) ListenAndServe() error {
	slog.Info(fmt.Sprintf("-- Starting Server on Port %s --", s.port))
	return http.ListenAndServe(
		fmt.Sprintf(":%s", s.port),
		s.router,
	)
}
