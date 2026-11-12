package routes

import (
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(api fiber.Router, h *handlers.Handlers) {
	admin := api.Group("/admin")

	// All admin routes require authentication and admin+ permissions
	admin.Use(middleware.Protected())
	admin.Use(middleware.ManagementOnly())

	// User management routes
	admin.Post("/users", h.UserHandler.Create)                   // POST /api/v1/admin/users
	admin.Get("/users", h.UserHandler.List)                      // GET /api/v1/admin/users
	admin.Get("/users/stats", h.UserHandler.GetStats)            // GET /api/v1/admin/users/stats
	admin.Get("/users/:id", h.UserHandler.GetByID)               // GET /api/v1/admin/users/:id
	admin.Put("/users/:id", h.UserHandler.Update)                // PUT /api/v1/admin/users/:id
	admin.Patch("/users/:id/status", h.UserHandler.UpdateStatus) // PATCH /api/v1/admin/users/:id/status
	admin.Patch("/users/:id/role", h.UserHandler.ChangeRole)     // PATCH /api/v1/admin/users/:id/role
	admin.Delete("/users/:id", h.UserHandler.Delete)             // DELETE /api/v1/admin/users/:id

	// Sales management routes
	admin.Get("/sales", h.SalesHandler.List)                        // GET /api/v1/admin/sales
	admin.Post("/sales", h.SalesHandler.Create)                     // POST /api/v1/admin/sales
	admin.Get("/sales/:id", h.SalesHandler.GetByID)                 // GET /api/v1/admin/sales/:id
	admin.Put("/sales/:id", h.SalesHandler.Update)                  // PUT /api/v1/admin/sales/:id
	admin.Delete("/sales/:id", h.SalesHandler.Delete)               // DELETE /api/v1/admin/sales/:id
	admin.Patch("/sales/:id/activate", h.SalesHandler.Activate)     // PATCH /api/v1/admin/sales/:id/activate
	admin.Patch("/sales/:id/deactivate", h.SalesHandler.Deactivate) // PATCH /api/v1/admin/sales/:id/deactivate

	// City management routes
	admin.Get("/cities", h.CityHandler.List)                    // GET /api/v1/admin/cities
	admin.Post("/cities", h.CityHandler.Create)                 // POST /api/v1/admin/cities
	admin.Get("/cities/:id", h.CityHandler.GetByID)             // GET /api/v1/admin/cities/:id
	admin.Put("/cities/:id", h.CityHandler.Update)              // PUT /api/v1/admin/cities/:id
	admin.Delete("/cities/:id", h.CityHandler.Delete)           // DELETE /api/v1/admin/cities/:id
	admin.Put("/cities/:id/status", h.CityHandler.UpdateStatus) // PUT /api/v1/admin/cities/:id/status

	// Tour management routes
	admin.Post("/tours", h.TourHandler.Create)                   // POST /api/v1/admin/tours
	admin.Get("/tours", h.TourHandler.List)                      // GET /api/v1/admin/tours (admin view)
	admin.Get("/tours/:id", h.TourHandler.GetByID)               // GET /api/v1/admin/tours/:id
	admin.Get("/tours/:id/periods", h.TourHandler.GetPeriods)    // GET /api/v1/admin/tours/:id/periods
	admin.Put("/tours/:id", h.TourHandler.Update)                // PUT /api/v1/admin/tours/:id
	admin.Delete("/tours/:id", h.TourHandler.Delete)             // DELETE /api/v1/admin/tours/:id
	admin.Patch("/tours/:id/status", h.TourHandler.UpdateStatus) // PATCH /api/v1/admin/tours/:id/status

	// Period management routes
	admin.Get("/periods", h.PeriodHandler.List)                             // GET /api/v1/admin/periods
	admin.Post("/periods", h.PeriodHandler.Create)                          // POST /api/v1/admin/periods
	admin.Get("/periods/:id", h.PeriodHandler.GetByID)                      // GET /api/v1/admin/periods/:id
	admin.Put("/periods/:id", h.PeriodHandler.Update)                       // PUT /api/v1/admin/periods/:id
	admin.Delete("/periods/:id", h.PeriodHandler.Delete)                    // DELETE /api/v1/admin/periods/:id
	admin.Put("/periods/:id/status", h.PeriodHandler.UpdateStatus)          // PUT /api/v1/admin/periods/:id/status
	admin.Get("/periods/:id/availability", h.PeriodHandler.GetAvailability) // GET /api/v1/admin/periods/:id/availability
	admin.Get("/periods/stats", h.PeriodHandler.GetStats)                   // GET /api/v1/admin/periods/stats
	//	admin.Get("/periods/test-deployment", h.PeriodHandler.TestDeployment)            // GET /api/v1/admin/periods/test-deployment
	//	admin.Get("/periods/:id/debug-booking-stats", h.PeriodHandler.DebugBookingStats) // GET /api/v1/admin/periods/:id/debug-booking-stats

	// Seat Hold Management Routes
	admin.Post("/periods/:id/hold-seats", h.SeatHoldHandler.CreateSeatHold)                           // POST /api/v1/admin/periods/:id/hold-seats
	admin.Get("/periods/:id/seat-holds", h.SeatHoldHandler.GetPeriodSeatHolds)                        // GET /api/v1/admin/periods/:id/seat-holds
	admin.Get("/periods/:id/seat-availability", h.SeatHoldHandler.CheckSeatAvailability)              // GET /api/v1/admin/periods/:id/seat-availability
	admin.Get("/seat-holds", h.SeatHoldHandler.GetActiveSeatHolds)                                    // GET /api/v1/admin/seat-holds
	admin.Get("/seat-holds/:hold_id", h.SeatHoldHandler.GetSeatHold)                                  // GET /api/v1/admin/seat-holds/:hold_id
	admin.Patch("/seat-holds/:hold_id", h.SeatHoldHandler.UpdateSeatHold)                             // PATCH /api/v1/admin/seat-holds/:hold_id
	admin.Delete("/seat-holds/:hold_id/release", h.SeatHoldHandler.ReleaseSeatHold)                   // DELETE /api/v1/admin/seat-holds/:hold_id/release
	admin.Post("/seat-holds/:hold_id/extend", h.SeatHoldHandler.ExtendSeatHold)                       // POST /api/v1/admin/seat-holds/:hold_id/extend
	admin.Delete("/seat-holds/:hold_id", h.SeatHoldHandler.DeleteSeatHold)                            // DELETE /api/v1/admin/seat-holds/:hold_id (super admin only)
	admin.Post("/seat-holds/:hold_id/convert-to-booking", h.SeatHoldHandler.ConvertSeatHoldToBooking) // POST /api/v1/admin/seat-holds/:hold_id/convert-to-booking
	admin.Post("/seat-holds/expire-old", h.SeatHoldHandler.ExpireOldSeatHolds)                        // POST /api/v1/admin/seat-holds/expire-old
	admin.Get("/admins/:admin_id/seat-holds", h.SeatHoldHandler.GetAdminSeatHolds)                    // GET /api/v1/admin/admins/:admin_id/seat-holds

	// Banner management routes
	admin.Get("/banners", h.BannerHandler.List)                    // GET /api/v1/admin/banners
	admin.Post("/banners", h.BannerHandler.Create)                 // POST /api/v1/admin/banners
	admin.Get("/banners/:id", h.BannerHandler.GetByID)             // GET /api/v1/admin/banners/:id
	admin.Put("/banners/:id", h.BannerHandler.Update)              // PUT /api/v1/admin/banners/:id
	admin.Delete("/banners/:id", h.BannerHandler.Delete)           // DELETE /api/v1/admin/banners/:id
	admin.Put("/banners/:id/status", h.BannerHandler.UpdateStatus) // PUT /api/v1/admin/banners/:id/status
	admin.Put("/banners/:id/order", h.BannerHandler.UpdateOrder)   // PUT /api/v1/admin/banners/:id/order

	// Booking management routes
	admin.Get("/bookings", h.BookingHandler.List)                                    // GET /api/v1/admin/bookings
	admin.Get("/bookings/stats", h.BookingHandler.GetStats)                          // GET /api/v1/admin/bookings/stats
	admin.Get("/bookings/:id", h.BookingHandler.GetByID)                             // GET /api/v1/admin/bookings/:id
	admin.Put("/bookings/:id/status", h.BookingHandler.UpdateStatus)                 // PUT /api/v1/admin/bookings/:id/status
	admin.Post("/bookings/:bookingId/payments", h.BookingHandler.CreateAdminPayment) // POST /api/v1/admin/bookings/:bookingId/payments
	admin.Post("/bookings/:bookingId/passport", h.BookingHandler.AdminUploadPassport) // POST /api/v1/admin/bookings/:bookingId/passport

	// Payment management routes
	admin.Get("/payments/pending", h.BookingHandler.GetPendingPayments)                        // GET /api/v1/admin/payments/pending
	admin.Post("/payment-slips/:paymentSlipId/verify", h.BookingHandler.VerifyPaymentSlip)     // POST /api/v1/admin/payment-slips/:paymentSlipId/verify
	admin.Post("/payment-slips/:paymentSlipId/reverse", h.BookingHandler.ReversePaymentSlip)   // POST /api/v1/admin/payment-slips/:paymentSlipId/reverse
	admin.Get("/payment-slips/:paymentSlipId/history", h.BookingHandler.GetPaymentSlipHistory) // GET /api/v1/admin/payment-slips/:paymentSlipId/history
	admin.Get("/payment-slips/:paymentSlipId", h.BookingHandler.GetPaymentSlipDetails)         // GET /api/v1/admin/payment-slips/:paymentSlipId

	// Company profile management routes
	admin.Get("/companies", h.CompanyHandler.ListProfiles)                     // GET /api/v1/admin/companies (pending profiles only)
	admin.Get("/companies/search", h.CompanyHandler.SearchProfiles)            // GET /api/v1/admin/companies/search
	admin.Get("/companies/:id", h.CompanyHandler.GetProfileByID)               // GET /api/v1/admin/companies/:id
	admin.Get("/companies/tax-id/:tax_id", h.CompanyHandler.GetProfileByTaxID) // GET /api/v1/admin/companies/tax-id/:tax_id
	admin.Put("/companies/:id", h.CompanyHandler.AdminUpdateProfile)           // PUT /api/v1/admin/companies/:id
	admin.Patch("/companies/:id/approve", h.CompanyHandler.ApproveProfile)     // PATCH /api/v1/admin/companies/:id/approve

	// Analytics and reporting routes
	admin.Get("/analytics/revenue", h.AnalyticsHandler.GetRevenue)       // GET /api/v1/admin/analytics/revenue
	admin.Get("/commissions", h.AnalyticsHandler.GetCommissions)         // GET /api/v1/admin/commissions
	admin.Get("/stats/basic", h.AnalyticsHandler.GetBasicStats)          // GET /api/v1/admin/stats/basic
	admin.Get("/reports/bookings", h.AnalyticsHandler.GetBookingReports) // GET /api/v1/admin/reports/bookings
}
