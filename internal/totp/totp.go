package totp

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateKey(username string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "OstoAssignment",
		AccountName: username,
	})
}

func VerifyOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}
