package response

//@Author: morris

// NAME               READY   UP-TO-DATE   AVAILABLE   AGE
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	//副本数
	Replicas int32 `json:"replicas"`
	//READY字段表示Deployment中正在运行的Pod副本的数量
	Ready int32 `json:"ready"`
	//UP-TO-DATE字段表示与Deployment所期望的副本数相比，有多少个Pod副本是最新的
	UpToDate int32 `json:"upToDate"`
	//AVAILABLE字段表示Deployment中可用的Pod副本数
	Available int32 `json:"available"`
	Age       int64 `json:"age"`
}
