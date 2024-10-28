package config

// 常量
const (
	ListenAddr = "iamIPaddr:8080"
)

// k8s区域关键词，配置文件关键词
var K8SRegion = [...]string{"cd-develop", "cd", "sh"}

// harbor镜像仓库关键词
var HarborKeyArray = [...]string{
	"iamIPaddr:8765",
	"iamIPaddr",
	"harbor.betack.com",
}
