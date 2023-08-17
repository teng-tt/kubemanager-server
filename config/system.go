package config

type Harbor struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	Host       string `json:"host" yaml:"host"`
	Scheme     string `json:"scheme" yaml:"scheme"`
	Enable     bool   `json:"enable" yaml:"enable"`
	CacertPath string `json:"cacertPath" yaml:"cacertPath"`
}

type K8sConfig struct {
	KubeConfig string `json:"kubeConfig" yaml:"kubeConfig"` // k8s集群权限，使用配置文件双向认证
	TokenFile  string `json:"tokenFile" yaml:"tokenFile"`   // k8s集群权限，使用token认证
	Host       string `json:"host" yaml:"host"`
	CacertPath string `json:"cacertPath" yaml:"cacertPath"`
}

type Prometheus struct {
	Host   string `json:"host" yaml:"host"`
	Scheme string `json:"scheme" yaml:"scheme"`
	Enable bool   `json:"enable" yaml:"enable"`
}

type System struct {
	Addr        string     `json:"addr" yaml:"addr"`
	K8sConfig   K8sConfig  `yaml:"k8SConfig" yaml:"k8SConfig"`
	Provisioner string     `json:"provisioner"`
	Harbor      Harbor     `json:"harbor" yaml:"harbor"`
	Prometheus  Prometheus `json:"prometheus" yaml:"prometheus"`
}
