package osint

import (
	"github.com/yash3605/spectre/internal/models"
	"github.com/yash3605/spectre/modules/osint/domain"
	"github.com/yash3605/spectre/modules/osint/ip"
)

func IPLookup(search string, activeModule int) models.Result {
	switch activeModule {
	case 0:
		return ip.Lookup(search)
	case 1:
		return domain.Lookup(search)
	default:
		return models.Result{
			Title:  "Not Implemented Yet.",
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}
}
