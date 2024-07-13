package database

import(
	"github.com/devmehta22/JWT-Auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"github.com/joho/godotenv"
)

func DotEnv(key string) string {
	
	err := godotenv.Load(".env")

	if err !=nil{
		panic(err)
	}

	return os.Getenv(key)
}

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error)  {
	dsn := DotEnv("URI")

	db,err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
		}

	DB = db
	db.AutoMigrate(&models.User{})

	return db,nil
	
}