package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"

	"backend/internal/models"
	"backend/internal/repositories"
	"github.com/google/uuid"
)

type WheelService interface {
	CreateWheel(title string, entryInputs []models.WheelEntry, settings models.WheelSettings) (*models.Wheel, string, error)
	GetWheel(shareCode string, requestToken string) (*models.Wheel, error)
	UpdateWheel(id uuid.UUID, title string, entryInputs []models.WheelEntry, settings models.WheelSettings, requestToken string) error
	DeleteWheel(id uuid.UUID, requestToken string) error
	RecordSpin(id uuid.UUID, entryID *uuid.UUID, resultLabel string) (*models.SpinHistory, error)
	GetHistory(id uuid.UUID) ([]models.SpinHistory, error)
	ClearHistory(id uuid.UUID, requestToken string) error
	DuplicateWheel(id uuid.UUID) (*models.Wheel, string, error)
	RecordVisit(ip string, userAgent string) (int64, error)
}

type wheelService struct {
	repo repositories.WheelRepository
}

func NewWheelService(repo repositories.WheelRepository) WheelService {
	return &wheelService{repo: repo}
}

// HashToken applies SHA-256 hashing to cleartext tokens before comparison or storage
func HashToken(token string) string {
	if token == "" {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

// generateRandomString returns a secure random alphanumeric string of length n
func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

func (s *wheelService) CreateWheel(title string, entryInputs []models.WheelEntry, settings models.WheelSettings) (*models.Wheel, string, error) {
	if title == "" {
		title = "Vòng Quay May Mắn"
	}

	shareCode, err := generateRandomString(10)
	if err != nil {
		return nil, "", err
	}

	editToken, err := generateRandomString(20)
	if err != nil {
		return nil, "", err
	}

	wheelID := uuid.New()
	wheel := &models.Wheel{
		ID:            wheelID,
		Title:         title,
		ShareCode:     shareCode,
		EditTokenHash: HashToken(editToken),
		Settings:      settings,
		Permission:    "edit",
	}

	entries := make([]models.WheelEntry, len(entryInputs))
	for i, entry := range entryInputs {
		entries[i] = models.WheelEntry{
			ID:       uuid.New(),
			WheelID:  wheelID,
			Label:    entry.Label,
			Weight:   entry.Weight,
			Color:    entry.Color,
			Position: i,
		}
	}
	wheel.Entries = entries

	if err := s.repo.Create(wheel); err != nil {
		return nil, "", err
	}

	return wheel, editToken, nil
}

func (s *wheelService) GetWheel(shareCode string, requestToken string) (*models.Wheel, error) {
	wheel, err := s.repo.GetByShareCode(shareCode)
	if err != nil {
		return nil, err
	}

	// Validate edit token to grant creator access
	if requestToken != "" && HashToken(requestToken) == wheel.EditTokenHash {
		wheel.Permission = "edit"
	} else {
		wheel.Permission = "spin"
	}

	return wheel, nil
}

func (s *wheelService) UpdateWheel(id uuid.UUID, title string, entryInputs []models.WheelEntry, settings models.WheelSettings, requestToken string) error {
	wheel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if HashToken(requestToken) != wheel.EditTokenHash {
		return errors.New("unauthorized: invalid edit token")
	}

	wheel.Title = title
	wheel.Settings = settings

	// Construct updated list of entries
	entries := make([]models.WheelEntry, len(entryInputs))
	for i, entry := range entryInputs {
		entries[i] = models.WheelEntry{
			ID:       uuid.New(),
			WheelID:  id,
			Label:    entry.Label,
			Weight:   entry.Weight,
			Color:    entry.Color,
			Position: i,
		}
	}
	wheel.Entries = entries

	return s.repo.Update(wheel)
}

func (s *wheelService) DeleteWheel(id uuid.UUID, requestToken string) error {
	wheel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if HashToken(requestToken) != wheel.EditTokenHash {
		return errors.New("unauthorized: invalid edit token")
	}

	return s.repo.Delete(id)
}

func (s *wheelService) RecordSpin(wheelID uuid.UUID, entryID *uuid.UUID, resultLabel string) (*models.SpinHistory, error) {
	history := &models.SpinHistory{
		ID:          uuid.New(),
		WheelID:     wheelID,
		EntryID:     entryID,
		ResultLabel: resultLabel,
	}

	if err := s.repo.CreateSpinHistory(history); err != nil {
		return nil, err
	}
	return history, nil
}

func (s *wheelService) GetHistory(wheelID uuid.UUID) ([]models.SpinHistory, error) {
	return s.repo.GetSpinHistory(wheelID)
}

func (s *wheelService) ClearHistory(wheelID uuid.UUID, requestToken string) error {
	wheel, err := s.repo.GetByID(wheelID)
	if err != nil {
		return err
	}

	if HashToken(requestToken) != wheel.EditTokenHash {
		return errors.New("unauthorized: invalid edit token")
	}

	return s.repo.ClearSpinHistory(wheelID)
}

func (s *wheelService) DuplicateWheel(id uuid.UUID) (*models.Wheel, string, error) {
	srcWheel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, "", err
	}

	return s.CreateWheel(srcWheel.Title+" (Bản sao)", srcWheel.Entries, srcWheel.Settings)
}

func (s *wheelService) RecordVisit(ip string, userAgent string) (int64, error) {
	return s.repo.RecordVisit(ip, userAgent)
}
