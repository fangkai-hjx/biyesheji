package entity

import (
	"strings"
)

/* ps: binding: "required" 好像没什么卵用 */
type Service struct {
	// service
	Id int64 `json:"id"`
	Name        string `json:"name" binding: "required"`
	// creator
	Creator     string `json:"creator"`
	Image       string `json:"image" binding: "required"`
	Ports       Port   `json:"ports"`
	Replicas    int32  `json:"replicas"`
	Description string `json:"description"`
	ServiceIp   string `json:"clusterIp"`
	CreatedAt   string `json:"createdAt"`
	Changer     string `json:"changer"`
	ChangedAt   string `json:"changedAt"`
}

type Port struct {
	ClusterPort int32  `json:"clusterPort"` // ClusterIP:ClusterPort
	TargetPort  int32  `json:"targetPort"`  // Pod/Container 暴露的端口
	NodePort    int32  `json:"nodePort"`    // NodeIP:NodePort
	Protocol    string `json:"protocol"`    // TCP/UDP/SCTP
}


type PublishServiceRequest struct {
	Service
}

//type ListServicesRequest struct {
//	ProjectId int64 `json:"projectId"`
//	PageSize  int64 `json:"pageSize"`
//	PageIndex int64 `json:"pageIndex"`
//}
//
//type GetServiceDetailRequest struct {
//	Id int64 `json:"id"`
//}
//
//type UpdateServiceRequest struct {
//	ProjectId int64  `json:"projectId"`
//	Id        int64  `json:"id"`
//	Name      string `json:"name"`
//	Image     string `json:"image"`
//	ChangerId int64  `json:"changerId"`
//}
//
//type DeleteServicesRequest struct {
//	ProjectId int64    `json:"projectId"`
//	Names     []string `json:"names"`
//	Ids       []int64  `json:"ids"`
//}

func (request *PublishServiceRequest) IsValid() bool {

	if strings.TrimSpace(request.Service.Name) == "" {
		return false
	}

	if request.Service.Image == "" || strings.TrimSpace(request.Service.Image) == "" {
		return false
	}

	return true
}

//func (request *ListServicesRequest) IsValid() bool {
//	if request.ProjectId == 0 {
//		return false
//	}
//
//	return true
//}
//
//func (request *GetServiceDetailRequest) IsValid() bool {
//	if request.Id <= 0 {
//		return false
//	}
//
//	return true
//}
//
//func (request *UpdateServiceRequest) IsValid() bool {
//	if request.ProjectId == 0 {
//		return false
//	}
//
//	if strings.TrimSpace(request.Name) == "" {
//		return false
//	}
//
//	if strings.TrimSpace(request.Image) == "" {
//		return false
//	}
//
//	return true
//}
//
//func (request *DeleteServicesRequest) IsValid() bool {
//	if request.ProjectId == 0 {
//		return false
//	}
//
//	for _, name := range request.Names {
//		if strings.TrimSpace(name) == "" {
//			return false
//		}
//	}
//
//	return true
//}
