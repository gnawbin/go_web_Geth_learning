package response

// @Author: morris
// NAME          CLASS    HOSTS         ADDRESS   PORTS   AGE
type Ingress struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Class     string `json:"class"`
	Hosts     string `json:"hosts"`
	Age       int64  `json:"age"`
}
