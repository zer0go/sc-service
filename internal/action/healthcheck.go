package action

import "net/http"

type HealthCheckAction struct{}

func NewHealthCheckAction() *HealthCheckAction {
	return new(HealthCheckAction)
}

func (a *HealthCheckAction) Handle(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"success":true}` + "\n"))
}
