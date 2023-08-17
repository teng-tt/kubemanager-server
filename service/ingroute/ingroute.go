package ingroute

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	ingroutes "kubemanager.com/model/ingroute"
	ingrouteReq "kubemanager.com/model/ingroute/request"
	ingrouteResp "kubemanager.com/model/ingroute/response"
	"kubemanager.com/utils"
	"strings"
)

type IngressRouteService struct {
}

func (i *IngressRouteService) CreateOrUpdateIngRoute(ingressRouteReq ingrouteReq.IngressRouteReq) error {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", ingressRouteReq.Namespace)
	// convert to k8s structure
	ingRoute := ingroutes.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "traefik.io/v1alpha1",
			Kind:       "IngressRoute",
		},
		Metadata: metav1.ObjectMeta{
			Name:      ingressRouteReq.Name,
			Namespace: ingressRouteReq.Namespace,
			Labels:    utils.ToMap(ingressRouteReq.Labels),
		},
		Spec: ingressRouteReq.IngressRouteSpec,
	}
	// 已存在则更新
	result, err := json.Marshal(ingRoute)
	if err != nil {
		return err
	}
	// 查询是否存在
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).Name(ingressRouteReq.Name).DoRaw(context.TODO())
	if err == nil {
		// 修改row
		var ingressRouteK8s ingroutes.IngressRoute
		err = json.Unmarshal(raw, &ingressRouteK8s)
		if err != nil {
			return err
		}
		// update
		ingressRouteK8s.Spec = ingRoute.Spec
		resultx, errMar := json.Marshal(ingressRouteK8s)
		if errMar != nil {
			return errMar
		}
		_, err = global.KubeConfigSet.RESTClient().Put().
			Name(ingressRouteK8s.Metadata.Name).
			AbsPath(url).
			Body(resultx).DoRaw(context.TODO())
	} else {
		// create
		_, err = global.KubeConfigSet.RESTClient().Post().AbsPath(url).Body(result).DoRaw(context.TODO())
	}
	return err
}

func (i *IngressRouteService) DeleteIngRoute(namespace, name string) error {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes/%s", namespace, name)
	_, err := global.KubeConfigSet.RESTClient().Delete().AbsPath(url).DoRaw(context.TODO())
	return err
}

func (i *IngressRouteService) GetIngRouteDetail(namespace, name string) (ingrouteReqs ingrouteReq.IngressRouteReq, err error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes/%s", namespace, name)
	// 查出得到K8S的结构 ——> response/request
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).
		DoRaw(context.TODO())
	if err != nil {
		return
	}
	var ingressRoute ingroutes.IngressRoute
	err = json.Unmarshal(raw, &ingressRoute)
	if err != nil {
		return
	}
	ingRtReq := ingrouteReq.IngressRouteReq{
		Name:             ingressRoute.Metadata.Name,
		Namespace:        ingressRoute.Metadata.Namespace,
		Labels:           utils.ToList(ingressRoute.Metadata.Labels),
		IngressRouteSpec: ingressRoute.Spec,
	}
	return ingRtReq, err
}

// GetIngRouteList https://github.com/kubernetes-client/python/blob/master/kubernetes/README.md
func (i *IngressRouteService) GetIngRouteList(namespace, keyword string) ([]ingrouteResp.IngressRouteRes, error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", namespace)
	// 查出得到K8S的结构 ——> response/request
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).
		DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}
	var ingressRouteList ingroutes.IngressRouteList
	err = json.Unmarshal(raw, &ingressRouteList)
	if err != nil {
		return nil, err
	}
	ingRtResp := make([]ingrouteResp.IngressRouteRes, 0)
	for _, item := range ingressRouteList.Items {
		if !strings.Contains(item.Metadata.Name, keyword) {
			continue
		}
		ingRtResp = append(ingRtResp, ingrouteResp.IngressRouteRes{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Age:       item.Metadata.CreationTimestamp.UnixMilli(),
		})
	}
	return ingRtResp, err
}

func (i *IngressRouteService) GetIngRouteMiddlewareList(namespace string) (mwList []string, err error) {
	//查询middleware 列表
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/middlewares", namespace)
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	mwList = make([]string, 0)
	var middlewareList ingroutes.MiddlewareList
	err = json.Unmarshal(raw, &middlewareList)
	if err != nil {
		return
	}
	for _, item := range middlewareList.Items {
		mwList = append(mwList, item.Metadata.Name)

	}
	return
}
