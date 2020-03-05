# 简要说明: 其他文件是干什么的?

[`api/v1/`](https://sigs.k8s.io/kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/api/v1)目录下除 `cronjob_types.go` 外还有另外两个文件：`groupversion_info.go`  和  `zz_generated.deepcopy.go`。

这两个文件都不需要编辑(前者保持不变，后者是自动生成的), 但是知道其中有什么是非常有用的。

## `groupversion_info.go`

`groupversion_info.go` 包含和 group-version 有关的元数据:

{{#literatego ./testdata/project/api/v1/groupversion_info.go}}

## `zz_generated.deepcopy.go`

`zz_generated.deepcopy.go` 包含[之前所说的](new-api.html#moment)由 `+kubebuilder:object:root`自动生成的 `runtime.Object` 接口的实现。

`runtime.Object` interface 的核心是一个 deep-copy 方法：`DeepCopyObject`。

controller-tools 中的 `object` 生成器也为每个 root type(`CronJob`) 和他的 sub-types(`CronJobList,CronJobSpec,CronJob1Status`) 都生成了两个方法：`DeepCopy`  和 `DeepCopyInto`

