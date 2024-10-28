package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

// 最外层结构体，用于封装排序、过滤、分页的数据类型
type DataSelector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// 接口，用于k8s各种资源的类型转换，后续进行排序、过滤、分页操作
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// 过滤、分页操作所需要的参数
type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

// 过滤参数
type FilterQuery struct {
	Name string
}

// 分页参数
type PaginateQuery struct {
	Limit int
	Page  int
}

// 结构体的自定义排序，重写Len、Swap、Less方法
// Len方法用于获取切片长度
func (d *DataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap方法用于切片中元素比较大小后的位置交换，升序或降序
// i和j是两个临近的数据元素的下标
func (d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less方法比较大小
func (d *DataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// 重写以上三个方法后，调用sort.Sort()对数据进行自定义排序
func (d *DataSelector) Sort() *DataSelector {
	sort.Sort(d)
	return d
}

// Filter方法用于过滤元素，比较元素的Name属性，若包含，在返回
func (d *DataSelector) Filter() *DataSelector {
	//判断参数是否正确
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	//若Name的传参不为空，则返回元素中包含Name的所有元素
	filterdList := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			matches = false
			continue
		}
		if matches {
			filterdList = append(filterdList, value)
		}
	}
	d.GenericDataList = filterdList

	return d
}

// Paginate方法用于数组分页，根据limit和page的传参，返回数据
func (d *DataSelector) Paginate() *DataSelector {
	//每一页多少数据
	limit := d.DataSelect.Paginate.Limit
	//第几页
	page := d.DataSelect.Paginate.Page
	//传参判断
	if limit <= 0 || page <= 0 {
		return d
	}

	//切片的分页原理是截取和重新组合
	//通过limit和page，计算得出startIndex和endIndex
	//举例：25个元素
	//limit=10, page=1, start=0, end=10, d.GenericDataList[0:10], 含头不含尾
	//limit=10, page=3, start=20, end=30
	startIndex := limit * (page - 1)
	endIndex := limit * page

	//处理最后一页
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]

	return d
}

// 类型重新定义
// 继承原生类型的所有属性和方法
// 可以在新类型上新增属性和方法
type namespacesCell corev1.Namespace

func (p namespacesCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p namespacesCell) GetName() string {
	return p.Name
}

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

type deploymentsCell appsv1.Deployment

func (d deploymentsCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentsCell) GetName() string {
	return d.Name
}

type statefulsetsCell appsv1.StatefulSet

func (d statefulsetsCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d statefulsetsCell) GetName() string {
	return d.Name
}
