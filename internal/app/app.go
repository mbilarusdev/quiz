package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mbilarusdev/quiz/internal/common"
	"github.com/mbilarusdev/quiz/internal/handler"
	"github.com/mbilarusdev/quiz/internal/repository"
	"github.com/mbilarusdev/quiz/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App interface {
	Start()
}

type QuizApp struct{}

func NewQuizApp() *QuizApp {
	return new(QuizApp)
}

func (app *QuizApp) Start() {
	// Logger
	fmt.Println("Init logger...")
	common.InitLogger()
	defer common.L.Sync()

	// Config
	fmt.Println("Parsing config...")
	common.Conf = common.NewQuizConfig()
	dsn := common.Conf.PostgresDsn

	// Postgres
	time.Sleep(time.Second * 5)
	fmt.Println("Connect postgres...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Can't open postgres connection: %v\n", err))
	}

	// Repositories
	questionRepo := repository.NewQuestionRepository(db)
	answerRepo := repository.NewAnswerRepository(db)

	// Services
	questionSrv := service.NewQuestionService(questionRepo, answerRepo)
	answerSrv := service.NewAnswerService(answerRepo, questionRepo, db)

	// Handlers
	questionHandler := handler.NewQuestionHandler(questionSrv)
	answerHandler := handler.NewAnswerHandler(answerSrv)

	router := mux.NewRouter()

	// Questions
	router.HandleFunc("/questions", questionHandler.FindAll).Methods(http.MethodGet)
	router.HandleFunc("/questions", questionHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/questions/{id}", questionHandler.FindOneDetailed).Methods(http.MethodGet)
	router.HandleFunc("/questions/{id}", questionHandler.Delete).Methods(http.MethodDelete)

	// Answers
	router.HandleFunc("/questions/{id}/answers", answerHandler.AddAnswer).Methods(http.MethodPost)
	router.HandleFunc("/answers/{id}", answerHandler.FindOne).Methods(http.MethodGet)
	router.HandleFunc("/answers/{id}", answerHandler.Delete).Methods(http.MethodDelete)

	fmt.Println("Starting server...")

	server := &http.Server{
		Addr:           common.Conf.Addr,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v\n", err))
	}

}
