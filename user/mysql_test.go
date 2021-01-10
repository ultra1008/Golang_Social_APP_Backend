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
		Sex:       "Мужчина",
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
		Sex:       "Мужчина",
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
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", "TestLogin", 1, "TestCity")

	mock.ExpectQuery("ELECT u.id").WillReturnRows(rows)
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
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name", "password"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", "TestLogin", 1, "TestCity", "TestPassword")

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
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name", "password"})

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
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name", "password"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", "TestLogin", nil, nil, "testPasswrod")

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

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name", "password"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", testLogin, 1, "TestCity", "testPassword")

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

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name", "password"})

	mock.ExpectQuery("SELECT u.id").WithArgs(testLogin).WillReturnRows(rows)

	user, err := repo.GetByLogin(testLogin)

	assert.Nil(t, err)
	assert.Nil(t, user)
}

func Test_mysql_AddFriend(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 2

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec("INSERT INTO friends ").WithArgs(testUserId, testFriendId).WillReturnResult(res)

	err = repo.AddFriend(testUserId, testFriendId)

	assert.Nil(t, err)
}

func Test_mysql_AddFriend_zeroLinesAdded(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 2
	expectedError := fmt.Errorf("test error")

	res := sqlmock.NewResult(0, 0)
	mock.ExpectExec("INSERT INTO friends ").WithArgs(testUserId, testFriendId).WillReturnResult(res)

	err = repo.AddFriend(testUserId, testFriendId)

	assert.NotNil(t, err)
	assert.Error(t, err, expectedError)
}

func Test_mysql_AddFriend_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 2
	expectedError := fmt.Errorf("test error")

	mock.ExpectExec("INSERT INTO friends ").WithArgs(testUserId, testFriendId).WillReturnError(expectedError)

	err = repo.AddFriend(testUserId, testFriendId)

	assert.NotNil(t, err)
	assert.Error(t, err, expectedError)
}

func Test_mysql_AddFriend_SameIDs(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 1
	expectedError := fmt.Errorf("user ID and friend ID are equal")

	err = repo.AddFriend(testUserId, testFriendId)

	assert.NotNil(t, err)
	assert.Error(t, err, expectedError)
}

func Test_mysql_DeleteFriend(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 2

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec("DELETE FROM friends ").WithArgs(testUserId, testFriendId).WillReturnResult(res)

	err = repo.DeleteFriend(testUserId, testFriendId)

	assert.Nil(t, err)
}

func Test_mysql_DeleteFriend_zeroLinesAdded(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 2
	expectedError := fmt.Errorf("test error")

	res := sqlmock.NewResult(0, 0)
	mock.ExpectExec("DELETE FROM friends ").WithArgs(testUserId, testFriendId).WillReturnResult(res)

	err = repo.DeleteFriend(testUserId, testFriendId)

	assert.NotNil(t, err)
	assert.Error(t, err, expectedError)
}

func Test_mysql_DeleteFriend_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	testUserId := 1
	testFriendId := 2
	expectedError := fmt.Errorf("test error")

	mock.ExpectExec("DELETE FROM friends ").WithArgs(testUserId, testFriendId).WillReturnError(expectedError)

	err = repo.DeleteFriend(testUserId, testFriendId)

	assert.NotNil(t, err)
	assert.Error(t, err, expectedError)
}
func Test_mysql_Friends_OneRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)
	testUserId := 1

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", "testLogin", 1, "TestCity")

	mock.ExpectQuery("SELECT u.id").WithArgs(testUserId).WillReturnRows(rows)

	friends, err := repo.Friends(testUserId)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(friends))
}

func Test_mysql_Friends_TwoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)
	testUserId := 1

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name"})
	rows.AddRow(1, "TestFirst", "TestLast", 12, "Мужчина", "testLogin", 1, "TestCity")
	rows.AddRow(2, "TestFirst2", "TestLast2", 16, "Мужчина", "testLogin2", 2, "TestCity2")

	mock.ExpectQuery("SELECT u.id").WithArgs(testUserId).WillReturnRows(rows)

	friends, err := repo.Friends(testUserId)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(friends))
}

func Test_mysql_Friends_ZeroRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)
	testUserId := 1

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "sex", "login", "city_id", "city_name"})

	mock.ExpectQuery("SELECT u.id").WithArgs(testUserId).WillReturnRows(rows)

	friends, err := repo.Friends(testUserId)

	assert.Nil(t, err)
	assert.NotNil(t, friends)
	assert.Equal(t, 0, len(friends))
}

func Test_mysql_Friends_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)
	testUserId := 1
	testError := fmt.Errorf("test error")

	mock.ExpectQuery("SELECT u.id").WithArgs(testUserId).WillReturnError(testError)

	friends, err := repo.Friends(testUserId)

	assert.Nil(t, friends)
	assert.Error(t, err, testError)
}
