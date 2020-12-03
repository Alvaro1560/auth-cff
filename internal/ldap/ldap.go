package ldap

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/parameters"

	ld "gopkg.in/ldap.v2"
)

// "LDAP" = t
// "LDAP_SERVER"
// "LDAP_PORT"
// "LDAP_DOMAIN"
// "LDAP_SSO" =  t  = Utiliza credenciales para ingreso SSO.
// "LDAP_USERNAME_SSO"
// "LDAP_PASSWORD_SSO"

func Authentication(username, bindUsername, bindPassword string) ([]string, error) {

	ldapDomain := parameters.GetParameter("LDAP_DOMAIN")
	ldapServer := parameters.GetParameter("LDAP_SERVER")
	ldapPort := parameters.GetParameter("LDAP_PORT")

	bindPassword = fmt.Sprintf("%s@%s", bindPassword, ldapDomain)
	var groups []string
	// conecta al directorio activo
	l, err := ld.Dial("tcp", fmt.Sprintf(`"%s":%s`, ldapServer, ldapPort))
	if err != nil {
		logger.Error.Printf("No se pudo conectar al Active Directory: %v", err)
		return groups, err
	}
	defer l.Close()

	// Autentica usuario
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		logger.Error.Printf("No se pudo autenticar al Active Directory: %v", err)
		return groups, err
	}
	domain := strings.Split(ldapDomain, ".")
	// Busqueda de un usuario por su account name
	searchRequest := ld.NewSearchRequest(
		fmt.Sprintf(`dc=%s,dc=%s`, domain[0], domain[1]),
		ld.ScopeWholeSubtree, ld.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(`(&(objectClass=*)(sAMAccountName=%s))`, username),
		[]string{"memberOf"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Error.Printf("No se pudo buscar el usuario en Active Directory: %v", err)
		return groups, err
	}

	if len(sr.Entries) != 1 {
		err = errors.New("No existe informacion de usuario en Active Directory")
		logger.Error.Printf("%v", err)
		return groups, err
	}

	for _, v := range sr.Entries {
		for _, k := range v.Attributes {
			for _, j := range k.Values {
				cn := strings.Split(j, ",")
				cnstr := strings.Replace(cn[0], "CN=", "", -1)
				groups = append(groups, cnstr)
			}
		}
	}
	return groups, nil
}
