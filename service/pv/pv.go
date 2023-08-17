package pv

import (
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	pv_req "kubemanager.com/model/pv/request"
	pv_resp "kubemanager.com/model/pv/response"
	"kubemanager.com/utils"
	"strconv"
	"strings"
)

type PvService struct {
}

func (p *PvService) CreatePV(pvReq pv_req.PersistentVolume) error {
	// 参数转换
	var volumeSource corev1.PersistentVolumeSource
	switch pvReq.VolumeSource.Type {
	case "nfs":
		volumeSource.NFS = &corev1.NFSVolumeSource{
			Server:   pvReq.VolumeSource.NfsVolumeSource.NfsServer,
			Path:     pvReq.VolumeSource.NfsVolumeSource.NfsPath,
			ReadOnly: pvReq.VolumeSource.NfsVolumeSource.NfsReadOnly,
		}
	default:
		return errors.New("不支持的存储卷类型")
	}
	pv := corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pvReq.Name,
			Labels: utils.ToMap(pvReq.Labels),
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: resource.MustParse(strconv.Itoa(int(pvReq.Capacity)) + "Mi"),
			},
			AccessModes:                   pvReq.AccessModes,
			PersistentVolumeReclaimPolicy: pvReq.ReClaimPolicy,
			PersistentVolumeSource:        volumeSource,
		},
	}
	ctx := context.TODO()
	_, err := global.KubeConfigSet.CoreV1().PersistentVolumes().Create(ctx, &pv, metav1.CreateOptions{})
	return err
}

func (p *PvService) DeletePV(_ string, name string) error {
	err := global.KubeConfigSet.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (p *PvService) GetPvList(keyword string) ([]pv_resp.PersistentVolume, error) {
	pvList, err := global.KubeConfigSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	pvResList := make([]pv_resp.PersistentVolume, 0)
	for _, item := range pvList.Items {
		// k8s -> response
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		claim := ""
		ref := item.Spec.ClaimRef
		if ref != nil {
			claim = ref.Name
		}
		pvRes := pv_resp.PersistentVolume{
			Name:             item.Name,
			Capacity:         int32(item.Spec.Capacity.Storage().Value() / (1024 * 1024)),
			Labels:           utils.ToList(item.Labels),
			AccessModes:      item.Spec.AccessModes,
			ReClaimPolicy:    item.Spec.PersistentVolumeReclaimPolicy,
			Status:           item.Status.Phase,
			Claim:            claim,
			StorageClassName: item.Spec.StorageClassName, // 当PV是通过SC创建时，会有该字段
			Age:              item.CreationTimestamp.UnixMilli(),
			Reason:           item.Status.Reason,
		}
		pvResList = append(pvResList, pvRes)
	}
	return pvResList, err
}
