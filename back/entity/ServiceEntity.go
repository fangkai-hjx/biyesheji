package entity

type Service struct {
	Name            string  `json:"name"`
	ClusterIP       string  `json:"cluster_ip"`
	SessionAffinity bool  `json:"session_affinity"`
	Status          string  `json:"status"`
	Success         int     `json:"success"`
	Fail            int     `json:"fail"`
	Total           int     `json:"total"`
	SuccessLu       float64 `json:"success_lu"`
	Runtime         int     `json:"runtime"`
	Pod             []Pod   `json:"pod"`
}

type Pod struct {
	Name string `json:"name"`
	//ServiceName       string            `json:"service_name"`
	Initialized     bool   `json:"initialized"`
	Ready           bool   `json:"ready"`
	ContainersReady bool   `json:"containersReady"`
	PodScheduled    bool   `json:"podScheduled"`
	Image           string `json:"image"`
	HostIP          string `json:"host_ip"`
	PodIP           string `json:"pod_ip"`
}
type Image struct {
	ImageName    string `json:"image_name"`
	UpdateTime   string `json:"update_time"`
	CreateTime   string `json:"create_time"`
	ProjectId    int64  `json:"project_id"`
	Description  string `json:"description"`
	RepositoryId int64  `json:"repository_id"`
	PullCount    int64  `json:"pull_count"`
	Tag          string `json:"tag"`
}
type Namespace struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type RepositoryAddr struct {
	ImageName   string `json:"imageName"`
	ProjectName string `json:"projectName"`
}
type HarborRepository struct {
	ProjectId   int64  `json:"project_id"`
	ProjectName string `json:"project_name"`
	CreateTime  string `json:"creation_time"`
	UpdateTime  string `json:"create_time"`
	OwnerName   string `json:"owner_name"`
	RepoCount   int64  `json:repo_count`
	Images      []Image
}
