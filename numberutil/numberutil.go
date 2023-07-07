package numberutil

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Uint64ToFloat64(num uint64) float64 {
	f, _ := strconv.ParseFloat(strconv.FormatUint(num, 10), 64)
	return f
}

// Round 按精度对浮点数四舍五入
// If places < 0, it will round the integer part to the nearest 10^(-places).
func Round(num float64, prec int32) float64 {
	f, _ := decimal.NewFromFloat(num).Round(prec).Float64()
	return f
}

// Precision 返回浮点数精度
func Precision(num float64) int {
	str := strconv.FormatFloat(num, 'f', -1, 64)
	i := strings.IndexByte(str, '.')
	if i == -1 {
		return 0
	}
	return len(str) - i - 1
}

// IsEqual 判断浮点数是否相等
func IsEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < 0.00000001
}

// IsEqualInPrec 按精度判断浮点数是否相等
func IsEqualInPrec(f1, f2 float64, prec int) bool {
	return math.Abs(f1-f2) < math.Pow10(-prec)
}

func MaxInt(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MinInt(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}

const (
	AddComputeType   = "add"
	SubComputeType   = "sub"
	MultiComputeType = "multi"
	DivComputeType   = "div"
)

func float64MCompute(computeNum1 float64, computeNum2 float64, computeType string) (float64, error) {
	value1 := decimal.NewFromFloat(computeNum1)
	value2 := decimal.NewFromFloat(computeNum2)

	var value decimal.Decimal

	switch computeType {
	case AddComputeType:
		value = value1.Add(value2)
	case SubComputeType:
		value = value1.Sub(value2)
	case MultiComputeType:
		value = value1.Mul(value2)
	case DivComputeType:
		value = value1.Div(value2)
	default:
		return 0, errors.New("compute type error")
	}

	// return value is float64, exact, exact is a bool indicating whether f represents d exactly
	float64Value, _ := value.Float64()

	return float64Value, nil
}

func AddFloat64(computeNum1 float64, computeNum2 float64) float64 {
	num, _ := float64MCompute(computeNum1, computeNum2, AddComputeType)
	return num
}

func SubFloat64(computeNum1 float64, computeNum2 float64) float64 {
	num, _ := float64MCompute(computeNum1, computeNum2, SubComputeType)
	return num
}

func MulFloat64(computeNum1 float64, computeNum2 float64) float64 {
	num, _ := float64MCompute(computeNum1, computeNum2, MultiComputeType)
	return num
}

func DivFloat64(computeNum1 float64, computeNum2 float64) float64 {
	num, _ := float64MCompute(computeNum1, computeNum2, DivComputeType)
	return num
}

// 千分位分隔符
func ThousandsInt(i int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", i)
}
