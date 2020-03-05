# 创建一个 API

创建一个新的 Kind (你还记得[上一章](./gvks.md#kinds-and-resources)的内容, 对吗?)和相应的 controller，我们可以使用` kubebuilder create api` 命令:

```bash
kubebuilder create api --group batch --version v1 --kind CronJob
```

当我们第一次执行这个命令时，它会为创建一个新的 `group-version` 目录

在这里，[`api/v1/`](https://sigs.k8s.io/kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/api/v1) 这个目录会被创建, 对应的 `group-version` 是 `batch.tutorial.kubebuilder.io/v1` (还记得我们一开始设置的[`--domain`
setting](cronjob-tutorial.md#scaffolding-out-our-project)) 参数么)

它也会创建 `api/v1/cronjob_types.go` 文件并添加  `CronJob Kind`, 每次我们运行这个命令但是指定不同的 `Kind` 时, 他都会为我们创建一个 `xxx_types.go` 文件。
 
先让我们看看我们得到了什么开箱即用的东西，然后我们就开始填写它。

{{#literatego ./testdata/emptyapi-cn.go}}

现在，我们已经认识了基础的结构体，下一步我们开始填写它们！
