package math

import "github.com/shopspring/decimal"

// AverageDecimals 计算多个 decimal.Decimal 的平均值
func AverageDecimals(decimals ...decimal.Decimal) decimal.Decimal {
	if len(decimals) == 0 {
		return decimal.Zero
	}
	var sum = decimal.Zero
	for _, d := range decimals {
		sum = sum.Add(d)
	}
	return sum.Div(decimal.NewFromInt(int64(len(decimals))))
}

// Sqrt 使用牛顿迭代法计算平方根
func Sqrt(x decimal.Decimal) decimal.Decimal {
	if x.LessThan(decimal.Zero) {
		return decimal.Zero
	}
	if x.Equal(decimal.Zero) {
		return decimal.Zero
	}

	// 初始猜测值
	z := x.Div(decimal.NewFromInt(2))
	// 迭代次数
	for i := 0; i < 10; i++ {
		z = z.Sub(z.Mul(z).Sub(x).Div(z.Mul(decimal.NewFromInt(2))))
	}
	return z
}

// StandardDeviation 计算多个 decimal.Decimal 的标准差
func StandardDeviation(decimals ...decimal.Decimal) decimal.Decimal {
	if len(decimals) == 0 {
		return decimal.Zero
	}

	// 计算平均值
	mean := AverageDecimals(decimals...)

	// 计算方差
	var variance = decimal.Zero
	for _, d := range decimals {
		diff := d.Sub(mean)
		variance = variance.Add(diff.Mul(diff))
	}
	variance = variance.Div(decimal.NewFromInt(int64(len(decimals))))

	// 计算标准差
	return Sqrt(variance)
}
