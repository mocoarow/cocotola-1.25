package gateway

import (
	"fmt"
	"time"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type BaseModelEntity struct {
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	UpdatedBy int
}

func (e *BaseModelEntity) ToBaseModel() (*libdomain.BaseModel, error) {
	model, err := libdomain.NewBaseModel(e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, fmt.Errorf("validate base model: %w", err)
	}

	return model, nil
}

type JunctionModelEntity struct {
	CreatedAt time.Time
	CreatedBy int
}
