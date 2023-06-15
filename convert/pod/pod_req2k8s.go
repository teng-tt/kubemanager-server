package pod

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	volume_empty = "emptyDir"
)

type PodConvert struct {
}

// PodReq2K8s 将 pod 的请求格式的数据 转换为 k8s 结构数据
func (p *PodConvert) PodReq2K8s(podReq pod_req.Pod) *corev1.Pod {
	labels := podReq.Base.Labels
	k8sLabels := p.getK8sLabels(labels)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podReq.Base.Name,
			Namespace: podReq.Base.NameSpace,
			Labels:    k8sLabels,
		},
		Spec: corev1.PodSpec{
			InitContainers: p.getK8sContainers(podReq.InitContainers),
			Containers:     p.getK8sContainers(podReq.Containers),
			Volumes:        p.getK8sVolumes(podReq.Volumes),
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: podReq.NetWorking.DnsConfig.Nameservers,
			},
			DNSPolicy:     corev1.DNSPolicy(podReq.NetWorking.DnsPolice),
			HostAliases:   p.getK8sHostAliases(podReq.NetWorking.HostAliases),
			Hostname:      podReq.NetWorking.HostName,
			RestartPolicy: corev1.RestartPolicy(podReq.Base.RestartPolicy),
		},
	}
}

func (p *PodConvert) getK8sHostAliases(podReqHostAliases []pod_req.ListMapItem) []corev1.HostAlias {
	podK8sHotsAliases := make([]corev1.HostAlias, 0)
	for _, item := range podReqHostAliases {
		podK8sHotsAliases = append(podK8sHotsAliases, corev1.HostAlias{
			IP:        item.Key,
			Hostnames: strings.Split(item.Value, ","),
		})
	}

	return podK8sHotsAliases
}

func (p *PodConvert) getK8sVolumes(podReqVilumes []pod_req.Volume) []corev1.Volume {
	podK8sVolumes := make([]corev1.Volume, 0)
	for _, volume := range podReqVilumes {
		if volume.Type != volume_empty {
			continue
		}
		source := corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		}
		podK8sVolumes = append(podK8sVolumes, corev1.Volume{
			Name:         volume.Name,
			VolumeSource: source,
		})
	}

	return podK8sVolumes

}

func (p *PodConvert) getK8sContainers(podReqContainers []pod_req.Container) []corev1.Container {
	podK8sContainers := make([]corev1.Container, 0)
	for _, item := range podReqContainers {
		podK8sContainers = append(podK8sContainers, p.getK8sContainer(item))

	}

	return podK8sContainers
}

func (p *PodConvert) getK8sContainer(podReqContainer pod_req.Container) corev1.Container {
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
		VolumeMounts:   p.getK8sVolumeMount(podReqContainer.VolumeMounts),
		StartupProbe:   p.getK8sContainerProbe(podReqContainer.StartProbe),
		LivenessProbe:  p.getK8sContainerProbe(podReqContainer.LivenessProbe),
		ReadinessProbe: p.getK8sContainerProbe(podReqContainer.ReadinessProbe),
		Resources:      p.getK8sResources(podReqContainer.Resources),
		Ports:          p.getK8sPort(podReqContainer.Port),
	}

	return k8sContainer
}

func (p *PodConvert) getK8sPort(podReqPort []pod_req.ContainerPort) []corev1.ContainerPort {
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

func (p *PodConvert) getK8sResources(podReqResources pod_req.Resources) corev1.ResourceRequirements {
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

func (p *PodConvert) getK8sContainerProbe(podReqProbe pod_req.ContainerProbe) *corev1.Probe {
	if podReqProbe.Enable {
		return nil
	}
	var k8sProbe corev1.Probe
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

func (p *PodConvert) getK8sVolumeMount(podReqMounts []pod_req.VolumeMount) []corev1.VolumeMount {
	podK8sVolumMounts := make([]corev1.VolumeMount, 0)
	for _, mount := range podReqMounts {
		podK8sVolumMounts = append(podK8sVolumMounts, corev1.VolumeMount{
			Name:      mount.Name,
			MountPath: mount.MountPath,
			ReadOnly:  mount.ReadOnly,
		})
	}

	return podK8sVolumMounts
}

func (p *PodConvert) getK8sEnv(podReqEnv []pod_req.ListMapItem) []corev1.EnvVar {
	podK8sEnvs := make([]corev1.EnvVar, 0)
	for _, item := range podReqEnv {
		podK8sEnvs = append(podK8sEnvs, corev1.EnvVar{
			Name:  item.Key,
			Value: item.Value,
		})
	}
	return podK8sEnvs
}

// Pod 请求 labels 转换为 k8s labels
func (p *PodConvert) getK8sLabels(podReqLabels []pod_req.ListMapItem) map[string]string {
	podK8sLabels := make(map[string]string)
	for _, label := range podReqLabels {
		podK8sLabels[label.Key] = label.Value
	}

	return podK8sLabels
}
