package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WheelHandler struct {
	service services.WheelService
}

func NewWheelHandler(service services.WheelService) *WheelHandler {
	return &WheelHandler{service: service}
}

// REST response helpers
func sendSuccess(c *fiber.Ctx, msg string, data interface{}, status ...int) error {
	code := fiber.StatusOK
	if len(status) > 0 {
		code = status[0]
	}
	return c.Status(code).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

func sendError(c *fiber.Ctx, msg string, status int, errors ...interface{}) error {
	var errs interface{}
	if len(errors) > 0 {
		errs = errors[0]
	} else {
		errs = []string{}
	}
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": msg,
		"errors":  errs,
	})
}

func (h *WheelHandler) HealthCheck(c *fiber.Ctx) error {
	return sendSuccess(c, "Server is healthy", fiber.Map{"status": "ok"})
}

type CreateRequest struct {
	Title    string               `json:"title"`
	Entries  []models.WheelEntry  `json:"entries"`
	Settings models.WheelSettings `json:"settings"`
}

func (h *WheelHandler) CreateWheel(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, "Invalid request payload", fiber.StatusBadRequest)
	}

	// Inputs validation
	if len(req.Entries) == 0 {
		return sendError(c, "Lựa chọn không được rỗng", fiber.StatusBadRequest)
	}
	if len(req.Entries) > 500 {
		return sendError(c, "Vượt quá giới hạn 500 lựa chọn", fiber.StatusBadRequest)
	}

	for _, entry := range req.Entries {
		if entry.Label == "" {
			return sendError(c, "Nhãn lựa chọn không được trống", fiber.StatusBadRequest)
		}
		if len(entry.Label) > 255 {
			return sendError(c, "Nhãn lựa chọn tối đa 255 ký tự", fiber.StatusBadRequest)
		}
		if entry.Weight <= 0 {
			return sendError(c, "Trọng số phải lớn hơn 0", fiber.StatusBadRequest)
		}
	}

	wheel, editToken, err := h.service.CreateWheel(req.Title, req.Entries, req.Settings)
	if err != nil {
		return sendError(c, "Failed to create wheel", fiber.StatusInternalServerError)
	}

	return sendSuccess(c, "Wheel created successfully", fiber.Map{
		"wheel":     wheel,
		"editToken": editToken,
	}, fiber.StatusCreated)
}

func (h *WheelHandler) GetWheel(c *fiber.Ctx) error {
	shareCode := c.Params("shareCode")
	if shareCode == "" {
		return sendError(c, "Share code is required", fiber.StatusBadRequest)
	}

	editToken := c.Get("X-Edit-Token")

	wheel, err := h.service.GetWheel(shareCode, editToken)
	if err != nil {
		return sendError(c, "Vòng quay không tồn tại hoặc đã bị xóa", fiber.StatusNotFound)
	}

	return sendSuccess(c, "Wheel loaded successfully", fiber.Map{
		"wheel":     wheel,
		"editToken": editToken, // return back edit token if matched
	})
}

type UpdateRequest struct {
	Title    string               `json:"title"`
	Entries  []models.WheelEntry  `json:"entries"`
	Settings models.WheelSettings `json:"settings"`
}

func (h *WheelHandler) UpdateWheel(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return sendError(c, "Invalid wheel ID format", fiber.StatusBadRequest)
	}

	editToken := c.Get("X-Edit-Token")
	if editToken == "" {
		return sendError(c, "Unauthorized: missing edit token", fiber.StatusUnauthorized)
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, "Invalid request payload", fiber.StatusBadRequest)
	}

	// Validation
	if len(req.Entries) == 0 {
		return sendError(c, "Lựa chọn không được rỗng", fiber.StatusBadRequest)
	}
	if len(req.Entries) > 500 {
		return sendError(c, "Vượt quá giới hạn 500 lựa chọn", fiber.StatusBadRequest)
	}

	for _, entry := range req.Entries {
		if entry.Label == "" {
			return sendError(c, "Nhãn lựa chọn không được trống", fiber.StatusBadRequest)
		}
		if len(entry.Label) > 255 {
			return sendError(c, "Nhãn lựa chọn tối đa 255 ký tự", fiber.StatusBadRequest)
		}
		if entry.Weight <= 0 {
			return sendError(c, "Trọng số phải lớn hơn 0", fiber.StatusBadRequest)
		}
	}

	err = h.service.UpdateWheel(id, req.Title, req.Entries, req.Settings, editToken)
	if err != nil {
		return sendError(c, err.Error(), fiber.StatusForbidden)
	}

	return sendSuccess(c, "Wheel updated successfully", nil)
}

func (h *WheelHandler) DeleteWheel(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return sendError(c, "Invalid wheel ID format", fiber.StatusBadRequest)
	}

	editToken := c.Get("X-Edit-Token")
	if editToken == "" {
		return sendError(c, "Unauthorized: missing edit token", fiber.StatusUnauthorized)
	}

	err = h.service.DeleteWheel(id, editToken)
	if err != nil {
		return sendError(c, err.Error(), fiber.StatusForbidden)
	}

	return sendSuccess(c, "Wheel deleted successfully", nil)
}

type SpinRequest struct {
	EntryID     *string `json:"entryId"`
	ResultLabel string  `json:"resultLabel"`
}

func (h *WheelHandler) RecordSpin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return sendError(c, "Invalid wheel ID format", fiber.StatusBadRequest)
	}

	var req SpinRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, "Invalid request payload", fiber.StatusBadRequest)
	}

	if req.ResultLabel == "" {
		return sendError(c, "Result label is required", fiber.StatusBadRequest)
	}

	var entryUUID *uuid.UUID
	if req.EntryID != nil && *req.EntryID != "" {
		u, err := uuid.Parse(*req.EntryID)
		if err == nil {
			entryUUID = &u
		}
	}

	history, err := h.service.RecordSpin(id, entryUUID, req.ResultLabel)
	if err != nil {
		return sendError(c, "Failed to record spin history", fiber.StatusInternalServerError)
	}

	return sendSuccess(c, "Spin history recorded successfully", history, fiber.StatusCreated)
}

func (h *WheelHandler) GetHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return sendError(c, "Invalid wheel ID format", fiber.StatusBadRequest)
	}

	history, err := h.service.GetHistory(id)
	if err != nil {
		return sendError(c, "Failed to load spin history", fiber.StatusInternalServerError)
	}

	return sendSuccess(c, "Spin history loaded successfully", history)
}

func (h *WheelHandler) ClearHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return sendError(c, "Invalid wheel ID format", fiber.StatusBadRequest)
	}

	editToken := c.Get("X-Edit-Token")
	if editToken == "" {
		return sendError(c, "Unauthorized: missing edit token", fiber.StatusUnauthorized)
	}

	err = h.service.ClearHistory(id, editToken)
	if err != nil {
		return sendError(c, err.Error(), fiber.StatusForbidden)
	}

	return sendSuccess(c, "Spin history cleared successfully", nil)
}

func (h *WheelHandler) DuplicateWheel(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return sendError(c, "Invalid wheel ID format", fiber.StatusBadRequest)
	}

	wheel, editToken, err := h.service.DuplicateWheel(id)
	if err != nil {
		return sendError(c, "Failed to duplicate wheel", fiber.StatusInternalServerError)
	}

	return sendSuccess(c, "Wheel duplicated successfully", fiber.Map{
		"wheel":     wheel,
		"editToken": editToken,
	}, fiber.StatusCreated)
}

func (h *WheelHandler) RecordVisit(c *fiber.Ctx) error {
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	count, err := h.service.RecordVisit(ip, userAgent)
	if err != nil {
		return sendError(c, "Failed to record visit", fiber.StatusInternalServerError)
	}

	return sendSuccess(c, "Visit recorded successfully", fiber.Map{
		"totalVisits": count,
	})
}
