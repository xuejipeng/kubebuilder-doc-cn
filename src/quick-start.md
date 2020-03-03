# 快速开始

这个 Quick Start 指南包括:

- [如何创建一个 project](#create-a-project)
- [如何创建一个 API](#create-an-api)
- [如何在本地运行 operator](#test-it-out)
- [如何在集群上运行 operator](#run-it-on-the-cluster)

## 环境准备

- [go](https://golang.org/dl/) version v1.13+.
- [docker](https://docs.docker.com/install/) version 17.03+.
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.11.3+.
- [kustomize](https://sigs.k8s.io/kustomize/docs/INSTALL.md) v3.1.0+
- Access to a Kubernetes v1.11.3+ cluster.

<span id="installation"></span>
## 安装

安装 [kubebuilder](https://sigs.k8s.io/kubebuilder):

```bash
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.2.0/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.2.0_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

<aside class="note">
<h1>使用最新版本</h1>

可以使用以下链接获取到最新版本 `https://go.kubebuilder.io/dl/latest/${os}/${arch}`.

</aside>

<span id="create-a-project"></span>
## 创建一个 Project

创建一个目录，然后运行 init 命令以初始化新项目

```bash
mkdir $GOPATH/src/example
cd $GOPATH/src/example
kubebuilder init --domain my.domain
```

<aside class="note">

<h1>项目不在 $GOPATH 下</h1>

如果创建的项目不在`GOPATH`中，为了让 kubebuilder 和 Go 识别导入路径，需要运行 ` go mod init <modulename> `

对于 `GOPATH` 的解释可以参考 [How to Write Go Code][how-to-write-go-code-golang-docs] 中的 [The GOPATH environment variable][GOPATH-golang-docs] 

</aside>

<aside class="note">

<h1>Go package 问题</h1>

如果报错 `cannot find package ... (from $GOROOT)` 请执行`$ export GO111MODULE=on` 开启 Go module 模式。

</aside>

<span id="create-an-api"></span>
## 创建一个 API
运行以下命令创建一个新的 API(group/version) -->`webapp/v1`, 并在此 API 上创建一个新的 Kind(CRD) --> `Guestbook `

<!--Run the following command to create a new API (group/version) as `webapp/v1` and the new Kind(CRD) `Guestbook` on it:
-->
```bash
kubebuilder create api --group webapp --version v1 --kind Guestbook
```

<aside class="note">
<h1>选项说明</h1>

`Create Resource [y/n]：y` 创建文件`api/v1/guestbook_types.go`，此文件用来定义 API

`Create Controller [y/n]：y` 创建文件 `controller/guestbook_controller.go`，此文件用来编写实现 Kind(CRD) 的业务逻辑


</aside>

**备注:** 对于如何设计 API 和如何实现业务逻辑可以参考 [Designing an API](/cronjob-tutorial/api-design.md) 和 [What's in
a Controller](cronjob-tutorial/controller-overview.md).


<details><summary>查看示例 `(api/v1/guestbook_types.go)` </summary>
<p>

```go
// GuestbookSpec defines the desired state of Guestbook
type GuestbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Quantity of instances
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Size int32 `json:"size"`

	// Name of the ConfigMap for GuestbookSpec's configuration
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	ConfigMapName string `json:"configMapName"`

	// +kubebuilder:validation:Enum=Phone;Address;Name
	Type string `json:"alias,omitempty"`
}

// GuestbookStatus defines the observed state of Guestbook
type GuestbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// PodName of the active Guestbook node.
	Active string `json:"active"`

	// PodNames of the standby Guestbook nodes.
	Standby []string `json:"standby"`
}

type Guestbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestbookSpec   `json:"spec,omitempty"`
	Status GuestbookStatus `json:"status,omitempty"`
}
```

</p>
</details>

<span id="test-it-out"></span>
## 在本地运行 operator
你需要一个 Kubernetes 测试集群来进行测试，你可以使用 [KIND](https://sigs.k8s.io/kind) 来启动一个本地或者远程的测试集群

<aside class="note">
<h1>Context 使用</h1>

Controller 会自动使用你当前 context 中的 kubeconfig 文件操作 kubernetes 集群(即：`kubectl cluster-info` 显示的集群)

</aside> 

在集群中安装 CRDs

```bash
make install
```

运行 controller (这个命令不会再后台运行 controller，所以建议新建一个 terminal 运行以下命令)

```bash
make run
```

## 安装 Custom Resources 实例

创建 Project时，对于 ` Create Resource [y/n] `，如果你选择 `y`, 将会在 `config/samples/` 目录中为你的 CRD 创建一个 CR 文件 xxx.yaml (如果你修改了 API 定义文件(api/v1/guestbook_types.go)，请务必修改此文件)

```bash
kubectl apply -f config/samples/
```

<span id="run-it-on-the-cluster"></span>
## 在集群上运行 operator

Build 并 Push 镜像 `IMG`:

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
```

使用指定镜像`IMG`将 operator 部署到集群中:

```bash
make deploy IMG=<some-registry>/<project-name>:tag
```

<aside class="note">
<h1>RBAC errors</h1>

如果你遇到 RBAC 相关错误，你需要使用 cluster-admin 权限或者使用 admin 登录，可以参考 [Prerequisites for using Kubernetes RBAC on GKE cluster v1.11.x and older][pre-rbc-gke]


</aside> 

## 卸载 CRDs

从集群中删除 CRDs：

```bash
make uninstall
```

## 下一步
现在，请继续学习 [CronJob tutorial][cronjob-tutorial] 教程，通过开发演示示例项目更好地了解其工作原理。


[pre-rbc-gke]:https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap
[cronjob-tutorial]: https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html
[GOPATH-golang-docs]: https://golang.org/doc/code.html#GOPATH
[how-to-write-go-code-golang-docs]: https://golang.org/doc/code.html 


