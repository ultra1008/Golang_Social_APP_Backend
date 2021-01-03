package city

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_mysql_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCityName := "Test"

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	mockRows := sqlmock.NewRows([]string{"id", "city_name", "created_by_user"}).AddRow(1, testCityName, 0)
	mock.ExpectQuery("INSERT INTO citys").WithArgs(testCityName).WillReturnRows(mockRows)

	city, err := repo.Create(testCityName)

	assert.Nil(t, err)
	assert.NotNil(t, city)
	assert.Equal(t, testCityName, city.Name)

}
func Test_mysql_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testError := fmt.Errorf("test error")
	testCityName := "TestCity"

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	mock.ExpectQuery("INSERT INTO citys").WithArgs(testCityName).WillReturnError(testError)

	city, err := repo.Create(testCityName)

	assert.NotNil(t, err)
	assert.Nil(t, city)
	if !strings.Contains(err.Error(), testError.Error()) {
		t.Errorf("got: %v; should contain %v", err, testError)
	}
}

func Test_mysql_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCityName := "Test"

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	mockRows := sqlmock.NewRows([]string{"id", "city_name", "created_by_user"}).AddRow(1, testCityName, 0)
	mock.ExpectQuery("SELECT id, city_name, created_by_user").WillReturnRows(mockRows)

	city, err := repo.List()

	assert.Nil(t, err)
	assert.NotNil(t, city)
	assert.Equal(t, 1, len(city))
	assert.Equal(t, testCityName, city[0].Name)
}

func Test_mysql_List_MultipleRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	tc := []struct {
		Name        string
		Rows        *sqlmock.Rows
		ExpectedLen int
	}{
		{
			Name:        "2 rows",
			Rows:        sqlmock.NewRows([]string{"id", "city_name", "created_by_user"}).AddRow(1, "Test1", 0).AddRow(2, "Test2", 0),
			ExpectedLen: 2,
		},
		{
			Name:        "3 rows",
			Rows:        sqlmock.NewRows([]string{"id", "city_name", "created_by_user"}).AddRow(1, "Test1", 0).AddRow(2, "Test2", 0).AddRow(3, "Test3", 0),
			ExpectedLen: 3,
		},
	}

	for _, cc := range tc {
		t.Run(cc.Name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, city_name, created_by_user").WillReturnRows(cc.Rows)

			citys, err := repo.List()
			assert.Nil(t, err)
			assert.Equal(t, cc.ExpectedLen, len(citys))
		})
	}
}

func Test_mysql_List_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testError := fmt.Errorf("test error")

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	mock.ExpectQuery("SELECT id, city_name, created_by_user").WillReturnError(testError)

	city, err := repo.List()

	assert.Nil(t, city)
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), testError.Error()) {
		t.Errorf("got: %v; should contain %v", err, testError)
	}
}

func Test_mysql_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCityID := 10

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	mockRows := sqlmock.NewRows([]string{"id", "city_name", "created_by_user"}).AddRow(testCityID, "Moscow", 0)
	mock.ExpectQuery("SELECT id, city_name").WithArgs(testCityID).WillReturnRows(mockRows)

	city, err := repo.GetByID(testCityID)

	assert.Nil(t, err)
	assert.NotNil(t, city)
	assert.Equal(t, testCityID, city.ID)
}

func Test_mysql_GetByID_ErrNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCityID := 10

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewRepository(sqlxDB)

	mockRows := sqlmock.NewRows([]string{"id", "city_name", "created_by_user"})
	mock.ExpectQuery("SELECT id, city_name").WithArgs(testCityID).WillReturnRows(mockRows)

	_, err = repo.GetByID(testCityID)

	assert.NotNil(t, err)
	assert.Error(t, err, sql.ErrNoRows)
}
