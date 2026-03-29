package conductor

import (
	"github.com/yash3605/spectre/internal/models"
	"github.com/yash3605/spectre/modules/osint"
)

func Search(query string, activeTab models.Tab, activeModule int) models.Result {
	switch activeTab {
	case models.OSINT:
		return osint.IPLookup(query, activeModule)
	case models.Infosys:
		return models.Result{
			Title:  "Science",
			Data:   make(map[string]string),
			Status: models.StateIdle,
		}
	case models.Entity:
		return models.Result{
			Title:  "Person",
			Data:   make(map[string]string),
			Status: models.StateIdle,
		}
	default:
		return models.Result{
			Title:  "",
			Data:   make(map[string]string),
			Status: models.StateError,
		}
	}
}
