package services

import (
	"regexp"
	"strings"
)

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

func IsValidNigerianPlate(plateNumber string) bool {
	plateNumber = strings.ToUpper(strings.TrimSpace(plateNumber))
	plateNumber = strings.ReplaceAll(plateNumber, " ", "")
	plateNumber = strings.ReplaceAll(plateNumber, "-", "")

	if len(plateNumber) == 0 {
		return false
	}

	if regexp.MustCompile(`^[A-Z]{2}[0-9]{3}[A-Z]{3}$`).MatchString(plateNumber) {
		stateCode := plateNumber[:2]
		if nigerianStates[stateCode] {
			return true
		}
	}

	if regexp.MustCompile(`^[A-Z]{3}[0-9]{3}[A-Z]{2}$`).MatchString(plateNumber) {
		stateCode := plateNumber[5:7]
		if nigerianStates[stateCode] {
			return true
		}
	}

	if regexp.MustCompile(`^[A-Z]{3}[0-9]{3}[A-Z]{1}$`).MatchString(plateNumber) {
		return true
	}

	if regexp.MustCompile(`^[A-Z]{3}[0-9]{3}[A-Z]{2,3}$`).MatchString(plateNumber) {
		return true
	}

	if regexp.MustCompile(`^CD[0-9]{3}[A-Z]{1}$`).MatchString(plateNumber) {
		return true
	}

	if regexp.MustCompile(`^(NA|AF|NN)[0-9]{3,4}[A-Z]?$`).MatchString(plateNumber) {
		return true
	}

	if regexp.MustCompile(`^[A-Z]{2,4}[0-9]{2,4}[A-Z]{1,3}$`).MatchString(plateNumber) {
		return true
	}

	return false
}
