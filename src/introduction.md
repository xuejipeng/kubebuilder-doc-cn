**注意:** 快速开始 [Quick
Start](quick-start.md).

**如果正在使用 Kubebuilder v1, 请查看 [Kubebuilder v1](https://book-v1.book.kubebuilder.io)。**

**英文原版: [book.kubebuilder.io](https://book.kubebuilder.io/introduction.html)**

## 哪些人适合看这个文档

#### Kubernetes 的使用者

Kubernetes 的使用者将通过学习 API 是如何设计和实现的，获得对 Kubernetes 更深入的了解。 本书将教读者如何开发自己的 Kubernetes API 以及实现 Kubernetes API 的核心原理。

包括:

- 如何构造 Kubernetes API 和 Resources
- 如何进行 API 版本控制
- 如何实现故障自愈
- 如何实现垃圾回收和 Finalizers
- 如何创建声明式和命令式 API
- 如何创建 Level-Based API 和 Edge-Base API
- 如何创建 Resources 和 Subresources


#### Kubernetes API extension developers

API extension developers 将学习实现标准的 Kubernetes API 的原理和概念，以及用于快速构建 API 的工具和库。本书涵盖了开发人员通常会遇到的陷阱和误解。

包括:

- 如何用一个 reconciliation 方法处理多个 events
-  如何定期执行 reconciliation 方法
- **将来会有**
	- 何时使用 lister cache 与 live lookups(实时查找)
	- 如何垃圾回收和 Finalizers
	- 如何使用 Declarative Validation 和 Webhook Validation
	- 如何实现 API 版本控制

## 相关资源

* 代码仓库: [sigs.k8s.io/kubebuilder](https://sigs.k8s.io/kubebuilder)

* Slack channel: [#kubebuilder](http://slack.k8s.io/#kubebuilder)

* Google Group:
  [kubebuilder@googlegroups.com](https://groups.google.com/forum/#!forum/kubebuilder)
