package caisses

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/appartment-app-api/database"
	"github.com/kgermando/appartment-app-api/models"
	"github.com/kgermando/appartment-app-api/utils"
)

// Paginate by SuperAdmin
func GetPaginatedCaissesSuperAdmin(c *fiber.Ctx) error {
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
		Where("type ILIKE ? OR signature ILIKE ? OR motif ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")

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
	dataQuery := db.
		Where("type ILIKE ? OR signature ILIKE ? OR motif ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")

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

	// Calculate totals for Income and Expense in USD and CDF
	var totalIncomeUSD, totalExpenseUSD, totalIncomeCDF, totalExpenseCDF float64

	// Build query for totals with same filters
	totalsQuery := db.Model(&models.Caisse{})

	// Add date range filter for totals
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			totalsQuery = totalsQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			totalsQuery = totalsQuery.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total Income USD
	totalsQuery.Where("type = ?", "Income").Select("COALESCE(SUM(device_usd), 0)").Scan(&totalIncomeUSD)

	// Calculate total Expense USD
	totalsQuery.Where("type = ?", "Expense").Select("COALESCE(SUM(device_usd), 0)").Scan(&totalExpenseUSD)

	// Reset query for CDF calculations
	totalsQuery = db.Model(&models.Caisse{})

	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			totalsQuery = totalsQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			totalsQuery = totalsQuery.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total Income CDF
	totalsQuery.Where("type = ?", "Income").Select("COALESCE(SUM(device_cdf), 0)").Scan(&totalIncomeCDF)

	// Calculate total Expense CDF
	totalsQuery.Where("type = ?", "Expense").Select("COALESCE(SUM(device_cdf), 0)").Scan(&totalExpenseCDF)

	// Prepare totals metadata
	totals := map[string]interface{}{
		"total_income_usd":  totalIncomeUSD,
		"total_expense_usd": totalExpenseUSD,
		"total_income_cdf":  totalIncomeCDF,
		"total_expense_cdf": totalExpenseCDF,
	}

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Caisses retrieved successfully",
		"data":       caisses,
		"pagination": pagination,
		"totals":     totals,
	})
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
		Where("type ILIKE ? OR signature ILIKE ? OR motif ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")

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
		Where("type ILIKE ? OR signature ILIKE ? OR motif ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")

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

	// Calculate totals for Income and Expense in USD and CDF
	var totalIncomeUSD, totalExpenseUSD, totalIncomeCDF, totalExpenseCDF float64

	// Build query for totals with same filters
	totalsQuery := db.Model(&models.Caisse{}).
		Where("appartment_uuid = ?", appartmentUUID)

	// Add date range filter for totals
	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			totalsQuery = totalsQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			totalsQuery = totalsQuery.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total Income USD
	totalsQuery.Where("type = ?", "Income").Select("COALESCE(SUM(device_usd), 0)").Scan(&totalIncomeUSD)

	// Calculate total Expense USD
	totalsQuery.Where("type = ?", "Expense").Select("COALESCE(SUM(device_usd), 0)").Scan(&totalExpenseUSD)

	// Reset query for CDF calculations
	totalsQuery = db.Model(&models.Caisse{}).
		Where("appartment_uuid = ?", appartmentUUID)

	if startDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDate); err == nil {
			totalsQuery = totalsQuery.Where("created_at >= ?", parsedStartDate)
		}
	}
	if endDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedEndDate.Add(24 * time.Hour)
			totalsQuery = totalsQuery.Where("created_at < ?", endDateTime)
		}
	}

	// Calculate total Income CDF
	totalsQuery.Where("type = ?", "Income").Select("COALESCE(SUM(device_cdf), 0)").Scan(&totalIncomeCDF)

	// Calculate total Expense CDF
	totalsQuery.Where("type = ?", "Expense").Select("COALESCE(SUM(device_cdf), 0)").Scan(&totalExpenseCDF)

	// Prepare totals metadata
	totals := map[string]interface{}{
		"total_income_usd":  totalIncomeUSD,
		"total_expense_usd": totalExpenseUSD,
		"total_income_cdf":  totalIncomeCDF,
		"total_expense_cdf": totalExpenseCDF,
	}

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Caisses retrieved successfully",
		"data":       caisses,
		"pagination": pagination,
		"totals":     totals,
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

	if p.AppartmentUUID == "" || p.Type == "" || p.Motif == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Form not complete - AppartmentUUID, Type and Motif are required",
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
		Motif:          p.Motif,
		Signature:      p.Signature,
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
		Motif          string  `json:"motif"`
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
	caisse.Motif = updateData.Motif
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
