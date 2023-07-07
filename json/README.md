# JSON

参考 https://github.com/gin-gonic/gin/tree/master/internal/json ，可在构建时选用不同的 JSON 包

uses `encoding/json` as default json package but you can change it by build from other tags.

[jsoniter](https://github.com/json-iterator/go)

```sh
go build -tags=jsoniter .
```

[go-json](https://github.com/goccy/go-json)

```sh
go build -tags=go_json .
```

[sonic](https://github.com/bytedance/sonic) (you have to ensure that your cpu support avx instruction.)

```sh
$ go build -tags="sonic avx" .
```
