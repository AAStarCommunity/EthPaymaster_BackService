package envirment

import "sync"

type JWT struct {
	Security string
	Realm    string
	IdKey    string
}

var jwt *JWT

var onceJwt sync.Once

// GetJwtKey represents jwt object
func GetJwtKey() *JWT {
	onceJwt.Do(func() {
		if jwt == nil {
			j := GetAppConf().Jwt
			jwt = &JWT{
				Security: j.Security,
				Realm:    j.Realm,
			}
		}
	})

	return jwt
}