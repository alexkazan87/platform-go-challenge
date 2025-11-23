// Package infra contains the services of the interface adapters
package infra

import (
	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/app/notification"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/akazantzidis/gwi-ass/internal/infra/http"
	"github.com/akazantzidis/gwi-ass/internal/infra/notification/console"
	"github.com/akazantzidis/gwi-ass/internal/infra/storage/memory"
)

// Services contains the exposed services of interface adapters
type Services struct {
	NotificationService notification.Service
	FavoriteRepository  favourite.Repository
	Server              *http.Server
}

// NewInfraProviders Instantiates the infra services
func NewInfraProviders() Services {
	return Services{
		NotificationService: console.NewNotificationService(),
		FavoriteRepository:  memory.NewRepo(),
	}
}

// NewHTTPServer creates a new server
func NewHTTPServer(services app.Services) *http.Server {
	return http.NewServer(services)
}
