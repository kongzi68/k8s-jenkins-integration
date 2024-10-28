package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Deployments deployments

type deployments struct{}

// 定义列表返回的内容
type DeploymentsListDTO struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// 获取deployments列表，支持过滤、排序、分页
func (n *deployments) GetDeployments(ctx *gin.Context, k8sRegion, filterName, namespaces string, limit, page int) (deploymentsListDTO *DeploymentsListDTO, err error) {
	//通过clientset获取原生deployments列表
	deploymentsList, err := K8s.cs[k8sRegion].AppsV1().Deployments(namespaces).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("获取deployments列表失败," + err.Error())
		return nil, errors.New("获取deployments列表失败," + err.Error())
	}
	//实例化dataselector
	selectableData := &DataSelector{
		GenericDataList: n.toCells(deploymentsList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	//再去排序和分页
	data := filtered.Sort().Paginate()
	deployments := n.fromCells(data.GenericDataList)

	return &DeploymentsListDTO{
		Items: deployments,
		Total: total,
	}, nil
}

// 获取deployments列表，仅返回名称，支持过滤、排序、分页
func (n *deployments) GetDeploymentsOnlyName(ctx *gin.Context, k8sRegion, filterName, namespaces string, limit, page int) (deployments []string, err error) {
	deployments = []string{}
	data, err := n.GetDeployments(ctx, k8sRegion, filterName, namespaces, limit, page)
	if err != nil {
		return nil, err
	}
	for _, deployment := range data.Items {
		deployments = append(deployments, deployment.Name)
	}
	return deployments, nil
}

func (*deployments) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentsCell(std[i])
	}
	return cells
}

func (*deployments) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentsCell))
	}
	return deployments
}
