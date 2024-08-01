package utils

import (
	"math"
	"strconv"
)

func ConvertTokenToNative(currentNativeUsdEquivalent, TokenAmount string) string {
	// Convert the strings to float64
	cusdAmountFloat, _ := strconv.ParseFloat(TokenAmount, 64)
	currentCeloUsdPriceFloat, _ := strconv.ParseFloat(currentNativeUsdEquivalent, 64)

	// Perform the conversion
	nativeAmount := cusdAmountFloat / currentCeloUsdPriceFloat
	nativeAmount = math.Round(nativeAmount*100) / 100

	// Convert the result back to a string
	nativeAmountStr := strconv.FormatFloat(nativeAmount, 'f', -1, 64)

	return nativeAmountStr

}

func CalculateOnePercent(amount string) string {
	// Convert the string to float64
	amountFloat, _ := strconv.ParseFloat(amount, 64)

	// Calculate 1% of the amount
	onePercent := amountFloat * 0.01

	// Convert the result back to a string
	onePercentStr := strconv.FormatFloat(onePercent, 'f', -1, 64)

	return onePercentStr
}

func PerformAmountCalculation(currentNativeUsdEquivalent, amount string) string {
	tokenAmount := CalculateOnePercent(amount)
	return ConvertTokenToNative(currentNativeUsdEquivalent, tokenAmount)

}

func ConvertAssetToFiat(currentEquivalent, amount string) string {
	// Convert the strings to float64
	amountFloat, _ := strconv.ParseFloat(amount, 64)
	currentEquivalentFloat, _ := strconv.ParseFloat(currentEquivalent, 64)

	// Perform the conversion
	fiatAmount := amountFloat * currentEquivalentFloat
	fiatAmount = math.Round(fiatAmount*100) / 100

	// Convert the result back to a string
	fiatAmountStr := strconv.FormatFloat(fiatAmount, 'f', -1, 64)

	return fiatAmountStr
}
