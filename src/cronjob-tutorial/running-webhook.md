# 部署 Admission Webhooks

## Kind Cluster

建议使用 [kind](../reference/kind.md) 集群开发 Webhook，以加快迭代速度。 为什么？

- 你可以在1分钟内在本地启动多节点群集。
- 你可以在几秒钟内将其拆除。
- 你无需将 images 推送到远程镜像仓库。

## Cert Manager

你需要 [按此](./cert-manager.md) 安装 cert manager bundle。你只需要安装就好，对于证书的申请 kubebuilder 会帮你做。

## Build your image

运行以下命令以在本地生成 image。

```bash
make docker-build
```

如果你使用 `kind` 创建的群集，则无需将 image 推送到远程镜像仓库。你可以将本地的 image 直接加载到 `kind` 创建的群集：

```bash
kind load docker-image your-image-name:your-tag
```

## 部署 Webhooks

您可以通过 kustomize 启动 webhook 和 cert manager 配置，将 `config/default/kustomization.yaml` 
改成如下所示：

```yaml
{{#include ./testdata/project/config/default/kustomization.yaml}}
```

现在你可以通过执行下面的命令将它们部署到你的集群中

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
make deploy IMG=<some-registry>/<project-name>:tag
```

稍等片刻，直到出现 webhook pod 启动并且证书被提供。它通常在1分钟内完成。

现在，您可以创建一个有效的 CronJob 来测试你的 webhooks，创建应该成功完成。

```bash
kubectl create -f config/samples/batch_v1_cronjob.yaml
```

您也可以尝试创建一个无效的 CronJob (例如，使用格式错误的 cron schedule 字段),您应该看到创建失败并显示  validation error。

<aside class="note warning">

<h1>The Bootstrapping Problem</h1>

如果你在集群中通过 Webhook 部署 pods, 请注意 Bootstrapping Problem，因为创建 Pod 的请求将会发给 Webhook, 但是在 Webhook 启动之前它还不能处理创建 Pod 的请求, 它们互相矛盾。

为了使其工作，如果你的 kubernetes 版本为 1.9+，则可以使用[namespaceSelector]；如果你的 kubernetes 版本为 1.15+，则可以使用 [objectSelector] 来忽略自己。

</aside>

[namespaceSelector]: https://github.com/kubernetes/api/blob/kubernetes-1.14.5/admissionregistration/v1beta1/types.go#L189-L233
[objectSelector]: https://github.com/kubernetes/api/blob/kubernetes-1.15.2/admissionregistration/v1beta1/types.go#L262-L274
