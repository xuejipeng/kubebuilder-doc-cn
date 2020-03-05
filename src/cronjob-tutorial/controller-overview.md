# controller 中有什么?

`Controllers` 是 operator 和 Kubernetes 的核心组件。

controller 的职责是确保实际的状态（包括群集状态，以及潜在的外部状态，例如正在运行的 Kubelet 容器和云提供商的 loadbalancers）与给定 object 期望的状态相匹配。
每个 controller 专注于一个 `root Kind`，但也可以与其他 `Kinds` 进行交互。

这种努力达到期望状态的过程，我们称之为 `reconciling`(调和，使...一直)。

在 controller-runtime 库中，实现 Kind reconciling 的逻辑我们称为
 [*Reconciler*](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile)。
 `reconciler` 获取对象的名称并返回是否需要重试(例如: 发生错误或是一些周期性的 controllers，像 HorizontalPodAutoscale)。

{{#literatego ./testdata/emptycontroller.go}}

现在，我们已经看到了 controller 的基本结构，下一步让我们来填写 `CronJob` 的逻辑
