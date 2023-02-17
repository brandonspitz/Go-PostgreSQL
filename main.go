package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brandonspitz/models"
	"github.com/brandonspitz/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateArtifacts(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}

type Artifact struct {
	Student string `json:"student"`
	Site    string `json:"site"`
	Type    string `json:"type"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_artifacts", r.CreateArtifact)
	api.Delete("delete_artifact/:id", r.DeleteArtifact)
	api.Get("/get_artifacts/:id", r.GetArtifactByID)
	api.Get("/artifacts", r.GetArtifacts)
}

func (r *Repository) CreateArtifact(context *fiber.Ctx) error {
	artifact := Artifact{}

	err := context.BodyParser(&artifact)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&artifact).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create artifact"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "artifact has been created"})

	return nil
}

func (r *Repository) GetArtifacts(context *fiber.Ctx) error {
	artifactModels := &[]models.Artifacts{}

	err := r.DB.Find(artifactModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get artifacts"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    artifactModels,
	})
	return nil
}

func (r *Repository) DeleteArtifact(context *fiber.Ctx) error {
	artifactModel := models.Artifacts{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(artifactModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete artifact",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "artifact deleted successfully",
	})
	return nil
}

func (r *Repository) GetArtifactByID(context *fiber.Ctx) error {
	id := context.Params("id")
	artifactModel := &models.Artifacts{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(artifactModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get artifact"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "artifact id fetched successfully",
		"data":    artifactModel,
	})
	return nil
}
