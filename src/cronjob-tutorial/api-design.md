# 设计一个 API

在 Kubernetes 中，我们有一些设计 API 的规范。
比如：所有可序列化字段必须是`camelCase`(驼峰) 格式，所以我们使用 JSON 的 tags 去指定它们，当一个字段可以为空时，我们会用 `omitempty` tag 去标记它。

字段类型大多是原始数据类型，Numbers 类型比较特殊：
出于兼容性的考虑，numbers 类型只接受三种数据类型：整形只能使用 `int32` 和 `int64` 声明, 
小数只能使用 `resource.Quantity` 声明

<details><summary>我的天, Quantity 又是什么鬼</summary>

Quantity 是十进制数字的一种特殊表示法，具有明确固定的表示形式，
使它们在计算机之间更易于移植。在Kubernetes中指定资源请求和Pod的限制时，您可能已经注意到它们。

从概念上讲，它们的工作方式类似于浮点数：它们具有有效位数，基数和指数。
它们的可序列化和人类可读格式使用整数和后缀来指定值，这与我们描述计算机存储的方式非常相似。

例如，该值2m表示0.002十进制表示法。 2Ki 表示2048十进制，而2K表示2000十进制。
如果要指定分数，请切换到一个后缀，该后缀使我们可以使用整数：2.5is 2500m。

支持两个基数：10和2（分别称为十进制和二进制）。
十进制基数用“常规” SI后缀（例如 M和K）表示，而二进制基数用“ mebi”表示法（例如Mi 和Ki）指定。想想[megabytes vs
mebibytes](https://en.wikipedia.org/wiki/Binary_prefix)。

</details>

我们还使用了另一种特殊类型：`metav1.Time`。它除了格式在 Kubernetes 中比较通用外，功能与 ` time.Time` 完全相同。

先让我们看一下 CronJob 应该是什么样子的！

{{#literatego ./testdata/project/api/v1/cronjob_types.go}}

现在，我们已经有了一个 API，然后我们需要编写一个 controller 来实现具体功能。