package interest

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_mysql_CreateIfNotExists(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)
	testInterest := &Interest{Name: "testInterest"}

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec("INSERT INTO interests").WithArgs(testInterest.Name).WillReturnResult(res)

	err := repo.CreateIfNotExists(testInterest)

	lastId, _ := res.LastInsertId()

	assert.Nil(t, err)
	assert.Equal(t, int(lastId), testInterest.ID)
}

func Test_mysql_CreateIfNotExists_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	testInterest := &Interest{Name: "testInterest"}
	testError := fmt.Errorf("test error")

	mock.ExpectExec("INSERT INTO interests").WithArgs(testInterest).WillReturnError(testError)

	err := repo.CreateIfNotExists(testInterest)

	assert.NotNil(t, err)
	assert.Error(t, err, testError)
}

func Test_mysql_List_OneRow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "testInterest1")

	mock.ExpectQuery("SELECT id , name FROM interests").WillReturnRows(rows)

	res, err := repo.List()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func Test_mysql_List_TwoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "testInterest1")
	rows.AddRow(2, "testInterest2")

	mock.ExpectQuery("SELECT id , name FROM interests").WillReturnRows(rows)

	res, err := repo.List()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
}

func Test_mysql_List_ZeroRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"})

	mock.ExpectQuery("SELECT id , name FROM interests").WillReturnRows(rows)

	res, err := repo.List()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(res))
}

func Test_mysql_List_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	testError := fmt.Errorf("test error")

	mock.ExpectQuery("SELECT id , name FROM interests").WillReturnError(testError)

	res, err := repo.List()

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, testError)
}

func Test_mysql_InterestsByUserId_OneRow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)
	testUserId := 1

	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "testInterest1")

	mock.ExpectQuery("SELECT ui.interest_id , i.name FROM user_interests").WithArgs(testUserId).WillReturnRows(rows)

	res, err := repo.InterestsByUserId(testUserId)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func Test_mysql_InterestsByUserId_TwoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)
	testUserId := 1

	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "testInterest1")
	rows.AddRow(2, "testInterest2")

	mock.ExpectQuery("SELECT ui.interest_id , i.name FROM user_interests").WithArgs(testUserId).WillReturnRows(rows)

	res, err := repo.InterestsByUserId(testUserId)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
}

func Test_mysql_InterestsByUserId_ZeroRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)
	testUserId := 1

	rows := sqlmock.NewRows([]string{"id", "name"})

	mock.ExpectQuery("SELECT ui.interest_id , i.name FROM user_interests").WithArgs(testUserId).WillReturnRows(rows)

	res, err := repo.InterestsByUserId(testUserId)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(res))
}

func Test_mysql_InterestsByUserId_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)
	testUserId := 1
	testError := fmt.Errorf("test error")

	mock.ExpectQuery("SELECT ui.interest_id , i.name FROM user_interests").WithArgs(testUserId).WillReturnError(testError)

	res, err := repo.InterestsByUserId(testUserId)

	assert.Nil(t, res)
	assert.Error(t, err, testError)
}
