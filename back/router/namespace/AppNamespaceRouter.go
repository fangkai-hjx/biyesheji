package namespace

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"t/back/entity"
	"t/back/utils"
)

// Router 服务发布
func Router(router *gin.RouterGroup) {
	app := router.Group("/ns/workspace")
	{
		app.GET("/all", getAllNamespace)
		app.PUT("/:namespace", addNamespace)
		app.DELETE("/:namespace", delNamespace)
		//cluster.PUT("/cluster",updateCluster )
		//cluster.PUT("/cluster",updateCluster )
	}
}
// 查询所有的工作空间 http://localhost:8080/pub_service/workspace/all
func getAllNamespace(c *gin.Context) {
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nill")
		return
	}
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	result := make([]entity.Namespace, 0)
	for _, v := range namespaceList.Items {
		var r = entity.Namespace{
			Name:   v.Name,
			Status: string(v.Status.Phase),
		}
		result = append(result, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"data":    result,
	})
}
// 创建指定的工作空间
func addNamespace(c *gin.Context) {
	//namespace := c.PostForm("namespace")
	namespace := c.Param("namespace")
	client := utils.GetK8sClient()
	if client == nil {
		fmt.Println("client is nill")
		return
	}
	ns := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: map[string]string{"name": namespace},
		},
	}
	_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "增加失败",
			"data":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "增加" + namespace + "成功",
		"data":    nil,
	})
	return
}
// 删除指定的工作空间
func delNamespace(c *gin.Context){
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