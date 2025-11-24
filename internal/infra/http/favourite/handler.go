package favourite

import (
	"encoding/json"
	"fmt"
	queries2 "github.com/akazantzidis/gwi-ass/internal/app/user/queries"
	"net/http"

	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/queries"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler Favorite http request Handler
type Handler struct {
	favoriteServices app.FavoriteServices
	userServices     app.UserServices
}

func NewHandler(app app.FavoriteServices, userApp app.UserServices) *Handler {
	return &Handler{favoriteServices: app, userServices: userApp}
}

// ErrorResponse standard JSON error response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// writeJSONError writes an error response as JSON
func writeJSONError(w http.ResponseWriter, status int, err error, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err.Error(),
		Details: details,
	})
}

// URL param constants
const (
	UserIDURLParam           = "userID"
	GetFavoriteIDURLParam    = "favoriteId"
	UpdateFavoriteIDURLParam = "favoriteId"
	DeleteFavoriteIDURLParam = "cragId"
)

// GetAllF returns all favorites for a given user
func (c Handler) GetAllF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"), nil)
		return
	}

	favorites, err := c.favoriteServices.Queries.GetAllFavoritesHandler.Handle(
		queries.GetAllFavoritesRequest{UserID: userID},
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

// GetByIDF returns a favorite for a given user
func (c Handler) GetByIDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"), nil)
		return
	}

	favoriteID, err := uuid.Parse(vars[GetFavoriteIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	favorite, err := c.favoriteServices.Queries.GetFavoriteHandler.Handle(
		queries.GetFavoriteRequest{
			UserID:     userID,
			FavoriteID: favoriteID,
		},
	)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}
	if favorite == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Errorf("favorite not found for this user"), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorite)
}

// CreateFavoriteRequestModel represents the request model expected for Add request
type CreateFavoriteRequestModel struct {
	Type        favourite.AssetType `json:"type"`
	Description string              `json:"description"`
	Data        json.RawMessage     `json:"data"`
}

// Create adds a new favorite for a user
func (c Handler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"), nil)
		return
	}

	var req CreateFavoriteRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"), nil)
		return
	}

	_, err = c.userServices.Queries.GetUserHandler.Handle(queries2.GetUserRequest{
		ID: userID,
	})

	if err != nil {
		writeJSONError(w, http.StatusNotFound, err, nil)
		return
	}

	err = c.favoriteServices.Commands.CreateFavoriteHandler.Handle(
		commands.AddFavoriteRequest{
			UserID:      userID,
			Type:        req.Type,
			Description: req.Description,
			Data:        req.Data,
		},
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateFavoriteRequestModel represents the request model of Update
type UpdateFavoriteRequestModel struct {
	ID          uuid.UUID           `json:"id"`
	Type        favourite.AssetType `json:"type"`
	Description string              `json:"description"`
	Data        json.RawMessage     `json:"data"`
}

// PatchFavoriteRequestModel represents the fields that can be partially updated
type PatchFavoriteRequestModel struct {
	Type        *favourite.AssetType `json:"type,omitempty"`
	Description *string              `json:"description,omitempty"`
	Data        *json.RawMessage     `json:"data,omitempty"`
}

// Patch handles partial updates for favorites
func (c Handler) Patch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"), nil)
		return
	}

	favID, err := uuid.Parse(vars[UpdateFavoriteIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	var req PatchFavoriteRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"), nil)
		return
	}

	reqCommand := commands.PatchFavoriteRequest{
		Type:        req.Type,
		Description: req.Description,
		Data:        req.Data,
	}

	fav, err := c.favoriteServices.Commands.UpdatePartialFavoriteHandler.HandlePartial(userID, favID, reqCommand)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fav)
}

// Update handles full updates
func (c Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"), nil)
		return
	}

	favoriteID, err := uuid.Parse(vars[UpdateFavoriteIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	var req UpdateFavoriteRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"), nil)
		return
	}

	err = c.favoriteServices.Commands.UpdateFavoriteHandler.Handle(
		commands.UpdateFavoriteRequest{
			UserID:      userID,
			ID:          favoriteID,
			Type:        req.Type,
			Description: req.Description,
			Data:        req.Data,
		},
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete deletes a favorite
func (c Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"), nil)
		return
	}

	favoriteID, err := uuid.Parse(vars[DeleteFavoriteIDURLParam])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	err = c.favoriteServices.Commands.DeleteFavoriteHandler.Handle(
		commands.DeleteFavoriteRequest{
			UserID:     userID,
			FavoriteID: favoriteID,
		},
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
