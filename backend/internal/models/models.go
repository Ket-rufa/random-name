package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type WheelSettings struct {
	Theme              string   `json:"theme"`
	CustomColors       []string `json:"customColors"`
	SpinDuration       int      `json:"spinDuration"`
	Volume             int      `json:"volume"`
	EnableSound        bool     `json:"enableSound"`
	EnableTickSound    bool     `json:"enableTickSound"`
	EnableVictorySound bool     `json:"enableVictorySound"`
	EnableConfetti     bool     `json:"enableConfetti"`
	AllowDuplicates    bool     `json:"allowDuplicates"`
	AutoRemoveWinner   bool     `json:"autoRemoveWinner"`
	EnableWeights      bool     `json:"enableWeights"`
	ShowTextOnWheel    bool     `json:"showTextOnWheel"`
	FontSize           int      `json:"fontSize"`
}

// Scan scans db value into WheelSettings struct (implements sql.Scanner)
func (ws *WheelSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, ws)
}

// Value serializes WheelSettings struct into driver Value (implements driver.Valuer)
func (ws WheelSettings) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

type Wheel struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title         string        `gorm:"type:varchar(255);not null" json:"title"`
	ShareCode     string        `gorm:"type:varchar(32);uniqueIndex;not null" json:"shareCode"`
	EditTokenHash string        `gorm:"type:text;not null" json:"-"`
	Settings      WheelSettings `gorm:"type:jsonb;not null" json:"settings"`
	Permission    string        `gorm:"type:varchar(32);default:'spin'" json:"permission"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Entries       []WheelEntry  `gorm:"foreignKey:WheelID;constraint:OnDelete:CASCADE" json:"entries"`
}

type WheelEntry struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WheelID   uuid.UUID `gorm:"type:uuid;not null;index" json:"wheelId"`
	Label     string    `gorm:"type:varchar(255);not null" json:"label"`
	Weight    float64   `gorm:"type:numeric;default:1" json:"weight"`
	Color     string    `gorm:"type:varchar(20)" json:"color"`
	Position  int       `gorm:"type:integer;not null" json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SpinHistory struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WheelID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"wheelId"`
	EntryID     *uuid.UUID `gorm:"type:uuid" json:"entryId"`
	ResultLabel string     `gorm:"type:varchar(255);not null" json:"resultLabel"`
	SpunAt      time.Time  `gorm:"autoCreateTime" json:"spunAt"`
}

type PageView struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IP        string    `gorm:"type:varchar(50)" json:"ip"`
	UserAgent string    `gorm:"type:text" json:"userAgent"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
