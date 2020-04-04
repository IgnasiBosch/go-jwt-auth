package seed

import (
	"fmt"
	"github.com/IgnasiBosch/go-jwt-auth/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"log"
)

var users = []models.User{
	models.User{
		ID:       uuid.New(),
		Nickname: "dummy",
		Email:    "test@example.com",
		Password: "1234",
	},
}

func Load(db *gorm.DB) {
	fmt.Println("Populating seed data")
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop the table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate the table: %v", err)
	}

	for _, u := range users {
		_, err = u.SaveUser(db)
		//err = db.Debug().Model(&models.User{}).Create(&u).Error
		if err != nil {

			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
