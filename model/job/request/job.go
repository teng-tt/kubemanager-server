package request

import (
	"kubmanager/model/base"
	podReq "kubmanager/model/pod/request"
)

type JobBase struct {
	Name        string             `json:"name"`
	Namespace   string             `json:"namespace"`
	Labels      []base.ListMapItem `json:"labels"`
	Completions int32              `json:"completions"` // job的pod副本数，全部副本运行成功，才能代表job运行成功

}

type Job struct {
	Base     JobBase    `json:"base"`
	Template podReq.Pod `json:"template"`
}
