package security

import (
	"crypto/subtle"
)

func SecureCompare(given, actual []byte) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare(given, actual) == 1
	}
	// Securely compare actual to itself to keep constant time, but always return false
	return subtle.ConstantTimeCompare(actual, actual) == 1 && false
}

func SecureCompareString(given, actual string) bool {
	// The following code is incorrect:
	// return SecureCompare([]byte(given), []byte(actual))

	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	}
	// Securely compare actual to itself to keep constant time, but always return false
	return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
}
