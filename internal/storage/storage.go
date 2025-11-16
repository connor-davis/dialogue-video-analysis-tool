package storage

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/common"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	Database() *gorm.DB
	Migrate() error
}

type storage struct {
	database *gorm.DB
}

func New() Storage {
	databaseDsn := common.EnvString("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=one_go port=5432 sslmode=disable TimeZone=Africa/Johannesburg")

	database, err := gorm.Open(postgres.Open(databaseDsn), &gorm.Config{})

	if err != nil {
		log.Errorf("🔥 Failed to establish connection with the database: %s", err.Error())
	}

	log.Info("✅ Connection established with the database.")

	return &storage{
		database: database,
	}
}

func (s *storage) Database() *gorm.DB {
	return s.database
}

func (s *storage) Migrate() error {
	if err := s.database.AutoMigrate(
		&models.User{},
		&models.Role{},
	); err != nil {
		return err
	}

	log.Info("✅ Database migration completed successfully.")

	return nil
}
