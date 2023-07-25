package fnOTP

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func Create(issuer, nm string) (key, secret string, err error) {
	var v *otp.Key
	if v, err = totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: nm,
	}); err != nil {
		return
	}
	return v.String(), v.Secret(), nil
}

func Verify(secret, key string) bool {
	return totp.Validate(key, secret)
}
