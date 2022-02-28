package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/s3f4/locationmatcher/internal/matching/client"
	"github.com/s3f4/locationmatcher/internal/matching/models"
	"github.com/s3f4/locationmatcher/internal/matching/server/middlewares"
	"github.com/s3f4/locationmatcher/pkg/apihelper"
	"github.com/s3f4/locationmatcher/pkg/log"
)

type httpServer struct {
	client client.APIClient
}

func (h *httpServer) Start(ctx context.Context) {
	service := os.Getenv("SERVICE")
	port := os.Getenv("PORT")

	var router *chi.Mux = chi.NewRouter()

	router.Route("/api/v1", func(router chi.Router) {
		router.Use(middlewares.AuthCtx)
		router.Post("/find_nearest", h.FindNearest)
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

func (h *httpServer) FindNearest(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	var query models.Query
	if err := apihelper.ParseAndValidate(r, &query); err != nil {
		apihelper.SendResponse(w, err.Code, apihelper.Response{
			Code: err.Code,
			Msg:  err.Msg,
		})
		return
	}

	response, err := h.client.FindNearest(context, "http://driverlocation:3001/api/v1/driver_locations/find_nearest", &query)
	if err != nil {
		log.Error(err)
		apihelper.Send500(w)
		return
	}

	if response.Code == http.StatusNotFound {
		apihelper.Send404(w)
		return
	}

	locationsData, ok := response.Data.(map[string]interface{})["locations"].([]interface{})
	if !ok {
		apihelper.Send400(w)
		return
	}

	apihelper.SendResponse(w, 200, locationsData[0])
}
