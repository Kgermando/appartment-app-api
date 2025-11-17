package dashboard

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/appartment-app-api/database"
	"github.com/kgermando/appartment-app-api/models"
	"gorm.io/gorm"
)

type DashboardStats struct {
	// Statistiques générales
	TotalAppartments      int64 `json:"total_apartments"`
	AvailableApartments   int64 `json:"available_apartments"`
	OccupiedApartments    int64 `json:"occupied_apartments"`
	MaintenanceApartments int64 `json:"maintenance_apartments"`

	// Statistiques financières
	TotalIncomeUSD  float64 `json:"total_income_usd"`
	TotalExpenseUSD float64 `json:"total_expense_usd"`
	NetBalanceUSD   float64 `json:"net_balance_usd"`

	// Statistiques de revenus
	MonthlyRevenueTarget float64 `json:"monthly_revenue_target"`
	ActualMonthlyRevenue float64 `json:"actual_monthly_revenue"`
	RevenuePercentage    float64 `json:"revenue_percentage"`

	// Top appartements
	TopApartmentsByRevenue []ApartmentRevenue `json:"top_apartments_by_revenue"`

	// Statistiques par manager
	ManagerStats []ManagerStats `json:"manager_stats"`
}

type ApartmentRevenue struct {
	UUID         string  `json:"uuid"`
	Name         string  `json:"name"`
	Number       string  `json:"number"`
	MonthlyRent  float64 `json:"monthly_rent"`
	TotalRevenue float64 `json:"total_revenue"`
	Status       string  `json:"status"`
	ManagerName  string  `json:"manager_name"`
}

type ManagerStats struct {
	ManagerUUID          string  `json:"manager_uuid"`
	ManagerName          string  `json:"manager_name"`
	TotalApartments      int64   `json:"total_apartments"`
	AvailableApartments  int64   `json:"available_apartments"`
	OccupiedApartments   int64   `json:"occupied_apartments"`
	TotalIncomeUSD       float64 `json:"total_income_usd"`
	TotalExpenseUSD      float64 `json:"total_expense_usd"`
	NetBalanceUSD        float64 `json:"net_balance_usd"`
	MonthlyRevenueTarget float64 `json:"monthly_revenue_target"`
}

type MonthlyTrend struct {
	Month      string  `json:"month"`
	Year       int     `json:"year"`
	IncomeUSD  float64 `json:"income_usd"`
	ExpenseUSD float64 `json:"expense_usd"`
}

// GetDashboardStats - Dashboard principal avec possibilité de filtrer par manager
func GetDashboardStats(c *fiber.Ctx) error {
	db := database.DB

	// Parse des paramètres de filtrage
	managerUUID := c.Query("manager_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var stats DashboardStats

	// Construire la requête de base pour les appartements
	apartmentQuery := db.Model(&models.Appartment{})
	if managerUUID != "" {
		apartmentQuery = apartmentQuery.Where("manager_uuid = ?", managerUUID)
	}

	// Statistiques des appartements
	apartmentQuery.Count(&stats.TotalAppartments)

	apartmentQuery.Where("status = ?", "available").Count(&stats.AvailableApartments)
	apartmentQuery.Where("status = ?", "occupied").Count(&stats.OccupiedApartments)
	apartmentQuery.Where("status = ?", "maintenance").Count(&stats.MaintenanceApartments)

	// Calculer le revenu mensuel cible
	var monthlyRevenueTarget float64
	db.Model(&models.Appartment{}).
		Select("COALESCE(SUM(monthly_rent), 0)").
		Where("manager_uuid = ? OR ? = ''", managerUUID, managerUUID).
		Row().Scan(&monthlyRevenueTarget)
	stats.MonthlyRevenueTarget = monthlyRevenueTarget

	// Construire la requête pour les caisses avec filtres de date
	caisseQuery := db.Model(&models.Caisse{})

	// Jointure avec appartments pour filtrer par manager si nécessaire
	if managerUUID != "" {
		caisseQuery = caisseQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", managerUUID)
	}

	// Ajouter les filtres de date
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			caisseQuery = caisseQuery.Where("caisses.created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			caisseQuery = caisseQuery.Where("caisses.created_at < ?", endDateTime)
		}
	}

	// Calculer les revenus et dépenses (USD uniquement)
	caisseQuery.Where("caisses.type = ?", "Income").
		Select("COALESCE(SUM(device_usd), 0)").
		Row().Scan(&stats.TotalIncomeUSD)

	caisseQuery.Where("caisses.type = ?", "Expense").
		Select("COALESCE(SUM(device_usd), 0)").
		Row().Scan(&stats.TotalExpenseUSD)

	// Calculer la balance nette
	stats.NetBalanceUSD = stats.TotalIncomeUSD - stats.TotalExpenseUSD

	// Calculer le revenu mensuel actuel (revenus du mois en cours)
	var actualMonthlyRevenue float64

	monthlyQuery := db.Model(&models.Caisse{}).
		Where("type = ? AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)", "Income")

	if managerUUID != "" {
		monthlyQuery = monthlyQuery.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", managerUUID)
	}

	monthlyQuery.Select("COALESCE(SUM(device_usd), 0)").
		Row().Scan(&actualMonthlyRevenue)

	stats.ActualMonthlyRevenue = actualMonthlyRevenue

	// Calculer le pourcentage de réalisation
	if stats.MonthlyRevenueTarget > 0 {
		stats.RevenuePercentage = (stats.ActualMonthlyRevenue / stats.MonthlyRevenueTarget) * 100
	}

	// Top appartements par revenus
	stats.TopApartmentsByRevenue = getTopApartmentsByRevenue(db, managerUUID, startDate, endDate)

	// Statistiques par manager (seulement si pas de filtre manager spécifique)
	if managerUUID == "" {
		stats.ManagerStats = getManagerStats(db, startDate, endDate)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Dashboard stats retrieved successfully",
		"data":    stats,
	})
}

// GetMonthlyTrends - Tendances mensuelles des revenus et dépenses
func GetMonthlyTrends(c *fiber.Ctx) error {
	db := database.DB

	managerUUID := c.Query("manager_uuid", "")
	months := c.Query("months", "12") // Par défaut 12 mois

	monthsInt, err := strconv.Atoi(months)
	if err != nil || monthsInt <= 0 {
		monthsInt = 12
	}

	var trends []MonthlyTrend

	// Calculer les tendances pour les X derniers mois
	for i := monthsInt - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, -i, 0)

		var trend MonthlyTrend
		trend.Month = date.Format("2006-01")
		trend.Year = date.Year()

		// Requête de base pour le mois
		query := db.Model(&models.Caisse{}).
			Where("DATE_TRUNC('month', created_at) = ?", date.Format("2006-01-01"))

		// Ajouter filtre manager si nécessaire
		if managerUUID != "" {
			query = query.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
				Where("appartments.manager_uuid = ?", managerUUID)
		}

		// Revenus du mois
		query.Where("caisses.type = ?", "Income").
			Select("COALESCE(SUM(device_usd), 0)").
			Row().Scan(&trend.IncomeUSD)

		// Dépenses du mois
		query.Where("caisses.type = ?", "Expense").
			Select("COALESCE(SUM(device_usd), 0)").
			Row().Scan(&trend.ExpenseUSD)

		trends = append(trends, trend)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Monthly trends retrieved successfully",
		"data":    trends,
	})
}

// GetManagerComparison - Comparaison entre tous les managers
func GetManagerComparison(c *fiber.Ctx) error {
	db := database.DB

	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	managerStats := getManagerStats(db, startDate, endDate)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Manager comparison retrieved successfully",
		"data":    managerStats,
	})
}

// GetApartmentPerformance - Performance détaillée des appartements
func GetApartmentPerformance(c *fiber.Ctx) error {
	db := database.DB

	managerUUID := c.Query("manager_uuid", "")
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")
	limit := c.Query("limit", "20")

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}

	apartmentPerformance := getTopApartmentsByRevenue(db, managerUUID, startDate, endDate)

	// Limiter les résultats
	if len(apartmentPerformance) > limitInt {
		apartmentPerformance = apartmentPerformance[:limitInt]
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Apartment performance retrieved successfully",
		"data":    apartmentPerformance,
	})
}

// GetFinancialSummary - Résumé financier détaillé
func GetFinancialSummary(c *fiber.Ctx) error {
	db := database.DB

	managerUUID := c.Query("manager_uuid", "")
	period := c.Query("period", "month") // month, quarter, year

	var summary fiber.Map

	switch period {
	case "month":
		summary = getMonthlyFinancialSummary(db, managerUUID)
	case "quarter":
		summary = getQuarterlyFinancialSummary(db, managerUUID)
	case "year":
		summary = getYearlyFinancialSummary(db, managerUUID)
	default:
		summary = getMonthlyFinancialSummary(db, managerUUID)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Financial summary retrieved successfully",
		"data":    summary,
	})
}

// GetOccupancyStats - Statistiques d'occupation
func GetOccupancyStats(c *fiber.Ctx) error {
	db := database.DB

	managerUUID := c.Query("manager_uuid", "")

	type OccupancyStats struct {
		TotalApartments       int64   `json:"total_apartments"`
		OccupiedApartments    int64   `json:"occupied_apartments"`
		AvailableApartments   int64   `json:"available_apartments"`
		MaintenanceApartments int64   `json:"maintenance_apartments"`
		OccupancyRate         float64 `json:"occupancy_rate"`
		AvailabilityRate      float64 `json:"availability_rate"`
		AverageRent           float64 `json:"average_rent"`
		TotalPotentialRevenue float64 `json:"total_potential_revenue"`
		LostRevenue           float64 `json:"lost_revenue"`
	}

	var stats OccupancyStats

	// Requête de base
	query := db.Model(&models.Appartment{})
	if managerUUID != "" {
		query = query.Where("manager_uuid = ?", managerUUID)
	}

	// Statistiques de base
	query.Count(&stats.TotalApartments)
	query.Where("status = ?", "occupied").Count(&stats.OccupiedApartments)
	query.Where("status = ?", "available").Count(&stats.AvailableApartments)
	query.Where("status = ?", "maintenance").Count(&stats.MaintenanceApartments)

	// Calculs des taux
	if stats.TotalApartments > 0 {
		stats.OccupancyRate = (float64(stats.OccupiedApartments) / float64(stats.TotalApartments)) * 100
		stats.AvailabilityRate = (float64(stats.AvailableApartments) / float64(stats.TotalApartments)) * 100
	}

	// Loyer moyen
	query.Select("COALESCE(AVG(monthly_rent), 0)").Row().Scan(&stats.AverageRent)

	// Revenus potentiels totaux
	query.Select("COALESCE(SUM(monthly_rent), 0)").Row().Scan(&stats.TotalPotentialRevenue)

	// Revenus perdus (appartements non occupés)
	lostRevenueQuery := db.Model(&models.Appartment{}).
		Where("status != ?", "occupied")
	if managerUUID != "" {
		lostRevenueQuery = lostRevenueQuery.Where("manager_uuid = ?", managerUUID)
	}
	lostRevenueQuery.Select("COALESCE(SUM(monthly_rent), 0)").Row().Scan(&stats.LostRevenue)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Occupancy stats retrieved successfully",
		"data":    stats,
	})
}

// GetTopManagers - Classement des meilleurs managers
func GetTopManagers(c *fiber.Ctx) error {
	db := database.DB

	period := c.Query("period", "month") // month, quarter, year
	limit := c.Query("limit", "10")

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	var dateFilter string
	switch period {
	case "month":
		dateFilter = "DATE_TRUNC('month', c.created_at) = DATE_TRUNC('month', CURRENT_DATE)"
	case "quarter":
		dateFilter = "DATE_TRUNC('quarter', c.created_at) = DATE_TRUNC('quarter', CURRENT_DATE)"
	case "year":
		dateFilter = "DATE_TRUNC('year', c.created_at) = DATE_TRUNC('year', CURRENT_DATE)"
	default:
		dateFilter = "DATE_TRUNC('month', c.created_at) = DATE_TRUNC('month', CURRENT_DATE)"
	}

	type TopManager struct {
		ManagerUUID    string  `json:"manager_uuid"`
		ManagerName    string  `json:"manager_name"`
		TotalRevenue   float64 `json:"total_revenue"`
		NetProfit      float64 `json:"net_profit"`
		ApartmentCount int64   `json:"apartment_count"`
		OccupancyRate  float64 `json:"occupancy_rate"`
		Efficiency     float64 `json:"efficiency"` // Net profit / Total apartments
	}

	var results []TopManager

	query := `
        SELECT 
            u.uuid as manager_uuid,
            u.fullname as manager_name,
            COALESCE(SUM(CASE WHEN c.type = 'Income' THEN c.device_usd ELSE 0 END), 0) as total_revenue,
            COALESCE(SUM(CASE WHEN c.type = 'Income' THEN c.device_usd ELSE 0 END) - 
                     SUM(CASE WHEN c.type = 'Expense' THEN c.device_usd ELSE 0 END), 0) as net_profit,
            COUNT(DISTINCT a.uuid) as apartment_count,
            CASE 
                WHEN COUNT(DISTINCT a.uuid) > 0 THEN 
                    (COUNT(DISTINCT CASE WHEN a.status = 'occupied' THEN a.uuid END)::float / COUNT(DISTINCT a.uuid)::float) * 100
                ELSE 0 
            END as occupancy_rate
        FROM users u
        LEFT JOIN appartments a ON u.uuid = a.manager_uuid
        LEFT JOIN caisses c ON a.uuid = c.appartment_uuid AND ` + dateFilter + `
        WHERE u.role IN ('Manager', 'Agent', 'Supervisor')
        GROUP BY u.uuid, u.fullname
        HAVING COUNT(DISTINCT a.uuid) > 0
        ORDER BY net_profit DESC
        LIMIT ?
    `

	rows, err := db.Raw(query, limitInt).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch top managers",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var result TopManager
		rows.Scan(
			&result.ManagerUUID,
			&result.ManagerName,
			&result.TotalRevenue,
			&result.NetProfit,
			&result.ApartmentCount,
			&result.OccupancyRate,
		)

		// Calculer l'efficacité
		if result.ApartmentCount > 0 {
			result.Efficiency = result.NetProfit / float64(result.ApartmentCount)
		}

		results = append(results, result)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Top managers retrieved successfully",
		"data":    results,
	})
}

// Fonctions utilitaires

func getMonthlyFinancialSummary(db *gorm.DB, managerUUID string) fiber.Map {
	currentMonth := time.Now().Format("2006-01")

	var incomeUSD, expenseUSD float64

	// Requête de base pour le mois en cours
	query := db.Model(&models.Caisse{}).
		Where("DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)")

	if managerUUID != "" {
		query = query.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", managerUUID)
	}

	// Revenus du mois
	query.Where("caisses.type = ?", "Income").
		Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&incomeUSD)

	// Dépenses du mois
	query.Where("caisses.type = ?", "Expense").
		Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&expenseUSD)

	return fiber.Map{
		"period":          "month",
		"month":           currentMonth,
		"income_usd":      incomeUSD,
		"expense_usd":     expenseUSD,
		"net_balance_usd": incomeUSD - expenseUSD,
	}
}

func getQuarterlyFinancialSummary(db *gorm.DB, managerUUID string) fiber.Map {
	currentQuarter := time.Now().Format("2006-Q1") // Approximation

	var incomeUSD, expenseUSD float64

	query := db.Model(&models.Caisse{}).
		Where("DATE_TRUNC('quarter', created_at) = DATE_TRUNC('quarter', CURRENT_DATE)")

	if managerUUID != "" {
		query = query.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", managerUUID)
	}

	query.Where("caisses.type = ?", "Income").
		Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&incomeUSD)
	query.Where("caisses.type = ?", "Expense").
		Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&expenseUSD)

	return fiber.Map{
		"period":          "quarter",
		"quarter":         currentQuarter,
		"income_usd":      incomeUSD,
		"expense_usd":     expenseUSD,
		"net_balance_usd": incomeUSD - expenseUSD,
	}
}

func getYearlyFinancialSummary(db *gorm.DB, managerUUID string) fiber.Map {
	currentYear := time.Now().Year()

	var incomeUSD, expenseUSD float64

	query := db.Model(&models.Caisse{}).
		Where("DATE_TRUNC('year', created_at) = DATE_TRUNC('year', CURRENT_DATE)")

	if managerUUID != "" {
		query = query.Joins("JOIN appartments ON caisses.appartment_uuid = appartments.uuid").
			Where("appartments.manager_uuid = ?", managerUUID)
	}

	query.Where("caisses.type = ?", "Income").
		Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&incomeUSD)
	query.Where("caisses.type = ?", "Expense").
		Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&expenseUSD)

	return fiber.Map{
		"period":          "year",
		"year":            currentYear,
		"income_usd":      incomeUSD,
		"expense_usd":     expenseUSD,
		"net_balance_usd": incomeUSD - expenseUSD,
	}
}

func getTopApartmentsByRevenue(db *gorm.DB, managerUUID, startDate, endDate string) []ApartmentRevenue {
	var results []ApartmentRevenue

	query := `
        SELECT 
            a.uuid,
            a.name,
            a.number,
            a.monthly_rent,
            COALESCE(SUM(c.device_usd), 0) as total_revenue,
            a.status,
            u.fullname as manager_name
        FROM appartments a
        LEFT JOIN caisses c ON a.uuid = c.appartment_uuid AND c.type = 'Income'
        LEFT JOIN users u ON a.manager_uuid = u.uuid
        WHERE 1=1
    `

	args := []interface{}{}

	if managerUUID != "" {
		query += " AND a.manager_uuid = ?"
		args = append(args, managerUUID)
	}

	if startDate != "" {
		query += " AND c.created_at >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND c.created_at < ?"
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			args = append(args, endDateTime)
		}
	}

	query += `
        GROUP BY a.uuid, a.name, a.number, a.monthly_rent, a.status, u.fullname
        ORDER BY total_revenue DESC
        LIMIT 10
    `

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		var result ApartmentRevenue
		rows.Scan(
			&result.UUID,
			&result.Name,
			&result.Number,
			&result.MonthlyRent,
			&result.TotalRevenue,
			&result.Status,
			&result.ManagerName,
		)
		results = append(results, result)
	}

	return results
}

func getManagerStats(db *gorm.DB, startDate, endDate string) []ManagerStats {
	var results []ManagerStats

	query := `
        SELECT 
            u.uuid as manager_uuid,
            u.fullname as manager_name,
            COUNT(DISTINCT a.uuid) as total_apartments,
            COUNT(DISTINCT CASE WHEN a.status = 'available' THEN a.uuid END) as available_apartments,
            COUNT(DISTINCT CASE WHEN a.status = 'occupied' THEN a.uuid END) as occupied_apartments,
            COALESCE(SUM(CASE WHEN c.type = 'Income' THEN c.device_usd ELSE 0 END), 0) as total_income_usd,
            COALESCE(SUM(CASE WHEN c.type = 'Expense' THEN c.device_usd ELSE 0 END), 0) as total_expense_usd,
            COALESCE(SUM(a.monthly_rent), 0) as monthly_revenue_target
        FROM users u
        LEFT JOIN appartments a ON u.uuid = a.manager_uuid
        LEFT JOIN caisses c ON a.uuid = c.appartment_uuid
        WHERE u.role IN ('Manager', 'Agent', 'Supervisor')
    `

	args := []interface{}{}

	if startDate != "" {
		query += " AND (c.created_at IS NULL OR c.created_at >= ?)"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND (c.created_at IS NULL OR c.created_at < ?)"
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			args = append(args, endDateTime)
		}
	}

	query += `
        GROUP BY u.uuid, u.fullname
        HAVING COUNT(DISTINCT a.uuid) > 0
        ORDER BY COALESCE(SUM(CASE WHEN c.type = 'Income' THEN c.device_usd ELSE 0 END), 0) DESC
    `

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		var result ManagerStats
		rows.Scan(
			&result.ManagerUUID,
			&result.ManagerName,
			&result.TotalApartments,
			&result.AvailableApartments,
			&result.OccupiedApartments,
			&result.TotalIncomeUSD,
			&result.TotalExpenseUSD,
			&result.MonthlyRevenueTarget,
		)

		// Calculer la balance nette
		result.NetBalanceUSD = result.TotalIncomeUSD - result.TotalExpenseUSD

		results = append(results, result)
	}

	return results
}
