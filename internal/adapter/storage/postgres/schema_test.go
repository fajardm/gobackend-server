package postgres_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fajardm/gobackend-server/internal/adapter/database"
	"github.com/fajardm/gobackend-server/internal/adapter/storage/postgres"
	"github.com/fajardm/gobackend-server/internal/adapter/storage/postgres/mock"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type SchemaTestSuite struct {
	suite.Suite
	ctx          context.Context
	mockCtrl     *gomock.Controller
	mockDB       *sql.DB
	mockSqlxDB   *sqlx.DB
	mockSQL      sqlmock.Sqlmock
	mockQuerier  *mock.MockSchemaQuerier
	mockPostgres *database.Postgres
	repo         *postgres.Schema
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}

func (t *SchemaTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.mockCtrl = gomock.NewController(t.T())
	t.mockDB, t.mockSQL, _ = sqlmock.New()
	t.mockSqlxDB = sqlx.NewDb(t.mockDB, "sqlmock")
	t.mockQuerier = mock.NewMockSchemaQuerier(t.mockCtrl)
	t.mockPostgres = &database.Postgres{DB: t.mockSqlxDB}
	t.repo = &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: json.Unmarshal,
		Querier:   postgres.GetSchemaQuerier(),
	}
}

func (t *SchemaTestSuite) TearDownTest() {
	t.mockDB.Close()
	t.mockCtrl.Finish()
}

func (t *SchemaTestSuite) TestAll_QuerierError() {
	mockError := errors.New("unexpected")

	t.mockQuerier.EXPECT().AllClasses().Return("", nil, mockError)

	repo := &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: json.Unmarshal,
		Querier:   t.mockQuerier,
	}
	res, err := repo.All(t.ctx)

	t.Nil(res)
	t.Equal(mockError, err)
}

func (t *SchemaTestSuite) TestAll_QueryxContextError() {
	mockError := errors.New("unexpected")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnError(mockError)

	res, err := t.repo.All(t.ctx)

	t.Nil(res)
	t.Equal(mockError, err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestAll_ScanError() {
	rows := sqlmock.NewRows([]string{"className"})
	rows.AddRow("")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	res, err := t.repo.All(t.ctx)

	t.Nil(res)
	t.Error(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestAll_UnmarshalError() {
	mockUnmarshal := func(data []byte, v interface{}) error { return errors.New("unexpected") }

	rows := sqlmock.NewRows([]string{"className", "schema"})
	rows.AddRow("", "")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	repo := &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: mockUnmarshal,
		Querier:   postgres.GetSchemaQuerier(),
	}
	res, err := repo.All(t.ctx)

	t.Nil(res)
	t.Error(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestAll_NoError() {
	rows := sqlmock.NewRows([]string{"className", "schema"})
	rows.AddRow("_User", `{
		"fields":{
		   "_rperm":{
			  "type":"Array"
		   },
		   "_wperm":{
			  "type":"Array"
		   },
		   "objectId":{
			  "type":"UUID"
		   },
		   "createdAt":{
			  "type":"Date"
		   },
		   "updatedAt":{
			  "type":"Date"
		   }
		},
		"className":"_User",
		"classLevelPermissions":{
		   "get":{
			  "*":true
		   },
		   "find":{
			  "*":true
		   },
		   "count":{
			  "*":true
		   },
		   "create":{
			  "*":true
		   },
		   "delete":{
			  "*":true
		   },
		   "update":{
			  "*":true
		   },
		   "addField":{
			  "*":true
		   },
		   "protectedFields":{
			  "*":[
				 
			  ]
		   }
		}
	}`)

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	res, err := t.repo.All(t.ctx)

	t.Equal(schema.ClassName("_User"), res["_User"].ClassName)
	t.NoError(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestFindByClassName_QuerierError() {
	mockError := errors.New("unexpected")
	mockClassName := schema.ClassName("xxxx")

	t.mockQuerier.EXPECT().FindClassByName(mockClassName).Return("", nil, mockError)

	repo := &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: json.Unmarshal,
		Querier:   t.mockQuerier,
	}
	res, err := repo.FindByClassName(t.ctx, mockClassName)

	t.Nil(res)
	t.Equal(mockError, err)
}

func (t *SchemaTestSuite) TestFindByClassName_QueryRowxContextError() {
	mockError := errors.New("unexpected")
	mockClassName := schema.ClassName("xxxx")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnError(mockError)

	res, err := t.repo.FindByClassName(t.ctx, mockClassName)

	t.Nil(res)
	t.Equal(mockError, err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestFindByClassName_NoRowsError() {
	mockClassName := schema.ClassName("xxxx")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnError(sql.ErrNoRows)

	res, err := t.repo.FindByClassName(t.ctx, mockClassName)

	t.Nil(res)
	t.Equal("error 3: class xxxx not exists", err.Error())
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestFindByClassName_ScanError() {
	mockClassName := schema.ClassName("xxxx")

	rows := sqlmock.NewRows([]string{"className"})
	rows.AddRow("xxxx")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	res, err := t.repo.FindByClassName(t.ctx, mockClassName)

	t.Nil(res)
	t.Error(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestFindByClassName_UnmarshalError() {
	mockClassName := schema.ClassName("xxxx")
	mockUnmarshal := func(data []byte, v interface{}) error { return errors.New("unexpected") }

	rows := sqlmock.NewRows([]string{"className", "schema"})
	rows.AddRow("", "")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	repo := &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: mockUnmarshal,
		Querier:   postgres.GetSchemaQuerier(),
	}
	res, err := repo.FindByClassName(t.ctx, mockClassName)

	t.Nil(res)
	t.Error(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestFindByClassName_NoError() {
	mockClassName := schema.ClassName("_User")

	rows := sqlmock.NewRows([]string{"className", "schema"})
	rows.AddRow("_User", `{
		"fields":{
		   "_rperm":{
			  "type":"Array"
		   },
		   "_wperm":{
			  "type":"Array"
		   },
		   "objectId":{
			  "type":"UUID"
		   },
		   "createdAt":{
			  "type":"Date"
		   },
		   "updatedAt":{
			  "type":"Date"
		   }
		},
		"className":"_User",
		"classLevelPermissions":{
		   "get":{
			  "*":true
		   },
		   "find":{
			  "*":true
		   },
		   "count":{
			  "*":true
		   },
		   "create":{
			  "*":true
		   },
		   "delete":{
			  "*":true
		   },
		   "update":{
			  "*":true
		   },
		   "addField":{
			  "*":true
		   },
		   "protectedFields":{
			  "*":[
				 
			  ]
		   }
		}
	}`)

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	res, err := t.repo.FindByClassName(t.ctx, mockClassName)

	t.Equal(mockClassName, res.ClassName)
	t.NoError(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestExists_QuerierError() {
	mockError := errors.New("unexpected")
	mockClassName := schema.ClassName("xxxx")

	t.mockQuerier.EXPECT().ExistsClass(mockClassName).Return("", nil, mockError)

	repo := &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: json.Unmarshal,
		Querier:   t.mockQuerier,
	}
	res, err := repo.Exists(t.ctx, mockClassName)

	t.False(res)
	t.Equal(mockError, err)
}

func (t *SchemaTestSuite) TestExists_QueryRowxContextError() {
	mockError := errors.New("unexpected")
	mockClassName := schema.ClassName("xxxx")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnError(mockError)

	res, err := t.repo.Exists(t.ctx, mockClassName)

	t.False(res)
	t.Equal(mockError, err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestExists_ScanError() {
	mockClassName := schema.ClassName("xxxx")

	rows := sqlmock.NewRows([]string{"exists", "xxxx"})
	rows.AddRow("FALSE", "")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	res, err := t.repo.Exists(t.ctx, mockClassName)

	t.False(res)
	t.Error(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestExists_NoError() {
	mockClassName := schema.ClassName("xxxx")

	rows := sqlmock.NewRows([]string{"exists"})
	rows.AddRow("TRUE")

	t.mockSQL.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	res, err := t.repo.Exists(t.ctx, mockClassName)

	t.True(res)
	t.NoError(err)
	t.NoError(t.mockSQL.ExpectationsWereMet())
}

func (t *SchemaTestSuite) TestCreate_BeginError() {
	mockError := errors.New("unexpected")
	mockData := schema.Schema{}

	t.mockSQL.ExpectBegin().WillReturnError(mockError)

	err := t.repo.Create(t.ctx, mockData)

	t.Equal(mockError, err)
}

func (t *SchemaTestSuite) TestCreate_txCreateSchemaTableIfNotExists_QuerierError() {
	mockError := errors.New("unexpected")
	mockData := schema.Schema{}

	t.mockSQL.ExpectBegin()
	t.mockQuerier.EXPECT().CreateSchemaTableIfNotExists().Return("", nil, mockError)
	t.mockSQL.ExpectRollback()

	repo := &postgres.Schema{
		DB:        t.mockPostgres,
		Unmarshal: json.Unmarshal,
		Querier:   t.mockQuerier,
	}
	err := repo.Create(t.ctx, mockData)

	t.Equal(mockError, err)
}

func (t *SchemaTestSuite) TestCreate_txCreateSchemaTableIfNotExists_ExecContextError() {
	mockError := errors.New("unexpected")
	mockData := schema.Schema{}

	t.mockSQL.ExpectBegin()
	t.mockSQL.ExpectExec("^CREATE (.+)").WillReturnError(mockError)
	t.mockSQL.ExpectRollback()

	err := t.repo.Create(t.ctx, mockData)

	t.Equal(mockError, err)
}
