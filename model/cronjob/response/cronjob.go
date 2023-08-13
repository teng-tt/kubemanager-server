package response

type CronJob struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Schedule     string `json:"schedule"`     // CronJob的执行时间表
	Suspend      *bool  `json:"suspend"`      // 于暂停CronJob的调度。当设置为true时，CronJob将停止生成新的Job，但保留现有的Job
	Active       int    `json:"active"`       // CronJob当前处于活跃状态的Job数量
	LastSchedule int64  `json:"lastSchedule"` // 这个字段记录了CronJob上一次成功调度的时间
	Age          int64  `json:"age"`          // 创建时间
}
