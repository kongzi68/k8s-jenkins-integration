package service

import (
	"bfJenkinsApi/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespaces namespaces

type namespaces struct{}

// 定义列表返回的内容
type NamesapcesListDTO struct {
	Items []corev1.Namespace `json:"items"`
	Total int                `json:"total"`
}

// 获取namespaces列表，支持过滤、排序、分页
func (n *namespaces) GetNamespaces(ctx *gin.Context, k8sRegion, filterName string, limit, page int) (namesapcesListDTO *NamesapcesListDTO, err error) {
	//通过clientset获取原生namespaces列表
	namespaceList, err := K8s.cs[k8sRegion].CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespaces列表失败," + err.Error())
		return nil, errors.New("获取namespaces列表失败," + err.Error())
	}
	//实例化dataselector
	selectableData := &DataSelector{
		GenericDataList: n.toCells(namespaceList.Items),
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
	namespaces := n.fromCells(data.GenericDataList)

	return &NamesapcesListDTO{
		Items: namespaces,
		Total: total,
	}, nil
}

// 获取namespaces列表，仅返回名称，支持过滤、排序、分页
func (n *namespaces) GetNamespacesOnlyName(ctx *gin.Context, k8sRegion, filterName string, limit, page int) (namespaces []string, err error) {
	namespaces = []string{}
	data, err := n.GetNamespaces(ctx, k8sRegion, filterName, limit, page)
	if err != nil {
		return nil, err
	}
	excluedNamespacesArray := []string{
		"artifactory-oss",
		"bf-metersphere",
		"default",
		"democratic-csi",
		"ops",
		"kube-ci",
		"kube-node-lease",
		"kube-public",
		"kube-system",
		"kuboard",
		"velero-system",
		"test-esm",
		"longhorn-system",
		"private-equity-data",
	}
	for _, namespace := range data.Items {
		if !utils.In(namespace.Name, excluedNamespacesArray) {
			namespaces = append(namespaces, namespace.Name)
		}
	}
	return namespaces, nil
}

func (*namespaces) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		//std[i]原生namespaces类型 -> namespacesCell -> DataCell
		cells[i] = namespacesCell(std[i])
	}
	return cells
}

func (*namespaces) fromCells(cells []DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i := range cells {
		//cells[i] DataCell接口类型 -> namespacesCell -> 原生namespaces类型
		namespaces[i] = corev1.Namespace(cells[i].(namespacesCell))
	}
	return namespaces
}
