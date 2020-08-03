# kubebuilder-doc-cn

目前文档还在翻译中，会尽量以通俗易懂的方式翻译，和原始文档会有些差异，欢迎大家提出建议。

## Running mdBook

本书使用 [mdBook](https://github.com/rust-lang-nursery/mdBook) 构建. 如果想在本地测试更改，请执行以下步骤:

1. 根据此链接 [https://github.com/rust-lang-nursery/mdBook#installation](https://github.com/rust-lang-nursery/mdBook#installation) 安装 mdBook.
2. cd 到项目目录
3. 运行 `mdbook serve`
4. 如果报错 `not found controller-gen in $PATH`, 执行以下命令 

  ```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5
  ```
5. 浏览器打开 [http://localhost:3000](http://localhost:3000)