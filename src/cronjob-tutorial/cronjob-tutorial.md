# 教程：构建 CronJob

很多教程都是从一些人为设置好的程序开始，或者是给你一些了解基础知识的小应用程序，这些不能让你深入了解更复杂的东西。 相反，本教程几乎带会您了解 Kubebuilder 的全部复杂性，从简单开始，然后逐步发展为功能齐全的产品。


现在让我们假设 Kubernetes 中实现的 `CronJob Controller` 并不能满足我们的需求，我们希望使用 Kubebuilder 重写它。

`CronJob controller` 会控制 kubernetes 集群上的 job 每隔一段时间运行一次，它是基于 `Job controller` 实现的，`Job controller` 的 job 只会执行任务一次。

通过重写 `Job controller`，我们可以更加了解如何与不属于集群的资源类型进行交互。


<aside class="note">

<h1>注意</h1>

请注意，本教程的大部分内容是由 literal Go 文件打包生成的：[docs/book/src/cronjob-tutorial/testdata][tutorial-source]，完整的可运行项目位于[project][tutorial-project-source]中，而打包生成的中间文件位于[testdata][tutorial-source]目录下。


[tutorial-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata

[tutorial-project-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project

</aside>

## 创建我们的项目

如[快速入门](../quick-start.md)中所述，我们需要创建一个新项目，请确保已经[安装Kubebuilder](../quick-start.md#installation) ，然后执行如下命令创建一个新项目：

```bash
# we'll use a domain of tutorial.kubebuilder.io,
# so all API groups will be <group>.tutorial.kubebuilder.io.
kubebuilder init --domain tutorial.kubebuilder.io
```

现在我们已经有了一个项目，让我们看一下到目前为止 Kubebuilder 为我们搭建的脚手架...
