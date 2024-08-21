package harbor

import "testing"

// @Author: morris
func TestInitHarbor(t *testing.T) {
	scheme := "https"
	host := "harbor.k8s-web"
	username := "admin"
	password := "123"
	_, err := InitHarbor(scheme, host, username, password)
	if err != nil {
		t.Error(err)
	}
	//harbor.GetProjects()
}
