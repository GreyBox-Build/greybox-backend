package utils

import (
	"backend/apis"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func ConvertTokenToNative(currentNativeUsdEquivalent, TokenAmount string) string {
	// Convert the strings to float64
	cusdAmountFloat, _ := strconv.ParseFloat(TokenAmount, 64)
	currentCeloUsdPriceFloat, _ := strconv.ParseFloat(currentNativeUsdEquivalent, 64)

	// Perform the conversion
	nativeAmount := cusdAmountFloat / currentCeloUsdPriceFloat
	nativeAmount = math.Round(nativeAmount*100) / 100
	fmt.Println("nativeAmount: ", nativeAmount)

	// Convert the result back to a string
	nativeAmountStr := strconv.FormatFloat(nativeAmount, 'f', -1, 64)

	return nativeAmountStr

}

// this now calculate 0.5 percent
func CalculateOnePercent(amount string) string {
	// Convert the string to float64
	amountFloat, _ := strconv.ParseFloat(amount, 64)

	// Calculate 0.5% of the amount
	onePercent := amountFloat * 0.005

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

func PerformDepositofNativeCalculation(amount, fiatCurrency, assetCurrency string) (string, error) {

	actualAmount := GetRemainingOnePercent(amount)

	resultChan := make(chan string)
	errChan := make(chan error)

	go apis.GetExchangeRate(fiatCurrency, assetCurrency, resultChan, errChan)

	select {
	case exchangeRate := <-resultChan:
		rate, err := strconv.ParseFloat(exchangeRate, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse exchange rate: %v", err)
		}
		actualAmountFloat, err := strconv.ParseFloat(actualAmount, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse actual amount: %v", err)
		}

		nativeAmount := actualAmountFloat / rate

		// Return the calculated amount as a string
		return fmt.Sprintf("%.8f", nativeAmount), nil

	case err := <-errChan:
		return "", err
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("timeout occurred while fetching exchange rate")
	}
}

func GetRemainingOnePercent(amount string) string {
	// Convert the string to float64
	amountFloat, _ := strconv.ParseFloat(amount, 64)

	// Calculate 1% of the amount
	onePercent := amountFloat / 100
	// Convert the result back to a string
	onePercentStr := strconv.FormatFloat(onePercent, 'f', -1, 64)

	return onePercentStr
}

func FormatAmountWithCommas(amount float64) string {
	amountStr := strconv.FormatFloat(amount, 'f', 2, 64)
	parts := strings.Split(amountStr, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	n := len(integerPart)
	if n <= 3 {
		return integerPart + "." + decimalPart
	}

	var result strings.Builder
	for i, digit := range integerPart {
		if (n-i)%3 == 0 && i != 0 {
			result.WriteByte(',')
		}
		result.WriteByte(byte(digit))
	}

	return result.String() + "." + decimalPart
}
