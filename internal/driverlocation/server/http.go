package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
	"github.com/s3f4/locationmatcher/internal/driverlocation/repository"
	"github.com/s3f4/locationmatcher/internal/driverlocation/server/middlewares"
	"github.com/s3f4/locationmatcher/pkg/apihelper"
	"github.com/s3f4/locationmatcher/pkg/log"
)

type httpServer struct {
	repository repository.Repository
}

// Start starts http server
func (h *httpServer) Start(ctx context.Context, repository repository.Repository) {
	service := os.Getenv("SERVICE")
	port := os.Getenv("PORT")

	h.repository = repository

	var router *chi.Mux = chi.NewRouter()

	router.Route("/api/v1/driver_locations", func(router chi.Router) {
		router.Use(middlewares.AuthCtx)
		router.Post("/", h.UpsertBulk)
		router.Post("/find_nearest", h.Find)
	})

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apihelper.Send404(w)
	}))

	server := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s HTTP server listen err: %s\n", service, err)
		}
	}()

	log.Infof("%s HTTP server started on port %s...\n", service, port)
	<-ctx.Done()
	log.Infof("%s HTTP server stopped. \n", service)
}

// UpsertBulk is used to create or/and update driver locations.
func (h *httpServer) UpsertBulk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var driverLocations []*models.DriverLocation
	if err := json.NewDecoder(r.Body).Decode(&driverLocations); err != nil {
		log.Error(err)
		apihelper.Send400(w)
		return
	}

	notValidResponse := apihelper.Response{
		Code: http.StatusBadRequest,
		Msg:  "provide valid driver locations",
	}

	if len(driverLocations) == 0 {
		apihelper.SendResponse(w, http.StatusBadRequest, notValidResponse)
		return
	} else {
		for _, driverLocation := range driverLocations {
			if err := driverLocation.Validate(); err != nil {
				log.Error(err)
				apihelper.SendResponse(w, http.StatusBadRequest, notValidResponse)
				return
			}
		}
	}

	if err := h.repository.UpsertBulk(ctx, driverLocations); err != nil {
		log.Error(err)
		apihelper.Send500(w)
		return
	}

	apihelper.SendResponse(w, http.StatusOK, driverLocations)
}

// Find returns nearest locations within the given query parameters
func (h *httpServer) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var query models.Query
	if err := apihelper.ParseAndValidate(r, &query); err != nil {
		log.Error(err)
		apihelper.SendResponse(w, err.Code, apihelper.Response{
			Code: err.Code,
			Msg:  err.Msg,
		})
		return
	}

	log.Debug(query)

	locations, err := h.repository.Find(ctx, &query)
	if err != nil {
		log.Error(err)
		apihelper.Send500(w)
		return
	}

	if len(locations) == 0 {
		apihelper.Send404(w)
		return
	}

	apihelper.SendResponse(w, http.StatusOK,
		apihelper.Response{
			Code: 200,
			Data: models.LocationsResponse{
				Total:     len(locations),
				Locations: locations,
			},
		},
	)
}
