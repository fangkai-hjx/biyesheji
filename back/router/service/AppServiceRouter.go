package svc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"t/back/entity"
	"t/back/utils"
	"time"
)

// Router 服务发布
func Router(router *gin.RouterGroup) {
	app := router.Group("/svc/workspace")
	{
		app.PUT("/:namespace", createService)
		app.DELETE("/:namespace", deleteService)
		app.POST("/:namespace", updateService)
		app.GET("/:namespace/all", getAllService)
		//app.GET("/:namespace/:svc_name", getService)
		app.GET("/:namespace/svc/status", getServiceStatusForNamespace)
		app.GET("/:namespace/pod/status", getPodStatusForNamespace)
		app.GET("/:namespace/:service_name/", getServiceDetail)
	}
}
func createService(c *gin.Context) {
	namespace := c.Param("namespace")
	service_name := c.PostForm("service_name")
	image := c.PostForm("image")
	//tag := "app=nginx-demo"
	port := strToInt32(c.PostForm("port"))
	//count := c.Param("count")
	count := strToInt32(c.PostForm("count"))
	fmt.Println(namespace, image, service_name, image, port)
	k8sClient := utils.GetK8sClient()
	if k8sClient == nil {
		fmt.Println("k8sClient is nil")
		return
	}
	// check service from db
	dbClient := utils.GetDBClient()
	if dbClient == nil {
		fmt.Println("dbClient is nil")
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": "获取数据库连接错误",
		})
		return
	}
	row := dbClient.QueryRow(`SELECT COUNT(*) AS total FROM tb_service WHERE delete_flag!=0 and service_name = ? and namespace= ?`, service_name, namespace)

	var r int64
	if err := row.Scan(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": "数据库查询失败",
		})
		return
	}
	if r != 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "服务已被创建，请勿重复创建",
		})
		return
	}
	// insert or update
	// if
	now := time.Now().Unix()
	row = dbClient.QueryRow(`SELECT COUNT(*) AS total FROM tb_service WHERE service_name = ? and namespace= ?`, service_name, namespace)
	if err := row.Scan(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": "数据库查询失败",
		})
		return
	}
	if r == 0 {
		_, err := dbClient.Exec(`INSERT INTO tb_service (service_name,namespace,create_time,replicas,image_name,env_vars,ports,description,creator,change_time) VALUES 
					(?,?,?,?,?,?,?,?,?,?)`,
			service_name, namespace, now, count, image, "ttt", port, "ttttt", "admin", now)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "数据库插入失败" + err.Error(),
			})
		}
		return
	} else {
		_, err := dbClient.Exec(`UPDATE tb_service SET delete_flag=? WHERE service_name = ? and namespace=?`, "1", service_name, namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "数据库更新失败" + err.Error(),
			})
		}
	}
	deploymentsClient := k8sClient.AppsV1().Deployments(namespace)
	servicesClient := k8sClient.CoreV1().Services(namespace)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: service_name,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port: 80,
				},
			},
			Selector: map[string]string{"app": "nginx-demo"},
		},
	}
	deployment := &appsv1.Deployment{ //
		ObjectMeta: metav1.ObjectMeta{
			Name: service_name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(count),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx-demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": service_name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  service_name,
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: port,
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("Creating deployment...")
	dep_result, err1 := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	svc_result, err2 := servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
	// 将服务ip 更新到数据库 svc_result.Spec.ClusterIP
	_, err := dbClient.Exec(`UPDATE tb_service SET cluster_ip=? WHERE service_name = ? and namespace= ?`, svc_result.Spec.ClusterIP, service_name, namespace)
	if err != nil {
		fmt.Println("服务IP更新到数据库失败")
		return
	}
	data := ""
	if err1 != nil {
		data += err1.Error()
	}
	if err2 != nil {
		data += err2.Error()
	}
	if data != "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "创建失败",
			"data":    data,
		})
		return
	}
	fmt.Printf("Created deployment %q.\n", dep_result.GetObjectMeta().GetName())
	fmt.Printf("Created service %q.\n", svc_result.GetObjectMeta().GetName())
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    nil,
	})
}

// 移除服务
func deleteService(c *gin.Context) {
	//namespace := c.PostForm("namespace")
	namespace := c.Param("namespace")
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nill")
		return
	}
	err := client.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "删除失败",
			"data":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除" + namespace + "成功",
		"data":    nil,
	})
}

// 修改服务参数
// 1 副本数
// 2 修改镜像
// 3 修改端口
func updateService(c *gin.Context) {
	namespace := c.Param("namespace")
	service_name := c.PostForm("service_name")
	replica := c.PostForm("replica")
	//image := c.PostForm("image")
	//port := c.PostForm("port")
	deploymentsClient := utils.GetK8sClient().AppsV1().Deployments(namespace)
	deployment, err := deploymentsClient.Get(context.TODO(), service_name, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "服务不存在" + err.Error(),
		})
		return
	}
	if deployment != nil {
		deployment.Spec.Replicas = int32Ptr(strToInt32(replica))
	}
	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "服务更新错误",
		})
		return
	}
	dbClient := utils.GetDBClient()
	if dbClient == nil {
		fmt.Println("dbClient is nil")
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": "获取数据库连接错误",
		})
		return
	}
	now := time.Now().Unix()
	_, err = dbClient.Exec(`UPDATE tb_service SET replicas=?,change_time =? WHERE service_name = ? and namespace= ?`, replica, now, service_name, namespace)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "服务删除失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "服务更新成功",
	})
}

type PodStatus struct {
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	Unknown   int `json:"unknown"`
}
type ServiceStatus struct {
	Healthy   int `json:"healthy"`
	UnHealthy int `json:"unhealthy"`
	Error     int `json:"error"`
}

func getServiceStatusForNamespace(c *gin.Context) {
	namespace := c.Param("namespace")
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nil")
		return
	}
	serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务失败",
			"data":    err.Error(),
		})
		return
	}
	Healthy := 0
	UnHealthy := 0
	Error := 0
	for _, v := range serviceList.Items {
		sss := getOneServiceStatus(namespace, v.Name)
		if sss == -1 {
			Error++
		} else if sss == 1 {
			Healthy++
		} else {
			UnHealthy++
		}
	}
	ss := ServiceStatus{
		Healthy:   Healthy,
		UnHealthy: UnHealthy,
		Error:     Error,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务成功",
		"data":    ss,
	})
}

// get service status[-1 error 1 healthy 2 unhealthy]
func getOneServiceStatus(namespace, service_name string) int {
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nil")
		return -1
	}
	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), service_name, metav1.GetOptions{})
	if err != nil {
		return -1
	}
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=" + service.Name})
	if err != nil {
		return -1
	}
	running_pod := 0
	error_pod := 0
	all_pod := len(podList.Items)
	result := -1
	if all_pod == 0 {
		return 1
	}
	for _, pod := range podList.Items {
		sc := getCondition(namespace, pod.Name)
		if sc.ContainersReady == "True" && sc.Ready == "True" &&
			sc.PodScheduled == "True" && sc.Initialized == "True" {
			running_pod++
		} else {
			error_pod++
		}
	}
	if running_pod == 0 {
		result = -1
	}
	if running_pod == all_pod {
		result = 1
	} else {
		result = 2
	}
	return result
}

// get pod status for namespace
func getPodStatusForNamespace(c *gin.Context) {
	namespace := c.Param("namespace")
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nil")
		return
	}
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务失败",
			"data":    err.Error(),
		})
		return
	}
	ps := PodStatus{}
	for _, v := range podList.Items {
		if string(v.Status.Phase) == "Succeeded"{
			ps.Succeeded ++
		}else if string(v.Status.Phase) == "Running" {
			ps.Running ++
		}else if string(v.Status.Phase) == "Pending"{
			ps.Pending ++
		}else if string(v.Status.Phase) == "Failed"{
			ps.Failed ++
		}else{
			ps.Unknown ++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务失败",
		"data":    ps,
	})
	return
	//ps := PodStatus{
	//	Healthy:   0,
	//	UnHealthy: 0,
	//	Error:     0,
	//}
}
// 通过服务名 获取服务的详细信息--->pod details
func getServiceDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	service_name := c.Param("service_name")
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nil")
		return
	}
	//查询所有的Pod
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=" + service_name})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务失败",
			"data":    err.Error(),
		})
		return
	}
	podItems := make([]entity.Pod, 0)
	for _, pod := range podList.Items {
		ppd := entity.Pod{
			Name:              pod.Name,
			ServiceConditions: getCondition(namespace, pod.Name),
			Image:             pod.Spec.Containers[0].Image,
		}
		podItems = append(podItems, ppd)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务失败",
		"data":    podItems,
	})
}
// 获取所有服务，包括获取服务下的实例
func getAllService(c *gin.Context) {
	namespace := c.Param("namespace")
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nil")
		return
	}
	serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务失败",
			"data":    err.Error(),
		})
		return
	}

	result := make([]entity.Service, 0)
	for _, v := range serviceList.Items {
		s := entity.Service{
			Name:      v.Name,
			ClusterIP: v.Spec.ClusterIP,
		}
		if string(v.Spec.SessionAffinity) == "None" {
			s.SessionAffinity = "false"
		} else {
			s.SessionAffinity = "true"
		}
		//查询所有的Pod
		podList, _ := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=" + v.Name})
		podItems := make([]entity.Pod, 0)
		running := 0
		all := len(podList.Items)
		for _, pod := range podList.Items {
			ppd := entity.Pod{
				Name:              pod.Name,
				//ServiceName:       v.Name,
				ServiceConditions: getCondition(namespace, pod.Name),
				Image:             pod.Spec.Containers[0].Image,
			}
			if ppd.ServiceConditions.ContainersReady == "True" && ppd.ServiceConditions.Ready == "True" &&
				ppd.ServiceConditions.PodScheduled == "True" && ppd.ServiceConditions.Initialized == "True" {
				running++
			}
			podItems = append(podItems, ppd)
		}
		s.Pod = podItems
		if all == 0 {
			s.SuccessLu = 0
			s.Total, s.Success, s.Fail = 0, 0, 0
		} else {
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(running)/float64(all)), 64)
			s.SuccessLu = value * 100
			s.Total, s.Success, s.Fail = all, running, all-running
		}
		s.Status = judgeSvcStatus(all, running)
		result = append(result, s)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务成功",
		"data":    result,
	})
}
func judgeSvcStatus(all, running int) string {
	if all == 0{
		return "ERROR"
	}
	if running == all {
		return "HEALTHY"
	} else if running == 0 {
		return "ERROR"
	} else {
		return "UNHEALTHY"
	}
}
func getService(c *gin.Context) {

}
func getCondition(namespace, svcName string) (sc entity.ServiceConditions) {
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nil")
		return sc
	}
	myPod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		return sc
	}
	for _, v := range myPod.Status.Conditions {
		if v.Type == "Initialized" {
			if v.Status == "True" {
				sc.Initialized = "True"
			} else {
				sc.Initialized = "False"
			}
			continue
		}
		if v.Type == "ContainersReady" {
			if v.Status == "True" {
				sc.ContainersReady = "True"
			} else {
				sc.ContainersReady = "False"
			}
			continue
		}
		if v.Type == "Ready" {
			if v.Status == "True" {
				sc.Ready = "True"
			} else {
				sc.Ready = "False"
			}
			continue
		}
		if v.Type == "PodScheduled" {
			if v.Status == "True" {
				sc.PodScheduled = "True"
			} else {
				sc.PodScheduled = "False"
			}
			continue
		}
	}
	return sc
}

func int32Ptr(i int32) *int32 { return &i }
func strToInt32(str string) int32 {
	j, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return -1
	}
	return int32(j)
}
