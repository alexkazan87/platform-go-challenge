// Package app contains the application layer of the application. It is the layer that contains the business logic of the application. It is the layer that interacts with the domain layer and the interface adapters. It is the layer that contains the use cases of the application
package app

import (
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/queries"
	"github.com/akazantzidis/gwi-ass/internal/app/notification"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/akazantzidis/gwi-ass/internal/pkg/time"
	"github.com/akazantzidis/gwi-ass/internal/pkg/uuid"
)

// Queries Contains all available query handlers of this app
type Queries struct {
	GetAllFavoritesHandler queries.GetAllFavoritesRequestHandler
	GetFavoriteHandler     queries.GetFavoriteRequestHandler
}

// Commands Contains all available command handlers of this app
type Commands struct {
	CreateFavoriteHandler        commands.CreateFavoriteRequestHandler
	UpdateFavoriteHandler        commands.UpdateFavoriteRequestHandler
	UpdatePartialFavoriteHandler commands.UpdatePartialFavoriteRequestHandler

	DeleteFavoriteHandler commands.DeleteFavoriteRequestHandler
}

// FavoriteServices Contains the grouped queries and commands of the app layer
type FavoriteServices struct {
	Queries  Queries
	Commands Commands
}

// Services contains all exposed services of the application layer
type Services struct {
	FavoriteServices FavoriteServices
}

// NewServices Bootstraps Application Layer dependencies
func NewServicesF(favoriteRepo favourite.Repository, ns notification.Service, _ uuid.Provider, _ time.Provider) Services {
	return Services{
		FavoriteServices: FavoriteServices{
			Queries: Queries{
				GetAllFavoritesHandler: queries.NewGetAllFavoritesRequestHandler(favoriteRepo),
				GetFavoriteHandler:     queries.NewGetFavoriteRequestHandler(favoriteRepo),
			},
			Commands: Commands{
				CreateFavoriteHandler:        commands.NewAddFavoriteRequestHandler(favoriteRepo, ns),
				UpdateFavoriteHandler:        commands.NewUpdateFavoriteRequestHandler(favoriteRepo),
				UpdatePartialFavoriteHandler: commands.NewUpdatePartialFavoriteRequestHandler(favoriteRepo),

				DeleteFavoriteHandler: commands.NewDeleteFavoriteRequestHandler(favoriteRepo),
			},
		},
	}
}
