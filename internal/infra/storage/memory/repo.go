// Package memory implements the Repository Interface to provide an in-memory storage provider
package memory

import (
	"fmt"

	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// Repo Implements the Repository Interface to provide an in-memory storage provider
type Repo struct {
	favourites map[string]map[string]favourite.Favorite
}

// NewRepo Constructor
func NewRepo() *Repo {
	return &Repo{
		favourites: make(map[string]map[string]favourite.Favorite),
	}
}

func (r *Repo) ensureUser(userID string) {
	if _, ok := r.favourites[userID]; !ok {
		r.favourites[userID] = make(map[string]favourite.Favorite)
	}
}

func (r Repo) GetByID(userID uuid.UUID, favoriteID uuid.UUID) (*favourite.Favorite, error) {
	userMap, ok := r.favourites[userID.String()]
	if !ok {
		return nil, nil
	}

	fav, ok := userMap[favoriteID.String()]
	if !ok {
		return nil, nil
	}
	return &fav, nil
}

// GetAll Returns all stored favourites
func (r Repo) GetAll(userID uuid.UUID) ([]favourite.Favorite, error) {
	userMap, ok := r.favourites[userID.String()]
	if !ok {
		return []favourite.Favorite{}, nil
	}

	values := make([]favourite.Favorite, 0, len(userMap))
	for _, fav := range userMap {
		values = append(values, fav)
	}
	return values, nil
}

func (r *Repo) Add(userID uuid.UUID, favorite favourite.Favorite) error {
	r.ensureUser(userID.String())

	r.favourites[userID.String()][favorite.ID.String()] = favorite
	return nil
}

// Update the provided favourite
func (r *Repo) Update(userID uuid.UUID, favorite favourite.Favorite) error {
	r.ensureUser(userID.String())

	r.favourites[userID.String()][favorite.ID.String()] = favorite
	return nil
}

// Delete the favourite with the provided id
func (r *Repo) Delete(userID uuid.UUID, favoriteID uuid.UUID) error {
	userMap, ok := r.favourites[userID.String()]
	if !ok {
		return fmt.Errorf("user not found")
	}

	if _, exists := userMap[favoriteID.String()]; !exists {
		return fmt.Errorf("favorite %v not found", favoriteID.String())
	}

	delete(userMap, favoriteID.String())
	return nil
}
