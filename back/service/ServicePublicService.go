package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"gee/util"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"gee/entity"
	//"gee/repository"
	//"service-manager/util"
	"strings"
)

type ApplicationService struct {
	k8sClient    *kubernetes.Clientset
	//appDao       *repository.ApplicationDao
	//flavorDao    *repository.FlavorDao
	//projectDao   *repository.ProjectDao
	//projectCache map[int64]*entity.ProjectInfoEntity
}

/* 创建 ServiceService 实例 */
func NewServiceService() *ApplicationService {
	cache := make(map[int64]*entity.ProjectInfoEntity, 10)
	return &ApplicationService{
		k8sClient:    util.K8sClient,
		//appDao:       repository.NewApplicationDao(),
		//flavorDao:    repository.NewFlavorDao(),
		//projectDao:   repository.NewProjectDao(),
		//projectCache: cache,
	}
}


/* 根据项目ID获取对应的 K8s 命名空间 */
func (service *ApplicationService) GetNamespace(projectId int64) (string, error) {
	service.getProjectAsNeeded(projectId)
	if _, ok := service.projectCache[projectId]; !ok {
		return "", errors.New("not found")
	}
	return service.projectCache[projectId].NameSpace, nil
}

/* 根据项目ID获取对应的镜像仓库地址 */
func (service *ApplicationService) GetImageRepository(image,version string) (string, error) {
	project, err := service.imageDao.GetOneImage(image,version)
	if err != nil {
		return "",errors.New("not found")
	}
	return service.projectCache[projectId].HarborAddress, nil
}

/* 查看资源库存是否充足 */
func (service *ApplicationService) CheckResourceInventory(svc *entity.Service) bool {
	return true
}

/* 发布服务 */
func (service *ApplicationService) PublishService(svc *entity.Service) error {
	namespace, err := service.GetNamespace(svc.ProjectId)
	if err != nil {
		return err
	}

	//flavor, err := service.flavorDao.GetFlavor(svc.FlavorId)
	//if err != nil {
	//	return err
	//}

	err = service.createK8sDeployment(namespace, flavor, svc)
	if err != nil {
		return err
	}

	k8sService, err := service.createK8sService(namespace, svc)
	if err != nil {
		service.deleteK8sDeployment(namespace, svc.Name)
		return err
	}

	svc.ClusterIp = k8sService.Spec.ClusterIP

	err = service.appDao.AddApplication(svc)
	if err != nil {
		return err
	}

	return nil
}

/* 获取服务列表 */
func (service *ApplicationService) ListServices(projectId, pageSize, pageIndex int64) ([]entity.Service, int64, error) {
	total, err := service.appDao.GetApplicationCount(projectId)
	if err != nil {
		return nil, -1, err
	}

	offset := pageSize * pageIndex

	apps, err := service.appDao.ListApplications(projectId, offset, pageSize)

	for i, _ := range apps {
		apps[i].Creator = "admin"
		apps[i].Changer = "admin"
	}

	return apps, total, err

	// 直接从 k8s 集群获取信息
	//namespace, err := GetNamespace(projectId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//deployments, err := service.listK8sDeployments(namespace)
	//if err != nil {
	//	return nil, err
	//}
	//
	//services, err := service.listK8sServices(namespace)
	//if err != nil {
	//	return nil, err
	//}
	//
	//serviceBriefs := make([]entity.Service, len(services.Items))
	//for i, svc := range services.Items {
	//	serviceBrief := entity.Service{}
	//	serviceBrief.Name = svc.Name
	//	serviceBrief.ClusterIp = svc.Spec.ClusterIP
	//	serviceBrief.CreatedAt = svc.CreationTimestamp.Format("2006-01-02 15:04:05")
	//
	//	for _, deployment := range deployments.Items {
	//		if deployment.Name == svc.Name {
	//			serviceBrief.Image = deployment.Spec.Template.Spec.Containers[0].Image
	//			break
	//		}
	//	}
	//
	//	serviceBrief.Creator = "admin"
	//	serviceBriefs[i] = serviceBrief
	//}
	//
	//return serviceBriefs, nil
}

/* 获取服务详情 */
func (service *ApplicationService) GetServiceDetail(appId int64) (*entity.Service, error) {
	app, err := service.appDao.GetApplication(appId)

	app.Creator = "admin"
	app.Changer = "admin"

	return app, err

	// 直接从 k8s 集群获取信息
	//namespace, err := GetNamespace(projectId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//deployment, err := service.getK8sDeployment(namespace, name)
	//if err != nil {
	//	return nil, err
	//}
	//
	//k8sService, err := service.getK8sService(namespace, name)
	//if err != nil {
	//	return nil, err
	//}
	//
	//serviceDetail := &entity.Service{}
	//serviceDetail.Name = name
	//serviceDetail.Image = deployment.Spec.Template.Spec.Containers[0].Image
	//serviceDetail.Replicas = *deployment.Spec.Replicas
	//serviceDetail.ClusterIp = k8sService.Spec.ClusterIP
	//serviceDetail.CreatedAt = k8sService.CreationTimestamp.Format("2006-01-02 15:04:05")
	//
	//ports := make([]entity.Port, len(k8sService.Spec.Ports))
	//for i, port := range k8sService.Spec.Ports {
	//	ports[i].ClusterPort = port.Port
	//	ports[i].NodePort = port.NodePort
	//	ports[i].TargetPort = port.TargetPort.IntVal
	//	ports[i].Protocol = string(port.Protocol)
	//}
	//serviceDetail.Ports = ports
	//
	//envVars := deployment.Spec.Template.Spec.Containers[0].Env
	//envs := make(map[string]string, len(envVars))
	//for _, env := range envVars {
	//	envs[env.Name] = env.Value
	//}
	//serviceDetail.Env = envs
	//
	//serviceDetail.Description = "service description"
	//serviceDetail.FlavorId = 1
	//serviceDetail.ProjectId = projectId
	//
	//return serviceDetail, nil
}

/* 更新服务 */
func (service *ApplicationService) UpdateService(request *entity.UpdateServiceRequest) error {
	namespace, err := service.GetNamespace(request.ProjectId)
	if err != nil {
		return err
	}

	err = service.patchK8sDeployment(namespace, request.Name, request.Image)
	if err != nil {
		return err
	}

	err = service.appDao.UpdateApplication(request.Id, request.Image, request.ChangerId)
	if err != nil {
		return err
	}

	return nil
}

/* 删除服务 */
func (service *ApplicationService) DeleteServices(projectId int64, names []string, ids []int64) error {
	namespace, err := service.GetNamespace(projectId)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = service.deleteK8sService(namespace, name)
		if err != nil {
			return err
		}

		err = service.deleteK8sDeployment(namespace, name)
		if err != nil {
			return err
		}
	}

	err = service.appDao.DeleteApplications(ids)
	if err != nil {
		return err
	}

	return nil
}

/* 删除命名空间 */
func (service *ApplicationService) DeleteNamespace(projectId int64) error {
	namespace, err := service.GetNamespace(projectId)
	if err != nil {
		return err
	}

	err = service.deleteK8sNamespace(namespace)
	return err
}

/* 删除 Registry Secret */
func (service *ApplicationService) DeleteRegistrySecret(projectId int64, name string) error {
	namespace, err := service.GetNamespace(projectId)
	if err != nil {
		return err
	}

	err = service.deleteK8sRegistrySecret(namespace, name)
	return err
}

/*************************************************项目管理的职责****************************************************/
/* 创建 命名空间 */
func (service *ApplicationService) CreateNamespace(name string) error {
	err := service.createK8sNamespace(name)
	return err
}

/* 创建 Registry Secret */
func (service *ApplicationService) CreateRegistrySecret(projectId int64, name, server, user, password string) error {
	namespace, err := service.GetNamespace(projectId)
	if err != nil {
		return err
	}

	err = service.createK8sRegistrySecret(namespace, name, server, user, password)
	return err
}

/* 创建 K8s 命名空间 */
func (service *ApplicationService) createK8sNamespace(name string) error {
	namespace := &corev1.Namespace{}
	namespace.Name = name

	_, err := service.k8sClient.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	return err
}

/* 创建 K8s Registry Secret */
func (service *ApplicationService) createK8sRegistrySecret(namespace, name, server, user, password string) error {
	data := make(map[string]string)

	secret := &corev1.Secret{}
	secret.Name = name
	secret.Type = corev1.SecretTypeDockerConfigJson

	content := `{
					"auths": {
						"%s": {
							"auth": "%s"
						}
					},
					"HttpHeaders": {
						"User-Agent": "Docker-Client/19.03.11 (linux)"
					}
				}`
	auth := fmt.Sprintf("%s:%s", user, password)
	auth = base64.StdEncoding.EncodeToString([]byte(auth))
	//fmt.Println(auth)
	content = fmt.Sprintf(content, server, auth)
	data[".dockerconfigjson"] = content
	secret.StringData = data

	_, err := service.k8sClient.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	return err
}

/*****************************************************************************************************/

/* 创建 K8s Service */
func (service *ApplicationService) createK8sService(namespace string, svc *entity.Service) (*corev1.Service, error) {
	k8sService := &corev1.Service{}
	k8sService.Name = svc.Name
	k8sService.Spec.Selector = map[string]string{"app": svc.Name}

	if svc.Ports != nil {
		servicePorts := make([]corev1.ServicePort, len(svc.Ports))

		for i, port := range svc.Ports {
			servicePort := corev1.ServicePort{}
			servicePort.Port = port.ClusterPort
			servicePort.TargetPort.IntVal = port.TargetPort
			servicePort.Protocol = corev1.Protocol(strings.ToUpper(port.Protocol))

			if port.NodePort > 0 {
				servicePort.NodePort = port.NodePort
				k8sService.Spec.Type = corev1.ServiceTypeNodePort
			}

			servicePorts[i] = servicePort
		}

		k8sService.Spec.Ports = servicePorts
	}

	k8sService, err := service.k8sClient.CoreV1().Services(namespace).Create(context.TODO(), k8sService, metav1.CreateOptions{})
	return k8sService, err
}

/* 创建 K8s Deployment */
func (service *ApplicationService) createK8sDeployment(namespace string, flavor *entity.Flavor, svc *entity.Service) error {
	secret := service.getK8sRegistrySecret(namespace)
	container := corev1.Container{}
	requests := corev1.ResourceList{}
	limits := corev1.ResourceList{}

	container.Name = svc.Name
	container.Image = svc.Image

	if svc.Ports != nil {
		containerPorts := make([]corev1.ContainerPort, len(svc.Ports))

		for i, port := range svc.Ports {
			containerPort := corev1.ContainerPort{}
			containerPort.ContainerPort = port.TargetPort
			containerPort.Protocol = corev1.Protocol(strings.ToUpper(port.Protocol))

			containerPorts[i] = containerPort
		}

		container.Ports = containerPorts
	}


	//quantity := resource.Quantity{}

	//quantity.Set(flavor.Cpu.Requests)
	//requests[corev1.ResourceCPU] = quantity

	//quantity.Set(flavor.Memory.Requests << 20)
	//requests[corev1.ResourceMemory] = quantity
	//
	//quantity.Set(flavor.Disk.Requests << 30)
	//requests[corev1.ResourceEphemeralStorage] = quantity
	//
	//quantity.Set(flavor.Cpu.Limits)
	//limits[corev1.ResourceCPU] = quantity

	//quantity.Set(flavor.Memory.Limits << 20)
	//limits[corev1.ResourceMemory] = quantity

	//quantity.Set(flavor.Disk.Limits << 30)
	//limits[corev1.ResourceEphemeralStorage] = quantity

	//container.Resources.Requests = requests
	//container.Resources.Limits = limits

	imagePullSecret := corev1.LocalObjectReference{}
	imagePullSecret.Name = secret

	deployment := &appsv1.Deployment{}

	deployment.Name = svc.Name
	deployment.Spec.Selector = &metav1.LabelSelector{MatchLabels: map[string]string{"app": svc.Name}}
	deployment.Spec.Replicas = &svc.Replicas
	deployment.Spec.Template.Labels = map[string]string{"app": svc.Name}
	deployment.Spec.Template.Spec.Containers = []corev1.Container{container}
	deployment.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{imagePullSecret}

	_, err := service.k8sClient.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

/* 根据 K8s 命名空间获取登录镜像仓库所用的 secret */
func (service *ApplicationService) getK8sRegistrySecret(namespace string) string {
	return "registry"
}

/* 获取 K8s Service 信息 */
func (service *ApplicationService) getK8sService(namespace string, name string) (*corev1.Service, error) {
	svc, err := service.k8sClient.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return svc, err
}

/* 获取 K8s Deployment 信息 */
func (service *ApplicationService) getK8sDeployment(namespace string, name string) (*appsv1.Deployment, error) {
	deployment, err := service.k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return deployment, err
}

/* 获取所有的 K8s Service */
func (service *ApplicationService) listK8sServices(namespace string) (*corev1.ServiceList, error) {
	services, err := service.k8sClient.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	return services, err
}

/* 获取所有的 K8s Deployment */
func (service *ApplicationService) listK8sDeployments(namespace string) (*appsv1.DeploymentList, error) {
	deployments, err := service.k8sClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	return deployments, err
}

/* 删除 K8s 命名空间 */
func (service *ApplicationService) deleteK8sNamespace(namespace string) error {
	err := service.k8sClient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	return err
}

/* 删除 K8s Registry secret */
func (service *ApplicationService) deleteK8sRegistrySecret(namespace, name string) error {
	err := service.k8sClient.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

/* 删除 K8s Service */
func (service *ApplicationService) deleteK8sService(namespace string, name string) error {
	err := service.k8sClient.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

/* 删除 K8s Deployment */
func (service *ApplicationService) deleteK8sDeployment(namespace string, name string) error {
	err := service.k8sClient.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

/* 回滚、升级 K8s Deployment */
func (service *ApplicationService) patchK8sDeployment(namespace string, name, image string) error {
	format := `{"spec":{"template":{"spec":{"containers":[{"name":"%s","image":"%s"}]}}}}`
	data := []byte(fmt.Sprintf(format, name, image))

	_, err := service.k8sClient.AppsV1().Deployments(namespace).Patch(context.TODO(), name, types.StrategicMergePatchType,
		data, metav1.PatchOptions{})
	return err
}
