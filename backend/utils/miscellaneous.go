package utils

import (
	"backend/apis"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ConvertTokenToNative(currentNativeUsdEquivalent, TokenAmount string) string {
	// Convert the strings to float64
	cusdAmountFloat, _ := strconv.ParseFloat(TokenAmount, 64)
	currentCeloUsdPriceFloat, _ := strconv.ParseFloat(currentNativeUsdEquivalent, 64)

	// Perform the conversion
	nativeAmount := cusdAmountFloat / currentCeloUsdPriceFloat
	nativeAmount = math.Round(nativeAmount*100) / 100
	log.Println("nativeAmount: ", nativeAmount)

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

func LastPart(url, sep string) string {
	parts := strings.Split(url, sep)
	return parts[len(parts)-1]
}

func BadRequest(c *gin.Context, err error, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   err.Error(),
		"message": msg,
	})
}

// Valid and Invalid Countries
func CreateValidCountryCodes() map[string]bool {
	validCodes := map[string]bool{
		"AF": true, "AX": true, "AL": true, "DZ": true, "AS": true, "AD": true, "AO": true, "AI": true,
		"AQ": true, "AG": true, "AR": true, "AM": true, "AW": true, "AU": true, "AT": true, "AZ": true,
		"BS": true, "BH": true, "BD": true, "BB": true, "BY": true, "BE": true, "BZ": true, "BJ": true,
		"BM": true, "BT": true, "BO": true, "BQ": true, "BA": true, "BW": true, "BV": true, "BR": true,
		"IO": true, "BN": true, "BG": true, "BF": true, "BI": true, "CV": true, "KH": true, "CM": true,
		"CA": true, "KY": true, "CF": true, "TD": true, "CL": true, "CN": true, "CX": true, "CC": true,
		"CO": true, "KM": true, "CG": true, "CD": true, "CK": true, "CR": true, "CI": true, "HR": true,
		"CU": true, "CW": true, "CY": true, "CZ": true, "DK": true, "DJ": true, "DM": true, "DO": true,
		"EC": true, "EG": true, "SV": true, "GQ": true, "ER": true, "EE": true, "ET": true, "FK": true,
		"FO": true, "FJ": true, "FI": true, "FR": true, "GF": true, "PF": true, "TF": true, "GA": true,
		"GM": true, "GE": true, "DE": true, "GH": true, "GI": true, "GR": true, "GL": true, "GD": true,
		"GP": true, "GU": true, "GT": true, "GG": true, "GN": true, "GW": true, "GY": true, "HT": true,
		"HM": true, "VA": true, "HN": true, "HK": true, "HU": true, "IS": true, "IN": true, "ID": true,
		"IR": true, "IQ": true, "IE": true, "IM": true, "IL": true, "IT": true, "JM": true, "JP": true,
		"JE": true, "JO": true, "KZ": true, "KE": true, "KI": true, "KP": true, "KR": true, "KW": true,
		"KG": true, "LA": true, "LV": true, "LB": true, "LS": true, "LR": true, "LY": true, "LI": true,
		"LT": true, "LU": true, "MO": true, "MK": true, "MG": true, "MW": true, "MY": true, "MV": true,
		"ML": true, "MT": true, "MH": true, "MQ": true, "MR": true, "MU": true, "YT": true, "MX": true,
		"FM": true, "MD": true, "MC": true, "MN": true, "ME": true, "MS": true, "MA": true, "MZ": true,
		"MM": true, "NA": true, "NR": true, "NP": true, "NL": true, "NC": true, "NZ": true, "NI": true,
		"NE": true, "NG": true, "NU": true, "NF": true, "MP": true, "NO": true, "OM": true, "PK": true,
		"PW": true, "PS": true, "PA": true, "PG": true, "PY": true, "PE": true, "PH": true, "PN": true,
		"PL": true, "PT": true, "PR": true, "QA": true, "RE": true, "RO": true, "RU": true, "RW": true,
		"BL": true, "SH": true, "KN": true, "LC": true, "MF": true, "PM": true, "VC": true, "WS": true,
		"SM": true, "ST": true, "SA": true, "SN": true, "RS": true, "SC": true, "SL": true, "SG": true,
		"SX": true, "SK": true, "SI": true, "SB": true, "SO": true, "ZA": true, "GS": true, "SS": true,
		"ES": true, "LK": true, "SD": true, "SR": true, "SJ": true, "SZ": true, "SE": true, "CH": true,
		"SY": true, "TW": true, "TJ": true, "TZ": true, "TH": true, "TL": true, "TG": true, "TK": true,
		"TO": true, "TT": true, "TN": true, "TR": true, "TM": true, "TC": true, "TV": true, "UG": true,
		"UA": true, "AE": true, "GB": true, "US": true, "UM": true, "UY": true, "UZ": true, "VU": true,
		"VE": true, "VN": true, "VG": true, "VI": true, "WF": true, "EH": true, "YE": true, "ZM": true,
		"ZW": true,
	}
	return validCodes
}

var blockedCountriesForVirtualAccounts = map[string]bool{
	"AF": true, // Afghanistan
	"BD": true, // Bangladesh
	"BI": true, // Burundi
	"BT": true, // Bhutan
	"CD": true, // Democratic Republic of the Congo
	"CN": true, // China
	"CU": true, // Cuba
	"DZ": true, // Algeria
	"GW": true, // Guinea-Bissau
	"HT": true, // Haiti
	"IQ": true, // Iraq
	"IR": true, // Iran
	"KE": true, // Kenya
	"KP": true, // North Korea
	"LB": true, // Lebanon
	"LY": true, // Libya
	"MA": true, // Morocco
	"MK": true, // North Macedonia
	"ML": true, // Mali
	"MM": true, // Myanmar
	"MZ": true, // Mozambique
	"NE": true, // Niger
	"NI": true, // Nicaragua
	"NP": true, // Nepal
	"PK": true, // Pakistan
	"QA": true, // Qatar
	"RU": true, // Russia
	"SD": true, // Sudan
	"SI": true, // Slovenia
	"SO": true, // Somalia
	"SS": true, // South Sudan
	"SY": true, // Syria
	"VE": true, // Venezuela
	"YE": true, // Yemen
	"ZW": true, // Zimbabwe
}

func IsBlockedCountry(alpha2Code string) bool {
	return blockedCountriesForVirtualAccounts[strings.ToUpper(alpha2Code)]
}