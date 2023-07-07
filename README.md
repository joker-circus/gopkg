# Gopkg

Go 常用工具包

## JSON

参考 https://github.com/gin-gonic/gin/tree/master/internal/json ，可在构建时选用不同的 JSON 包

构建：

```shell
go build -tags=jsoniter main.go

// 或者
go run -tags=jsoniter main.go
```



## Gin 工具包

### Bind &validator

Gin 默认 validator：将 validator 和 binding 绑定 [v8_to_v9](https://github.com/go-playground/validator/blob/master/_examples/gin-upgrading-overriding/v8_to_v9.go)

另一种校验器：[govalidator](https://github.com/asaskevich/govalidator)

以下展示的是 Gin validator 结合错误中文翻译

```go
package main

import (
    "github.com/gin-gonic/gin/binding"
    "github.com/joker-circus/gopkg/ginutil"
)
func main() {

	binding.Validator = new(ginutil.DefaultValidator)

	// regular gin logic
}
```

最后使用，只需要使用 tag 就行

```go
type Test struct {
	ID          int    `validate:"required"`             //数字确保不为0
	Name        string `validate:"required,min=1,max=8"` //字符串确保不为""，且长度 >=1 && <=8 （min=1,max=8等于gt=0,lt=9）
	Value       string `validate:"required,gte=1,lte=8"` //字符串确保不为""，且长度 >=1 && <=8
	Status      int    `validate:"min=1,max=10"`         //最小为0，最大为10（min=0,max=10等于gt=0,lt=11）
	PhoneNumber string `validate:"required,len=11"`      //不为""且长度为11
	Time        string `validate:"datetime=2006-01-02"`  //必须如2006-01-02的datetime格式
	Color       string `validate:"oneof=red green"`      //是能是red或者green
	Size        int    `validate:"oneof=37 39 41"`       //是能是37或者39或者41
	Email       string `validate:"email"`                //必须邮件格式
	JSON        string `validate:"json"`                 //必须json格式
	URL         string `validate:"url"`                  //必须url格式
	UUID        string `validate:"uuid"`                 //必须uuid格式
}

// ……
// 用 *gin.Context 的 Bind 方法即可
```
