package user

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/niklod/highload-social-network/user/city"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	cityDb, cityMock, _ := sqlmock.New()

	testUser := &User{
		ID:        1,
		FirstName: "testFirstName",
		Lastname:  "testLastName",
		Age:       12,
		Sex:       "Женщина",
		City: city.City{
			ID:   1,
			Name: "Москва",
		},
		Login:    "TestLogin",
		Password: "TestPassword",
	}

	repo := NewRepository(db)
	cityRepo := city.NewRepository(cityDb)
	citySvc := city.NewService(cityRepo)
	userSvc := NewService(repo, citySvc)

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name"})
	cityRows := sqlmock.NewRows([]string{"id"}).AddRow(testUser.City.ID)

	mock.ExpectQuery("SELECT u.id").WithArgs(testUser.Login).WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(int64(testUser.ID), 1))
	cityMock.ExpectQuery("INSERT INTO citys").WillReturnRows(cityRows)

	u, err := userSvc.Create(testUser)

	assert.Nil(t, err)
	assert.Equal(t, u, testUser)
}

func TestService_Create_UserAlreadyExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testUser := &User{
		ID:        1,
		FirstName: "testFirstName",
		Lastname:  "testLastName",
		Age:       12,
		Sex:       "Мужчина",
		City: city.City{
			ID:   1,
			Name: "Москва",
		},
		Login:    "TestLogin",
		Password: "TestPassword",
	}

	repo := NewRepository(db)
	cityRepo := city.NewRepository(db)
	citySvc := city.NewService(cityRepo)
	userSvc := NewService(repo, citySvc)
	expectedErrorString := "user already exist"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", testUser.Login, 1, "TestCity")

	mock.ExpectQuery("SELECT u.id").WithArgs(testUser.Login).WillReturnRows(rows)

	_, err = userSvc.Create(testUser)

	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), expectedErrorString) {
		t.Errorf("got %v; should contain %v", err.Error(), expectedErrorString)
	}
}
