package harbor

import "testing"

func TestInitHarbor(t *testing.T) {
	scheme := "https"
	host := "harbor.kubemanager.com"
	username := "admin"
	password := "123123"
	harbor, err := InitHarbor(scheme, host, username, password)
	if err != nil {
		t.Error(err)
	}
	harbor.GetProjects()
}
