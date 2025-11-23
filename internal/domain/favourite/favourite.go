// Package favourite contains the Crag model.
package favourite

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	AssetChart    AssetType = "chart"
	AssetInsight  AssetType = "insight"
	AssetAudience AssetType = "audience"
)

type Favorite struct {
	ID          uuid.UUID
	Type        AssetType
	Description string
	Data        json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
