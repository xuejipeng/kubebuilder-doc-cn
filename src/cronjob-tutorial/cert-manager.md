# 部署 cert manager

我们建议使用 [cert manager](https://github.com/jetstack/cert-manager) 为 Webhook 服务器配置证书，
如果使用其他方法，将证书放在正确的位置，它们也会起作用。

您可以根据 
[the cert manager documentation](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html) 
安装 cert manager。

Cert manager 还具有一个名为 CA injector 的组件，
该组件负责将 CA bundle 注入到 Mutating | ValidatingWebhookConfiguration 中。

为此，您需要在 Mutating | ValidatingWebhookConfiguration 对象中使用 key 为 
` certmanager.k8s.io/inject-ca-from`  的 annotation，
annotation 的 value 应该以 `<certificate-namespace>/<certificate-name>` 的格式指向现有证书 CR 实例。

这是带有 Mutating | ValidatingWebhookConfiguration 对象的 
[kustomize](https://github.com/kubernetes-sigs/kustomize) 文件补丁。

```yaml
{{#include ./testdata/project/config/default/webhookcainjection_patch.yaml}}
```
