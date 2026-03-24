package country

import (
	"strings"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, names []string) error {
	if len(names) == 0 {
		return nil
	}

	clean := make([]string, 0, len(names))
	for _, n := range names {
		v := strings.TrimSpace(n)
		if v != "" {
			clean = append(clean, v)
		}
	}
	if len(clean) == 0 {
		return nil
	}

	for _, name := range clean {
		var existing schema.CountryManagement
		if err := db.Where("name_country = ?", name).First(&existing).Error; err == nil {
			continue
		}
		if err := db.Create(&schema.CountryManagement{NameCountry: name}).Error; err != nil {
			return err
		}
	}

	return nil
}
