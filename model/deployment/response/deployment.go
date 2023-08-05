package response

// NAME               READY   UP-TO-DATE   AVAILABLE   AGE

type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int32  `json:"replicas"`  // 副本数
	Ready     int32  `json:"ready"`     // READY字段表示Deployment中正在运行的Pod副本的数量
	UpToDate  int32  `json:"upToDate"`  // UP-TO-DATE字段表示与Deployment所期望的副本数相比，有多少个Pod副本是最新的
	Available int32  `json:"available"` // AVAILABLE字段表示Deployment中可用的Pod副本数
	Age       int64  `json:"age"`
}
