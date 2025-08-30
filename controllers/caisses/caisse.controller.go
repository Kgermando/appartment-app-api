package caisses

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/appartment-app-api/database"
	"github.com/kgermando/appartment-app-api/models"
	"github.com/kgermando/appartment-app-api/utils"
)

// getCurrentExchangeRate utilise les taux par défaut définis dans utils
func getCurrentExchangeRate(fromCurrency, toCurrency string) float64 {
	return utils.GetDefaultExchangeRate(fromCurrency, toCurrency)
}

// Paginate
func GetPaginatedCaisses(c *fiber.Ctx) error {
	db := database.DB

	appartmentUUID := c.Params("appartment_uuid")

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	// Parse search query
	search := c.Query("search", "")

	// Parse date range filters
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var caisses []models.Caisse
	var totalRecords int64

	// Build query with date filter
	query := db.Model(&models.Caisse{}).
		Where("appartment_uuid = ?", appartmentUUID).
		Where("type ILIKE ? OR signature ILIKE ?", "%"+search+"%", "%"+search+"%")

	// Add date range filter if provided
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add 24 hours to include the entire end date
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endDateTime)
		}
	}

	// Count total records matching the search and date filters
	query.Count(&totalRecords)

	// Apply the same filters for fetching data
	dataQuery := db.Where("appartment_uuid = ?", appartmentUUID).
		Where("type ILIKE ? OR signature ILIKE ?", "%"+search+"%", "%"+search+"%")

	// Add date range filter for data query
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			dataQuery = dataQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			dataQuery = dataQuery.Where("created_at < ?", endDateTime)
		}
	}

	err = dataQuery.
		Offset(offset).
		Limit(limit).
		Order("caisses.updated_at DESC").
		Preload("Appartment").
		Find(&caisses).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch Caisses",
			"error":   err.Error(),
		})
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	//  Prepare pagination metadata
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Caisses retrieved successfully",
		"data":       caisses,
		"pagination": pagination,
	})
}

func GetPaginatedCaissesManagerGeneral(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	// Parse search query
	search := c.Query("search", "")

	// Parse date range filters
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var caisses []models.Caisse
	var totalRecords int64

	// Build query with date filter
	query := db.Model(&models.Caisse{}).
		Where("type ILIKE ? OR signature ILIKE ?", "%"+search+"%", "%"+search+"%")

	// Add date range filter if provided
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add 24 hours to include the entire end date
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endDateTime)
		}
	}

	// Count total records matching the search and date filters
	query.Count(&totalRecords)

	// Apply the same filters for fetching data
	dataQuery := db.Where("type ILIKE ? OR signature ILIKE ?", "%"+search+"%", "%"+search+"%")

	// Add date range filter for data query
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			dataQuery = dataQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			dataQuery = dataQuery.Where("created_at < ?", endDateTime)
		}
	}

	err = dataQuery.
		Offset(offset).
		Limit(limit).
		Order("caisses.updated_at DESC").
		Preload("Appartment").
		Find(&caisses).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch Caisses",
			"error":   err.Error(),
		})
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	//  Prepare pagination metadata
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Caisses retrieved successfully",
		"data":       caisses,
		"pagination": pagination,
	})
}

// query all data
func GetAllCaisses(c *fiber.Ctx) error {
	db := database.DB
	var caisses []models.Caisse
	db.Preload("Appartment").Find(&caisses)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses",
		"data":    caisses,
	})
}

func GetAllCaissesByAppartmentUUID(c *fiber.Ctx) error {
	db := database.DB
	appartmentUUID := c.Params("appartment_uuid")

	var caisses []models.Caisse
	db.Where("appartment_uuid = ?", appartmentUUID).Preload("Appartment").Find(&caisses)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses",
		"data":    caisses,
	})
}

// Get one data
func GetCaisse(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var caisse models.Caisse
	db.Where("uuid = ?", uuid).Preload("Appartment").First(&caisse)
	if caisse.UUID == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Caisse found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Caisse found",
			"data":    caisse,
		},
	)
}

// Create data
func CreateCaisse(c *fiber.Ctx) error {
	p := &models.Caisse{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	if p.AppartmentUUID == "" || p.Type == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Form not complete - AppartmentUUID and Type are required",
				"data":    nil,
			},
		)
	}

	// Validate Type enum
	if p.Type != "Income" && p.Type != "Expense" {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Type must be either 'Income' or 'Expense'",
				"data":    nil,
			},
		)
	}

	caisse := &models.Caisse{
		AppartmentUUID: p.AppartmentUUID,
		Type:           p.Type,
		DeviceCDF:      p.DeviceCDF,
		DeviceUSD:      p.DeviceUSD,
		Signature:      p.Signature,
	}

	if err := utils.ValidateStruct(*caisse); err != nil {
		c.Status(400)
		return c.JSON(err)
	}

	caisse.UUID = utils.GenerateUUID()

	database.DB.Create(caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Caisse Created success",
			"data":    caisse,
		},
	)
}

// Update data
func UpdateCaisse(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateDataInput struct {
		AppartmentUUID string  `json:"appartment_uuid"`
		Type           string  `json:"type"`
		DeviceCDF      float64 `json:"device_cdf"`
		DeviceUSD      float64 `json:"device_usd"`
		Signature      string  `json:"signature"`
	}

	var updateData UpdateDataInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    nil,
			},
		)
	}

	// Validate Type enum
	if updateData.Type != "Income" && updateData.Type != "Expense" {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Type must be either 'Income' or 'Expense'",
				"data":    nil,
			},
		)
	}

	caisse := new(models.Caisse)

	db.Where("uuid = ?", uuid).First(&caisse)
	caisse.AppartmentUUID = updateData.AppartmentUUID
	caisse.Type = updateData.Type
	caisse.DeviceCDF = updateData.DeviceCDF
	caisse.DeviceUSD = updateData.DeviceUSD
	caisse.Signature = updateData.Signature

	db.Save(&caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Caisse updated success",
			"data":    caisse,
		},
	)
}

// Delete data
func DeleteCaisse(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var caisse models.Caisse
	db.Where("uuid = ?", uuid).First(&caisse)
	if caisse.UUID == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Caisse found",
				"data":    nil,
			},
		)
	}

	db.Delete(&caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Caisse deleted success",
			"data":    nil,
		},
	)
}

// Get balance for an appartment
func GetAppartmentBalance(c *fiber.Ctx) error {
	db := database.DB
	appartmentUUID := c.Params("appartment_uuid")

	// Parse date range filters
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var totalIncomeCDF, totalExpenseCDF, totalIncomeUSD, totalExpenseUSD float64

	// Build base queries with appartment filter
	incomeQuery := db.Model(&models.Caisse{}).Where("appartment_uuid = ? AND type = ?", appartmentUUID, "Income")
	expenseQuery := db.Model(&models.Caisse{}).Where("appartment_uuid = ? AND type = ?", appartmentUUID, "Expense")

	// Add date range filter if provided
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			incomeQuery = incomeQuery.Where("created_at >= ?", parsedStartDate)
			expenseQuery = expenseQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			incomeQuery = incomeQuery.Where("created_at < ?", endDateTime)
			expenseQuery = expenseQuery.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total income CDF
	incomeQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalIncomeCDF)

	// Calculate total expense CDF
	expenseQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalExpenseCDF)

	// Recreate queries for USD calculations
	incomeQueryUSD := db.Model(&models.Caisse{}).Where("appartment_uuid = ? AND type = ?", appartmentUUID, "Income")
	expenseQueryUSD := db.Model(&models.Caisse{}).Where("appartment_uuid = ? AND type = ?", appartmentUUID, "Expense")

	// Add date range filter for USD queries
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			incomeQueryUSD = incomeQueryUSD.Where("created_at >= ?", parsedStartDate)
			expenseQueryUSD = expenseQueryUSD.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			incomeQueryUSD = incomeQueryUSD.Where("created_at < ?", endDateTime)
			expenseQueryUSD = expenseQueryUSD.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total income USD
	incomeQueryUSD.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalIncomeUSD)

	// Calculate total expense USD
	expenseQueryUSD.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalExpenseUSD)

	balanceCDF := totalIncomeCDF - totalExpenseCDF
	balanceUSD := totalIncomeUSD - totalExpenseUSD

	// Get current exchange rates from database or default values
	currentRateUSDToCDF := getCurrentExchangeRate("USD", "CDF")
	currentRateCDFToUSD := getCurrentExchangeRate("CDF", "USD")

	// Convert for comparison
	totalIncomeCDFInUSD := utils.ConvertCDFToUSD(totalIncomeCDF, currentRateCDFToUSD)
	totalExpenseCDFInUSD := utils.ConvertCDFToUSD(totalExpenseCDF, currentRateCDFToUSD)
	totalIncomeUSDInCDF := utils.ConvertUSDToCDF(totalIncomeUSD, currentRateUSDToCDF)
	totalExpenseUSDInCDF := utils.ConvertUSDToCDF(totalExpenseUSD, currentRateUSDToCDF)

	balance := map[string]interface{}{
		"appartment_uuid":   appartmentUUID,
		"total_income_cdf":  totalIncomeCDF,
		"total_expense_cdf": totalExpenseCDF,
		"balance_cdf":       balanceCDF,
		"total_income_usd":  totalIncomeUSD,
		"total_expense_usd": totalExpenseUSD,
		"balance_usd":       balanceUSD,
		"conversions": map[string]interface{}{
			"income_cdf_in_usd":  totalIncomeCDFInUSD,
			"expense_cdf_in_usd": totalExpenseCDFInUSD,
			"income_usd_in_cdf":  totalIncomeUSDInCDF,
			"expense_usd_in_cdf": totalExpenseUSDInCDF,
		},
		"exchange_rates": map[string]interface{}{
			"usd_to_cdf": currentRateUSDToCDF,
			"cdf_to_usd": currentRateCDFToUSD,
		},
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Appartment balance retrieved successfully",
		"data":    balance,
	})
}

// Get global totals for all Income and Expense
func GetGlobalTotals(c *fiber.Ctx) error {
	db := database.DB

	// Parse date range filters
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var totalIncomeCDF, totalExpenseCDF, totalIncomeUSD, totalExpenseUSD float64

	// Build base queries with date filters
	incomeQuery := db.Model(&models.Caisse{}).Where("type = ?", "Income")
	expenseQuery := db.Model(&models.Caisse{}).Where("type = ?", "Expense")

	// Add date range filter if provided
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			incomeQuery = incomeQuery.Where("created_at >= ?", parsedStartDate)
			expenseQuery = expenseQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			incomeQuery = incomeQuery.Where("created_at < ?", endDateTime)
			expenseQuery = expenseQuery.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total income CDF
	incomeQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalIncomeCDF)

	// Calculate total expense CDF
	expenseQuery.Select("COALESCE(SUM(device_cdf), 0)").Row().Scan(&totalExpenseCDF)

	// Recreate queries for USD calculations
	incomeQueryUSD := db.Model(&models.Caisse{}).Where("type = ?", "Income")
	expenseQueryUSD := db.Model(&models.Caisse{}).Where("type = ?", "Expense")

	// Add date range filter for USD queries
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			incomeQueryUSD = incomeQueryUSD.Where("created_at >= ?", parsedStartDate)
			expenseQueryUSD = expenseQueryUSD.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			incomeQueryUSD = incomeQueryUSD.Where("created_at < ?", endDateTime)
			expenseQueryUSD = expenseQueryUSD.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total income USD
	incomeQueryUSD.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalIncomeUSD)

	// Calculate total expense USD
	expenseQueryUSD.Select("COALESCE(SUM(device_usd), 0)").Row().Scan(&totalExpenseUSD)

	balanceCDF := totalIncomeCDF - totalExpenseCDF
	balanceUSD := totalIncomeUSD - totalExpenseUSD

	// Get current exchange rates from database or default values
	currentRateUSDToCDF := getCurrentExchangeRate("USD", "CDF")
	currentRateCDFToUSD := getCurrentExchangeRate("CDF", "USD")

	// Convert for comparison
	totalIncomeCDFInUSD := utils.ConvertCDFToUSD(totalIncomeCDF, currentRateCDFToUSD)
	totalExpenseCDFInUSD := utils.ConvertCDFToUSD(totalExpenseCDF, currentRateCDFToUSD)
	totalIncomeUSDInCDF := utils.ConvertUSDToCDF(totalIncomeUSD, currentRateUSDToCDF)
	totalExpenseUSDInCDF := utils.ConvertUSDToCDF(totalExpenseUSD, currentRateUSDToCDF)

	balance := map[string]interface{}{
		"total_income_cdf":  totalIncomeCDF,
		"total_expense_cdf": totalExpenseCDF,
		"balance_cdf":       balanceCDF,
		"total_income_usd":  totalIncomeUSD,
		"total_expense_usd": totalExpenseUSD,
		"balance_usd":       balanceUSD,
		"conversions": map[string]interface{}{
			"income_cdf_in_usd":  totalIncomeCDFInUSD,
			"expense_cdf_in_usd": totalExpenseCDFInUSD,
			"income_usd_in_cdf":  totalIncomeUSDInCDF,
			"expense_usd_in_cdf": totalExpenseUSDInCDF,
		},
		"exchange_rates": map[string]interface{}{
			"usd_to_cdf": currentRateUSDToCDF,
			"cdf_to_usd": currentRateCDFToUSD,
		},
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Appartment balance retrieved successfully",
		"data":    balance,
	})
}

// Get totals by manager
func GetTotalsByManager(c *fiber.Ctx) error {
	db := database.DB
	managerUUID := c.Params("manager_uuid")

	// Parse date range filters
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	var totalIncomeCDF, totalExpenseCDF, totalIncomeUSD, totalExpenseUSD float64

	// Build base queries with manager filter
	baseQuery := "JOIN appartments ON caisses.appartment_uuid = appartments.uuid WHERE appartments.manager_uuid = ?"

	// Build date filter condition
	dateFilter := ""
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			dateFilter += " AND caisses.created_at >= '" + parsedStartDate.Format("2006-01-02 15:04:05") + "'"
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			dateFilter += " AND caisses.created_at < '" + endDateTime.Format("2006-01-02 15:04:05") + "'"
		}
	}

	// Join with appartment table to filter by manager
	db.Table("caisses").
		Joins(baseQuery, managerUUID).
		Where("caisses.type = ?"+dateFilter, "Income").
		Select("COALESCE(SUM(caisses.device_cdf), 0)").
		Row().Scan(&totalIncomeCDF)

	db.Table("caisses").
		Joins(baseQuery, managerUUID).
		Where("caisses.type = ?"+dateFilter, "Expense").
		Select("COALESCE(SUM(caisses.device_cdf), 0)").
		Row().Scan(&totalExpenseCDF)

	db.Table("caisses").
		Joins(baseQuery, managerUUID).
		Where("caisses.type = ?"+dateFilter, "Income").
		Select("COALESCE(SUM(caisses.device_usd), 0)").
		Row().Scan(&totalIncomeUSD)

	db.Table("caisses").
		Joins(baseQuery, managerUUID).
		Where("caisses.type = ?"+dateFilter, "Expense").
		Select("COALESCE(SUM(caisses.device_usd), 0)").
		Row().Scan(&totalExpenseUSD)

	// Calculate net balances
	netBalanceCDF := totalIncomeCDF - totalExpenseCDF
	netBalanceUSD := totalIncomeUSD - totalExpenseUSD

	// Get current exchange rate for conversions
	currentRateUSDToCDF := getCurrentExchangeRate("USD", "CDF")
	currentRateCDFToUSD := getCurrentExchangeRate("CDF", "USD")

	// Convert totals for comparison
	totalIncomeCDFInUSD := utils.ConvertCDFToUSD(totalIncomeCDF, currentRateCDFToUSD)
	totalExpenseCDFInUSD := utils.ConvertCDFToUSD(totalExpenseCDF, currentRateCDFToUSD)
	totalIncomeUSDInCDF := utils.ConvertUSDToCDF(totalIncomeUSD, currentRateUSDToCDF)
	totalExpenseUSDInCDF := utils.ConvertUSDToCDF(totalExpenseUSD, currentRateUSDToCDF)

	// Calculate grand totals in both currencies
	grandTotalIncomeCDF := totalIncomeCDF + totalIncomeUSDInCDF
	grandTotalExpenseCDF := totalExpenseCDF + totalExpenseUSDInCDF
	grandTotalIncomeUSD := totalIncomeUSD + totalIncomeCDFInUSD
	grandTotalExpenseUSD := totalExpenseUSD + totalExpenseCDFInUSD

	grandNetBalanceCDF := grandTotalIncomeCDF - grandTotalExpenseCDF
	grandNetBalanceUSD := grandTotalIncomeUSD - grandTotalExpenseUSD

	totals := map[string]interface{}{
		"manager_uuid": managerUUID,
		"income_totals": map[string]interface{}{
			"cdf_total":       totalIncomeCDF,
			"usd_total":       totalIncomeUSD,
			"cdf_in_usd":      totalIncomeCDFInUSD,
			"usd_in_cdf":      totalIncomeUSDInCDF,
			"grand_total_cdf": grandTotalIncomeCDF,
			"grand_total_usd": grandTotalIncomeUSD,
		},
		"expense_totals": map[string]interface{}{
			"cdf_total":       totalExpenseCDF,
			"usd_total":       totalExpenseUSD,
			"cdf_in_usd":      totalExpenseCDFInUSD,
			"usd_in_cdf":      totalExpenseUSDInCDF,
			"grand_total_cdf": grandTotalExpenseCDF,
			"grand_total_usd": grandTotalExpenseUSD,
		},
		"net_balances": map[string]interface{}{
			"net_balance_cdf":       netBalanceCDF,
			"net_balance_usd":       netBalanceUSD,
			"grand_net_balance_cdf": grandNetBalanceCDF,
			"grand_net_balance_usd": grandNetBalanceUSD,
		},
		"exchange_rates": map[string]interface{}{
			"usd_to_cdf": currentRateUSDToCDF,
			"cdf_to_usd": currentRateCDFToUSD,
		},
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Manager totals retrieved successfully",
		"data":    totals,
	})
}

// Convert amount between currencies
func ConvertCurrency(c *fiber.Ctx) error {
	type ConversionRequest struct {
		Amount       float64 `json:"amount" validate:"required,gt=0"`
		FromCurrency string  `json:"from_currency" validate:"required"`
		ToCurrency   string  `json:"to_currency" validate:"required"`
		Rate         float64 `json:"rate,omitempty"` // Taux optionnel, sinon utilise le taux de la DB ou par défaut
	}

	var req ConversionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"error":   err,
		})
	}

	// Validate currencies
	validCurrencies := []string{"USD", "CDF"}
	fromValid := false
	toValid := false
	for _, currency := range validCurrencies {
		if req.FromCurrency == currency {
			fromValid = true
		}
		if req.ToCurrency == currency {
			toValid = true
		}
	}

	if !fromValid || !toValid {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Supported currencies are: USD, CDF",
			"data":    nil,
		})
	}

	var exchangeRate float64
	if req.Rate > 0 {
		// Utiliser le taux fourni par l'utilisateur
		exchangeRate = req.Rate
	} else {
		// Utiliser le taux depuis la base de données ou par défaut
		exchangeRate = getCurrentExchangeRate(req.FromCurrency, req.ToCurrency)
	}

	convertedAmount := utils.ConvertCurrency(req.Amount, exchangeRate)

	conversion := map[string]interface{}{
		"original_amount":  req.Amount,
		"from_currency":    req.FromCurrency,
		"to_currency":      req.ToCurrency,
		"converted_amount": convertedAmount,
		"exchange_rate":    exchangeRate,
		"conversion_time":  time.Now(),
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Currency converted successfully",
		"data":    conversion,
	})
}
