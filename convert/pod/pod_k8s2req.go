package pod

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	pod_req "kubmanager/model/pod/request"
	pod_res "kubmanager/model/pod/response"
	"strings"
)

const (
	volume_type_emptydir = "emptyDir"
)

type K8s2RqeConver struct {
}

func (k *K8s2RqeConver) PodK8s2Req(podK8s corev1.Pod) pod_req.Pod {
	return pod_req.Pod{
		Base:           k.getReqBase(podK8s),
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

func (k *K8s2RqeConver) getReqContainerEnvs(envsK8s []corev1.EnvVar) []pod_req.ListMapItem {
	envsReq := make([]pod_req.ListMapItem, 0)
	for _, item := range envsK8s {
		envsReq = append(envsReq, pod_req.ListMapItem{
			Key:   item.Name,
			Value: item.Value,
		})
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
		volumesReq = append(volumesReq, pod_req.VolumeMount{
			MountName: item.Name,
			MountPath: item.MountPath,
			ReadOnly:  item.ReadOnly,
		})
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
			headersReq := make([]pod_req.ListMapItem, 0)
			for _, headerK8s := range httpGet.HTTPHeaders {
				headersReq = append(headersReq, pod_req.ListMapItem{
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
	for _, volume := range volumes {
		if volume.EmptyDir == nil {
			continue
		}
		volumesReq = append(volumesReq, pod_req.Volume{
			Type: volume_type_emptydir,
			Name: volume.Name,
		})
	}
	return volumesReq
}

func (k *K8s2RqeConver) getReqHostAliases(hostAlias []corev1.HostAlias) []pod_req.ListMapItem {
	hostAliasReq := make([]pod_req.ListMapItem, 0)
	for _, alias := range hostAlias {
		hostAliasReq = append(hostAliasReq, pod_req.ListMapItem{
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

func (k *K8s2RqeConver) getReqLabels(data map[string]string) []pod_req.ListMapItem {
	labels := make([]pod_req.ListMapItem, 0)
	for k, v := range data {
		labels = append(labels, pod_req.ListMapItem{
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
