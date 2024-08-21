package response

//@Author: morris

type IngressRoute struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Age       int64  `json:"age"`
}
