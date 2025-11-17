package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/appartment-app-api/controllers/appartments"
	"github.com/kgermando/appartment-app-api/controllers/auth"
	"github.com/kgermando/appartment-app-api/controllers/caisses"
	"github.com/kgermando/appartment-app-api/controllers/dashboard"
	"github.com/kgermando/appartment-app-api/controllers/users"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {

	api := app.Group("/api", logger.New())

	// Authentification controller
	a := api.Group("/auth")
	a.Post("/register", auth.Register)
	a.Post("/login", auth.Login)
	a.Post("/create-admin", auth.CreateAdminUser) // Nouveau endpoint pour créer un admin
	a.Post("/forgot-password", auth.ForgotPassword)
	a.Post("/reset/:token", auth.ResetPassword)

	// app.Use(middlewares.IsAuthenticated)

	a.Get("/user", auth.AuthUser)
	a.Put("/profil/info", auth.UpdateInfo)
	a.Put("/change-password", auth.ChangePassword)
	a.Post("/logout", auth.Logout)

	// Users controller
	u := api.Group("/users")
	u.Get("/all/paginate", users.GetPaginatedUsers) // Route statique en premier
	u.Get("/all/:uuid", users.GetAllUsersByUUID)    // Route dynamique après
	u.Get("/all", users.GetAllUsers)
	u.Get("/get/:uuid", users.GetUser)
	u.Post("/create", users.CreateUser)
	u.Put("/update/:uuid", users.UpdateUser)
	u.Delete("/delete/:uuid", users.DeleteUser)

	// Appartments controller
	ap := api.Group("/appartments")
	ap.Get("/all/paginate", appartments.GetPaginatedAppartmentsManagerGeneral) // Route statique en premier
	ap.Get("/all/:manager_uuid/paginate", appartments.GetPaginatedAppartments) // Route avec paramètre + suffixe
	ap.Get("/all/:manager_uuid", appartments.GetAllAppartmentsByManagerUUID)   // Route dynamique seule
	ap.Get("/all", appartments.GetAllAppartments)
	ap.Get("/stats/:uuid", appartments.GetAppartmentStats) // Route statique "stats" avant "get"
	ap.Get("/get/:uuid", appartments.GetAppartment)
	ap.Post("/create", appartments.CreateAppartment)
	ap.Put("/update/:uuid", appartments.UpdateAppartment)
	ap.Delete("/delete/:uuid", appartments.DeleteAppartment)

	// Caisses controller
	c := api.Group("/caisses")
	c.Get("/all/paginate", caisses.GetPaginatedCaissesManagerGeneral)     // Route statique en premier
	c.Get("/all/:appartment_uuid/paginate", caisses.GetPaginatedCaisses)  // Route avec paramètre + suffixe
	c.Get("/all/:appartment_uuid", caisses.GetAllCaissesByAppartmentUUID) // Route dynamique seule
	c.Get("/all", caisses.GetAllCaisses)
	c.Get("/get/:uuid", caisses.GetCaisse)
	c.Post("/create", caisses.CreateCaisse)
	c.Put("/update/:uuid", caisses.UpdateCaisse)
	c.Delete("/delete/:uuid", caisses.DeleteCaisse)

	// Financial endpoints
	c.Get("/totals/global", caisses.GetGlobalTotals)                   // Route statique en premier
	c.Get("/totals/manager/:manager_uuid", caisses.GetTotalsByManager) // Route avec paramètre après
	c.Get("/balance/:appartment_uuid", caisses.GetAppartmentBalance)
	c.Post("/convert", caisses.ConvertCurrency)

	// Dashboard controller
	d := api.Group("/dashboard")
	d.Get("/stats", dashboard.GetDashboardStats)                        // Dashboard principal avec filtres
	d.Get("/trends", dashboard.GetMonthlyTrends)                        // Tendances mensuelles
	d.Get("/managers", dashboard.GetManagerComparison)                  // Comparaison entre managers
	d.Get("/apartments/performance", dashboard.GetApartmentPerformance) // Performance des appartements
	d.Get("/financial", dashboard.GetFinancialSummary)                  // Résumé financier détaillé
	d.Get("/occupancy", dashboard.GetOccupancyStats)                    // Statistiques d'occupation
	d.Get("/top-managers", dashboard.GetTopManagers)                    // Classement des meilleurs managers

}
