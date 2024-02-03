package main

import (
	"log"
	"net/http"
	"os"

	"github.com/methlox/daisy/models"
	"github.com/methlox/daisy/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"time"
)

type Form struct {
	Title			*string 	`json:"title"`
	Description		*string 	`json:"desc"`
	Created_at 		time.Time   `json:"created_at"`

}

type Question struct {
	Question_text	*string 	`json:"ques_text"`
	Question_order	*string 	`json:"ques_order"`
	Created_at 		time.Time   `json:"created_at"`
}

type FormResponse struct {
	Responded_at 	time.Time   `json:"responded_at"`
}

type Response struct {
	Answer			*string 	`json:"answer"`
}
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateForm(context *fiber.Ctx) error {
	form := Form{}

	err := context.BodyParser(&form)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&form).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create form"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "form has been created"})
	return nil
}

func (r *Repository) GetForm(context *fiber.Ctx) error {
	formModel := &[]models.Form{}

	err := r.DB.Find(formModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get form"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "form fetched successfully",
		"data": formModel,
	})
	return nil
}

func (r *Repository) GetQuestion(context *fiber.Ctx) error {
	formQues := &[]models.Question{}

	err := r.DB.Find(formQues).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get question"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "questions fetched successfully",
		"data": formQues,
	})
	return nil
}

func (r *Repository) CreateQuestion(context *fiber.Ctx) error {
	ques := Question{}

	err := context.BodyParser(&ques)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&ques).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create question"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "question has been created"})
	return nil
}

func (r *Repository) GetAllFormResponses(context *fiber.Ctx) error {
	formRes := &[]models.FormResponse{}

	err := r.DB.Find(formRes).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get form responses"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "form responses fetched successfully",
		"data": formRes,
	})
	return nil
}

func (r *Repository) GetResponse(context *fiber.Ctx) error {
	res := &[]models.Response{}

	err := r.DB.Find(res).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get response"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "response fetched successfully",
		"data": res,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_form", r.CreateForm)
	api.Get("/get_form/", r.GetForm)
	api.Get("/get_ques", r.GetQuestion)
	api.Post("/create_ques", r.CreateQuestion)
	api.Get("/get_all_responses", r.GetAllFormResponses)
	api.Get("/get_response", r.GetResponse)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateForm(db)
	if err != nil {
		log.Fatal("could not migrate form db")
	}
	err = models.MigrateFormResponse(db)
	if err != nil {
		log.Fatal("could not migrate form response db")
	}
	err = models.MigrateQuestion(db)
	if err != nil {
		log.Fatal("could not migrate ques db")
	}
	err = models.MigrateResponse(db)
	if err != nil {
		log.Fatal("could not migrate response db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}