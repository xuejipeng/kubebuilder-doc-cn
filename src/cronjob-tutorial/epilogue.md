# 结语

至此，我们已经有了 CronJob controller 的相当完整的实现，并利用了KubeBuilder 的大多数功能。

如果需要更多内容，请转到 [Multi-Version
指南l](/multiversion-tutorial/tutorial.md) 以了解如何向项目添加新的 API version。

此外，您可以自行尝试执行以下步骤 -- 很快我们将为它们提供一个教程：

- 编写单元/整合测试(使用  [envtest])
- 为 `kubectl get` 命令[添加额外字段][printer-columns] 

[envtest]: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/envtest

[printer-columns]: /reference/generating-crd.md#additional-printer-columns
