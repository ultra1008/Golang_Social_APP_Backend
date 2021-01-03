package user

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/niklod/highload-social-network/user/city"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	cityDb, cityMock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(cityDb, "mysql")

	testUser := &User{
		ID:        1,
		FirstName: "testFirstName",
		Lastname:  "testLastName",
		Age:       12,
		Sex: Sex{
			ID:   1,
			Name: "Мужской",
		},
		City: city.City{
			ID:   1,
			Name: "Москва",
		},
		Login:    "TestLogin",
		Password: "TestPassword",
	}

	repo := NewRepository(db)
	cityRepo := city.NewRepository(sqlxDB)
	citySvc := city.NewService(cityRepo)
	userSvc := NewService(repo, citySvc)

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})
	cityRows := sqlmock.NewRows([]string{"id", "city_name", "created_by_user"}).AddRow(testUser.City.ID, testUser.City.Name, 0)

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
	sqlxDB := sqlx.NewDb(db, "mysql")

	testUser := &User{
		ID:        1,
		FirstName: "testFirstName",
		Lastname:  "testLastName",
		Age:       12,
		Sex: Sex{
			ID:   1,
			Name: "Мужской",
		},
		City: city.City{
			ID:   1,
			Name: "Москва",
		},
		Login:    "TestLogin",
		Password: "TestPassword",
	}

	repo := NewRepository(db)
	cityRepo := city.NewRepository(sqlxDB)
	citySvc := city.NewService(cityRepo)
	userSvc := NewService(repo, citySvc)
	expectedErrorString := "user already exist"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, 1, "Мужчина", testUser.Login, 1, "TestCity")

	mock.ExpectQuery("SELECT u.id").WithArgs(testUser.Login).WillReturnRows(rows)

	_, err = userSvc.Create(testUser)

	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), expectedErrorString) {
		t.Errorf("got %v; should contain %v", err.Error(), expectedErrorString)
	}
}
