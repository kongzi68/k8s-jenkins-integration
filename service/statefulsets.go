package service

import (
	"bfJenkinsApi/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Statefulsets statefulsets

type statefulsets struct{}

// 定义列表返回的内容
type StatefulsetsListDTO struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// 获取statefulsets列表，支持过滤、排序、分页
func (n *statefulsets) GetStatefulsets(ctx *gin.Context, k8sRegion, filterName, namespaces string, limit, page int) (statefulsetsListDTO *StatefulsetsListDTO, err error) {
	//通过clientset获取原生statefulsets列表
	statefulsetsList, err := K8s.cs[k8sRegion].AppsV1().StatefulSets(namespaces).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("获取statefulsets列表失败," + err.Error())
		return nil, errors.New("获取statefulsets列表失败," + err.Error())
	}
	//实例化dataselector
	selectableData := &DataSelector{
		GenericDataList: n.toCells(statefulsetsList.Items),
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
	statefulsets := n.fromCells(data.GenericDataList)

	return &StatefulsetsListDTO{
		Items: statefulsets,
		Total: total,
	}, nil
}

// 获取statefulsets列表，仅返回名称，支持过滤、排序、分页
func (n *statefulsets) GetStatefulsetsOnlyName(ctx *gin.Context, k8sRegion, filterName, namespaces string, limit, page int) (statefulsets []string, err error) {
	statefulsets = []string{}
	data, err := n.GetStatefulsets(ctx, k8sRegion, filterName, namespaces, limit, page)
	if err != nil {
		return nil, err
	}
	excluedStsArray := []string{
		"etcd",
		"mysql57-statefuleset",
		"mysql80-statefuleset",
	}
	for _, sts := range data.Items {
		if !utils.In(sts.Name, excluedStsArray) {
			statefulsets = append(statefulsets, sts.Name)
		}
	}
	return statefulsets, nil
}

func (*statefulsets) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulsetsCell(std[i])
	}
	return cells
}

func (*statefulsets) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulsets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulsets[i] = appsv1.StatefulSet(cells[i].(statefulsetsCell))
	}
	return statefulsets
}
