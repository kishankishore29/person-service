package handlers

import (
	"encoding/json"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"net/http"
	"net/http/httptest"
	"person-service/internal/models"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestGetPersons Tests the GetPersons method.
func TestGetPersons(t *testing.T) {

	var err error

	mockServer := APIServer{
		Database: nil,
	}

	// This acts as a mock for the gin test context and will be used to check the response
	recorder := httptest.NewRecorder()

	// This test context is passed to the handler function
	testContext, _ := gin.CreateTestContext(recorder)

	// This will help mus mock the database dependency.
	db, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	mockServer.Database, err = gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic("Failed to create test database for the test")
	}

	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()

	testId, _ := uuid.NewV4()

	want := models.Person{
		Id:      testId.String(),
		Name:    "Kishan",
		Age:     10,
		Email:   "abc@gmail.com",
		Country: "India",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "people"`)).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"Id", "Name", "Age", "Email", "Country"}).
			AddRow(want.Id, want.Name, want.Age, want.Email, want.Country))

	mockServer.GetPersons(testContext)
	mock.ExpectationsWereMet()

	var got []models.Person
	err = json.Unmarshal(recorder.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	// Assert that the response was a 200 OK
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, reflect.DeepEqual(want, got[0]))

}
