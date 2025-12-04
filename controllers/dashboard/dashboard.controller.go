package dashboard

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/appartment-app-api/database"
	"github.com/kgermando/appartment-app-api/models"
)

func GetDashboardStats(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters (tous optionnels)
	userUUID := c.Query("user_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var stats models.DashboardStats

	// 1. Statistiques générales des appartements
	apartmentBaseQuery := db.Model(&models.Appartment{})
	if userUUID != "" {
		apartmentBaseQuery = apartmentBaseQuery.Where("manager_uuid = ?", userUUID)
	}

	// Clone queries for each status
	apartmentBaseQuery.Count(&stats.TotalAppartments)

	availableQuery := db.Model(&models.Appartment{}).Where("status = ?", "available")
	if userUUID != "" {
		availableQuery = availableQuery.Where("manager_uuid = ?", userUUID)
	}
	availableQuery.Count(&stats.AvailableApartments)

	occupiedQuery := db.Model(&models.Appartment{}).Where("status = ?", "occupied")
	if userUUID != "" {
		occupiedQuery = occupiedQuery.Where("manager_uuid = ?", userUUID)
	}
	occupiedQuery.Count(&stats.OccupiedApartments)

	maintenanceQuery := db.Model(&models.Appartment{}).Where("status = ?", "maintenance")
	if userUUID != "" {
		maintenanceQuery = maintenanceQuery.Where("manager_uuid = ?", userUUID)
	}
	maintenanceQuery.Count(&stats.MaintenanceApartments)

	// 2. Statistiques financières en USD et CDF avec filtres optionnels
	var totalIncomeUSD, totalExpenseUSD, totalIncomeCDF, totalExpenseCDF float64

	// Build income query (filtre par user_uuid et dates si fournis)
	incomeQuery := db.Model(&models.Caisse{}).Where("type = ?", "Income")
	if userUUID != "" {
		incomeQuery = incomeQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", userUUID)
	}
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			incomeQuery = incomeQuery.Where("caisses.created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			incomeQuery = incomeQuery.Where("caisses.created_at < ?", endDateTime)
		}
	}
	incomeQuery.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalIncomeUSD)

	// Income CDF
	incomeCDFQuery := db.Model(&models.Caisse{}).Where("type = ?", "Income")
	if userUUID != "" {
		incomeCDFQuery = incomeCDFQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", userUUID)
	}
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			incomeCDFQuery = incomeCDFQuery.Where("caisses.created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			incomeCDFQuery = incomeCDFQuery.Where("caisses.created_at < ?", endDateTime)
		}
	}
	incomeCDFQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalIncomeCDF)

	// Build expense query (filtre par user_uuid et dates si fournis)
	expenseQuery := db.Model(&models.Caisse{}).Where("type = ?", "Expense")
	if userUUID != "" {
		expenseQuery = expenseQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", userUUID)
	}
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			expenseQuery = expenseQuery.Where("caisses.created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			expenseQuery = expenseQuery.Where("caisses.created_at < ?", endDateTime)
		}
	}
	expenseQuery.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalExpenseUSD)

	// Expense CDF
	expenseCDFQuery := db.Model(&models.Caisse{}).Where("type = ?", "Expense")
	if userUUID != "" {
		expenseCDFQuery = expenseCDFQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", userUUID)
	}
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			expenseCDFQuery = expenseCDFQuery.Where("caisses.created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			expenseCDFQuery = expenseCDFQuery.Where("caisses.created_at < ?", endDateTime)
		}
	}
	expenseCDFQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalExpenseCDF)

	stats.TotalIncomeUSD = totalIncomeUSD
	stats.TotalIncomeCDF = totalIncomeCDF
	stats.TotalExpenseUSD = totalExpenseUSD
	stats.TotalExpenseCDF = totalExpenseCDF

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Dashboard statistics retrieved successfully",
		"data":    stats,
	})
}

// GetApartmentRevenues returns revenue statistics for each apartment
func GetApartmentRevenues(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters
	userUUID := c.Query("user_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	// Build apartment query
	apartmentQuery := db.Preload("Manager")
	if userUUID != "" {
		apartmentQuery = apartmentQuery.Where("manager_uuid = ?", userUUID)
	}

	var apartments []models.Appartment
	apartmentQuery.Find(&apartments)

	var revenues []models.ApartmentRevenue

	for _, apt := range apartments {
		var totalIncomeUSD, totalIncomeCDF, totalExpenseUSD, totalExpenseCDF float64

		// Build income query with date filters
		incomeQuery := db.Model(&models.Caisse{}).
			Where("appartment_uuid = ? AND type = ?", apt.UUID, "Income")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				incomeQuery = incomeQuery.Where("created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				incomeQuery = incomeQuery.Where("created_at < ?", endDateTime)
			}
		}

		incomeQuery.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalIncomeUSD)

		// Income CDF
		incomeCDFQuery := db.Model(&models.Caisse{}).
			Where("appartment_uuid = ? AND type = ?", apt.UUID, "Income")
		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				incomeCDFQuery = incomeCDFQuery.Where("created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				incomeCDFQuery = incomeCDFQuery.Where("created_at < ?", endDateTime)
			}
		}
		incomeCDFQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalIncomeCDF)

		// Build expense query with date filters
		expenseQuery := db.Model(&models.Caisse{}).
			Where("appartment_uuid = ? AND type = ?", apt.UUID, "Expense")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				expenseQuery = expenseQuery.Where("created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				expenseQuery = expenseQuery.Where("created_at < ?", endDateTime)
			}
		}

		expenseQuery.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalExpenseUSD)

		// Expense CDF
		expenseCDFQuery := db.Model(&models.Caisse{}).
			Where("appartment_uuid = ? AND type = ?", apt.UUID, "Expense")
		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				expenseCDFQuery = expenseCDFQuery.Where("created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				expenseCDFQuery = expenseCDFQuery.Where("created_at < ?", endDateTime)
			}
		}
		expenseCDFQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalExpenseCDF)

		revenue := models.ApartmentRevenue{
			UUID:            apt.UUID,
			Name:            apt.Name,
			Number:          apt.Number,
			MonthlyRent:     apt.MonthlyRent,
			TotalIncomeUSD:  totalIncomeUSD,
			TotalIncomeCDF:  totalIncomeCDF,
			TotalExpenseUSD: totalExpenseUSD,
			TotalExpenseCDF: totalExpenseCDF,
			Status:          apt.Status,
			ManagerName:     apt.Manager.Fullname,
		}

		revenues = append(revenues, revenue)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Apartment revenues retrieved successfully",
		"data":    revenues,
	})
}

// GetManagerStats returns statistics grouped by manager
func GetManagerStats(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters
	userUUID := c.Query("user_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	// Build manager query
	managerQuery := db.Where("role IN ?", []string{"Manager"})
	if userUUID != "" {
		managerQuery = managerQuery.Where("uuid = ?", userUUID)
	}

	var managers []models.User
	managerQuery.Find(&managers)

	var managerStats []models.ManagerStats

	for _, manager := range managers {
		var stats models.ManagerStats
		stats.ManagerUUID = manager.UUID
		stats.ManagerName = manager.Fullname

		// Count apartments by status for this manager
		db.Model(&models.Appartment{}).
			Where("manager_uuid = ?", manager.UUID).
			Count(&stats.TotalApartments)

		db.Model(&models.Appartment{}).
			Where("manager_uuid = ? AND status = ?", manager.UUID, "available").
			Count(&stats.AvailableApartments)

		db.Model(&models.Appartment{}).
			Where("manager_uuid = ? AND status = ?", manager.UUID, "occupied").
			Count(&stats.OccupiedApartments)

		// Calculate financial stats for this manager with date filters
		var totalIncomeUSD, totalIncomeCDF, totalExpenseUSD, totalExpenseCDF float64

		// Income USD query
		incomeQuery := db.Table("caisses").
			Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ? AND caisses.type = ?", manager.UUID, "Income")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				incomeQuery = incomeQuery.Where("caisses.created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				incomeQuery = incomeQuery.Where("caisses.created_at < ?", endDateTime)
			}
		}
		incomeQuery.Select("COALESCE(SUM(caisses.device_usd), 0)").Row().Scan(&totalIncomeUSD)

		// Income CDF query
		incomeCDFQuery := db.Table("caisses").
			Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ? AND caisses.type = ?", manager.UUID, "Income")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				incomeCDFQuery = incomeCDFQuery.Where("caisses.created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				incomeCDFQuery = incomeCDFQuery.Where("caisses.created_at < ?", endDateTime)
			}
		}
		incomeCDFQuery.Select("COALESCE(SUM(caisses.device_cdf), 0)").Row().Scan(&totalIncomeCDF)

		// Expense USD query
		expenseQuery := db.Table("caisses").
			Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ? AND caisses.type = ?", manager.UUID, "Expense")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				expenseQuery = expenseQuery.Where("caisses.created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				expenseQuery = expenseQuery.Where("caisses.created_at < ?", endDateTime)
			}
		}
		expenseQuery.Select("COALESCE(SUM(caisses.device_usd), 0)").Row().Scan(&totalExpenseUSD)

		// Expense CDF query
		expenseCDFQuery := db.Table("caisses").
			Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ? AND caisses.type = ?", manager.UUID, "Expense")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				expenseCDFQuery = expenseCDFQuery.Where("caisses.created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				expenseCDFQuery = expenseCDFQuery.Where("caisses.created_at < ?", endDateTime)
			}
		}
		expenseCDFQuery.Select("COALESCE(SUM(caisses.device_cdf), 0)").Row().Scan(&totalExpenseCDF)

		stats.TotalIncomeUSD = totalIncomeUSD
		stats.TotalIncomeCDF = totalIncomeCDF
		stats.TotalExpenseUSD = totalExpenseUSD
		stats.TotalExpenseCDF = totalExpenseCDF

		managerStats = append(managerStats, stats)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Manager statistics retrieved successfully",
		"data":    managerStats,
	})
}

// GetMonthlyTrends returns income and expense trends by month
func GetMonthlyTrends(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters
	userUUID := c.Query("user_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	// Get year from query parameter, default to current year
	year := time.Now().Year()
	if yearParam := c.Query("year"); yearParam != "" {
		if parsedYear, err := strconv.Atoi(yearParam); err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	var trends []models.MonthlyTrend

	// Loop through each month
	for month := 1; month <= 12; month++ {
		monthStartDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		monthEndDate := monthStartDate.AddDate(0, 1, 0).Add(-time.Second)

		// Apply custom date range if provided
		actualStartDate := monthStartDate
		actualEndDate := monthEndDate

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				if parsedStartDate.After(monthStartDate) {
					actualStartDate = parsedStartDate
				}
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				parsedEndDateTime := parsedEndDate.Add(24 * time.Hour)
				if parsedEndDateTime.Before(monthEndDate) {
					actualEndDate = parsedEndDateTime
				}
			}
		}

		// Skip month if it's outside the date range
		if actualStartDate.After(actualEndDate) {
			continue
		}

		var incomeUSD, expenseUSD float64

		// Income query
		incomeQuery := db.Model(&models.Caisse{}).
			Where("type = ? AND created_at >= ? AND created_at <= ?", "Income", actualStartDate, actualEndDate)

		if userUUID != "" {
			incomeQuery = incomeQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
				Where("appartments.manager_uuid = ?", userUUID)
		}
		incomeQuery.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&incomeUSD)

		// Expense query
		expenseQuery := db.Model(&models.Caisse{}).
			Where("type = ? AND created_at >= ? AND created_at <= ?", "Expense", actualStartDate, actualEndDate)

		if userUUID != "" {
			expenseQuery = expenseQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
				Where("appartments.manager_uuid = ?", userUUID)
		}
		expenseQuery.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&expenseUSD)

		trend := models.MonthlyTrend{
			Month:      time.Month(month).String(),
			Year:       year,
			IncomeUSD:  incomeUSD,
			ExpenseUSD: expenseUSD,
		}

		trends = append(trends, trend)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Monthly trends retrieved successfully",
		"data":    trends,
	})
}

// GetOccupancyStats returns detailed occupancy statistics
func GetOccupancyStats(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters
	userUUID := c.Query("user_uuid", "")

	var stats models.OccupancyStats

	// Build base query with user filter
	baseQuery := db.Model(&models.Appartment{})
	if userUUID != "" {
		baseQuery = baseQuery.Where("manager_uuid = ?", userUUID)
	}

	// Count apartments by status
	baseQuery.Count(&stats.TotalApartments)

	occupiedQuery := db.Model(&models.Appartment{}).Where("status = ?", "occupied")
	if userUUID != "" {
		occupiedQuery = occupiedQuery.Where("manager_uuid = ?", userUUID)
	}
	occupiedQuery.Count(&stats.OccupiedApartments)

	availableQuery := db.Model(&models.Appartment{}).Where("status = ?", "available")
	if userUUID != "" {
		availableQuery = availableQuery.Where("manager_uuid = ?", userUUID)
	}
	availableQuery.Count(&stats.AvailableApartments)

	maintenanceQuery := db.Model(&models.Appartment{}).Where("status = ?", "maintenance")
	if userUUID != "" {
		maintenanceQuery = maintenanceQuery.Where("manager_uuid = ?", userUUID)
	}
	maintenanceQuery.Count(&stats.MaintenanceApartments)

	// Calculate occupancy and availability rates
	if stats.TotalApartments > 0 {
		stats.OccupancyRate = (float64(stats.OccupiedApartments) / float64(stats.TotalApartments)) * 100
		stats.AvailabilityRate = (float64(stats.AvailableApartments) / float64(stats.TotalApartments)) * 100
	}

	// Calculate average rent
	avgRentQuery := db.Model(&models.Appartment{})
	if userUUID != "" {
		avgRentQuery = avgRentQuery.Where("manager_uuid = ?", userUUID)
	}
	avgRentQuery.Select("COALESCE(AVG(monthly_rent), 0)").Row().Scan(&stats.AverageRent)

	// Calculate total potential revenue (all apartments)
	potentialRevenueQuery := db.Model(&models.Appartment{})
	if userUUID != "" {
		potentialRevenueQuery = potentialRevenueQuery.Where("manager_uuid = ?", userUUID)
	}
	potentialRevenueQuery.Select("COALESCE(SUM(monthly_rent), 0)").Row().Scan(&stats.TotalPotentialRevenue)

	// Calculate lost revenue from vacant apartments
	lostRevenueQuery := db.Model(&models.Appartment{}).Where("status != ?", "occupied")
	if userUUID != "" {
		lostRevenueQuery = lostRevenueQuery.Where("manager_uuid = ?", userUUID)
	}
	lostRevenueQuery.Select("COALESCE(SUM(monthly_rent), 0)").Row().Scan(&stats.LostRevenue)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Occupancy statistics retrieved successfully",
		"data":    stats,
	})
}

// GetTopManagers returns the top performing managers
func GetTopManagers(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters
	userUUID := c.Query("user_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	// Build manager query
	managerQuery := db.Where("role IN ?", []string{"Manager"})
	if userUUID != "" {
		managerQuery = managerQuery.Where("uuid = ?", userUUID)
	}

	var managers []models.User
	managerQuery.Find(&managers)

	var topManagers []models.TopManager

	for _, manager := range managers {
		var topMgr models.TopManager
		topMgr.ManagerUUID = manager.UUID
		topMgr.ManagerName = manager.Fullname

		// Count apartments
		db.Model(&models.Appartment{}).
			Where("manager_uuid = ?", manager.UUID).
			Count(&topMgr.ApartmentCount)

		// Calculate financial stats with date filters
		var totalRevenue, totalExpense float64
		var occupiedCount int64

		// Revenue query with date filters
		revenueQuery := db.Table("caisses").
			Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ? AND caisses.type = ?", manager.UUID, "Income")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				revenueQuery = revenueQuery.Where("caisses.created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				revenueQuery = revenueQuery.Where("caisses.created_at < ?", endDateTime)
			}
		}
		revenueQuery.Select("COALESCE(SUM(caisses.device_usd), 0)").Row().Scan(&totalRevenue)

		// Expense query with date filters
		expenseQuery := db.Table("caisses").
			Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ? AND caisses.type = ?", manager.UUID, "Expense")

		if startDate != "" {
			if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
				expenseQuery = expenseQuery.Where("caisses.created_at >= ?", parsedStartDate)
			}
		}
		if endDate != "" {
			if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
				endDateTime := parsedEndDate.Add(24 * time.Hour)
				expenseQuery = expenseQuery.Where("caisses.created_at < ?", endDateTime)
			}
		}
		expenseQuery.Select("COALESCE(SUM(caisses.device_usd), 0)").Row().Scan(&totalExpense)

		db.Model(&models.Appartment{}).
			Where("manager_uuid = ? AND status = ?", manager.UUID, "occupied").
			Count(&occupiedCount)

		topMgr.TotalRevenue = totalRevenue
		topMgr.NetProfit = totalRevenue - totalExpense

		// Calculate occupancy rate
		if topMgr.ApartmentCount > 0 {
			topMgr.OccupancyRate = (float64(occupiedCount) / float64(topMgr.ApartmentCount)) * 100
		}

		// Calculate efficiency (net profit per apartment)
		if topMgr.ApartmentCount > 0 {
			topMgr.Efficiency = topMgr.NetProfit / float64(topMgr.ApartmentCount)
		}

		topManagers = append(topManagers, topMgr)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Top managers retrieved successfully",
		"data":    topManagers,
	})
}

// Get appartment payment statistics by month
func GetAppartmentStats(c *fiber.Ctx) error {

	db := database.DB

	userUUID := c.Query("user_uuid", "")

	// Obtenir l'année courante ou depuis les paramètres de requête
	year := time.Now().Year()
	if yearParam := c.Query("year"); yearParam != "" {
		if parsedYear, err := strconv.Atoi(yearParam); err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Récupérer les appartements avec filtre optionnel
	var appartments []models.Appartment
	appartmentQuery := db.Model(&models.Appartment{})
	if userUUID != "" {
		appartmentQuery = appartmentQuery.Where("manager_uuid = ?", userUUID)
	}

	if err := appartmentQuery.Find(&appartments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch appartments",
			"error":   err.Error(),
		})
	}

	if len(appartments) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No appartments found",
			"data":    nil,
		})
	}

	// Initialiser les statistiques pour les 12 mois
	monthlyStats := make(map[string]map[string]float64)
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	for _, month := range months {
		monthlyStats[month] = map[string]float64{
			"income_cdf":  0.0,
			"income_usd":  0.0,
			"expense_cdf": 0.0,
			"expense_usd": 0.0,
		}
	}

	// Récupérer toutes les entrées (Income et Expense) de la caisse pour ces appartements
	var caisses []models.Caisse
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	caisseQuery := db.Where("created_at >= ? AND created_at <= ?", startDate, endDate)

	// Filtrer par appartements si nécessaire
	if userUUID != "" {
		var appartmentUUIDs []string
		for _, apt := range appartments {
			appartmentUUIDs = append(appartmentUUIDs, apt.UUID)
		}
		caisseQuery = caisseQuery.Where("appartment_uuid IN ?", appartmentUUIDs)
	}

	err := caisseQuery.Find(&caisses).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch payment statistics",
			"error":   err.Error(),
		})
	}

	// Calculer les totaux par mois
	for _, caisse := range caisses {
		monthIndex := int(caisse.CreatedAt.Month()) - 1
		if monthIndex >= 0 && monthIndex < 12 {
			monthName := months[monthIndex]
			switch caisse.Type {
			case "Income":
				monthlyStats[monthName]["income_cdf"] += caisse.DeviceCDF
				monthlyStats[monthName]["income_usd"] += caisse.DeviceUSD
			case "Expense":
				monthlyStats[monthName]["expense_cdf"] += caisse.DeviceCDF
				monthlyStats[monthName]["expense_usd"] += caisse.DeviceUSD
			}
		}
	}

	// Calculer les totaux mensuels et annuels
	var totalYearIncomeCDF, totalYearIncomeUSD, totalYearExpenseCDF, totalYearExpenseUSD float64
	for _, month := range months {
		// Calculer les totaux pour chaque mois
		monthlyStats[month]["total_income_cdf"] = monthlyStats[month]["income_cdf"]
		monthlyStats[month]["total_income_usd"] = monthlyStats[month]["income_usd"]
		monthlyStats[month]["total_expense_cdf"] = monthlyStats[month]["expense_cdf"]
		monthlyStats[month]["total_expense_usd"] = monthlyStats[month]["expense_usd"]

		// Calculer les totaux annuels
		totalYearIncomeCDF += monthlyStats[month]["income_cdf"]
		totalYearIncomeUSD += monthlyStats[month]["income_usd"]
		totalYearExpenseCDF += monthlyStats[month]["expense_cdf"]
		totalYearExpenseUSD += monthlyStats[month]["expense_usd"]
	}

	// Préparer la réponse
	response := map[string]interface{}{
		"year":          year,
		"monthly_stats": monthlyStats,
		"yearly_totals": map[string]float64{
			"total_income_cdf":  totalYearIncomeCDF,
			"total_income_usd":  totalYearIncomeUSD,
			"total_expense_cdf": totalYearExpenseCDF,
			"total_expense_usd": totalYearExpenseUSD,
		},
		"appartment_count": len(appartments),
		"filter_applied":   userUUID != "",
		"currency_info": map[string]string{
			"cdf": "Francs Congolais",
			"usd": "US Dollars",
		},
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Appartment payment statistics retrieved successfully",
		"data":    response,
	})
}
