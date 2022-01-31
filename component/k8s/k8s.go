package k8s

import (
	"bot/common"
	"bytes"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"text/template"
	"time"
)

// NodeAllocatedResources Refer: https://github.com/kubernetes/dashboard/blob/d3b1a14cd87791c0e1b217669846738d53229027/src/app/backend/resource/node/detail.go#L37
type NodeAllocatedResources struct {
	// CPURequests is number of allocated milicores.
	CPURequests int64 `json:"cpuRequests"`

	// CPURequestsFraction is a fraction of CPU, that is allocated.
	CPURequestsFraction float64 `json:"cpuRequestsFraction"`

	// CPULimits is defined CPU limit.
	CPULimits int64 `json:"cpuLimits"`

	// CPULimitsFraction is a fraction of defined CPU limit, can be over 100%, i.e.
	// overcommitted.
	CPULimitsFraction float64 `json:"cpuLimitsFraction"`

	// CPUCapacity is specified node CPU capacity in milicores.
	CPUCapacity int64 `json:"cpuCapacity"`

	// MemoryRequests is a fraction of memory, that is allocated.
	MemoryRequests int64 `json:"memoryRequests"`

	// MemoryRequestsFraction is a fraction of memory, that is allocated.
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`

	// MemoryLimits is defined memory limit.
	MemoryLimits int64 `json:"memoryLimits"`

	// MemoryLimitsFraction is a fraction of defined memory limit, can be over 100%, i.e.
	// overcommitted.
	MemoryLimitsFraction float64 `json:"memoryLimitsFraction"`

	// MemoryCapacity is specified node memory capacity in bytes.
	MemoryCapacity int64 `json:"memoryCapacity"`

	// AllocatedPods in number of currently allocated pods on the node.
	AllocatedPods int `json:"allocatedPods"`

	// PodCapacity is maximum number of pods, that can be allocated on the node.
	PodCapacity int64 `json:"podCapacity"`

	// PodFraction is a fraction of pods, that can be allocated on given node.
	PodFraction float64 `json:"podFraction"`
}

type ClusterAllocatedResources struct {
	// CPURequests is number of allocated milicores.
	CPURequests int64 `json:"cpuRequests"`

	// CPURequestsFraction is a fraction of CPU, that is allocated.
	CPURequestsFraction float64 `json:"cpuRequestsFraction"`

	// CPULimits is defined CPU limit.
	CPULimits int64 `json:"cpuLimits"`

	// CPULimitsFraction is a fraction of defined CPU limit, can be over 100%, i.e.
	// overcommitted.
	CPULimitsFraction float64 `json:"cpuLimitsFraction"`

	// CPUCapacity is specified node CPU capacity in milicores.
	CPUCapacity int64 `json:"cpuCapacity"`

	// MemoryRequests is a fraction of memory, that is allocated.
	MemoryRequests int64 `json:"memoryRequests"`

	// MemoryRequestsFraction is a fraction of memory, that is allocated.
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`

	// MemoryLimits is defined memory limit.
	MemoryLimits int64 `json:"memoryLimits"`

	// MemoryLimitsFraction is a fraction of defined memory limit, can be over 100%, i.e.
	// overcommitted.
	MemoryLimitsFraction float64 `json:"memoryLimitsFraction"`

	// MemoryCapacity is specified node memory capacity in bytes.
	MemoryCapacity int64 `json:"memoryCapacity"`
}

//type Report struct {
//	CountBacsiFrontend         int
//	CountBacsiBackend          int
//	CountSehatFrontend         int
//	CountSehatBackend          int
//	Count400mCpu               int
//	Count200mCpu               int
//	CountAll                   int
//	DateTime                   string
//	ClusterCPURequests         int64
//	ClusterCPURequestsFraction string
//	ClusterCPUAllocatable      int64
//	ClusterCPUCapacity         int64
//	CountNodes                 int
//	CountNodePools             map[string]int
//}

func Authenticate(kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	// Authenticate from inside of k8s
	config, err = rest.InClusterConfig()
	if err != nil {
		// Authenticate from outside of k8s
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	log.Println("Kubernetes authenticate successfully")
	return clientSet, nil
}

func Analyze(clientSet *kubernetes.Clientset) (string, error) {
	log.Println("analyzing kubenetes...")
	var err error
	var tmpl *template.Template
	tmpl = template.Must(template.ParseFiles("templates/k8s-counting-pods.txt"))

	var report common.Report
	//init the loc
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	//set timezone
	now := time.Now().In(loc)
	report.DateTime = now.Format("02-01-2006 15:04:05 MST")

	if err = countPods(clientSet, &report); err != nil {
		return "", err
	}

	if err = getClusterAllocatedResources(clientSet, &report); err != nil {
		return "", err
	}

	var output bytes.Buffer
	// Execute template
	if err = tmpl.Execute(&output, report); err != nil {
		return "", err
	}
	return output.String(), nil
}

func listNamespaces(clientSet *kubernetes.Clientset) {
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))
	for i := 0; i < len(namespaces.Items); i++ {
		fmt.Println(namespaces.Items[i].Name)
	}
}

func listNodes(clientSet *kubernetes.Clientset) (*v1.NodeList, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// https://github.com/kubernetes/dashboard/blob/d3b1a14cd87791c0e1b217669846738d53229027/src/app/backend/resource/node/detail.go#L171
func getNodeAllocatedResources(clientSet *kubernetes.Clientset, node *v1.Node) (NodeAllocatedResources, error) {
	pods, err := listRunningPodsOnNode(clientSet, node.Name)
	if err != nil {
		return NodeAllocatedResources{}, err
	}

	var nodeCpuRequests int64
	for _, pod := range pods.Items {
		var podCpuRequest int64
		for _, container := range pod.Spec.Containers {
			podCpuRequest += container.Resources.Requests.Cpu().MilliValue()
		}
		nodeCpuRequests += podCpuRequest
	}

	var nodeCpuRequestsFraction float64
	if capacity := float64(node.Status.Allocatable.Cpu().MilliValue()); capacity > 0 {
		nodeCpuRequestsFraction = float64(nodeCpuRequests) / capacity * 100
	}

	return NodeAllocatedResources{
		CPURequests:         nodeCpuRequests,
		CPURequestsFraction: nodeCpuRequestsFraction,
	}, nil
}

// namespace = ""       -> all namespaces
// labelSelector = ""   -> don't use label selector
func listPods(clientSet *kubernetes.Clientset, namespace string, labelSelector string) (*v1.PodList, error) {
	pods, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func listRunningPods(clientSet *kubernetes.Clientset, namespace string, labelSelector string) (*v1.PodList, error) {
	pods, err := listPods(clientSet, namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	var result v1.PodList
	for _, pod := range pods.Items {
		if pod.Status.Phase != v1.PodFailed && pod.Status.Phase != v1.PodSucceeded {
			result.Items = append(result.Items, pod)
		}
	}

	return &result, nil
}

func listFailedPods(clientSet *kubernetes.Clientset, namespace string, labelSelector string) (*v1.PodList, error) {
	pods, err := listPods(clientSet, namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	var result v1.PodList
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodFailed {
			result.Items = append(result.Items, pod)
		}
	}

	return &result, nil
}

// cpuRequest: 200m, 400m, etc.
func listPodsWithCpuRequest(clientSet *kubernetes.Clientset, namespace string, labelSelector string, cpuRequest string) (*v1.PodList, error) {
	pods, err := listPods(clientSet, namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	var result v1.PodList
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			for _, container := range pod.Spec.Containers {
				if container.Resources.Requests.Cpu().String() == cpuRequest {
					result.Items = append(result.Items, pod)
					break
				}
			}
		}

	}

	return &result, nil
}

func listRunningPodsOnNode(clientSet *kubernetes.Clientset, nodeName string) (*v1.PodList, error) {
	pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, err
	}

	var result v1.PodList
	for _, pod := range pods.Items {
		if pod.Status.Phase != v1.PodFailed && pod.Status.Phase != v1.PodSucceeded {
			result.Items = append(result.Items, pod)
		}
	}
	return &result, nil
}

// namespace = ""       -> all namespaces
// labelSelector = ""   -> don't use label selector
func countRunningPods(clientSet *kubernetes.Clientset, namespace string, labelSelector string) (int, error) {
	pods, err := listRunningPods(clientSet, namespace, labelSelector)
	if err != nil {
		return 0, err
	}

	return len(pods.Items), nil
}

func countFailedPods(clientSet *kubernetes.Clientset, namespace string, labelSelector string) (int, error) {
	pods, err := listFailedPods(clientSet, namespace, labelSelector)
	if err != nil {
		return 0, err
	}

	return len(pods.Items), nil
}

// cpuRequest: 200m, 400m, etc.
func countPodsWithCpuRequest(clientSet *kubernetes.Clientset, namespace string, labelSelector string, cpuRequest string) (int, error) {
	pods, err := listPodsWithCpuRequest(clientSet, namespace, labelSelector, cpuRequest)
	if err != nil {
		return 0, err
	}

	return len(pods.Items), nil
}

func countPods(clientSet *kubernetes.Clientset, report *common.Report) error {
	var err error

	report.CountBacsiFrontend, err = countRunningPods(clientSet, "production", "app=discover-fe-bacsi")
	if err != nil {
		return err
	}

	report.CountBacsiBackend, err = countRunningPods(clientSet, "production", "app=discover-be-bacsi")
	if err != nil {
		return err
	}

	report.CountSehatFrontend, err = countRunningPods(clientSet, "production", "app=discover-fe-sehat")
	if err != nil {
		return err
	}

	report.CountSehatBackend, err = countRunningPods(clientSet, "production", "app=discover-be-sehat")
	if err != nil {
		return err
	}

	report.Count400mCpu, err = countPodsWithCpuRequest(clientSet, "", "", "400m")
	if err != nil {
		return err
	}

	report.Count200mCpu, err = countPodsWithCpuRequest(clientSet, "", "", "200m")
	if err != nil {
		return err
	}

	report.CountAll, err = countRunningPods(clientSet, "", "")
	if err != nil {
		return err
	}

	return nil
}

func getClusterAllocatedResources(clientSet *kubernetes.Clientset, report *common.Report) error {
	nodes, err := listNodes(clientSet)
	if err != nil {
		return err
	}

	var clusterAllocatedResources ClusterAllocatedResources
	var clusterAllocatableCPU int64
	report.CountNodePools = make(map[string]int)

	for _, node := range nodes.Items {
		nodeAllocatedResources, err := getNodeAllocatedResources(clientSet, &node)
		if err != nil {
			return err
		}
		clusterAllocatedResources.CPURequests += nodeAllocatedResources.CPURequests

		clusterAllocatableCPU += node.Status.Allocatable.Cpu().MilliValue()
		clusterAllocatedResources.CPUCapacity += node.Status.Capacity.Cpu().MilliValue()

		report.CountNodePools[node.Labels["cloud.google.com/gke-nodepool"]]++
	}

	var clusterCpuRequestsFraction float64
	capacity := float64(clusterAllocatableCPU)
	if capacity > 0 {
		clusterCpuRequestsFraction = float64(clusterAllocatedResources.CPURequests) / capacity * 100
	}

	// Copy result to Report
	report.ClusterCPURequests = clusterAllocatedResources.CPURequests
	report.ClusterCPURequestsFraction = fmt.Sprintf("%.0f", clusterCpuRequestsFraction)
	report.ClusterCPUCapacity = clusterAllocatedResources.CPUCapacity
	report.ClusterCPUAllocatable = clusterAllocatableCPU

	report.CountNodes = len(nodes.Items)

	return nil
}
