package response

// @Author: morris
// NAME                READY   AGE
type StatefulSet struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace"`
	//ready的pod的个数
	Ready int32 `json:"ready"`
	//副本数
	Replicas int32 `json:"replicas"`
	Age      int64 `json:"age"`
}
