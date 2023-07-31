package sc

import (
	"context"
	"fmt"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubmanager/global"
	scReq "kubmanager/model/sc/request"
	scResp "kubmanager/model/sc/response"
	"kubmanager/utils"
	"strings"
)

type SCService struct {
}

func (s *SCService) CreateSc(scReqs scReq.StorageClass) error {
	// 判断provisioner是否在系统支持列表
	provisionerList := strings.Split(global.CONF.System.Provisioner, ",")
	var flag bool
	for _, val := range provisionerList {
		if scReqs.Provisioner == val {
			flag = true
			break
		}
	}
	if !flag {
		err := fmt.Errorf("%s 当前K8S只支持未支持", scReqs.Provisioner)
		return err
	}

	sc := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   scReqs.Name,
			Labels: utils.ToMap(scReqs.Labels),
		},
		Provisioner:          scReqs.Provisioner,
		MountOptions:         scReqs.MountOptions,
		VolumeBindingMode:    &scReqs.VolumeBindIngMode,
		ReclaimPolicy:        &scReqs.ReclaimPolicy,
		AllowVolumeExpansion: &scReqs.AllowVolumeExpansion,
		Parameters:           utils.ToMap(scReqs.Parameters),
	}
	ctx := context.TODO()
	_, err := global.KubeConfigSet.StorageV1().StorageClasses().
		Create(ctx, &sc, metav1.CreateOptions{})
	return err
}

func (s *SCService) DeleteSc(name string) error {
	err := global.KubeConfigSet.StorageV1().StorageClasses().
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (s *SCService) GetScList(keyword string) ([]scResp.StorageClass, error) {
	scList, err := global.KubeConfigSet.StorageV1().StorageClasses().
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	scRespList := make([]scResp.StorageClass, 0)
	for _, item := range scList.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		// item -> response
		var allowVolumeExpansion bool
		if item.AllowVolumeExpansion != nil {
			allowVolumeExpansion = *item.AllowVolumeExpansion
		}
		mountOptions := make([]string, 0)
		if item.MountOptions != nil {
			mountOptions = item.MountOptions
		}
		scRespItme := scResp.StorageClass{
			Name:                 item.Name,
			Labels:               utils.ToList(item.Labels),
			Provisioner:          item.Provisioner,
			MountOptions:         mountOptions,
			Parameters:           utils.ToList(item.Parameters),
			ReclaimPolicy:        *item.ReclaimPolicy,
			AllowVolumeExpansion: allowVolumeExpansion,
			Age:                  item.CreationTimestamp.UnixMilli(),
			VolumeBindIngMode:    *item.VolumeBindingMode,
		}
		scRespList = append(scRespList, scRespItme)
	}
	return scRespList, err
}
