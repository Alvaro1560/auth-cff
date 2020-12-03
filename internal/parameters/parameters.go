package parameters

import "strings"

// APP_VALIDATE_IP = t(habilita validacion ip permitidas en login ) -  f(deshabilita validacion IP en login )
func GetParameter(parameter string) string {
	if strings.ToUpper(parameter) == "APP_VALIDATE_IP" {
		return "f"
	}
	if strings.ToUpper(parameter) == "RSA" {
		return `MIICXQIBAAKBgQCl6RCbisYho3oREByilGdicL3QA3iujpyQIDwd7f/wQv72cd9t\\n7dwc4h3MReo5a02BGGKLV83PTybPvepPsuiuyxuNDp1zvX3MH7GK8Ms3zqWWW/6b\\nMBlCWcxP4wLFwh804Ii+oYd6pyMmyfewI4/0EbAp09u9acSYYbZVZK73WQIDAQAB\\nAoGAWkcY11pK95Dp4hD/U+Qm4WTxlBffejR1suMncy3HX6hE7jsGVd3hMQJFLps4\\nmWfu83keXi43+j9aoh34OsfiXiIjF8g8mXBqx2mrgMWAx4YqrREmA/jQNap6660V\\nZr2o0AZBAre/rO22ghaPYSVRbm+ZNOSn9LZJhFTiOj0bwCUCQQDTJImclP0hikHK\\njBeOhuyRGKl1w349hAG+Sv+o79lw7bZcOrZ4lcik4zzOuaC1MMU88Im1MMkb1N4B\\njrefPUj/AkEAySh32mFLhb76o+s8QiyKE3mBWzjFnqt1uZY+0YXKee21Oay4IiZ2\\nEFLwYmPUX8JNGKW4nVKnTci258kJMOqnpwJBALfH0c/tDvems+VtUwPIBRm2caoD\\nY1qAEFRmS2nse0OEZXqZ0EHdfiunb0Iw6OVNciC87eA8epAzFJoec02ztlkCQQCK\\n7SYJbpTIJCPaPcZ6NWSPGqWaKNVjRiuiJv/vmKVEHEXDNWReQY3crEtUyHmOQRUZ\\n5qpgDQt7DxozboaogAeVAkAZNQa8dXFfMdFT8YoFNLR8z3EhEujl+Nx4boJ6O4oH\\nYAKafqdT/JTyZsvuKnqX7psEFcZDuy9U/QPt4cjEJstx`
	}
	return "Se debe implementar este metodo consultando la coleccion paramters"
}
