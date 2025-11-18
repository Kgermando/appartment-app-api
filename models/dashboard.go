package models

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

type TopManager struct {
	ManagerUUID    string  `json:"manager_uuid"`
	ManagerName    string  `json:"manager_name"`
	TotalRevenue   float64 `json:"total_revenue"`
	NetProfit      float64 `json:"net_profit"`
	ApartmentCount int64   `json:"apartment_count"`
	OccupancyRate  float64 `json:"occupancy_rate"`
	Efficiency     float64 `json:"efficiency"` // Net profit / Total apartments
}
