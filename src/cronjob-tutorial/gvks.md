# Groups、Versions 和 Kinds 之间的关系

在开始编写 API 之前，让我们讨论一些关键术语

当描述 Kubernetes 中的 API 时，我们经常用到的四个术语是：`groups`、
`versions`、`kinds` 和 `resources`

## Groups 和 Versions

Kubernetes 中的 `API Group` 只是相关功能的集合, **每个 `group` 包含一个或者多个 `versions`**, 这样的关系可以让我们随着时间的推移，通过创建不同的 `versions ` 来更改 API 的工作方式。

## Kinds 和 Resources
<h3>Kinds</h3>

**每个 `API group-version` 包含一个或者多个 API 类型, 我们称之为 `Kinds`**.
虽然同类型 `Kind` 在不同的 `version` 之间的表现形式可能不同，但是同类型 `Kind` 必须能够存储其他 `Kind` 的全部数据，也就是说同类型 `Kind` 之间必须是互相兼容的(我们可以把数据存到 fields 或者 annotations)，这样当你使用老版本的 `API group-version` 时不会造成丢失或损坏, 有关更多信息，请参见[Kubernetes API指南](https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md)。

<aside class="note">
<h1>自己的理解</h1>

通常部署到 k8s 集群中的 YAML 文件前几行格式如下, 其中`extensions/v1beta1` 就是 `API group-version`、`extensions` 是 `Group`、 `v1beta1` 是 `version`、 `Deployment` 是 `Kind`。

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
```
</aside>
<h3>Resources</h3>

你应该听别人提到过 `Resources`， `Resource` 是 `Kind` 在 API 中的标识，通常情况下 `Kind` 和 `Resource` 是一一对应的, 但是有时候相同的 `Kind` 可能对应多个 `Resources`, 比如 Scale Kind 可能对应很多 Resources：deployments/scale 或者 replicasets/scale, 但是在 CRD 中，每个 `Kind` 只会对应一种 `Resource`。

请注意，`Resource` 始终是小写形式，并且通常情况下是 `Kind` 的小写形式。
具体对应关系可以查看 [resource type](https://links.jianshu.com/go?to=https%3A%2F%2Fyq.aliyun.com%2Fgo%2FarticleRenderRedirect%3Furl%3Dhttps%253A%252F%252Fkubernetes.io%252Fdocs%252Freference%252Fkubectl%252Foverview%252F%2523resource-types)。

<aside class="note">
<h1>自己的理解</h1>

当我们使用 kubectl 操作 API 时，操作的就是 `Resource`，比如 `kubectl get pods`, 这里的 `pods` 就是指 `Resource`。

而我们在编写 YAML 文件时，会编写类似 `Kind: Pod` 这样的内容，这里 `Pod` 就是 `Kind`
</aside>

## 那么以上类型在框架中是如何定义的呢

当我们在一个特定的 `group-version`中使用 `Kind` 时，我们称它为 `GroupVersionKind`, 简称 `GVK`, 同样的 `resources` 我们称它为 `GVR`。
稍后我们将看到每个 `GVK` 都对应一个 root Go type (比如：Deployment 就关联着 K8s 源码里面 k8s.io/api/apps/v1 package 中的 Deployment struct)。

现在我们已经熟悉了一些关键术语，那么我们可以开始创建我们的 API 了！

## 额, 但是 Scheme 是什么东东?

`Scheme` 提供了 `GVK` 与对应 Go types(struct) 之间的映射（请不要和 [godocs](https://godoc.org/k8s.io/apimachinery/pkg/runtime#Scheme) 中的 Scheme 混淆）

也就是说给定 Go type 就可以获取到它的 GVK，给定 GVK 可以获取到它的 Go type

例如，让我们假设 `tutorial.kubebuilder.io/api/v1/cronjob_types.go` 中的 `CronJob` 结构体在 `batch.tutorial.kubebuilder.io/v1`  API group 中（也就是说，假设这个 API group 有一种 Kind：`CronJob`，并且已经注册到 Api Server 中）

那么我们可以从 Api Server 获取到数据并反序列化至 `&CronJob{}`中，那么结果会是如下格式:

```json
{
    "kind": "CronJob",
    "apiVersion": "batch.tutorial.kubebuilder.io/v1",
    ...
}
```
*ps：这里翻译的不好，请继续向后看，慢慢就理解了 ：）*
