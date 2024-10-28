package service

import (
	"bfJenkinsApi/config"
	"bfJenkinsApi/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct{}

// 定义列表返回的内容
type PodListDTO struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(ctx *gin.Context, k8sRegion, filterName, namespace string, limit, page int) (podListDTO *PodListDTO, err error) {
	//通过clientset获取原生pod列表
	podList, err := K8s.cs[k8sRegion].CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("获取Pod列表失败," + err.Error())
		return nil, errors.New("获取Pod列表失败," + err.Error())
	}
	//实例化dataselector
	selectableData := &DataSelector{
		GenericDataList: p.toCells(podList.Items), //这里就要把原生pod转成podcell
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
	pods := p.fromCells(data.GenericDataList)

	return &PodListDTO{
		Items: pods,
		Total: total,
	}, nil
}

// 获取所有k8s集群中，所有pod中容器的镜像
func (p *pod) GetPodsContainersImages(ctx *gin.Context, filterName string, limit, page int) (images []string, err error) {
	images = []string{}
	for _, k8sRegion := range config.K8SRegion {
		data, err := p.GetPods(ctx, k8sRegion, filterName, "", limit, page)
		if err != nil {
			return nil, err
		}
		for _, val := range data.Items {
			for item := range val.Spec.Containers {
				image := val.Spec.Containers[item].Image
				if utils.RegMatchStringArray(image, config.HarborKeyArray) {
					images = append(images, image)
				}
			}
		}
	}
	tmpRet := utils.Unique(images)
	return tmpRet, nil
}

func (*pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		//std[i]原生pod类型 -> podCell -> DataCell
		cells[i] = podCell(std[i])
	}
	return cells
}

func (*pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i] DataCell接口类型 -> podCell -> 原生pod类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
