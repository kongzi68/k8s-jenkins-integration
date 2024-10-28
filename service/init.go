package service

import (
	"fmt"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	cs map[string]*kubernetes.Clientset
}

// 初始化clientset，并把clientset赋值到k8s结构体中
func (k *k8s) Init(k8sRegion string) {
	// 在 kubeconfig 中使用当前上下文
	// path-to-kubeconfig -- 例如 /root/.kube/config
	kubeConfigPath := fmt.Sprintf("config/%s-k8s-kube.config", k8sRegion)
	logger.Info(kubeConfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err)
	}
	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	logger.Info("创建k8s clientset成功")

	// 初始化 k.cs
	if k.cs == nil {
		k.cs = make(map[string]*kubernetes.Clientset)
	}

	k.cs[k8sRegion] = clientset

	// 访问 API 以列出 Pod
	//pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	//fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
