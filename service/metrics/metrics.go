package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"hash/fnv"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubemanager.com/global"
	metricsK8s "kubemanager.com/model/metrics/k8s"
	metricsResp "kubemanager.com/model/metrics/response"
	"strconv"
	"time"
)

type MetricsService struct {
}

// 基于字符串的哈希码生成 RGB 字符串
func generateHashBaseRGB(str string) string {
	hash := hashString(str)    // 计算字符串的哈希码
	r, g, b := hashToRGB(hash) // 将哈希码转换为 RGB 分量
	return strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b)
}

// 计算字符串的哈希码
func hashString(str string) uint32 {
	h := fnv.New32()
	h.Write([]byte(str))
	return h.Sum32()
}

// 将哈希码转换为 RGB 分量
func hashToRGB(hash uint32) (r, g, b int) {
	r = int(hash & 0xFF)         // 取低 8 位作为红色分量
	g = int((hash >> 8) & 0xFF)  // 取中间 8 位作为绿色分量
	b = int((hash >> 16) & 0xFF) // 取最高 8 位作为蓝色分量
	return
}

func getFormatTimeByUnix(createTime int64) string {
	if createTime == 0 {
		return "Unknown"
	}
	// 计算时间
	currentTime := time.Now()
	timestampTime := time.Unix(createTime, 0)
	// 当前时间-创建时间 然后转换为天数得到创建天数， 然后算出创建 年、月、日
	days := int(currentTime.Sub(timestampTime).Hours() / 24)
	years := days / 365         // 计算年份
	remainingDays := days % 365 // 剩余的天数

	months := remainingDays / 30       // 计算月份
	remainingDays = remainingDays % 30 // 剩余的天数

	result := ""
	if years > 0 {
		result += fmt.Sprintf("%d年", years)
	}
	if months > 0 {
		result += fmt.Sprintf("%d月", months)
	}
	if remainingDays > 0 {
		result += fmt.Sprintf("%d天", remainingDays)
	}
	return result
}

// GetClusterInfo 获取集群信息
func (m *MetricsService) GetClusterInfo() []metricsResp.MetricsItem {
	metricsList := make([]metricsResp.MetricsItem, 0)

	// k8s类型
	metricsList = append(metricsList, metricsResp.MetricsItem{
		Title: "Cluster",
		Value: "k8s",
		Logo:  "k8s",
	})

	// k8s版本
	serverVersion, err := global.KubeConfigSet.ServerVersion()
	if err == nil {
		metricsList = append(metricsList, metricsResp.MetricsItem{
			Title: "Kubernetes Version",
			Value: fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor),
			Logo:  "k8s",
		})
	}
	list, err := global.KubeConfigSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	// k8s集群初始化时间
	if err == nil {
		var k8sCreateTime int64 = 0
		for _, item := range list.Items {
			if _, ok := item.Labels["node-role.kubernetes.io/control-plane"]; ok {
				// 当k8sCreateTime未被赋值 直接赋值
				if k8sCreateTime == 0 {
					k8sCreateTime = item.CreationTimestamp.Unix()
				}
				// 当找到一个节点的初始化时间更早 就依据当前的结点的初始化时间作为集群初始化时间
				if k8sCreateTime > 0 && item.CreationTimestamp.Unix() < k8sCreateTime {
					k8sCreateTime = item.CreationTimestamp.Unix()
				}
			}
		}
		formatTime := getFormatTimeByUnix(k8sCreateTime)
		metricsList = append(metricsList, metricsResp.MetricsItem{
			Title: "Created",
			Value: formatTime,
			Logo:  "k8s",
		})
	}

	// k8s node数量
	if err == nil {
		metricsList = append(metricsList, metricsResp.MetricsItem{
			Title: "Nodes",
			Value: strconv.Itoa(len(list.Items)),
			Logo:  "k8s",
		})
	}
	for index, item := range metricsList {
		metricsList[index].Color = generateHashBaseRGB(item.Title)
	}
	return metricsList
}

// GetResource 获取集群资源统计信息
func (m *MetricsService) GetResource() []metricsResp.MetricsItem {
	metricsItemList := make([]metricsResp.MetricsItem, 0)
	ctx := context.TODO()
	//namespace
	list, err := global.KubeConfigSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(list.Items)),
			Logo:  "k8s",
			Title: "Namespaces",
		})
	}
	//pods
	podList, err := global.KubeConfigSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(podList.Items)),
			Logo:  "pod",
			Title: "Pods",
		})
	}

	cmList, err := global.KubeConfigSet.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(cmList.Items)),
			Logo:  "cm",
			Title: "ConfigMaps",
		})
	}

	//secret
	secretList, err := global.KubeConfigSet.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(secretList.Items)),
			Logo:  "secret",
			Title: "Secrets",
		})
	}

	//pv
	pvList, err := global.KubeConfigSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(pvList.Items)),
			Logo:  "pv",
			Title: "PV",
		})
	}

	//pvc
	pvcList, err := global.KubeConfigSet.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(pvcList.Items)),
			Logo:  "pvc",
			Title: "PVC",
		})
	}
	//sc
	scList, err := global.KubeConfigSet.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(scList.Items)),
			Logo:  "sc",
			Title: "StorageClass",
		})
	}

	//service
	serviceList, err := global.KubeConfigSet.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(serviceList.Items)),
			Logo:  "svc",
			Title: "Services",
		})
	}

	//ingesses
	ingressList, err := global.KubeConfigSet.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(ingressList.Items)),
			Logo:  "ingress",
			Title: "Ingresses",
		})
	}

	//deployment
	deploymentList, err := global.KubeConfigSet.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(deploymentList.Items)),
			Logo:  "pod",
			Title: "Deployments",
		})
	}

	//DaemonSets
	daemonSetsList, err := global.KubeConfigSet.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(daemonSetsList.Items)),
			Logo:  "pod",
			Title: "DaemonSets",
		})
	}

	//StatefulSets
	statefulSetsList, err := global.KubeConfigSet.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(statefulSetsList.Items)),
			Logo:  "ingress",
			Title: "StatefulSets",
		})
	}

	//Jobs
	jobList, err := global.KubeConfigSet.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(jobList.Items)),
			Logo:  "pod",
			Title: "Jobs",
		})
	}

	//CronJobs
	cronJobsList, err := global.KubeConfigSet.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(cronJobsList.Items)),
			Logo:  "pod",
			Title: "CronJobs",
		})
	}

	//ServiceAccounts
	saList, err := global.KubeConfigSet.CoreV1().ServiceAccounts("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(saList.Items)),
			Logo:  "secret",
			Title: "ServiceAccounts",
		})
	}

	//roles
	rolesList, err := global.KubeConfigSet.RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(rolesList.Items)),
			Logo:  "secret",
			Title: "Roles",
		})
	}

	//clusterrole
	clusterRoleList, err := global.KubeConfigSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(clusterRoleList.Items)),
			Logo:  "secret",
			Title: "ClusterRoles",
		})
	}

	//rolesbinding
	rolesBindingList, err := global.KubeConfigSet.RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(rolesBindingList.Items)),
			Logo:  "secret",
			Title: "RoleBindings",
		})
	}

	//clusterrole
	clusterRoleBindingList, err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Value: strconv.Itoa(len(clusterRoleBindingList.Items)),
			Logo:  "secret",
			Title: "CRBindings",
		})
	}

	for index, item := range metricsItemList {
		metricsItemList[index].Color = generateHashBaseRGB(item.Value)
	}
	return metricsItemList
}

// GetClusterUsage 集群使用情况
func (m *MetricsService) GetClusterUsage() []metricsResp.MetricsItem {
	metricsItemList := make([]metricsResp.MetricsItem, 0)
	url := "/apis/metrics.k8s.io/v1beta1/nodes"
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	if err != nil {
		return metricsItemList
	}
	var nodeMetricsList metricsK8s.NodeMetricsList
	err = json.Unmarshal(raw, &nodeMetricsList)
	if err != nil {
		return metricsItemList
	}
	nodeList, err := global.KubeConfigSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return metricsItemList
	}
	if len(nodeList.Items) != len(nodeMetricsList.Items) {
		return metricsItemList
	}
	var cpuUsage, cpuTotal int64
	var memUsage, memTotal int64
	podList, err := global.KubeConfigSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return metricsItemList
	}
	var podUsage, podTotal int64 = int64(len(podList.Items)), 0
	for i, item := range nodeList.Items {
		// pod cpu mem 使用情况
		cpuUsage += nodeMetricsList.Items[i].Usage.Cpu().MilliValue()
		memUsage += nodeMetricsList.Items[i].Usage.Memory().Value()
		cpuTotal += item.Status.Capacity.Cpu().MilliValue()
		memTotal += item.Status.Capacity.Memory().Value()
		podTotal += item.Status.Capacity.Pods().Value()
	}
	// 每一项使用的值和我们k8s值系统总的值除 得到 pod
	podUsageFormat := fmt.Sprintf("%.2f", float64(podUsage)/float64(podTotal)*100)
	metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
		Value: podUsageFormat,
		Title: "Pod使用占比",
	})
	// 每一项使用的值和我们k8s值系统总的值除 得到 cpu
	cpuUsageFormat := fmt.Sprintf("%.2f", float64(cpuUsage)/float64(cpuTotal)*100)
	metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
		Value: cpuUsageFormat,
		Title: "CPU使用占比",
		Label: "cluster_cpu",
	})
	// 每一项使用的值和我们k8s值系统总的值除 得到 mem
	memUsageFormat := fmt.Sprintf("%.2f", float64(memUsage)/float64(memTotal)*100)
	metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
		Value: memUsageFormat,
		Title: "MEM使用占比",
		Label: "cluster_mem",
	})
	return metricsItemList
}

func (m *MetricsService) getMetricsFromPrometheus(metricsName string) (string, error) {
	if !global.CONF.System.Prometheus.Enable {
		err := fmt.Errorf("prometheus未开启")
		return "", err
	}
	resultMap := make(map[string][]string)
	scheme := global.CONF.System.Prometheus.Scheme
	host := global.CONF.System.Prometheus.Host
	addr := fmt.Sprintf("%s://%s", scheme, host)
	client, err := promapi.NewClient(promapi.Config{
		Address: addr,
	})
	if err != nil {
		return "", err
	}
	promApi := promv1.NewAPI(client)
	now := time.Now()
	start, end := now.Add(-time.Hour*24), now
	// 查询范围，过去24小时
	r := promv1.Range{
		Start: start,
		End:   end,
		Step:  5 * time.Minute, // 每5分钟一个点
	}
	queryRange, _, err := promApi.QueryRange(context.TODO(), metricsName, r)
	if err != nil {
		return "", err
	}
	matrix := queryRange.(model.Matrix)
	if len(matrix) == 0 {
		err = fmt.Errorf("prometheus查询数据为空")
		return "", err
	}
	x := make([]string, 0)
	y := make([]string, 0)
	for _, value := range matrix[0].Values {
		formatTime := value.Timestamp.Time().Format("15:04")
		x = append(x, formatTime)
		y = append(y, value.Value.String())
	}
	resultMap["x"] = x
	resultMap["y"] = y
	marsha1, _ := json.Marshal(resultMap)
	return string(marsha1), nil

}

func (m *MetricsService) GetClusterUsageRange() []metricsResp.MetricsItem {
	metricsItemList := make([]metricsResp.MetricsItem, 0)
	// 去prometheus 查询数据
	promData, err := m.getMetricsFromPrometheus("cluster_cpu")
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Title: "CPU变化趋势",
			Value: promData,
		})
	}
	promData, err = m.getMetricsFromPrometheus("cluster_mem")
	if err == nil {
		metricsItemList = append(metricsItemList, metricsResp.MetricsItem{
			Title: "内存变化趋势",
			Value: promData,
		})
	}
	return metricsItemList
}
