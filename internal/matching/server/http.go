package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/s3f4/locationmatcher/internal/matching/client"
	"github.com/s3f4/locationmatcher/internal/matching/models"
	"github.com/s3f4/locationmatcher/pkg/apihelper"
	"github.com/s3f4/locationmatcher/pkg/log"
)

type httpServer struct {
}

func (h *httpServer) Start(ctx context.Context) {
	service := os.Getenv("SERVICE")
	port := os.Getenv("PORT")

	var router *chi.Mux = chi.NewRouter()

	router.Post("/api/v1/find_nearest", h.FindNearest)
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apihelper.Send404(w)
	}))

	server := &http.Server{
		Handler: router,
		Addr:    port,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s HTTP server listen err: %s\n", service, err)
		}
	}()

	log.Infof("%s HTTP server started on port %s...\n", service, port)
	<-ctx.Done()
}

func (h *httpServer) FindNearest(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	if err := apihelper.ParseAndValidate(r, &query); err != nil {
		log.Error(err)
		apihelper.Send400(w)
		return
	}

	response, err := client.GetAPIClient().FindNearest("http://driverlocation:3001/api/v1/driver_locations/find_nearest", &query)
	if err != nil {
		log.Error(err)
		apihelper.Send500(w)
		return
	}

	fmt.Printf("%#v\n", response.Data.(map[string]interface{}))
	log.Debug(response.Data)

	apihelper.SendResponse(w, 200, response)
}