package postgres

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"person/models"
	"person/user"
	"testing"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository user.Repository
	person     *models.User
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.repository = &UserRepository{s.DB}
}

func (s *Suite) TestPersonRepository_GetById() {
	var (
		id      = 1
		name    = "test-name"
		address = "test-address"
		work    = "test-work"
		age     = 10
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT *`)).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "address", "work", "age"}).
				AddRow(id, name, address, work, age))

	res, err := s.repository.GetUser(context.Background(), id)
	require.NoError(s.T(), err)

	expected := models.User{Id: id, Name: name, Address: address, Work: work, Age: age}
	require.Nil(s.T(), deep.Equal(&expected, res))
}

func (s *Suite) TestPersonRepository_GetAll() {
	var (
		persons = []*models.User{
			{Id: 1, Name: "name1", Address: "add1", Work: "work1", Age: 10},
			{Id: 2, Name: "name2", Address: "add2", Work: "work2", Age: 20},
		}
	)
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "address", "work", "age"}).
				AddRow(persons[0].Id, persons[0].Name, persons[0].Address, persons[0].Work, persons[0].Age).
				AddRow(persons[1].Id, persons[1].Name, persons[1].Address, persons[1].Work, persons[1].Age))

	res, err := s.repository.GetAll(context.Background())
	require.NoError(s.T(), err)

	require.Nil(s.T(), deep.Equal(persons, res))
}

func (s *Suite) TestPersonRepository_Create() {
	var (
		id      = 1
		name    = "test-name"
		address = "test-address"
		work    = "test-work"
		age     = 10
	)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(name, address, work, age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectCommit()

	created := models.User{Name: name, Address: address, Work: work, Age: age}
	res, err := s.repository.CreateUser(context.Background(), &created)
	require.NoError(s.T(), err)

	require.Equal(s.T(), res, id)
}

func (s *Suite) TestPersonRepository_Update() {
	var (
		id      = 1
		name    = "test-name"
		address = "test-address"
		work    = "test-work"
		age     = 10
	)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WithArgs(name, address, work, age, id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	s.mock.ExpectCommit()

	updated := models.User{Name: name, Address: address, Work: work, Age: age}
	res, err := s.repository.ChangeUser(context.Background(), &updated, id)
	require.NoError(s.T(), err)
	updated.Id = id

	require.Nil(s.T(), deep.Equal(res, &updated))
}

func (s *Suite) TestPersonRepository_Delete() {
	var (
		id      = 1
		name    = "test-name"
		address = "test-address"
		work    = "test-work"
		age     = 10
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "address", "work", "age"}).
			AddRow(id, name, address, work, age))

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	s.mock.ExpectCommit()

	err := s.repository.DeleteUser(context.Background(), id)

	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestRun(t *testing.T) {
	suite.Run(t, new(Suite))
}