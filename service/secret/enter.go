package secret

import "kubemanager.com/convert"

type ServicerGroup struct {
	SecretService
}

var secretConvert = convert.ConvertGroupApp.SecretConvert
