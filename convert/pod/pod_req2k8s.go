package pod

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kubmanager/model/base"
	pod_req "kubmanager/model/pod/request"
	"strconv"
	"strings"
)

const (
	probe_http = "http"
	probe_tcp  = "tcp"
	probe_exec = "exec"
)

const (
	volume_empty     = "emptyDir"
	volume_configMap = "configMap"
	volume_secret    = "secret"
	volume_hostPath  = "hostPath"
	volume_downward  = "downward"
	volume_pvc       = "pvc"
)

const (
	scheduling_nodename     = "nodeName"
	scheduling_nodselector  = "nodeSelector"
	scheduling_nodeaffinity = "nodeAffinity"
	scheduling_nodeany      = "nodeAny"
)

const (
	ref_type_configMap = "configMap"
	ref_type_secret    = "secret"
)

type Req2K8sConvert struct {
}

func (p *Req2K8sConvert) getNodeK8sScheduling(podReq pod_req.Pod) (affinity *corev1.Affinity, nodeSelector map[string]string, nodeName string) {
	nodeScheduling := podReq.NodeScheduling
	switch nodeScheduling.Type {
	case scheduling_nodename:
		nodeName = nodeScheduling.NodeName
	case scheduling_nodselector:
		nodeSelectoMap := make(map[string]string)
		for _, item := range nodeScheduling.NodeSelector {
			nodeSelectoMap[item.Key] = item.Value
		}
		nodeSelector = nodeSelectoMap
	case scheduling_nodeaffinity:
		nodeSelectorTermExpressions := nodeScheduling.NodeAffinity
		matchExpression := make([]corev1.NodeSelectorRequirement, 0)
		for _, expression := range nodeSelectorTermExpressions {
			matchExpression = append(matchExpression, corev1.NodeSelectorRequirement{
				Key:      expression.Key,
				Values:   strings.Split(expression.Value, ","),
				Operator: expression.Operator,
			})
		}
		affinity = &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: matchExpression,
						},
					},
				},
			},
		}
	case scheduling_nodeany:
		// do nothing
	default:
		// do nothing
	}
	return
}

// PodReq2K8s 将 pod 的请求格式的数据 转换为 k8s 结构数据
func (p *Req2K8sConvert) PodReq2K8s(podReq pod_req.Pod) *corev1.Pod {
	nodeAffinity, nodeSelector, nodeName := p.getNodeK8sScheduling(podReq)
	labels := podReq.Base.Labels
	k8sLabels := p.getK8sLabels(labels)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podReq.Base.Name,
			Namespace: podReq.Base.NameSpace,
			Labels:    k8sLabels,
		},
		Spec: corev1.PodSpec{
			NodeName:       nodeName,
			NodeSelector:   nodeSelector,
			Affinity:       nodeAffinity,
			Tolerations:    podReq.Tolerations,
			InitContainers: p.getK8sContainers(podReq.InitContainers),
			Containers:     p.getK8sContainers(podReq.Containers),
			Volumes:        p.getK8sVolumes(podReq.Volumes),
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: podReq.NetWorking.DnsConfig.Nameservers,
			},
			DNSPolicy:     corev1.DNSPolicy(podReq.NetWorking.DnsPolicy),
			HostAliases:   p.getK8sHostAliases(podReq.NetWorking.HostAliases),
			Hostname:      podReq.NetWorking.HostName,
			RestartPolicy: corev1.RestartPolicy(podReq.Base.RestartPolicy),
		},
	}
}

func (p *Req2K8sConvert) getK8sHostAliases(podReqHostAliases []base.ListMapItem) []corev1.HostAlias {
	podK8sHotsAliases := make([]corev1.HostAlias, 0)
	for _, item := range podReqHostAliases {
		podK8sHotsAliases = append(podK8sHotsAliases, corev1.HostAlias{
			IP:        item.Key,
			Hostnames: strings.Split(item.Value, ","),
		})
	}

	return podK8sHotsAliases
}

func (p *Req2K8sConvert) getK8sVolumes(podReqVilumes []pod_req.Volume) []corev1.Volume {
	podK8sVolumes := make([]corev1.Volume, 0)
	for _, volume := range podReqVilumes {
		//if volume.Type != volume_empty {
		//	continue
		//}
		source := corev1.VolumeSource{}
		switch volume.Type {
		case volume_empty:
			source = corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			}
		case volume_hostPath:
			source = corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Type: &volume.HostPathVolume.Type,
					Path: volume.HostPathVolume.Path,
				},
			}
		case volume_downward:
			items := make([]corev1.DownwardAPIVolumeFile, 0)
			for _, item := range volume.DownwardAPIVolume.Items {
				items = append(items, corev1.DownwardAPIVolumeFile{
					// 容器内的文件访问路径
					Path: item.Path,
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: item.FiledRefPath,
					},
				})
			}
			source = corev1.VolumeSource{
				DownwardAPI: &corev1.DownwardAPIVolumeSource{
					Items: items,
				},
			}
		case volume_pvc:
			source = corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volume.PVCVolume.Name,
				},
			}
		case volume_configMap:
			optional := volume.ConfiMapRefVolume.Optional
			source = corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: volume.ConfiMapRefVolume.Name,
					},
					Optional: &optional,
				},
			}
		case volume_secret:
			optional := volume.SecretRefVolume.Optional
			source = corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: volume.SecretRefVolume.Name,
					Optional:   &optional,
				},
			}
		default:
			continue
		}
		podK8sVolumes = append(podK8sVolumes, corev1.Volume{
			Name:         volume.Name,
			VolumeSource: source,
		})
	}

	return podK8sVolumes

}

func (p *Req2K8sConvert) getK8sContainers(podReqContainers []pod_req.Container) []corev1.Container {
	podK8sContainers := make([]corev1.Container, 0)
	for _, item := range podReqContainers {
		podK8sContainers = append(podK8sContainers, p.getK8sContainer(item))

	}

	return podK8sContainers
}

func (p *Req2K8sConvert) getK8sEnvFrom(podReqEnvsFrom []pod_req.EnvsVarFromResource) []corev1.EnvFromSource {
	podK8sEnvsFrom := make([]corev1.EnvFromSource, 0)
	for _, fromResource := range podReqEnvsFrom {
		// 前缀通用
		envFromResource := corev1.EnvFromSource{
			Prefix: fromResource.Prefix,
		}
		switch fromResource.RefType {
		case ref_type_configMap:
			envFromResource.ConfigMapRef = &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: fromResource.Name,
				},
			}
			podK8sEnvsFrom = append(podK8sEnvsFrom, envFromResource)
		case ref_type_secret:
			envFromResource.SecretRef = &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: fromResource.Name,
				},
			}
			podK8sEnvsFrom = append(podK8sEnvsFrom, envFromResource)
		}
	}
	return podK8sEnvsFrom
}

func (p *Req2K8sConvert) getK8sContainer(podReqContainer pod_req.Container) corev1.Container {
	k8sContainer := corev1.Container{
		Name:            podReqContainer.Name,
		Image:           podReqContainer.Image,
		ImagePullPolicy: corev1.PullPolicy(podReqContainer.ImagePullPolicy),
		TTY:             podReqContainer.Tty,
		Command:         podReqContainer.Command,
		Args:            podReqContainer.Args,
		WorkingDir:      podReqContainer.WorkingDir,
		SecurityContext: &corev1.SecurityContext{
			Privileged: &podReqContainer.Privileged,
		},
		Env:            p.getK8sEnv(podReqContainer.Envs),
		EnvFrom:        p.getK8sEnvFrom(podReqContainer.EnvsFrom),
		VolumeMounts:   p.getK8sVolumeMount(podReqContainer.VolumeMounts),
		StartupProbe:   p.getK8sContainerProbe(podReqContainer.StartupProbe),
		LivenessProbe:  p.getK8sContainerProbe(podReqContainer.LivenessProbe),
		ReadinessProbe: p.getK8sContainerProbe(podReqContainer.ReadinessProbe),
		Resources:      p.getK8sResources(podReqContainer.Resources),
		Ports:          p.getK8sPort(podReqContainer.Port),
	}

	return k8sContainer
}

func (p *Req2K8sConvert) getK8sPort(podReqPort []pod_req.ContainerPort) []corev1.ContainerPort {
	podK8sContainerPort := make([]corev1.ContainerPort, 0)
	for _, item := range podReqPort {
		podK8sContainerPort = append(podK8sContainerPort, corev1.ContainerPort{
			Name:          item.Name,
			HostPort:      item.HostPort,
			ContainerPort: item.ContainerPort,
		})
	}
	return podK8sContainerPort
}

func (p *Req2K8sConvert) getK8sResources(podReqResources pod_req.Resources) corev1.ResourceRequirements {
	var k8sPodResources corev1.ResourceRequirements
	if !podReqResources.Enable {
		return k8sPodResources
	}
	k8sPodResources.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(strconv.Itoa(int(podReqResources.CpuRequest)) + "m"),
		corev1.ResourceMemory: resource.MustParse(strconv.Itoa(int(podReqResources.MemRequest)) + "Mi"),
	}
	k8sPodResources.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(strconv.Itoa(int(podReqResources.CpuLimit)) + "m"),
		corev1.ResourceMemory: resource.MustParse(strconv.Itoa(int(podReqResources.MemLimit)) + "Mi"),
	}
	return k8sPodResources
}

func (p *Req2K8sConvert) getK8sContainerProbe(podReqProbe pod_req.ContainerProbe) *corev1.Probe {
	if podReqProbe.Enable {
		k8sProbe := corev1.Probe{
			InitialDelaySeconds: podReqProbe.InitialDelaySeconds,
			PeriodSeconds:       podReqProbe.PeriodSeconds,
			TimeoutSeconds:      podReqProbe.TimeOutSeconds,
			SuccessThreshold:    podReqProbe.SuccessThreshold,
			FailureThreshold:    podReqProbe.FailureThreshold,
		}
		switch podReqProbe.Type {
		case probe_http:
			httpGet := podReqProbe.HttpGet
			K8sHttpHeaders := make([]corev1.HTTPHeader, 0)
			for _, header := range httpGet.HttpHeaders {
				K8sHttpHeaders = append(K8sHttpHeaders, corev1.HTTPHeader{
					Name:  header.Key,
					Value: header.Value,
				})
			}
			k8sProbe.HTTPGet = &corev1.HTTPGetAction{
				Scheme:      corev1.URIScheme(httpGet.Scheme),
				Host:        httpGet.Host,
				Port:        intstr.FromInt(int(httpGet.Port)),
				Path:        httpGet.Path,
				HTTPHeaders: K8sHttpHeaders,
			}
		case probe_tcp:
			tcpSocket := podReqProbe.TcpSocket
			k8sProbe.TCPSocket = &corev1.TCPSocketAction{
				Host: tcpSocket.Host,
				Port: intstr.FromInt(int(tcpSocket.Port)),
			}
		case probe_exec:
			exec := podReqProbe.Exec
			k8sProbe.Exec = &corev1.ExecAction{
				Command: exec.Command,
			}
		}
		return &k8sProbe
	}
	return nil
}

func (p *Req2K8sConvert) getK8sVolumeMount(podReqMounts []pod_req.VolumeMount) []corev1.VolumeMount {
	podK8sVolumMounts := make([]corev1.VolumeMount, 0)
	for _, mount := range podReqMounts {
		podK8sVolumMounts = append(podK8sVolumMounts, corev1.VolumeMount{
			Name:      mount.MountName,
			MountPath: mount.MountPath,
			ReadOnly:  mount.ReadOnly,
		})
	}

	return podK8sVolumMounts
}

func (p *Req2K8sConvert) getK8sEnv(podReqEnv []pod_req.EnvVar) []corev1.EnvVar {
	podK8sEnvs := make([]corev1.EnvVar, 0)
	for _, item := range podReqEnv {
		envVar := corev1.EnvVar{
			Name: item.Name,
		}
		switch item.Type {
		case ref_type_configMap:
			envVar.ValueFrom = &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					Key: item.Value,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: item.Name,
					},
				},
			}
		case ref_type_secret:
			envVar.ValueFrom = &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: item.Value,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: item.Name,
					},
				},
			}
		default:
			envVar.Value = item.Value
		}
		podK8sEnvs = append(podK8sEnvs, envVar)
	}
	return podK8sEnvs
}

// Pod 请求 labels 转换为 k8s labels
func (p *Req2K8sConvert) getK8sLabels(podReqLabels []base.ListMapItem) map[string]string {
	podK8sLabels := make(map[string]string)
	for _, label := range podReqLabels {
		podK8sLabels[label.Key] = label.Value
	}

	return podK8sLabels
}
