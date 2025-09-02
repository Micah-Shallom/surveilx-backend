package services

import (
	"regexp"
	"strings"
)

// Nigerian states and their codes
var nigerianStates = map[string]bool{
	"AB": true, // Abia
	"AD": true, // Adamawa
	"AK": true, // Akwa Ibom
	"AN": true, // Anambra
	"BA": true, // Bauchi
	"BY": true, // Bayelsa
	"BE": true, // Benue
	"BO": true, // Borno
	"CR": true, // Cross River
	"DE": true, // Delta
	"EB": true, // Ebonyi
	"ED": true, // Edo
	"EK": true, // Ekiti
	"EN": true, // Enugu
	"GO": true, // Gombe
	"IM": true, // Imo
	"JI": true, // Jigawa
	"KD": true, // Kaduna
	"KN": true, // Kano
	"KT": true, // Katsina
	"KE": true, // Kebbi
	"KO": true, // Kogi
	"KW": true, // Kwara
	"LA": true, // Lagos
	"NA": true, // Nasarawa
	"NI": true, // Niger
	"OG": true, // Ogun
	"ON": true, // Ondo
	"OS": true, // Osun
	"OY": true, // Oyo
	"PL": true, // Plateau
	"RI": true, // Rivers
	"SO": true, // Sokoto
	"TA": true, // Taraba
	"YO": true, // Yobe
	"ZA": true, // Zamfara
	"FC": true, // Federal Capital Territory (Abuja)
}

// IsValidNigerianPlate validates if a given string is a valid Nigerian license plate
func IsValidNigerianPlate(plateNumber string) bool {
	// Clean the input
	plateNumber = strings.ToUpper(strings.TrimSpace(plateNumber))
	plateNumber = strings.ReplaceAll(plateNumber, " ", "")
	plateNumber = strings.ReplaceAll(plateNumber, "-", "")

	if len(plateNumber) == 0 {
		return false
	}

	// Pattern 1: New format - AA123BCD (2 letters + 3 digits + 3 letters)
	if regexp.MustCompile(`^[A-Z]{2}[0-9]{3}[A-Z]{3}$`).MatchString(plateNumber) {
		stateCode := plateNumber[:2]
		if nigerianStates[stateCode] {
			return true
		}
	}

	// Pattern 2: Old format - ABC123DE (3 letters + 3 digits + 2 letters)
	if regexp.MustCompile(`^[A-Z]{3}[0-9]{3}[A-Z]{2}$`).MatchString(plateNumber) {
		stateCode := plateNumber[5:7]
		if nigerianStates[stateCode] {
			return true
		}
	}

	// Pattern 3: Commercial/Government format - ABC123D (3 letters + 3 digits + 1 letter)
	if regexp.MustCompile(`^[A-Z]{3}[0-9]{3}[A-Z]{1}$`).MatchString(plateNumber) {
		return true
	}

	// Pattern 4: Mixed format - ABC123DE or similar variations (3 letters + 3 digits + 2+ letters)
	if regexp.MustCompile(`^[A-Z]{3}[0-9]{3}[A-Z]{2,3}$`).MatchString(plateNumber) {
		return true
	}

	// Pattern 5: Diplomatic format - CD123A (CD + 3 digits + 1 letter)
	if regexp.MustCompile(`^CD[0-9]{3}[A-Z]{1}$`).MatchString(plateNumber) {
		return true
	}

	// Pattern 6: Military format variations
	if regexp.MustCompile(`^(NA|AF|NN)[0-9]{3,4}[A-Z]?$`).MatchString(plateNumber) {
		return true
	}

	// Pattern 7: General Nigerian plate format - flexible validation
	if regexp.MustCompile(`^[A-Z]{2,4}[0-9]{2,4}[A-Z]{1,3}$`).MatchString(plateNumber) {
		return true
	}

	// Only return false if none of the patterns matched
	return false
}
