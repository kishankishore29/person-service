module person-service

// +heroku goVersion go1.16
go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/gin-gonic/gin v1.7.4
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/jaswdr/faker v1.8.0
	github.com/lib/pq v1.10.2
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.17.0
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.16
)
