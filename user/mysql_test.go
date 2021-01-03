package user

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_mysql_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(10, 1))

	testUser := &User{
		FirstName: "TestFirst",
		Lastname:  "TestLast",
		Age:       12,
		Sex: Sex{
			ID: 1,
		},
	}

	res, err := repo.Create(testUser)
	if err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, res.ID, 0)
	assert.Equal(t, res.ID, 10)
}

func Test_mysql_Create_WrongCityID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	dbError := fmt.Errorf("FK constraint error")

	mock.ExpectExec("INSERT INTO users").WillReturnError(dbError)

	testUser := &User{
		FirstName: "TestFirst",
		Lastname:  "TestLast",
		Age:       12,
		Sex: Sex{
			ID: 1,
		},
	}

	res, err := repo.Create(testUser)

	assert.NotNil(t, err)
	assert.Nil(t, res)
	if !strings.Contains(err.Error(), dbError.Error()) {
		t.Errorf("got: %v; shoult contain %v", err, dbError)
	}
}

func Test_mysql_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, 1, "Мужчина", "TestLogin", 1, "TestCity")

	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)
	users, err := repo.List()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
}

func Test_mysql_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, 1, "Мужчина", "TestLogin", 1, "TestCity")

	mock.ExpectQuery("SELECT u.id").WithArgs(1).WillReturnRows(rows)
	user, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
}

func Test_mysql_GetByID_ErrNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})

	mock.ExpectQuery("SELECT u.id").WithArgs(1).WillReturnRows(rows)
	_, err = repo.GetByID(1)

	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
		t.Errorf("got %v; should contain %v", err, sql.ErrNoRows)
	}
}

func Test_mysql_GetByID_EmptyCityData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, 1, "Мужчина", "TestLogin", nil, nil)

	mock.ExpectQuery("SELECT u.id").WithArgs(1).WillReturnRows(rows)
	res, err := repo.GetByID(1)

	assert.Nil(t, err)
	assert.NotNil(t, res.City)
}

func Test_mysql_GetByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testLogin := "TestLogin"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, 1, "Мужчина", testLogin, 1, "TestCity")

	mock.ExpectQuery("SELECT u.id").WithArgs(testLogin).WillReturnRows(rows)

	res, err := repo.GetByLogin(testLogin)

	assert.Nil(t, err)
	assert.Equal(t, res.Login, testLogin)
}

func Test_mysql_GetByLogin_ErrNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testLogin := "TestLogin"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex_id", "sex", "login", "city_id", "city_name"})

	mock.ExpectQuery("SELECT u.id").WithArgs(testLogin).WillReturnRows(rows)

	user, err := repo.GetByLogin(testLogin)

	assert.Nil(t, err)
	assert.Nil(t, user)
}
