# k8s api 给 Jenkins pipeline 用

## 接口列表

```go
// pod操作
router.GET("/api/k8s/pods", Pod.GetPods)
router.GET("/api/k8s/pods/containers/all-images", Pod.GetPodsContainersImages)
// namespaces操作
router.GET("/api/k8s/namespaces", Namespaces.GetNamespaces)
router.GET("/api/k8s/namespaces/only-name", Namespaces.GetNamespacesOnlyName)
// deployments操作
router.GET("/api/k8s/deployments/only-name", Deployments.GetDeploymentsOnlyName)
// statefulsets操作
router.GET("/api/k8s/statefulsets/only-name", Statefulsets.GetStatefulsetsOnlyName)
```

## /api/k8s/statefulsets/only-name

```bash
http://iamIPaddr:8080/api/k8s/statefulsets/only-name?k8s_region=cd&namespaces=saasdata-dev-1

[root@localhost ~]# curl -s "http://iamIPaddr:8080/api/k8s/statefulsets/only-name?k8s_region=cd&namespaces=saasdata-dev-1" | jq
{
  "code": 0,
  "data": [
    "saas-data-server-dev-1"
  ],
  "msg": "获取Statefulsets列表成功"
}
```

## 部署

### 推镜像

```bash
[root@localhost study-go]# docker image list bf-jenkins-api
REPOSITORY       TAG       IMAGE ID       CREATED          SIZE
bf-jenkins-api   latest    iamIPaddreb7   10 minutes ago   128MB
[root@localhost study-go]# docker image tag bf-jenkins-api:latest iamIPaddr:8765/bf-devops/bf-jenkins-api:viamIPaddr
[root@localhost study-go]# docker image push iamIPaddr:8765/bf-devops/bf-jenkins-api:viamIPaddr
The push refers to repository [iamIPaddr:8765/bf-devops/bf-jenkins-api]
caiamIPaddr: Pushed
f60721fe925b: Pushed
c5bec285958d: Pushed
950d1cd21157: Layer already exists
viamIPaddr: digest: shaiamIPaddrac3aced7b2d29caiamIPaddr5569d57abiamIPaddr3b0a0bb2f17e4 size: 1156
```

## Jenkins pipeline 使用场景

Jenkins动态选择参数

```groovy
  parameters {
    string description: '服务需要部署到生产环境shanghai k8s集群，请输入docker镜像TAG。\n比如：docker镜像 iamIPaddr:8765/bar/bar-svc-api-app:6b2cfiamIPaddr 中冒号后面的 6b2cfiamIPaddr 是镜像TAG。',
           name: 'SVC_IMAGE_TAG'
    extendedChoice defaultValue: '请选择',
                    description: '选择K8S所在的区域：成都开发测试环境、上海生产环境',
                    multiSelectDelimiter: ',',
                    name: 'K8S_REGION',
                    quoteValue: false,
                    saveJSONParameterToFile: false,
                    type: 'PT_SINGLE_SELECT',
                    value: '请选择,上海prod-k8s,成都-k8s',
                    visibleItemCount: 5
    reactiveChoice choiceType: 'PT_SINGLE_SELECT', 
        description: '选择需要更新的K8S命名空间环境',
        filterLength: 1,
        filterable: true,
        name: 'K8S_NAMESPACES',
        randomName: 'choice-parameter-iamIPaddr76',
        referencedParameters: 'K8S_REGION',
            script: groovyScript(
                fallbackScript: [classpath: [], oldScript: '', sandbox: false, script: 'return ["请选择"]'],
                script: [classpath: [], oldScript: '', sandbox: false, script: '''
                    import groovy.json.JsonSlurper
                    switch(K8S_REGION) {
                        case "上海prod-k8s":
                            k8sRegion = "sh"
                        break
                        case "成都-k8s":
                            k8sRegion = "cd"
                        break
                        default:
                            return []
                        break
                    }
                    def apiURL = "http://k8s-api-devops-in-svc:8080/api/k8s/namespaces/only-name?k8s_region=" + k8sRegion
                    def pkgObject = ["curl", apiURL].execute().text
                    def jsonSlurper = new JsonSlurper()
                    def artifactsJsonObject = jsonSlurper.parseText(pkgObject)
                    retList = artifactsJsonObject.data
                    return retList.sort()
                '''])
    reactiveChoice choiceType: 'PT_SINGLE_SELECT',
        description: '选择需要更新的服务',
        filterLength: 1,
        filterable: true,
        name: 'K8S_SERVICES',
        randomName: 'choice-parameter-iamIPaddr91',
        referencedParameters: 'K8S_REGION,K8S_NAMESPACES',
            script: groovyScript(
                fallbackScript: [classpath: [], oldScript: '', sandbox: false, script: 'return ["请选择"]'],
                script: [classpath: [], oldScript: '', sandbox: false, script: '''
                    import groovy.json.JsonSlurper
                    switch(K8S_REGION) {
                        case "上海prod-k8s":
                            k8sRegion = "sh"
                        break
                        case "成都-k8s":
                            k8sRegion = "cd"
                        break
                        default:
                            return []
                        break
                    }
                    def apiK8S = "http://k8s-api-devops-in-svc:8080/api/k8s/"
                    def workloadTypeList = ["deployments", "statefulsets"]
                    def retList = []
                    for (workloadType in workloadTypeList) {
                        // http://k8s-api-devops-in-svc:8080/api/k8s/deployments/only-name?k8s_region=cd&namespaces=saasdata-dev-1
                        def apiURL =  apiK8S + workloadType  + "/only-name?k8s_region=" + k8sRegion  + "&namespaces=" +  K8S_NAMESPACES
                        // println(apiURL)
                        def pkgObject = ["curl", apiURL].execute().text
                        def jsonSlurper = new JsonSlurper()
                        def artifactsJsonObject = jsonSlurper.parseText(pkgObject)
                        retList.add(artifactsJsonObject.data)
                    }
                    ret = retList[0] + retList[1]
                    return ret.sort()
                '''])
  }
```
