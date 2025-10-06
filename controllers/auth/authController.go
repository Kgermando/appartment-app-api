package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kgermando/appartment-app-api/database"
	"github.com/kgermando/appartment-app-api/models"
	"github.com/kgermando/appartment-app-api/utils"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// ensureAdminUser vérifie si un utilisateur admin existe, sinon en crée un
func ensureAdminUser() error {
	var adminUser models.User

	// Vérifier si un utilisateur avec le rôle Administrator existe
	result := database.DB.Where("role = ?", "Admin").First(&adminUser)

	if result.Error == nil {
		// Un admin existe déjà
		return nil
	}

	// Créer un utilisateur admin par défaut
	defaultAdmin := &models.User{
		UUID:       uuid.New().String(),
		Fullname:   "Super Admin",
		Email:      "admin@appartment-app.com",
		Telephone:  "+243000000000",
		Role:       "Admin",
		Permission: "ALL",
		Status:     true,
		Signature:  "Super Administrator",
	}

	// Définir un mot de passe par défaut (vous devriez le changer)
	defaultAdmin.SetPassword("Admin@123")

	// Sauvegarder l'admin en base de données
	if err := database.DB.Create(defaultAdmin).Error; err != nil {
		return fmt.Errorf("erreur lors de la création de l'admin par défaut: %v", err)
	}

	fmt.Println("Utilisateur admin créé avec succès:")
	fmt.Printf("Email: %s\n", defaultAdmin.Email)
	fmt.Printf("Téléphone: %s\n", defaultAdmin.Telephone)
	fmt.Println("Mot de passe: Admin@123")
	fmt.Println("⚠️  IMPORTANT: Changez le mot de passe par défaut après la première connexion!")

	return nil
}

// CreateAdminUser endpoint pour créer manuellement un utilisateur admin
func CreateAdminUser(c *fiber.Ctx) error {
	type AdminInput struct {
		Fullname  string `json:"fullname" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Telephone string `json:"telephone" validate:"required"`
		Password  string `json:"password" validate:"required,min=6"`
	}

	adminInput := new(AdminInput)

	if err := c.BodyParser(&adminInput); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := utils.ValidateStruct(*adminInput); err != nil {
		c.Status(400)
		return c.JSON(err)
	}

	// Vérifier si un admin existe déjà avec cet email ou téléphone
	var existingUser models.User
	result := database.DB.Where("email = ? OR telephone = ?", adminInput.Email, adminInput.Telephone).First(&existingUser)
	if result.Error == nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Un utilisateur avec cet email ou téléphone existe déjà",
		})
	}

	// Créer le nouvel admin
	newAdmin := &models.User{
		UUID:       uuid.New().String(),
		Fullname:   adminInput.Fullname,
		Email:      adminInput.Email,
		Telephone:  adminInput.Telephone,
		Role:       "Admin",
		Permission: "ALL",
		Status:     true,
		Signature:  "Admin",
	}

	newAdmin.SetPassword(adminInput.Password)

	if err := database.DB.Create(newAdmin).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Erreur lors de la création de l'admin",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Utilisateur admin créé avec succès",
		"data": fiber.Map{
			"uuid":      newAdmin.UUID,
			"fullname":  newAdmin.Fullname,
			"email":     newAdmin.Email,
			"telephone": newAdmin.Telephone,
			"role":      newAdmin.Role,
		},
	})
}

func Register(c *fiber.Ctx) error {

	nu := new(models.User)

	if err := c.BodyParser(&nu); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if nu.Password != nu.PasswordConfirm {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	u := &models.User{
		UUID:       uuid.New().String(),
		Fullname:   nu.Fullname,
		Email:      nu.Email,
		Telephone:  nu.Telephone,
		Role:       nu.Role,
		Permission: nu.Permission,
		Status:     nu.Status,
		Signature:  nu.Signature,
	}

	u.SetPassword(nu.Password)

	if err := utils.ValidateStruct(*u); err != nil {
		c.Status(400)
		return c.JSON(err)
	}

	database.DB.Create(u)

	return c.JSON(fiber.Map{
		"message": "user account created",
		"data":    u,
	})
}

func Login(c *fiber.Ctx) error {
	// S'assurer qu'un utilisateur admin existe
	if err := ensureAdminUser(); err != nil {
		fmt.Printf("Erreur lors de la vérification/création de l'admin: %v\n", err)
		// On continue le processus de login même si la création de l'admin échoue
	}

	lu := new(models.Login)

	if err := c.BodyParser(&lu); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := utils.ValidateStruct(*lu); err != nil {
		c.Status(400)
		return c.JSON(err)
	}

	u := &models.User{}

	result := database.DB.Where("email = ? OR telephone = ?", lu.Identifier, lu.Identifier).
		First(&u)

	if result.Error != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "invalid email or telephone 😰",
		})
	}

	if err := u.ComparePassword(lu.Password); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "mot de passe incorrect! 😰",
		})
	}

	if !u.Status {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "vous n'êtes pas autorisé de se connecter 😰",
		})
	}

	token, err := utils.GenerateJwt(u.UUID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    token,
	})

}

func AuthUser(c *fiber.Ctx) error {

	token := c.Query("token")

	fmt.Println("token", token)

	// cookie := c.Cookies("token")
	UserUUID, _ := utils.VerifyJwt(token)

	fmt.Println("UserUUID", UserUUID)

	u := models.User{}

	database.DB.Where("users.uuid = ?", UserUUID).
		First(&u)
	r := &models.UserResponse{
		UUID:       u.UUID,
		Fullname:   u.Fullname,
		Email:      u.Email,
		Telephone:  u.Telephone,
		Role:       u.Role,
		Permission: u.Permission,
		Status:     u.Status,
		Signature:  u.Signature,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
	return c.JSON(r)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // 1 day ,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
		"Logout":  "success",
	})

}

// User bioprofile
func UpdateInfo(c *fiber.Ctx) error {
	type UpdateDataInput struct {
		Fullname  string `json:"fullname"`
		Email     string `json:"email"`
		Telephone string `json:"telephone"`
		Signature string `json:"signature"`
	}
	var updateData UpdateDataInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error(),
		})
	}

	cookie := c.Cookies("token")

	UserUUID, _ := utils.VerifyJwt(cookie)

	user := new(models.User)

	db := database.DB

	// Utiliser UUID au lieu de convertir en int
	result := db.Where("uuid = ?", UserUUID).First(&user)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Utilisateur non trouvé",
		})
	}

	user.Fullname = updateData.Fullname
	user.Email = updateData.Email
	user.Telephone = updateData.Telephone
	user.Signature = updateData.Signature

	db.Save(&user)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully updated",
		"data":    user,
	})

}

func ChangePassword(c *fiber.Ctx) error {
	type UpdateDataInput struct {
		OldPassword     string `json:"old_password"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
	}
	var updateData UpdateDataInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error(),
		})
	}

	// Utiliser la même logique que AuthUser - récupérer le token depuis les query params
	token := c.Query("token")

	fmt.Println("token", token)

	UserUUID, err := utils.VerifyJwt(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Token invalide ou expiré",
		})
	}

	fmt.Println("UserUUID", UserUUID)

	user := new(models.User)

	// Utiliser UUID au lieu de id car c'est la clé primaire du modèle User
	result := database.DB.Where("uuid = ?", UserUUID).First(&user)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Utilisateur non trouvé",
		})
	}

	if err := user.ComparePassword(updateData.OldPassword); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "votre mot de passe n'est pas correct! 😰",
		})
	}

	if updateData.Password != updateData.PasswordConfirm {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "passwords do not match",
		})
	}

	// Utiliser la méthode SetPassword du modèle au lieu de utils.HashPassword
	user.SetPassword(updateData.Password)

	db := database.DB
	db.Save(&user)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Mot de passe modifié avec succès",
	})
}
