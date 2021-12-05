package seed

import (
	"fmt"
	"log"
	"person-service/internal/models"

	"github.com/jaswdr/faker"
	"gorm.io/gorm"
)

// Denotes the number of entries per batch
const recordsPerBatch = 100

//LoadRandomPersonData Creates entries in the person table for test usage.
func LoadRandomPersonData(numberOfRecords int32, databaseHandle *gorm.DB) {

	// We need to seed data only if the person table is completely empty.
	var person models.Person
	var count int64
	databaseHandle.Model(&person).Count(&count)

	if count > 0 {
		log.Println(fmt.Sprintf("Found %d rows already in the database! Not seeding any data!", count))
		return
	}

	// This slice will contain all the Person entries.
	var persons []models.Person

	// The root fake data generator.
	fakeDataGenerator := faker.New()

	// Create numberOfRecords number of entries.
	for i := 0; i < int(numberOfRecords); i++ {

		// This will act as a fake data generator for person.
		person := fakeDataGenerator.Person()

		// Generate a random entry for a new person and append it to the persons array.
		persons = append(persons, models.Person{
			Id:      fakeDataGenerator.UUID().V4(),
			Name:    person.Name(),
			Age:     fakeDataGenerator.Int32Between(1, 100),
			Email:   fakeDataGenerator.Internet().Email(),
			Country: person.Faker.Address().Country(),
		})
	}

	// The batch size needs to be dynamic.
	batchSize := numberOfRecords / recordsPerBatch

	// Create entries in the database in batches.
	transaction := databaseHandle.CreateInBatches(&persons, int(batchSize))

	if transaction.Error != nil {
		log.Println(transaction.Error)
		log.Fatalln("Unable to create seed data! Check error above for more details!")
	}

}
