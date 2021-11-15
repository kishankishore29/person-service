package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"person-service/internal/models"
	"reflect"
	"regexp"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestSuite Contains the variables to be used by the tests
type TestSuite struct {
	suite.Suite                            // testify suite struct
	mockServer  APIServer                  // The mockserver to be used to call the gin handler methods
	mockHandle  sqlmock.Sqlmock            // The mock handler to be used for mocking database calls.
	recorder    *httptest.ResponseRecorder // Recorder to assert the response of the HTTP request.
	testContext *gin.Context               // Context to be passed to the handler methods
	testPerson  models.Person              // Test person struct instance to use in the tests
}

// TestSetup This is the setup method for the test suite. It declares all the common variables which will be used by the test in this file
func (t *TestSuite) SetupTest() {
	t.mockServer = APIServer{
		Database: nil,
	}

	// This acts as a mock for the gin test context and will be used to check the response
	t.recorder = httptest.NewRecorder()

	// This test context is passed to the handler function
	t.testContext, _ = gin.CreateTestContext(t.recorder)

	// This will help us mock the database dependency.
	db, mock, _ := sqlmock.New()
	mock.MatchExpectationsInOrder(false)

	t.mockHandle = mock

	// The dialector helps us pass a custom database object
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	// This needs to be declared explcitly because of the following assignment
	var err error
	t.mockServer.Database, err = gorm.Open(dialector, &gorm.Config{})

	testId, _ := uuid.NewV4()

	t.testPerson = models.Person{
		Id:      testId.String(),
		Name:    "Test",
		Age:     10,
		Email:   "abc@gmail.com",
		Country: "USA",
	}

	if err != nil {
		panic("Failed to create test database for the test")
	}
}

// TestGetPersons Tests the GetPersons method.
func (ts *TestSuite) TestGetPersons() {

	assert.True(ts.T(), true)

	ts.mockHandle.ExpectBegin()

	want := ts.testPerson

	ts.mockHandle.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "people"`)).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Name", "Age", "Email", "Country"}).
			AddRow(want.Id, want.Name, want.Age, want.Email, want.Country))

	ts.mockServer.GetPersons(ts.testContext)
	ts.mockHandle.ExpectationsWereMet()

	var got []models.Person
	err := json.Unmarshal(ts.recorder.Body.Bytes(), &got)

	if err != nil {
		ts.T().Fatal(err)
	}

	// Assert that the response was a 200 OK
	assert.Equal(ts.T(), http.StatusOK, ts.recorder.Code)

	// Asseert that the expected and the actual response are matching.
	assert.True(ts.T(), reflect.DeepEqual(want, got[0]))

}

// TestCreatePersons Tests the create person method
func (ts *TestSuite) TestCreatePersons() {
	ts.mockHandle.ExpectBegin()

	want := ts.testPerson

	// Mock the corresponding query.
	ts.mockHandle.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "people" ("name","age","email","country") VALUES ($1,$2,$3,$4) RETURNING "id"`)).WithArgs(want.Name, want.Age, want.Email, want.Country).WillReturnRows(
		sqlmock.NewRows([]string{"Id"}).
			AddRow(want.Id))

	ts.mockHandle.ExpectCommit()

	// Create the request body
	request := fmt.Sprintf(`{
		"Name": "%s",
		"Age": %d,
		"Email": "%s",
		"Country": "%s"
	}`, want.Name, want.Age, want.Email, want.Country)

	// Set the request body.
	ts.testContext.Request = &http.Request{
		Header: make(http.Header),
		Body:   io.NopCloser(strings.NewReader(request)),
	}

	ts.testContext.Request.Method = "POST"
	ts.testContext.Request.Header.Set("Content-Type", "application/json")

	ts.mockServer.CreatePerson(ts.testContext)
	ts.mockHandle.ExpectationsWereMet()

	var got models.Person
	err := json.Unmarshal(ts.recorder.Body.Bytes(), &got)

	if err != nil {
		ts.T().Fatal(err)
	}

	// Assert that the response was a 200 OK
	assert.Equal(ts.T(), http.StatusCreated, ts.recorder.Code)

	// Assert that the expected and actual response is matching
	assert.True(ts.T(), reflect.DeepEqual(want, got))

}

// TestRunSuite Runs all the tests in the suite.
func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
