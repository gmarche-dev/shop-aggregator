package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"shop-aggregator/internal/config"
	"shop-aggregator/internal/db/postgresql"
	"shop-aggregator/internal/handler"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
	"shop-aggregator/internal/router"
	"shop-aggregator/internal/usecase"
	"shop-aggregator/tools/migrations"
	"strconv"
	"strings"
	"testing"
)

type HandlerRepositories struct {
	Users       *postgresql.User
	Auth        *postgresql.Auth
	Company     *postgresql.Company
	Brand       *postgresql.Brand
	Store       *postgresql.Store
	Bill        *postgresql.Bill
	UserProduct *postgresql.UserProduct
	Product     *postgresql.Product
}

type HandlerUseCases struct {
	AuthUseCase        handler.AuthUsecase
	UserUseCase        handler.UserUsecase
	BrandUseCase       handler.BrandUseCase
	CompanyUseCase     handler.CompanyUseCase
	BillUseCase        handler.BillUseCase
	StoreUseCase       handler.StoreUseCase
	ProductUseCase     handler.ProductUseCase
	ProductUserProduct handler.UserProductUseCase
}

type Handlers struct {
	Auth           *handler.Auth
	User           *handler.User
	Brand          *handler.Brand
	Company        *handler.Company
	Bill           *handler.Bill
	Store          *handler.Store
	Product        *handler.Product
	UserProduct    *handler.UserProduct
	Initialisation *handler.Initialisation
}

type HandlerTestSuite struct {
	suite.Suite
	DB                  *postgresql.Client
	Container           testcontainers.Container
	ctx                 context.Context
	HandlerRepositories HandlerRepositories
	HandlerUseCases     HandlerUseCases
	Handlers            Handlers
	router              *gin.Engine
}

func (s *HandlerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "shopdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	s.Container = postgresContainer
	host, err := postgresContainer.Host(s.ctx)
	s.Require().NoError(err)
	portString, err := postgresContainer.MappedPort(s.ctx, "5432")
	s.Require().NoError(err)

	parts := strings.Split(string(portString), "/")
	if len(parts) == 0 {
		s.T().Fatal("invalid port")
	}
	i, err := strconv.Atoi(parts[0])
	s.Require().NoError(err)

	dbConfig := &config.DatabaseConfig{
		Host:     host,
		Port:     i,
		Username: "user",
		Password: "password",
		DBName:   "shopdb",
	}

	db, err := postgresql.NewDB(s.ctx, dbConfig)
	s.Require().NoError(err)

	s.DB = db

	// migration
	err = migrations.Run(s.ctx, s.DB, "../../migrations/deploy")
	s.Require().NoError(err)

	// load db requests
	s.HandlerRepositories.Users = postgresql.NewUsers(s.DB)
	s.HandlerRepositories.Auth = postgresql.NewAuth(s.DB)
	s.HandlerRepositories.Company = postgresql.NewCompany(s.DB)
	s.HandlerRepositories.Brand = postgresql.NewBrand(s.DB)
	s.HandlerRepositories.Store = postgresql.NewStore(s.DB)
	s.HandlerRepositories.Bill = postgresql.NewBill(s.DB)
	s.HandlerRepositories.UserProduct = postgresql.NewUserProduct(s.DB)
	s.HandlerRepositories.Product = postgresql.NewProduct(s.DB)

	// load usecases
	s.HandlerUseCases.AuthUseCase = usecase.NewAuth(s.HandlerRepositories.Auth, s.HandlerRepositories.Users)
	s.HandlerUseCases.UserUseCase = usecase.NewUsers(s.HandlerRepositories.Users)
	s.HandlerUseCases.BrandUseCase = usecase.NewBrand(s.HandlerRepositories.Brand)
	s.HandlerUseCases.CompanyUseCase = usecase.NewCompany(s.HandlerRepositories.Company)
	s.HandlerUseCases.BillUseCase = usecase.NewBill(s.HandlerRepositories.Bill, s.HandlerRepositories.Store, s.HandlerRepositories.Company, s.HandlerRepositories.UserProduct)
	s.HandlerUseCases.StoreUseCase = usecase.NewStore(s.HandlerRepositories.Store, s.HandlerRepositories.Company)
	s.HandlerUseCases.ProductUseCase = usecase.NewProduct(s.HandlerRepositories.Product, s.HandlerRepositories.Brand)
	s.HandlerUseCases.ProductUserProduct = usecase.NewUserProduct(s.HandlerRepositories.UserProduct)

	// load handlers
	s.Handlers.User = handler.NewUser(s.HandlerUseCases.UserUseCase)
	s.Handlers.Auth = handler.NewAuth(s.HandlerUseCases.AuthUseCase)
	s.Handlers.Brand = handler.NewBrand(s.HandlerUseCases.BrandUseCase)
	s.Handlers.Company = handler.NewCompany(s.HandlerUseCases.CompanyUseCase)
	s.Handlers.Bill = handler.NewBill(s.HandlerUseCases.BillUseCase)
	s.Handlers.Store = handler.NewStore(s.HandlerUseCases.StoreUseCase)
	s.Handlers.Product = handler.NewProduct(s.HandlerUseCases.ProductUseCase)
	s.Handlers.UserProduct = handler.NewUserProduct(s.HandlerUseCases.ProductUserProduct)
	s.Handlers.Initialisation = handler.NewInitialisation()

	s.router = gin.New()
	s.router = router.NewRouter(
		s.router,
		s.HandlerRepositories.Auth,
		s.Handlers.Auth,
		s.Handlers.User,
		s.Handlers.Brand,
		s.Handlers.Company,
		s.Handlers.Bill,
		s.Handlers.Store,
		s.Handlers.Product,
		s.Handlers.UserProduct,
		s.Handlers.Initialisation,
	)
}

func (s *HandlerTestSuite) TearDownSuite() {
	s.DB.Close()
	err := s.Container.Terminate(s.ctx)
	s.Require().NoError(err)
}

func (s *HandlerTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE store")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE bill")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE product")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE user_product")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE users")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE auth")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE brand")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE company")
	s.Require().NoError(err)
}

func (s *HandlerTestSuite) request(requestType, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(requestType, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)
	return w
}

func (s *HandlerTestSuite) requestWithToken(requestType, path, token string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(requestType, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	s.router.ServeHTTP(w, req)
	return w
}

func (s *HandlerTestSuite) createUserAndGenerateToken(login, password, email string) string {
	// create user
	s.createUser(login, password, email)

	// login
	return s.login(login, password)
}

func (s *HandlerTestSuite) createUser(login, password, email string) {
	// create user
	bodyCreate := request.CreateUser{
		Login:    login,
		Password: password,
		Email:    email,
	}
	bodyCreateBytes, err := json.Marshal(bodyCreate)
	s.Require().NoError(err)

	wCreate := s.request("POST", "/create-user", bodyCreateBytes)
	s.T().Log(fmt.Sprintf("createUser response body: %s", wCreate.Body.String()))
	s.Equal(200, wCreate.Code)
	user, err := s.HandlerRepositories.Users.GetUserByLogin(s.ctx, login)
	s.NoError(err)
	s.NotNil(user)
}

func (s *HandlerTestSuite) login(login, password string) string {
	bodyLogin := request.Login{
		Login:    login,
		Password: password,
	}
	bodyLoginBytes, err := json.Marshal(bodyLogin)
	s.Require().NoError(err)
	wLogin := s.request("POST", "/login", bodyLoginBytes)
	s.Equal(200, wLogin.Code)

	var loginResponse response.Login
	s.NoError(json.Unmarshal(wLogin.Body.Bytes(), &loginResponse))

	return loginResponse.Token
}

func (s *HandlerTestSuite) logout(token string) {
	wLogout := s.requestWithToken("POST", "/user/logout", token, nil)
	s.Equal(200, wLogout.Code)
}

func TestUserProductTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
