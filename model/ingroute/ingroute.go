package ingroute

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type IngressRouteSpec struct {
	EntryPoints []string `json:"entryPoints"`
	Routes      []struct {
		Kind     string `json:"kind"`
		Match    string `json:"match"`
		Services []struct {
			Name string `json:"name"`
			Port int32  `json:"port"`
		} `json:"services"`
	} `json:"routes"`
	// 配置TLS证书，指针类型如果为空 tls-nil ->k8s忽略
	Tls *struct {
		SecretName string `json:"secretName"`
	} `json:"tls"`
}
type IngressRoute struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
	Spec            IngressRouteSpec  `json:"spec"`
}
type IngressRouteList struct {
	Items           []IngressRoute `json:"items"`
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

type Middleware struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
}

type MiddlewareList struct {
	Items           []Middleware `json:"items"`
	metav1.TypeMeta `json:",inline"`
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}
