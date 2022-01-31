package common

type Report struct {
	CountBacsiFrontend         int
	CountBacsiBackend          int
	CountSehatFrontend         int
	CountSehatBackend          int
	Count400mCpu               int
	Count200mCpu               int
	CountAll                   int
	DateTime                   string
	ClusterCPURequests         int64
	ClusterCPURequestsFraction string
	ClusterCPUAllocatable      int64
	ClusterCPUCapacity         int64
	CountNodes                 int
	CountNodePools             map[string]int
	UsdtPrice                  float64
}
