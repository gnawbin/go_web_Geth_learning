package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"hash/fnv"
	"k8s-web/global"
	metrics_k8s "k8s-web/model/metrics/k8s"
	metrics_res "k8s-web/model/metrics/response"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"time"
)

// @Author: morris
type MetricsService struct {
}

func (MetricsService) getMetricsFromProm(metricsName string) (string, error) {
	if !global.CONF.System.Prometheus.Enable {
		err := fmt.Errorf("prometheus未开启！")
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
	r := promv1.Range{
		Start: start,
		End:   end,
		Step:  5 * time.Minute,
	}
	queryRange, _, err := promApi.QueryRange(context.TODO(), metricsName, r)
	if err != nil {
		return "", err
	}
	matrix := queryRange.(model.Matrix)
	if len(matrix) == 0 {
		err = fmt.Errorf("prometheus查询数据为空！")
		return "", err
	}
	x := make([]string, 0)
	y := make([]string, 0)
	for _, value := range matrix[0].Values {
		format := value.Timestamp.Time().Format("15:04")
		x = append(x, format)
		y = append(y, value.Value.String())
	}
	resultMap["x"] = x
	resultMap["y"] = y
	raw, _ := json.Marshal(resultMap)
	return string(raw), nil
}

func (m MetricsService) GetClusterUsageRange() []metrics_res.MetricsItem {
	metricsItemList := make([]metrics_res.MetricsItem, 0)
	//去prometheus 查询数据
	promData, err := m.getMetricsFromProm("cluster_cpu")
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Title: "CPU变化趋势",
			Value: promData,
		})
	}
	promData, err = m.getMetricsFromProm("cluster_mem")
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Title: "内存变化趋势",
			Value: promData,
		})
	}
	return metricsItemList
}

func (MetricsService) GetClusterUsage() []metrics_res.MetricsItem {
	metricsItemList := make([]metrics_res.MetricsItem, 0)
	url := "/apis/metrics.k8s.io/v1beta1/nodes"
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	if err != nil {
		return metricsItemList
	}
	var nodeMetricsList metrics_k8s.NodeMetricsList
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
		//pod cpu mem 使用情况
		cpuUsage += nodeMetricsList.Items[i].Usage.Cpu().MilliValue() //!!!重点
		memUsage += nodeMetricsList.Items[i].Usage.Memory().Value()
		cpuTotal += item.Status.Capacity.Cpu().MilliValue()
		memTotal += item.Status.Capacity.Memory().Value()
		podTotal += item.Status.Capacity.Pods().Value()
	}
	//每一项使用的值和我们k8s值系统总的值除 得到 21.10%
	podUsageFormat := fmt.Sprintf("%.2f", (float64(podUsage)/float64(podTotal))*100)
	metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
		Value: podUsageFormat,
		Title: "Pod使用占比",
	})
	//每一项使用的值和我们k8s值系统总的值除 得到 21.10%
	cpuUsageFormat := fmt.Sprintf("%.2f", (float64(cpuUsage)/float64(cpuTotal))*100)
	metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
		Value: cpuUsageFormat,
		Label: "cluster_cpu",
		Title: "CPU使用占比",
	})
	//每一项使用的值和我们k8s值系统总的值除 得到 21.10%
	memUsageFormat := fmt.Sprintf("%.2f", (float64(memUsage)/float64(memTotal))*100)
	metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
		Value: memUsageFormat,
		Label: "cluster_mem",
		Title: "内存使用占比",
	})
	return metricsItemList
}

func (MetricsService) GetResource() []metrics_res.MetricsItem {
	metricsItemList := make([]metrics_res.MetricsItem, 0)
	ctx := context.TODO()
	//namespace
	list, err := global.KubeConfigSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(list.Items)),
			Logo:  "k8s",
			Title: "Namespaces",
		})
	}
	//pods
	podlist, err := global.KubeConfigSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(podlist.Items)),
			Logo:  "pod",
			Title: "Pods",
		})
	}

	cmlist, err := global.KubeConfigSet.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(cmlist.Items)),
			Logo:  "cm",
			Title: "ConfigMaps",
		})
	}

	//secret
	secretlist, err := global.KubeConfigSet.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(secretlist.Items)),
			Logo:  "secret",
			Title: "Secrets",
		})
	}

	//pv
	pvlist, err := global.KubeConfigSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(pvlist.Items)),
			Logo:  "pv",
			Title: "PV",
		})
	}

	//pvc
	pvclist, err := global.KubeConfigSet.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(pvclist.Items)),
			Logo:  "pvc",
			Title: "PVC",
		})
	}
	//sc
	sclist, err := global.KubeConfigSet.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(sclist.Items)),
			Logo:  "sc",
			Title: "StorageClass",
		})
	}

	//service
	servicelist, err := global.KubeConfigSet.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(servicelist.Items)),
			Logo:  "svc",
			Title: "Services",
		})
	}

	//ingesses
	ingresslist, err := global.KubeConfigSet.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(ingresslist.Items)),
			Logo:  "ingress",
			Title: "Ingresses",
		})
	}

	//deployment
	deploymentlist, err := global.KubeConfigSet.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(deploymentlist.Items)),
			Logo:  "pod",
			Title: "Deployments",
		})
	}

	//DaemonSets
	daemonsetslist, err := global.KubeConfigSet.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(daemonsetslist.Items)),
			Logo:  "pod",
			Title: "DaemonSets",
		})
	}

	//StatefulSets
	statefulSetslist, err := global.KubeConfigSet.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(statefulSetslist.Items)),
			Logo:  "ingress",
			Title: "StatefulSets",
		})
	}

	//Jobs
	joblist, err := global.KubeConfigSet.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(joblist.Items)),
			Logo:  "pod",
			Title: "Jobs",
		})
	}

	//CronJobs
	cronjobslist, err := global.KubeConfigSet.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(cronjobslist.Items)),
			Logo:  "pod",
			Title: "CronJobs",
		})
	}

	//ServiceAccounts
	salist, err := global.KubeConfigSet.CoreV1().ServiceAccounts("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(salist.Items)),
			Logo:  "secret",
			Title: "ServiceAccounts",
		})
	}

	//roles
	roleslist, err := global.KubeConfigSet.RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(roleslist.Items)),
			Logo:  "secret",
			Title: "Roles",
		})
	}

	//clusterrole
	clusterrolelist, err := global.KubeConfigSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(clusterrolelist.Items)),
			Logo:  "secret",
			Title: "ClusterRoles",
		})
	}

	//rolesbinding
	rolesbindinglist, err := global.KubeConfigSet.RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(rolesbindinglist.Items)),
			Logo:  "secret",
			Title: "RoleBindings",
		})
	}

	//clusterrole
	clusterrolebindinglist, err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err == nil {
		metricsItemList = append(metricsItemList, metrics_res.MetricsItem{
			Value: strconv.Itoa(len(clusterrolebindinglist.Items)),
			Logo:  "secret",
			Title: "CRBindings",
		})
	}

	for index, item := range metricsItemList {
		metricsItemList[index].Color = generateHashBasedRGB(item.Value)
	}
	return metricsItemList
}

// 获取集群信息
func (MetricsService) GetClusterInfo() []metrics_res.MetricsItem {
	metricsList := make([]metrics_res.MetricsItem, 0)
	//k8s类型
	metricsList = append(metricsList, metrics_res.MetricsItem{
		Title: "Cluster",
		Value: "K8S",
		Logo:  "k8s",
	})
	//k8s版本
	serverVersion, err := global.KubeConfigSet.ServerVersion()
	if err == nil {
		metricsList = append(metricsList, metrics_res.MetricsItem{
			Title: "Kubernetes Version",
			Value: fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor),
			Logo:  "k8s",
		})
	}
	list, err := global.KubeConfigSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	//k8s集群初始化时间
	if err == nil {
		var k8sCreateTime int64 = 0
		for _, item := range list.Items {
			if _, ok := item.Labels["node-role.kubernetes.io/control-plane"]; ok {
				//当k8sCreateTime未被赋值 直接赋值
				if k8sCreateTime == 0 {
					k8sCreateTime = item.CreationTimestamp.Unix()
				}
				//当找到一个节点的初始化时间更早 就依据当前的节点的初始化时间作为集群初始化时间
				if k8sCreateTime > 0 && item.CreationTimestamp.Unix() < k8sCreateTime {
					k8sCreateTime = item.CreationTimestamp.Unix()
				}
			}
		}
		formatTime := getFormatTimeByUnix(k8sCreateTime)
		metricsList = append(metricsList, metrics_res.MetricsItem{
			Title: "Created",
			Value: formatTime,
			Logo:  "k8s",
		})
	}

	//K8s node数量
	if err == nil {
		metricsList = append(metricsList, metrics_res.MetricsItem{
			Title: "Nodes",
			Value: strconv.Itoa(len(list.Items)),
			Logo:  "k8s",
		})
	}
	for i, item := range metricsList {
		metricsList[i].Color = generateHashBasedRGB(item.Title)
	}
	return metricsList
}

// 基于字符串的哈希码生成 RGB 字符串
func generateHashBasedRGB(str string) string {
	hash := hashString(str)    // 计算字符串的哈希码
	r, g, b := hashToRGB(hash) // 将哈希码转换为 RGB 分量

	return strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b)
}

// 计算字符串的哈希码
func hashString(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

// 将哈希码转换为 RGB 分量
func hashToRGB(hash uint32) (r, g, b int) {
	r = int(hash & 0xFF)         // 取低 8 位作为红色分量
	g = int((hash >> 8) & 0xFF)  // 取中间 8 位作为绿色分量
	b = int((hash >> 16) & 0xFF) // 取高 8 位作为蓝色分量
	return
}

func getFormatTimeByUnix(createTime int64) string {
	if createTime == 0 {
		return "Unknown"
	}
	//计算时间
	currentTime := time.Now()
	timestampTime := time.Unix(createTime, 0)
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
