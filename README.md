# kubebuilder-doc-cn

目前文档还在翻译中，欢迎大家提出建议。

## Running mdBook

The kubebuilder book is served using [mdBook](https://github.com/rust-lang-nursery/mdBook). If you want to test changes to the book locally, follow these directions:

1. Follow the instructions at [https://github.com/rust-lang-nursery/mdBook#installation](https://github.com/rust-lang-nursery/mdBook#installation) to
  install mdBook.
2. cd into the repo directory
3. Run `mdbook serve`
4. If you  get Error `not found controller-gen in $PATH`, you can do 

  ```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5
  ```
5. Visit [http://localhost:3000](http://localhost:3000)