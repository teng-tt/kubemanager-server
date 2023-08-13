package request

import (
	batchv1 "k8s.io/api/batch/v1"
	"kubmanager/model/base"
	podReq "kubmanager/model/pod/request"
)

type CronJobBase struct {
	Name                       string                    `json:"name"`
	Namespace                  string                    `json:"namespace"`
	Labels                     []base.ListMapItem        `json:"labels"`
	Schedule                   string                    `json:"schedule"`          // cron 表达式
	Suspend                    bool                      `json:"suspend"`           // 是否暂停cronjob
	ConcurrencyPolicy          batchv1.ConcurrencyPolicy `json:"concurrencyPolicy"` // 并发策略
	SuccessfulJobsHistoryLimit int32                     `json:"successfulJobsHistoryLimit"`
	FailedJobsHistoryLimit     int32                     `json:"failedJobsHistoryLimit"`
	Selector                   []base.ListMapItem        `json:"selector"`
	JobBase                    JobBase                   `json:"jobBase"`
}

type JobBase struct {
	Completions  int32 `json:"completions"`
	BackoffLimit int32 `json:"backoffLimit"`
}

type CronJob struct {
	Base     CronJobBase `json:"base"`
	Template podReq.Pod  `json:"template"`
}
