package request

import (
	corev1 "k8s.io/api/core/v1"
	"kubmanager/model/base"
)

type Base struct {
	Name          string             `json:"name"`          // 名称
	Labels        []base.ListMapItem `json:"labels"`        // 标签
	NameSpace     string             `json:"nameSpace"`     // 命名空间
	RestartPolicy string             `json:"restartPolicy"` // 重启策略 Always | Never | On-Failure
}

type Volume struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// hostNetwork: false
// # 可选值：Default|ClusterFirst|ClusterFirstWithHostNet|None
// dnsPolicy: "Default"
// # dns配置
// dnsConfig:
//
//	nameservers:
//	- 8.8.8.8
//
// # 域名映射
// hostAliases:
//   - ip: 192.168.1.18
//     hostnames:
//   - "foo.local"
//   - "bar.local"
type DnsConfig struct {
	Nameservers []string `json:"nameservers"`
}

type NetWorking struct {
	HostNetwork bool               `json:"hostNetwork"`
	HostName    string             `json:"hostName"`
	DnsPolicy   string             `json:"dnsPolicy"`
	DnsConfig   DnsConfig          `json:"dnsConfig"`
	HostAliases []base.ListMapItem `json:"hostAliases"`
}

type Resources struct {
	Enable     bool  `json:"enable"`     // 是否配置容器配额
	MemRequest int32 `json:"memRequest"` // 内存 Mi
	MemLimit   int32 `json:"memLimit"`
	CpuRequest int32 `json:"cpuRequest"` // cpu m
	CpuLimit   int32 `json:"cpuLimit"`
}

type VolumeMount struct {
	MountName string `json:"mountName"` // 挂载卷名称
	MountPath string `json:"mountPath"` // 挂载卷->对应的容器内的路径
	ReadOnly  bool   `json:"readOnly"`  // 是否只读
}

type ProbeHttpGet struct {
	Scheme      string             `json:"scheme"`      // 请求协议 http|https
	Host        string             `json:"host"`        // 请求主机，如果为空，那么就是Pod内请求
	Path        string             `json:"path"`        // 请求路径
	Port        int32              `json:"port"`        // 请求端口
	HttpHeaders []base.ListMapItem `json:"httpHeaders"` // 请求头
}

type ProbeCommand struct {
	Command []string `json:"command"` // 执行命令： cat /test/test.txt
}

type ProbeTcpSocket struct {
	Host string `json:"host"` // 请求主机，如果为空，那么就是Pod内请求
	Port int32  `json:"port"` // 请求端口
}

type ProbeTime struct {
	InitialDelaySeconds int32 `json:"initialDelaySeconds"` // 初始化时间 初始化若干秒之后才开始探针
	PeriodSeconds       int32 `json:"periodSeconds"`       // 每隔若干秒之后去探针
	TimeOutSeconds      int32 `json:"timeOutSeconds"`      // 探针等待时间 等待若干秒之后还没有返回，那么就是探测失败
	SuccessThreshold    int32 `json:"successThreshold"`    // 探针若干次成功了 才认为这次探测成功
	FailureThreshold    int32 `json:"failureThreshold"`    // 探针若干次失败了 才认为这次探测失败
}

type ContainerProbe struct {
	Enable    bool           `json:"enable"` // 是否打开探针
	Type      string         `json:"type"`   // 探针类型 http/tcp/exec
	HttpGet   ProbeHttpGet   `json:"httpGet"`
	Exec      ProbeCommand   `json:"exec"`
	TcpSocket ProbeTcpSocket `json:"tcpSocket"`
	ProbeTime
}

type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
	HostPort      int32  `json:"hostPort"`
}
type Container struct {
	Name            string             `json:"name"`            // 容器的名称
	Image           string             `json:"image"`           // 容器点镜像
	ImagePullPolicy string             `json:"imagePullPolicy"` // 镜像的拉取策略
	Tty             bool               `json:"tty"`             // 是否开启伪终端
	Port            []ContainerPort    `json:"port"`            // 映射端口
	WorkingDir      string             `json:"workingDir"`      // 工作目录
	Command         []string           `json:"command"`         // 执行命令
	Args            []string           `json:"args"`            // 命令行参数
	Envs            []base.ListMapItem `json:"envs"`            // 环境变量
	Privileged      bool               `json:"privileged"`      // 是否开启特权模式
	Resources       Resources          `json:"resources"`       // 容器申请配额
	VolumeMounts    []VolumeMount      `json:"volumeMounts"`    // 容器挂载卷
	StartupProbe    ContainerProbe     `json:"startupProbe"`    // 启动探针
	LivenessProbe   ContainerProbe     `json:"livenessProbe"`   // 存活探针
	ReadinessProbe  ContainerProbe     `json:"readinessProbe"`  // 就绪探针
}

type NodeSelectorTermExpressions struct {
	Key      string                      `json:"key"`
	Operator corev1.NodeSelectorOperator `json:"operator"`
	Value    string                      `json:"value"`
}

type NodeScheduling struct {
	Type         string                        `json:"type"` // 调度类型 nodeName | nodeSelector | nodeAffinity
	NodeName     string                        `json:"nodeName"`
	NodeSelector []base.ListMapItem            `json:"nodeSelector"`
	NodeAffinity []NodeSelectorTermExpressions `json:"nodeAffinity"`
}
type Pod struct {
	Base           Base                `json:"base"`           // 基础定义信息
	Tolerations    []corev1.Toleration `json:"tolerations"`    // 容忍
	NodeScheduling NodeScheduling      `json:"nodeScheduling"` // 容器调度方式
	Volumes        []Volume            `json:"volumes"`        // 卷
	NetWorking     NetWorking          `json:"netWorking"`     // 网络相关
	InitContainers []Container         `json:"initContainers"` // init containers,一定不以守护进程方式运行，用于容器数据配置初始化
	Containers     []Container         `json:"containers"`     // containers 守护进程方式运行、也可不以守护进程方式运行
}
