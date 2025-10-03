package sql

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity interface {
	TableName() string
	RepositoryName() string
}

type BaseEntity struct {
	ID        uuid.UUID  `json:"id" gorm:"primarykey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID != uuid.Nil {
		return nil
	}

	b.ID, err = uuid.NewV7()
	return err
}
