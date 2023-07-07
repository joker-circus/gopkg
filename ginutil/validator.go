package ginutil

/**
设置 Gin 在 Bind 数据时，默认的数据格式校验，添加中文错误提示
*/

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
				for _, e := range errs {
					sliceErrs = append(sliceErrs, e.Translate(trans))
				}
				return errors.New(strings.Join(sliceErrs, "；") + "。")
			}

			return error(err)
		}
	}

	return nil
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
