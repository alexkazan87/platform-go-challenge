// package main contains the entry point of the application.
package main

import (
	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/infra"
	"github.com/akazantzidis/gwi-ass/internal/pkg/time"
	"github.com/akazantzidis/gwi-ass/internal/pkg/uuid"
)

func main() {
	infraProviders := infra.NewInfraProviders()
	tp := time.NewTimeProvider()
	up := uuid.NewUUIDProvider()
	appServices := app.NewServicesF(infraProviders.FavoriteRepository, infraProviders.NotificationService, up, tp)

	infraHTTPServer := infra.NewHTTPServer(appServices)
	infraHTTPServer.ListenAndServe(":8080")
}
