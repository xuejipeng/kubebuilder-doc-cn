# Summary

[引言](./introduction.md)

[快速入门](quick-start.md)

---

- [教程: 构建 CronJob](cronjob-tutorial/cronjob-tutorial.md)

  - [基础项目中什么?](./cronjob-tutorial/basic-project.md)
  - [每个程序都需要一个 main 入口](./cronjob-tutorial/empty-main.md)
  - [Groups、Versions 和 Kinds 之间的关系](./cronjob-tutorial/gvks.md)
  - [创建一个 API](./cronjob-tutorial/new-api.md)
  - [设计一个 API](./cronjob-tutorial/api-design.md)

      - [简要说明: 其他文件是干什么的?](./cronjob-tutorial/other-api-files.md)

  - [controller 中有什么?](./cronjob-tutorial/controller-overview.md)
  - [实现一个 controller](./cronjob-tutorial/controller-implementation.md)

    - [main 文件修改什么?](./cronjob-tutorial/main-revisited.md)

  - [实现 defaulting webhooks 和 validating webhooks](./cronjob-tutorial/webhook-implementation.md)
  - [运行和部署 controller](./cronjob-tutorial/running.md)

    - [部署 cert manager](./cronjob-tutorial/cert-manager.md)
    - [部署 webhooks](./cronjob-tutorial/running-webhook.md)

  - [结语](./cronjob-tutorial/epilogue.md)

- [教程: 多版本 API](./multiversion-tutorial/tutorial.md)

  - [Changing things up](./multiversion-tutorial/api-changes.md)
  - [Hubs, spokes, and other wheel metaphors](./multiversion-tutorial/conversion-concepts.md)
  - [Implementing conversion](./multiversion-tutorial/conversion.md)

      - [and setting up the webhooks](./multiversion-tutorial/webhooks.md)

  - [Deployment and Testing](./multiversion-tutorial/deployment.md)

---

- [Migrations](./migrations.md)

  - [Kubebuilder v1 vs v2](./migration/v1vsv2.md)

      - [Migration Guide](./migration/guide.md)

  - [Single Group to Multi-Group](./migration/multi-group.md)

---

- [Reference](./reference/reference.md)

  - [Generating CRDs](./reference/generating-crd.md)
  - [Using Finalizers](./reference/using-finalizers.md)
  - [Kind cluster](reference/kind.md)
  - [What's a webhook?](reference/webhook-overview.md)
    - [Admission webhook](reference/admission-webhook.md)
    - [Webhooks for Core Types](reference/webhook-for-core-types.md)
  - [Markers for Config/Code Generation](./reference/markers.md)

      - [CRD Generation](./reference/markers/crd.md)
      - [CRD Validation](./reference/markers/crd-validation.md)
      - [CRD Processing](./reference/markers/crd-processing.md)
      - [Webhook](./reference/markers/webhook.md)
      - [Object/DeepCopy](./reference/markers/object.md)
      - [RBAC](./reference/markers/rbac.md)

  - [controller-gen CLI](./reference/controller-gen.md)
  - [Artifacts](./reference/artifacts.md)
  - [Writing controller tests](./reference/writing-tests.md)

    - [Using envtest in integration tests](./reference/testing/envtest.md)

  - [Metrics](./reference/metrics.md)

---

[Appendix: The TODO Landing Page](./TODO.md)
