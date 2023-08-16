package config

type Harbor struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	Host       string `json:"host" yaml:"host"`
	Scheme     string `json:"scheme" yaml:"scheme"`
	Enable     bool   `json:"enable" yaml:"enable"`
	CacertPath string `json:"cacertPath" yaml:"cacertPath"`
}
type System struct {
	Addr        string `json:"addr" yaml:"addr"`
	Provisioner string `json:"provisioner"`
	Harbor      Harbor `json:"harbor" yaml:"harbor"`
}
