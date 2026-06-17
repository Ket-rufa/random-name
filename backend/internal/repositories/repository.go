package repositories

import (
	"backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WheelRepository interface {
	Create(wheel *models.Wheel) error
	GetByShareCode(code string) (*models.Wheel, error)
	GetByID(id uuid.UUID) (*models.Wheel, error)
	Update(wheel *models.Wheel) error
	Delete(id uuid.UUID) error
	DeleteEntries(wheelID uuid.UUID) error
	CreateSpinHistory(history *models.SpinHistory) error
	GetSpinHistory(wheelID uuid.UUID) ([]models.SpinHistory, error)
	ClearSpinHistory(wheelID uuid.UUID) error
	RecordVisit(ip string, userAgent string) (int64, error)
}

type wheelRepository struct {
	db *gorm.DB
}

func NewWheelRepository(db *gorm.DB) WheelRepository {
	return &wheelRepository{db: db}
}

func (r *wheelRepository) Create(wheel *models.Wheel) error {
	return r.db.Create(wheel).Error
}

func (r *wheelRepository) GetByShareCode(code string) (*models.Wheel, error) {
	var wheel models.Wheel
	err := r.db.Preload("Entries").First(&wheel, "share_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &wheel, nil
}

func (r *wheelRepository) GetByID(id uuid.UUID) (*models.Wheel, error) {
	var wheel models.Wheel
	err := r.db.Preload("Entries").First(&wheel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wheel, nil
}

func (r *wheelRepository) Update(wheel *models.Wheel) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Clean up old entries
		if err := tx.Where("wheel_id = ?", wheel.ID).Delete(&models.WheelEntry{}).Error; err != nil {
			return err
		}
		// Write the new updated entries
		if len(wheel.Entries) > 0 {
			if err := tx.Create(&wheel.Entries).Error; err != nil {
				return err
			}
		}
		// Save wheel metadata (omitting entries relation handled manually above)
		return tx.Omit("Entries").Save(wheel).Error
	})
}

func (r *wheelRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Wheel{}, "id = ?", id).Error
}

func (r *wheelRepository) DeleteEntries(wheelID uuid.UUID) error {
	return r.db.Where("wheel_id = ?", wheelID).Delete(&models.WheelEntry{}).Error
}

func (r *wheelRepository) CreateSpinHistory(history *models.SpinHistory) error {
	return r.db.Create(history).Error
}

func (r *wheelRepository) GetSpinHistory(wheelID uuid.UUID) ([]models.SpinHistory, error) {
	var history []models.SpinHistory
	err := r.db.Where("wheel_id = ?", wheelID).Order("spun_at desc").Limit(100).Find(&history).Error
	return history, err
}

func (r *wheelRepository) ClearSpinHistory(wheelID uuid.UUID) error {
	return r.db.Where("wheel_id = ?", wheelID).Delete(&models.SpinHistory{}).Error
}

func (r *wheelRepository) RecordVisit(ip string, userAgent string) (int64, error) {
	visit := models.PageView{
		IP:        ip,
		UserAgent: userAgent,
	}
	if err := r.db.Create(&visit).Error; err != nil {
		return 0, err
	}

	var count int64
	if err := r.db.Model(&models.PageView{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
