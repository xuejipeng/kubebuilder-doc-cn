# 运行和部署 controller

为了测试，我们可以在本地运行 controller。像[快速入门](..//quick-start.md)中说的，在运行 controller 之前，我们需要先安装 CRD。

一下命令会通过 controller-tools 更新我们的 YAML manifests 文件。

```bash
make install
```

现在，我们已经安装了 CRD，我们可以针对我们的集群运行 controller了，运行 controller 的证书和我们连接 k8s 集群的证书是同一个 ，因此我们现在不必担心 RBAC 权限问题。

<aside class="note"> 

<h1>在本地运行 webhooks</h1>

如果你想在本地运行 Webhooks，则必须生成用于服务 Webhooks 的证书，并把它放到正确的目录中(`/tmp/k8s-webhook-server/serving-certs/tls.{crt,key}`，默认情况下）。

如果你没有运行本地 API server，则还需要清楚如何将流量从远程群集代理到本地 Webhook 服务器。因此，在本地做 code-run-test 时，建议你关闭 webhooks ，如下所示。

</aside>

在另一个终端中，运行：

```bash
make run ENABLE_WEBHOOKS=false
```

您应该会从 controller 中看到有关启动的日志，但是它目前还没有做任何事。

此时，我们需要创建一个 CronJob 进行测试。让我们写个简单的例子放到 `config/samples/batch_v1_cronjob.yaml` 中，并运行该例子：

```yaml
{{#include ./testdata/project/config/samples/batch_v1_cronjob.yaml}}
```

```bash
kubectl create -f config/samples/batch_v1_cronjob.yaml
```

执行下面的命令，你应该会看到一连串的信息。如果你使用 `-w` 参数，应该会看到你的 `cronjob` 正在运行，并且正在更新状态：

```bash
kubectl get cronjob.batch.tutorial.kubebuilder.io -o yaml
kubectl get job
```

现在, 我们已经可以在集群中运行 CronJob 了。 停掉 `make run`  命令，然后运行

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
make deploy IMG=<some-registry>/<project-name>:tag
```

此时 `Webhook` 并不能正确运行，如果想 `Webhook` 正确运行，请参[运行 Webhook](./running-webhook.md)。

### 注意

1. 为了解决镜像拉不到的问题，可以将 `Dockerfile` 中的 
`FROM gcr.io/distroless/static:nonroot` 修改为 `FROM golang:1.13` 
或任意可正常拉取到的镜像。将 `/config/default/manager_auth_proxy_patch.yaml` 文件中的
 `gcr.io/kubebuilder/kube-rbac-proxy:v0.4.1` 改为 `xuejipeng/learn:kube-rbac-proxy-v0.4.1`

2. 为了加速构建，我们可以把 Makefile 中 `docker-build` 后面的 `test` 去掉



如果像以前一样再次get cronjobs ，我们应该可以看到 controller 再次运行！
