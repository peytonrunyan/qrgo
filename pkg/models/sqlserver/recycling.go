package sqlserver

import (
	"database/sql"

	"qrgo/pkg/models"
)

// Model for MaterialGuidelinesResults
type MGRModel struct {
	DB *sql.DB
}

func (m *MGRModel) Get(cID string) (*models.MaterialGuidelineResults, error) {
	return nil, nil
}
