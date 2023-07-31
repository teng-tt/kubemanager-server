package pod

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"kubmanager/model/base"
	pod_req "kubmanager/model/pod/request"
	pod_res "kubmanager/model/pod/response"
	"strings"
)

const (
	volume_type_emptydir = "emptyDir"
)

type K8s2RqeConver struct {
	volumeMap map[string]string
}

func (k *K8s2RqeConver) getNodeReqScheduling(podK8s corev1.Pod) pod_req.NodeScheduling {
	nodeScheduling := pod_req.NodeScheduling{
		Type: scheduling_nodeany,
	}
	if podK8s.Spec.NodeSelector != nil {
		nodeScheduling.Type = scheduling_nodselector
		labels := make([]base.ListMapItem, 0)
		for k, v := range podK8s.Spec.NodeSelector {
			labels = append(labels, base.ListMapItem{
				Key:   k,
				Value: v,
			})
		}
		nodeScheduling.NodeSelector = labels
		return nodeScheduling
	}
	if podK8s.Spec.Affinity != nil {
		nodeScheduling.Type = scheduling_nodeaffinity
		term := podK8s.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0]
		matchExpressions := make([]pod_req.NodeSelectorTermExpressions, 0)
		for _, expression := range term.MatchExpressions {
			matchExpressions = append(matchExpressions, pod_req.NodeSelectorTermExpressions{
				Key:      expression.Key,
				Value:    strings.Join(expression.Values, ","),
				Operator: expression.Operator,
			})
		}
		nodeScheduling.NodeAffinity = matchExpressions
		return nodeScheduling
	}
	if podK8s.Spec.NodeName != "" {
		nodeScheduling.Type = scheduling_nodename
		nodeScheduling.NodeName = podK8s.Spec.NodeName
		return nodeScheduling
	}
	return nodeScheduling
}

func (k *K8s2RqeConver) PodK8s2Req(podK8s corev1.Pod) pod_req.Pod {
	return pod_req.Pod{
		Base:           k.getReqBase(podK8s),
		Tolerations:    podK8s.Spec.Tolerations,
		NodeScheduling: k.getNodeReqScheduling(podK8s),
		NetWorking:     k.getReqNetworking(podK8s),
		Volumes:        k.getReqVolumes(podK8s.Spec.Volumes),
		Containers:     k.getReqContainers(podK8s.Spec.Containers),
		InitContainers: k.getReqContainers(podK8s.Spec.InitContainers),
	}
}

func (k *K8s2RqeConver) getReqContainers(containersK8s []corev1.Container) []pod_req.Container {
	podReqContainers := make([]pod_req.Container, 0)
	for _, item := range containersK8s {
		// container转换
		reqContainer := k.getReqContainer(item)
		podReqContainers = append(podReqContainers, reqContainer)
	}
	return podReqContainers
}

func (k *K8s2RqeConver) getReqContainerPorts(portsK8s []corev1.ContainerPort) []pod_req.ContainerPort {
	portsReq := make([]pod_req.ContainerPort, 0)
	for _, item := range portsK8s {
		portsReq = append(portsReq, pod_req.ContainerPort{
			Name:          item.Name,
			ContainerPort: item.ContainerPort,
			HostPort:      item.HostPort,
		})
	}

	return portsReq
}

func (k *K8s2RqeConver) getReqContainerEnvs(envsK8s []corev1.EnvVar) []pod_req.EnvVar {
	envsReq := make([]pod_req.EnvVar, 0)
	for _, item := range envsK8s {
		envVar := pod_req.EnvVar{
			Name: item.Name,
		}
		if item.ValueFrom != nil {
			if item.ValueFrom.ConfigMapKeyRef != nil {
				envVar.Type = ref_type_configMap
				envVar.Value = item.ValueFrom.ConfigMapKeyRef.Key
				envVar.RefName = item.ValueFrom.ConfigMapKeyRef.Name
			}
			if item.ValueFrom.SecretKeyRef != nil {
				envVar.Type = ref_type_secret
				envVar.Value = item.ValueFrom.SecretKeyRef.Key
				envVar.RefName = item.ValueFrom.SecretKeyRef.Name
			}
		} else {
			envVar.Type = "default"
			envVar.Value = item.Value
		}
		envsReq = append(envsReq, envVar)
	}
	return envsReq
}

func (k *K8s2RqeConver) getReqContainerPrivileged(ctx *corev1.SecurityContext) (privileged bool) {
	if ctx != nil {
		privileged = *ctx.Privileged
	}
	return privileged
}

func (k *K8s2RqeConver) getReqContainerResources(requirements corev1.ResourceRequirements) pod_req.Resources {
	reqResources := pod_req.Resources{
		Enable: false,
	}
	requests := requirements.Requests
	limits := requirements.Limits
	if requests != nil {
		reqResources.Enable = true
		reqResources.CpuRequest = int32(requests.Cpu().MilliValue()) // 单位 m
		reqResources.MemRequest = int32(requests.Memory().Value())   // 单位MiB 1024*1024 bytes
	}
	if limits != nil {
		reqResources.Enable = true
		reqResources.CpuLimit = int32(limits.Cpu().MilliValue())
		reqResources.MemLimit = int32(limits.Memory().Value())
	}

	return reqResources
}

func (k *K8s2RqeConver) getReqContainerVolumeMounts(volumeMountsK8s []corev1.VolumeMount) []pod_req.VolumeMount {
	volumesReq := make([]pod_req.VolumeMount, 0)
	for _, item := range volumeMountsK8s {
		// 非 emptyDir 过滤
		_, ok := k.volumeMap[item.Name]
		if ok {
			volumesReq = append(volumesReq, pod_req.VolumeMount{
				MountName: item.Name,
				MountPath: item.MountPath,
				ReadOnly:  item.ReadOnly,
			})
		}
	}
	return volumesReq
}

func (k *K8s2RqeConver) getReqContainerProbe(probeK8s *corev1.Probe) pod_req.ContainerProbe {
	containerProbe := pod_req.ContainerProbe{
		Enable: false,
	}
	// 先判断探针是否为空
	if probeK8s != nil {
		containerProbe.Enable = true
		// 在判断探针具体是什么类型 exec | http | tcp
		if probeK8s.Exec != nil {
			containerProbe.Type = probe_exec
			containerProbe.Exec.Command = probeK8s.Exec.Command
		} else if probeK8s.HTTPGet != nil {
			containerProbe.Type = probe_http
			httpGet := probeK8s.HTTPGet
			headersReq := make([]base.ListMapItem, 0)
			for _, headerK8s := range httpGet.HTTPHeaders {
				headersReq = append(headersReq, base.ListMapItem{
					Key:   headerK8s.Name,
					Value: headerK8s.Value,
				})
			}
			containerProbe.HttpGet = pod_req.ProbeHttpGet{
				Scheme:      string(httpGet.Scheme),
				Host:        httpGet.Host,
				Path:        httpGet.Path,
				Port:        httpGet.Port.IntVal,
				HttpHeaders: headersReq,
			}
		} else if probeK8s.TCPSocket != nil {
			containerProbe.Type = probe_tcp
			containerProbe.TcpSocket = pod_req.ProbeTcpSocket{
				Host: probeK8s.TCPSocket.Host,
				Port: probeK8s.TCPSocket.Port.IntVal,
			}
		} else {
			// 不支持探针，特殊处理
			containerProbe.Type = probe_http
			return containerProbe
		}
		containerProbe.InitialDelaySeconds = probeK8s.InitialDelaySeconds
		containerProbe.PeriodSeconds = probeK8s.PeriodSeconds
		containerProbe.TimeOutSeconds = probeK8s.TimeoutSeconds
		containerProbe.SuccessThreshold = probeK8s.SuccessThreshold
		containerProbe.FailureThreshold = probeK8s.FailureThreshold
	}
	return containerProbe
}

func (k *K8s2RqeConver) getReqContainerEnvsFrom(envsFromK8s []corev1.EnvFromSource) []pod_req.EnvsVarFromResource {
	podReqEnvsFromList := make([]pod_req.EnvsVarFromResource, 0)
	for _, envK8sItem := range envsFromK8s {
		// 前缀通用
		podReqEnvsFrom := pod_req.EnvsVarFromResource{
			Prefix: envK8sItem.Prefix,
		}
		if envK8sItem.ConfigMapRef != nil {
			podReqEnvsFrom.RefType = ref_type_configMap
			podReqEnvsFrom.Name = envK8sItem.ConfigMapRef.Name
		}
		if envK8sItem.SecretRef != nil {
			podReqEnvsFrom.RefType = ref_type_secret
			podReqEnvsFrom.Name = envK8sItem.SecretRef.Name
		}
		podReqEnvsFromList = append(podReqEnvsFromList, podReqEnvsFrom)
	}
	return podReqEnvsFromList
}

func (k *K8s2RqeConver) getReqContainer(containerK8s corev1.Container) pod_req.Container {
	return pod_req.Container{
		Name:            containerK8s.Name,
		Image:           containerK8s.Image,
		ImagePullPolicy: string(containerK8s.ImagePullPolicy),
		Tty:             containerK8s.TTY,
		Port:            k.getReqContainerPorts(containerK8s.Ports),
		WorkingDir:      containerK8s.WorkingDir,
		Command:         containerK8s.Command,
		Args:            containerK8s.Args,
		Envs:            k.getReqContainerEnvs(containerK8s.Env),
		EnvsFrom:        k.getReqContainerEnvsFrom(containerK8s.EnvFrom),
		Privileged:      k.getReqContainerPrivileged(containerK8s.SecurityContext),
		Resources:       k.getReqContainerResources(containerK8s.Resources),
		VolumeMounts:    k.getReqContainerVolumeMounts(containerK8s.VolumeMounts),
		StartupProbe:    k.getReqContainerProbe(containerK8s.StartupProbe),
		LivenessProbe:   k.getReqContainerProbe(containerK8s.LivenessProbe),
		ReadinessProbe:  k.getReqContainerProbe(containerK8s.ReadinessProbe),
	}
}

func (k *K8s2RqeConver) getReqVolumes(volumes []corev1.Volume) []pod_req.Volume {
	volumesReq := make([]pod_req.Volume, 0)
	if k.volumeMap == nil {
		k.volumeMap = make(map[string]string)
	}
	for _, volume := range volumes {
		//if volume.EmptyDir == nil {
		//	continue
		//}
		var volumeReq *pod_req.Volume

		if volume.EmptyDir != nil {
			volumeReq = &pod_req.Volume{
				Type: volume_empty,
				Name: volume.Name,
			}
		}
		if volume.ConfigMap != nil {
			var optional bool
			if volume.ConfigMap.Optional != nil {
				optional = *volume.ConfigMap.Optional
			}
			volumeReq = &pod_req.Volume{
				Type: volume_configMap,
				Name: volume.Name,
				ConfiMapRefVolume: pod_req.ConfiMapRefVolume{
					Name:     volume.ConfigMap.Name,
					Optional: optional,
				},
			}
		}
		if volume.Secret != nil {
			var optional bool
			if volume.Secret.Optional != nil {
				optional = *volume.Secret.Optional
			}
			volumeReq = &pod_req.Volume{
				Type: volume_secret,
				Name: volume.Name,
				SecretRefVolume: pod_req.SecretRefVolume{
					Name:     volume.Secret.SecretName,
					Optional: optional,
				},
			}
		}
		if volume.HostPath != nil {
			volumeReq = &pod_req.Volume{
				Name: volume.Name,
				Type: volume_hostPath,
				HostPathVolume: pod_req.HostPathVolume{
					Path: volume.HostPath.Path,
					Type: *volume.HostPath.Type,
				},
			}
		}
		if volume.PersistentVolumeClaim != nil {
			volumeReq = &pod_req.Volume{
				Name: volume.Name,
				Type: volume_pvc,
				PVCVolume: pod_req.PVCVolume{
					Name: volume.PersistentVolumeClaim.ClaimName,
				},
			}
		}
		if volume.DownwardAPI != nil {
			items := make([]pod_req.DownwardAPIVolumeItem, 0)
			for _, item := range volume.DownwardAPI.Items {
				items = append(items, pod_req.DownwardAPIVolumeItem{
					Path:         item.Path,
					FiledRefPath: item.FieldRef.FieldPath,
				})
			}
			volumeReq = &pod_req.Volume{
				Type: volume_downward,
				Name: volume.Name,
				DownwardAPIVolume: pod_req.DownwardAPIVolume{
					Items: items,
				},
			}
		}
		if volumeReq == nil {
			continue
		}
		k.volumeMap[volume.Name] = ""
		volumesReq = append(volumesReq, *volumeReq)
	}
	return volumesReq
}

func (k *K8s2RqeConver) getReqHostAliases(hostAlias []corev1.HostAlias) []base.ListMapItem {
	hostAliasReq := make([]base.ListMapItem, 0)
	for _, alias := range hostAlias {
		hostAliasReq = append(hostAliasReq, base.ListMapItem{
			Key:   alias.IP,
			Value: strings.Join(alias.Hostnames, ","),
		})
	}

	return hostAliasReq
}

func (k *K8s2RqeConver) getReqDnsConfig(dnsConfigK8s *corev1.PodDNSConfig) pod_req.DnsConfig {
	var dnsConfigReq pod_req.DnsConfig
	if dnsConfigK8s != nil {
		dnsConfigReq.Nameservers = dnsConfigK8s.Nameservers
	}
	return dnsConfigReq
}

func (k *K8s2RqeConver) getReqNetworking(pod corev1.Pod) pod_req.NetWorking {
	return pod_req.NetWorking{
		HostNetwork: pod.Spec.HostNetwork,
		HostName:    pod.Spec.Hostname,
		DnsPolicy:   string(pod.Spec.DNSPolicy),
		DnsConfig:   k.getReqDnsConfig(pod.Spec.DNSConfig),
		HostAliases: k.getReqHostAliases(pod.Spec.HostAliases),
	}
}

func (k *K8s2RqeConver) getReqLabels(data map[string]string) []base.ListMapItem {
	labels := make([]base.ListMapItem, 0)
	for k, v := range data {
		labels = append(labels, base.ListMapItem{
			Key:   k,
			Value: v,
		})
	}

	return labels
}

func (k *K8s2RqeConver) getReqBase(pod corev1.Pod) pod_req.Base {
	return pod_req.Base{
		Name:          pod.Name,
		NameSpace:     pod.Namespace,
		Labels:        k.getReqLabels(pod.Labels),
		RestartPolicy: string(pod.Spec.RestartPolicy),
	}
}

func (k *K8s2RqeConver) PodK8s2ItemRes(pod corev1.Pod) pod_res.PodListItem {
	var totalC, readyC, restartC int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			readyC++
		}
		restartC += containerStatus.RestartCount
		totalC++
	}
	var podStatus string
	if pod.Status.Phase != "Running" {
		podStatus = "Error"
	} else {
		podStatus = "Running"
	}
	return pod_res.PodListItem{
		Name:    pod.Name,
		Ready:   fmt.Sprintf("%d/%d", readyC, totalC),
		Status:  podStatus,
		Restart: restartC,
		Age:     pod.CreationTimestamp.Unix(),
		Ip:      pod.Status.PodIP,
		Node:    pod.Spec.Hostname,
	}
}
