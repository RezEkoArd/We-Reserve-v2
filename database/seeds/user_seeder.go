package seeds

import (
	"wereserve/internal/core/domain/model"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	bytes, err := bcrypt.GenerateFromPassword([]byte("admin123"),14)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating password hash")
	}

	admin :=  model.User{
		Name: "Admin",
		Email: "admin@wereserve.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@wereserve.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Error seeding Admin Account")
	} else {
		log.Info().Msg("Admin Role Has been Seed")
	}
}