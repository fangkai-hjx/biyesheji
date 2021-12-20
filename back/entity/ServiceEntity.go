package entity

type Service struct {
	Name            string `json:"name"`
	ClusterIP       string `json:"cluster_ip"`
	SessionAffinity string `json:"session_affinity"`
	Status          string `json:"status"`
	SuccessLu		float32 `json:"success_lu"`
	Pod             []Pod  `json:"pod"`
}
type ServiceConditions struct {
	Initialized     string `json:"initialized"`
	Ready           string `json:"ready"`
	ContainersReady string `json:"containersReady"`
	PodScheduled    string `json:"podScheduled"`
}
type Pod struct {
	Name              string            `json:"name"`
	ServiceName       string            `json:"service_name"`
	ServiceConditions ServiceConditions `json:"service_conditions"`
	Image             string            `json:"image"`
}
type Image struct {
	ImageName    string   `json:"image_name"`
	UpdateTime   string   `json:"update_time"`
	CreateTime   string   `json:"create_time"`
	ProjectId    int64    `json:"project_id"`
	Description  string   `json:"description"`
	RepositoryId int64    `json:"repository_id"`
	PullCount    int64    `json:"pull_count"`
	Tag          []string `json:"tag"`
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
