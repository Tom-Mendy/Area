package schemas

import (
	"time"
)

// GithubToken represents the GithubToken entity in the database.
type Area struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"           json:"id,omitempty"`
	ServiceId  uint64    `                                          json:"-"` // Foreign key for Service
	Service    Service   `gorm:"foreignKey:ServiceId;references:Id" json:"service_id,omitempty" binding:"required"`
	ActionId   uint64    `                                          json:"-"` // Foreign key for Action
	Action     Action    `gorm:"foreignKey:ActionId;references:Id" json:"action_id,omitempty" binding:"required"`
	ReactionId uint64    `                                          json:"-"` // Foreign key for Reaction
	Reaction   Reaction  `gorm:"foreignKey:ReactionId;references:Id" json:"reaction_id,omitempty" binding:"required"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"          json:"createdAt"`
	UpdateAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"          json:"update_at"`
}
