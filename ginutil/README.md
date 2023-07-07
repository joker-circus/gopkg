# Gin

学习文献：[Gin 文档-中文翻译 | Go 技术论坛](https://learnku.com/docs/gin-gonic/2018/gin-readme/3819)、[Gin 中文文档 | 看云](https://www.kancloud.cn/shuangdeyu/gin_book/949411)

Go 常用 Web 框架：[go语言几个最快最好运用最广的web框架比较](https://www.cnblogs.com/desmond123/p/9821687.html)



## http 请求

```go
g.GET("/ping", func(ctx *gin.Context) {
  ctx.JSON(http.StatusOK, gin.H{
    "message": "pong",
  })
})
g.POST("/ping", func(ctx *gin.Context) {
  buf := new(bytes.Buffer)
  body, err := ctx.Request.GetBody()
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
    return
  }
  _, _ = buf.ReadFrom(body)
  fmt.Println(buf.String())
  ctx.JSON(http.StatusOK, gin.H{"message": "success",})
})
```



## http 参数

```go
// path路径参数
value := c.Param("key")

// post body 参数
body, err := ioutil.ReadAll(c.Request.Body)
if err != nil {
    c.JSON(http.StatusInternalServerError, utils.StdErrorf("fail to read the body of request, err: %s", err.Error()))
    return
}
role := gjson.GetBytes(body, "role")
if !role.Exists() {
    c.JSON(http.StatusBadRequest, utils.StdError("role is required!"))
    return
}
```



## Bind & Validator

Gin 默认 validator：将 validator 和 binding 绑定 [v8_to_v9](https://github.com/go-playground/validator/blob/master/_examples/gin-upgrading-overriding/v8_to_v9.go)

另一种校验器：[govalidator](https://github.com/asaskevich/govalidator)

以下展示的是 Gin validator 结合错误中文翻译

```go
import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var trans ut.Translator

var _ binding.StructValidator = &defaultValidator{}

// ValidateStruct 如果接收到的类型是一个结构体或指向结构体的指针，则执行验证。
func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		// 如果传递不合规则的值，则返回InvalidValidationError，否则返回nil。
		// 如果返回err != nil，可通过err.(validator.ValidationErrors)来访问错误数组。
		if err := v.validate.Struct(obj); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				sliceErrs := make([]string, 0, len(errs))
				
				//	for _, e := range errs {
				//		sliceErrs = append(sliceErrs, e.Translate(trans))
				//	}

				for k, v := range translate(err, obj) {
					sliceErrs = append(sliceErrs, k + "：" + v)
				}
				return errors.New(strings.Join(sliceErrs, "；") + "。")
			}

			return error(err)
		}
	}

	return nil
}

//translate 翻译工具
func translate(err error, s interface{}) map[string]string {
	r := make(map[string]string)
	t := reflect.TypeOf(s).Elem()
	for _, err := range err.(validator.ValidationErrors) {
		//使用反射方法获取struct种的json标签作为key --重点2
		var k string
		if field, ok := t.FieldByName(err.StructField()); ok {
			k = field.Tag.Get("json")
		}
		if k == "" {
			k = err.StructField()
		}
		r[k] = err.Translate(trans)
	}
	return r
}

// Engine 返回支持`StructValidator`实现的底层验证引擎
func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()

		zhLoc := zh.New()
		enLoc := en.New()
		uni := ut.New(enLoc, zhLoc)
		trans, _ = uni.GetTranslator("zh")
		err := zhtranslations.RegisterDefaultTranslations(v.validate, trans)
		if err != nil {
			panic(err)
		}

		// 设置 Json Tag 字段代替 StructFields 作为错误中的提示变量
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			jsonName := fld.Tag.Get("json")
			if len(jsonName) == 0 {
				return fld.Name
			}

			name := strings.SplitN(jsonName, ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})
    
		// 新项目使用最新 Tag validate，未使用旧 Tag binding
		//v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

```

在 mian.go 文件中进行添加

```go
package main

import "github.com/gin-gonic/gin/binding"

func main() {

	binding.Validator = new(defaultValidator)

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



## 路由参数过滤器

使用

```go
func NewRouter (prefix *gin.RouterGroup) *Router {
	router := Router{Prefix:prefix}
	decorator := NewDecorator(prefix.BasePath())
	router.Prefix.Use(decorator.CheckHeaderAuth(), decorator.CheckParamsIsValid())
	return &router
}
```



过滤器

```go
package manage

import (
	FUtils "git.ucloudadmin.com/sre/firefang/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Decorator struct {
	WhiteList map[string][]string
}

func NewDecorator (prefix string) *Decorator {
	whiteList := map[string][]string {
		prefix + "/login": {"GET"},
		prefix + "/user":  {"GET"},
	}
	return &Decorator{WhiteList:whiteList}
}

// 检测身份白名单
func (d *Decorator) checkAuthWhiteList (method string, path string) bool {
	methodList, ok := d.WhiteList[path]
	if ok {
		for _, v := range methodList {
			if method == v {
				return true
			}
		}
	}
	return false
}

// 检测身份
func (d *Decorator) CheckHeaderAuth () gin.HandlerFunc{
	return func(c *gin.Context) {
		user := FUtils.GetUserByHeader(c.Request.Header)
		if 0 == len(user) && !d.checkAuthWhiteList(c.Request.Method, c.FullPath()){
			c.JSON(http.StatusBadRequest, FUtils.StdError("user not found"))
			c.Abort()
			return
		}
	}
}

// 检测参数是否合法
func (d *Decorator) CheckParamsIsValid () gin.HandlerFunc{
	return func(c *gin.Context) {
		for _, value := range c.Params {
			// 如果参数包含空格
			if strings.Contains(value.Value, " ") {
				c.JSON(http.StatusBadRequest, FUtils.StdErrorf("[%s] cannot contain spaces!", value.Key))
				c.Abort()
				return
			}
			// 如果参数是id，校验uint64
			if strings.HasSuffix(value.Key, "id") {
				_, err := strconv.ParseUint(value.Value, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, FUtils.StdErrorf("[%s] not valid", value.Key))
					c.Abort()
					return
				}
			}
		}
	}
}

```



## 跨域中间件

参考：[跨域共享CORS详解及Gin配置跨域](https://www.cnblogs.com/you-men/p/14054348.html)、[跨域资源共享 CORS 详解 - 阮一峰](http://www.ruanyifeng.com/blog/2016/04/cors.html)

民用

```shell
package middlewares

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			// 可将将* 替换为指定的域名
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}


// 严格控制
func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method
        origin := c.Request.Header.Get("Origin") //请求头部
        if origin != "" {
            //接收客户端发送的origin （重要！）
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin) 
            //服务器支持的所有跨域请求的方法
            c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") 
            //允许跨域设置可以返回其他子段，可以自定义字段
            c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
            // 允许浏览器（客户端）可以解析的头部 （重要）
            c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers") 
            //设置缓存时间
            c.Header("Access-Control-Max-Age", "172800") 
            //允许客户端传递校验信息比如 cookie (重要)
            c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                                                                                          
        }

        //允许类型校验 
        if method == "OPTIONS" {
            c.JSON(http.StatusOK, "ok!")
        }

        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic info is: %v", err)
            }
        }()

        c.Next()
    }
}
```

[官方库](https://github.com/gin-gonic/contrib)

```go
import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	app.Use(cors.Default())

}
```



