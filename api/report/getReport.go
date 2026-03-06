package report

import (
	"encoding/json"
	"go-inventory/repository"
	"net/http"
	"sort"
)

// blm ada service, tapi readonly jadi pake repo
type ReportAPIService struct {
	Repo *repository.LogInMemory
}

func (s *ReportAPIService) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report := s.Repo.GetLogs(r.Context())
	sort.Slice(report, func(i, j int) bool {
		return report[i].CreatedAt.Before(report[j].CreatedAt)
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}
