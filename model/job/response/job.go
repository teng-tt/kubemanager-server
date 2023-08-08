package response

type Job struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Completions int32  `json:"completions"` // 控制Job成功完成的实例数目的  当指定的实例数目达到Completions字段所设定的值时，Job将被标记为成功完成
	Succeeded   int32  `json:"succeeded"`   // 就绪的Job个数
	Duration    int64  `json:"duration"`    // Job 持续时间
	Age         int64  `json:"age"`
}
