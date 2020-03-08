# 实现 defaulting webhooks 和 validating webhooks

如果你想为你的 CRD 实现 [admission webhooks](../reference/admission-webhook.md)，你只需要实现 `Defaulter` 和 (或)  `Validator`  接口即可。

其余的东西 Kubebuilder 会为你实现，比如：

1.  创建一个 webhook server
2.  确保这个 server 被添加到 manager 中
3.  为你的 webhooks 创建一个 handlers
4.  将每个 handler 以 path 形式注册到你的 server 中

首先，让我们为 CRD（CronJob）创建 webhooks，我们需要用到 `--defaulting` 和 `--programmatic-validation` 参数(因为我们的测试项目将使用 defaulting webhooks 和 validating webhooks)：

```bash
kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation
```
这将创建 Webhook 功能相关的方法，并在 `main.go` 中注册 Webhook 到你的 manager 中。

{{#literatego ./testdata/project/api/v1/cronjob_webhook.go}}
