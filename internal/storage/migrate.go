package storage

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/gofiber/fiber/v3/log"
)

func (s *storage) Migrate() error {
	if err := s.database.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Organization{},
	); err != nil {
		return err
	}

	log.Info("✅ Database migration completed successfully.")

	return nil
}
