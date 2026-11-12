package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupBankAccountRoutes(api fiber.Router, h *handlers.Handlers) {
	bankAccounts := api.Group("/bank-accounts")

	// Public endpoint — agent ใช้ดูบัญชีที่เปิดรับโอนเงิน
	bankAccounts.Get("/active", h.BankAccountHandler.ListActive)

	// Admin only — CRUD บัญชีธนาคารบริษัท
	admin := bankAccounts.Group("/")
	admin.Use(middleware.Protected())
	admin.Use(middleware.AdminOnly())

	admin.Get("/", h.BankAccountHandler.List)
	admin.Post("/", h.BankAccountHandler.Create)
	admin.Get("/:id", h.BankAccountHandler.GetByID)
	admin.Put("/:id", h.BankAccountHandler.Update)
	admin.Delete("/:id", h.BankAccountHandler.Delete)
}
