package schemas

import (
	"time"
)

type AreaMessage struct {
	ActionOption   string `json:"action_option"   binding:"required"`
	ActionId       uint64 `json:"action_id"` // Foreign key for Action
	ReactionOption string `json:"reaction_option" binding:"required"`
	ReactionId     uint64 `json:"reaction_id"` // Foreign key for Reaction
}

type Area struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"            json:"id,omitempty"`
	UserId         uint64    `                                           json:"-"` // Foreign key for User
	User           User      `gorm:"foreignKey:UserId;references:Id"     json:"user,omitempty"     binding:"required"`
	ActionOption   string    `                                           json:"action_option"         binding:"required"`
	ActionId       uint64    `                                           json:"-"` // Foreign key for Action
	Action         Action    `gorm:"foreignKey:ActionId;references:Id"   json:"action,omitempty"   binding:"required"`
	ReactionOption string    `                                           json:"reaction_option"       binding:"required"`
	ReactionId     uint64    `                                           json:"-"` // Foreign key for Reaction
	Reaction       Reaction  `gorm:"foreignKey:ReactionId;references:Id" json:"reaction,omitempty" binding:"required"`
	Enable         bool      `gorm:"default:true"                        json:"enable"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"           json:"created_at"`
	UpdateAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"           json:"update_at"`
}
