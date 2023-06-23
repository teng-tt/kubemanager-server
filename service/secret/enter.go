package secret

import "kubmanager/convert"

type ServicerGroup struct {
	SecretService
}

var secretConvert = convert.ConvertGroupApp.SecretConvert
