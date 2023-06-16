package response

type PodListItem struct {
	Name    string `json:"name"`    // 名称
	Ready   string `json:"ready"`   // 状态 1/2
	Status  string `json:"status"`  // Running/Error
	Restart int32  `json:"restart"` // 重启 n 次
	Age     int64  `json:"age"`     // 运行时间
	Ip      string `json:"IP"`      // pod ip
	Node    string `json:"node"`    // pod 被调度到那个node
}
