package storage

import "gorm.io/gorm"

func (s *storage) Database() *gorm.DB {
	return s.database
}
