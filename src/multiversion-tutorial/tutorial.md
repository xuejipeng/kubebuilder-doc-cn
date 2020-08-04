# 教程: 多版本 API

大多数项目都是从 Alpha API 开始的，经过不断迭代演变成 stable API。但是，一旦 API 稳定了，就很难对其进行重大更改。这就是 API 版本的作用。

让我们对 `CronJob` API spec 进行一些更改，并确保我们的 CronJob 项目支持所有不同版本的 API。


如果还没有 `CronJob` API ，请完成 [CronJob 基础教程](/cronjob-tutorial/cronjob-tutorial.md)。

<aside class="note">

<h1>注意</h1>

请注意，本教程的大部分内容都是从构成可运行项目的 Go 语言文件生成的，它们位于书的源目录中: [docs/book/src/multiversion-tutorial/testdata/project][tutorial-source]。

[tutorial-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/multiversion-tutorial/testdata/project

</aside>

<aside class="note warning">

<h1>Kubernetes 最低版本！</h1>

CRD 转换支持在 Kubernetes 1.13 中作为 Alpha 功能引入（这意味着它默认情况下未启用，并且需要通过Feature Gate启用），并在Kubernetes 1.15中成为 Beta 版本（这意味着默认情况下处于启用状态）。

如果您使用的是Kubernetes 1.13-1.14，请确保启用此功能。如果您使用的是Kubernetes 1.12或更低版本，则需要创建一个新的集群。请查看 [KinD instructions](/reference/kind.md)，查看如何搭建一个 all-in-one 集群。

</aside>

接下来，让我们看看我们需要修改哪些内容...

[kube-feature-gates]: https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/ "Kubernetes Feature Gates"
