package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubmanager/global"
	"kubmanager/model/base"
	"kubmanager/response"
	"kubmanager/utils"
	"strings"
)

type IngressRouteApi struct {
}

func (i *IngressRouteApi) CreateOrUpdateIngRoute(c *gin.Context) {

}

func (i *IngressRouteApi) DeleteIngRoute(c *gin.Context) {}

// GetIngRouteDetailOrList https://github.com/kubernetes-client/python/blob/master/kubernetes/README.md
func (i *IngressRouteApi) GetIngRouteDetailOrList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	name := c.Query("name")

	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", namespace)
	// 查出得到K8S的结构 ——> response/request
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
		Tls struct {
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
	type IngressRouteReq struct {
		Name      string             `json:"name"`
		Namespace string             `json:"namespace"`
		Labels    []base.ListMapItem `json:"labels"`
		IngressRouteSpec
	}

	type IngressRouteRes struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Age       int64  `json:"age"`
	}
	if name == "" {
		raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).
			DoRaw(context.TODO())
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		var ingressRouteList IngressRouteList
		err = json.Unmarshal(raw, &ingressRouteList)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		ingRtResp := make([]IngressRouteRes, 0)
		for _, item := range ingressRouteList.Items {
			if !strings.Contains(item.Metadata.Name, keyword) {
				continue
			}
			ingRtResp = append(ingRtResp, IngressRouteRes{
				Name:      item.Metadata.Name,
				Namespace: item.Metadata.Namespace,
				Age:       item.Metadata.CreationTimestamp.UnixMilli(),
			})
		}
		response.SuccessWithDetailed(c, "获取IngressRoute列表成功!", ingRtResp)
	} else {
		url = url + "/" + name
		raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).
			DoRaw(context.TODO())
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		var ingressRoute IngressRoute
		err = json.Unmarshal(raw, &ingressRoute)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		ingRtReq := IngressRouteReq{
			Name:             ingressRoute.Metadata.Name,
			Namespace:        ingressRoute.Metadata.Namespace,
			Labels:           utils.ToList(ingressRoute.Metadata.Labels),
			IngressRouteSpec: ingressRoute.Spec,
		}
		response.SuccessWithDetailed(c, "获取Ingress详情成功!", ingRtReq)
	}

}

func (i *IngressRouteApi) GetIngRouteMiddlewareList(c *gin.Context) {}
