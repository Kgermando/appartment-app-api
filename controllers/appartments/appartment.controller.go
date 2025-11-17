package appartments

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/appartment-app-api/database"
	"github.com/kgermando/appartment-app-api/models"
	"github.com/kgermando/appartment-app-api/utils"
)

// Paginate
func GetPaginatedAppartments(c *fiber.Ctx) error {
	db := database.DB

	managerUUID := c.Params("manager_uuid")

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

	var appartments []models.Appartment
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Appartment{}).
		Where("manager_uuid = ?", managerUUID).
		Where("name ILIKE ? OR number ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("manager_uuid = ?", managerUUID).
		Where("name ILIKE ? OR number ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("appartments.updated_at DESC").
		Preload("Manager").
		Preload("Caisses").
		Find(&appartments).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch Appartments",
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
		"message":    "Appartments retrieved successfully",
		"data":       appartments,
		"pagination": pagination,
	})
}

func GetPaginatedAppartmentsManagerGeneral(c *fiber.Ctx) error {
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

	var appartments []models.Appartment
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Appartment{}).
		Where("name ILIKE ? OR number ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Where("name ILIKE ? OR number ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("appartments.updated_at DESC").
		Preload("Manager").
		Preload("Caisses").
		Find(&appartments).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch Appartments",
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
		"message":    "Appartments retrieved successfully",
		"data":       appartments,
		"pagination": pagination,
	})
}

// query all data
func GetAllAppartments(c *fiber.Ctx) error {
	db := database.DB
	var appartments []models.Appartment
	db.Preload("Manager").Preload("Caisses").Find(&appartments)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All appartments",
		"data":    appartments,
	})
}

func GetAllAppartmentsByManagerUUID(c *fiber.Ctx) error {
	db := database.DB
	managerUUID := c.Params("manager_uuid")

	var appartments []models.Appartment
	db.Where("manager_uuid = ?", managerUUID).Preload("Manager").Preload("Caisses").Find(&appartments)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All appartments",
		"data":    appartments,
	})
}

// Get one data
func GetAppartment(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var appartment models.Appartment
	db.Where("uuid = ?", uuid).Preload("Manager").Preload("Caisses").First(&appartment)
	if appartment.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Appartment name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Appartment found",
			"data":    appartment,
		},
	)
}

// Get appartment payment statistics by month
func GetAppartmentStats(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	// Vérifier si l'appartement existe
	var appartment models.Appartment
	if err := db.Where("uuid = ?", uuid).First(&appartment).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Appartment not found",
			"data":    nil,
		})
	}

	// Obtenir l'année courante ou depuis les paramètres de requête
	year := time.Now().Year()
	if yearParam := c.Query("year"); yearParam != "" {
		if parsedYear, err := strconv.Atoi(yearParam); err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Initialiser les statistiques pour les 12 mois
	monthlyStats := make(map[string]map[string]float64)
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	for _, month := range months {
		monthlyStats[month] = map[string]float64{
			"total_cdf": 0.0,
			"total_usd": 0.0,
		}
	}

	// Récupérer toutes les entrées (Income) de la caisse pour cet appartement
	var caisses []models.Caisse
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	err := db.Where("appartment_uuid = ? AND type = ? AND created_at >= ? AND created_at <= ?",
		uuid, "Income", startDate, endDate).Find(&caisses).Error

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
			monthlyStats[monthName]["total_cdf"] += caisse.DeviceCDF
			monthlyStats[monthName]["total_usd"] += caisse.DeviceUSD
		}
	}

	// Calculer les totaux annuels
	var totalYearCDF, totalYearUSD float64
	for _, month := range months {
		totalYearCDF += monthlyStats[month]["total_cdf"]
		totalYearUSD += monthlyStats[month]["total_usd"]
	}

	// Préparer la réponse
	response := map[string]interface{}{
		"appartment_info": map[string]interface{}{
			"uuid":           appartment.UUID,
			"created_at":     appartment.CreatedAt,
			"updated_at":     appartment.UpdatedAt,
			"name":           appartment.Name,
			"number":         appartment.Number,
			"surface":        appartment.Surface,
			"rooms":          appartment.Rooms,
			"bathrooms":      appartment.Bathrooms,
			"balcony":        appartment.Balcony,
			"furnished":      appartment.Furnished,
			"monthly_rent":   appartment.MonthlyRent,
			"garantie_month": appartment.GarantieMonth,
			"garantie":       appartment.Garantie,
			"echeance":       appartment.Echeance,
			"status":         appartment.Status,
			"manager_uuid":   appartment.ManagerUUID,
		},
		"year":          year,
		"monthly_stats": monthlyStats,
		"yearly_totals": map[string]float64{
			"total_cdf": totalYearCDF,
			"total_usd": totalYearUSD,
		},
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

// Create data
func CreateAppartment(c *fiber.Ctx) error {
	// Define input struct with string for date field
	type CreateAppartmentInput struct {
		Name          string    `json:"name"`
		Number        string    `json:"number"`
		Surface       float64   `json:"surface"`
		Rooms         int       `json:"rooms"`
		Bathrooms     int       `json:"bathrooms"`
		Balcony       bool      `json:"balcony"`
		Furnished     bool      `json:"furnished"`
		MonthlyRent   float64   `json:"monthly_rent"`
		GarantieMonth float64   `json:"garantie_month"`
		Garantie      float64   `json:"garantie_montant"`
		Echeance      time.Time `json:"echeance"`
		Status        string    `json:"status"`
		ManagerUUID   string    `json:"manager_uuid"`
	}

	var input CreateAppartmentInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	appartment := &models.Appartment{
		Name:          input.Name,
		Number:        input.Number,
		Surface:       input.Surface,
		Rooms:         input.Rooms,
		Bathrooms:     input.Bathrooms,
		Balcony:       input.Balcony,
		Furnished:     input.Furnished,
		MonthlyRent:   input.MonthlyRent,
		GarantieMonth: input.GarantieMonth,
		Garantie:      input.Garantie,
		Echeance:      input.Echeance,
		Status:        input.Status,
		ManagerUUID:   input.ManagerUUID,
	}

	appartment.UUID = utils.GenerateUUID()

	database.DB.Create(appartment)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Appartment Created success",
			"data":    appartment,
		},
	)
}

// Update data
func UpdateAppartment(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateDataInput struct {
		Name          string    `json:"name"`
		Number        string    `json:"number"`
		Surface       float64   `json:"surface"`
		Rooms         int       `json:"rooms"`
		Bathrooms     int       `json:"bathrooms"`
		Balcony       bool      `json:"balcony"`
		Furnished     bool      `json:"furnished"`
		MonthlyRent   float64   `json:"monthly_rent"`
		GarantieMonth float64   `json:"garantie_month"`
		Garantie      float64   `json:"garantie_montant"`
		Echeance      time.Time `json:"echeance"` // Accept as string
		Status        string    `json:"status"`
		ManagerUUID   string    `json:"manager_uuid"`
	}

	var updateData UpdateDataInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    nil,
			},
		)
	}

	appartment := new(models.Appartment)

	db.Where("uuid = ?", uuid).First(&appartment)
	appartment.Name = updateData.Name
	appartment.Number = updateData.Number
	appartment.Surface = updateData.Surface
	appartment.Rooms = updateData.Rooms
	appartment.Bathrooms = updateData.Bathrooms
	appartment.Balcony = updateData.Balcony
	appartment.Furnished = updateData.Furnished
	appartment.MonthlyRent = updateData.MonthlyRent
	appartment.GarantieMonth = updateData.GarantieMonth
	appartment.Garantie = updateData.Garantie
	appartment.Echeance = updateData.Echeance
	appartment.Status = updateData.Status
	appartment.ManagerUUID = updateData.ManagerUUID

	db.Save(&appartment)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Appartment updated success",
			"data":    appartment,
		},
	)
}

// Delete data
func DeleteAppartment(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var appartment models.Appartment
	db.Where("uuid = ?", uuid).First(&appartment)
	if appartment.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Appartment name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&appartment)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Appartment deleted success",
			"data":    nil,
		},
	)
}
