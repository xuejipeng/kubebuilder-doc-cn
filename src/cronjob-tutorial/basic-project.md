# 基本项目中有什么？

在创建新项目时，Kubebuilder 为我们提供了一些基本的模板。

## 基础设施模板


首先，会创建一些用于构建项目的 basic infrastructure：

<details> <summary><code>go.mod</code>: 一个与项目匹配，包含最基本依赖关系的 go module 文件</summary>

```go
{{#include ./testdata/project/go.mod}}
```
</details>

<details><summary><code>Makefile</code>: 用于构建和部署 controller</summary>

```makefile
{{#include ./testdata/project/Makefile}}
```
</details>

<details><summary><code>PROJECT</code>: 用于创建新组件的 Kubebuilder 元数据</summary>

```yaml
{{#include ./testdata/project/PROJECT}}
```
</details>

## 启动配置项模板

我们能在 [`config/`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config) 目录下找到运行 operator 所需的所有配置文件，现在它只包含运行 controller 所需要的 [Kustomize](https://sigs.k8s.io/kustomize) YAML 配置文件，后续我们编写 operator 时，这个目录还会包含 CustomResourceDefinitions(CRD)、RBAC 和 Webhook 等相关的配置文件

[`config/default`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/default) 包含 [Kustomize base](https://github.com/kubernetes-sigs/kubebuilder/blob/master/docs/book/src/cronjob-tutorial/testdata/project/config/default/kustomization.yaml) 文件，用于以标准配置启动 controller。

[`config/`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config)目录下每个目录都包含不同的配置：

- [`config/manager`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/manager): 包含在 k8s 集群中以 pod 形式运行 controller 的 YAML 配置文件 

- [`config/rbac`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/rbac): 包含运行 controller 所需最小权限的配置文件

## 程序入口

最后，但同样重要的是，Kubebuilder 创建了我们项目的程序入口：main.go。下面让我们看看这个文件...
