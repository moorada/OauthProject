package main

import (
	"encoding/json"
	"log"

	"github.com/dgrijalva/jwt-go"

	"github.com/MicahParks/keyfunc"
)

func main() {

	// Get the JWKS as JSON.
	var jwksJSON json.RawMessage = []byte(`{"keys":[{"use":"sig","kty":"RSA","kid":"public:44155750-30e3-457d-b59d-585b2fd71fbd","alg":"RS256","n":"4WkukAKi35Ep-hQogVzsQUgp0-U1VzepITA9VVyNfwL5EyAagsZg_36CzmuCkB93qn9gusrqIjleRHWA0feJBQrOAKuNU2L2ZNwPvogUSffU0hl4axPjmg_HOgq3qtSrQuVlHq3p0LhtvHDbJ5cZggwsZ-KTfA7-29iu3LNrjCddTi1msxydKKlrqXv1Ct5zyb7OrXUUXIzvodx8_MK8xQs9JAAlEF5N1K-CE0ewNq20QgMC_dVOxoAgoLnKKNR0m_TqhHnJPSmG39zP2r9089p4DBODYJvae99UnF6vI8wdthkXBZFP3VvlJ_WSElGTuxeGFjhOvHEVrWtH1Fx15s4ev8Q_bPSxOrYkLgRH2L4yzkgH0tsxy6CHwfVV3NO3FFuBQTswGGSnt7x91d3SsUZBa_8coXv7jkrntTw3OFDzxkTvWizNma2b8FEshxA3dsi0EP9ZDE0uP1M7OsDRycyl5lvIpgNwN0hgwUKJzKs26xIwSF8Mcmzbx_6caXHQOeQdIWttukUZaw7zsHvf4Kot2hM1IANLqsFtjmoWkBolKelvnaoT1_BVeyTbZ0XxHPiSJ0PLjlkPRByK9lYydg50uT8du5OKg_XZ4k5p-hk_caNfh3-hZ8f83X0LyYRtpCYid3KSRXYNoqvF3rM4Jdn5o_wMqDjYb7gHlqe4tQM","e":"AQAB"}]}`)

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.New(jwksJSON)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}

	// Get a JWT to parse.
	jwtB64 := "eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzo0NDE1NTc1MC0zMGUzLTQ1N2QtYjU5ZC01ODViMmZkNzFmYmQifQ.eyJhdF9oYXNoIjoiMUwtRVN0X3R6UEp0S0x6WnhFX3NxZyIsImF1ZCI6WyJteWNsaWVudE9wZW5JZCJdLCJhdXRoX3RpbWUiOjE2MjY3MjQ3MzUsImV4cCI6MTYyNjcyODMzNiwiZXh0cmFfdmFycyI6eyJjb2dub21lIjoiUm9idXN0ZWxsYSIsImVtYWlsIjoiYWxlc3Npb0Byb2J1c3RlbGxhLmNvbSIsImlkIjoiUzk2MDIyOCIsIm5vbWUiOiJBbGVzc2lvIn0sImlhdCI6MTYyNjcyNDczNiwiaXNzIjoiaHR0cDovL2h5ZHJhLnRlc3QvIiwianRpIjoiOWRjNjQ4Y2QtNGJiNi00YjA1LThkNTctMTNmNGM1NTFhZTQyIiwicmF0IjoxNjI2NzI0NzMzLCJzaWQiOiI2MWY3NTJkNy02MjJmLTQyYzYtYWNmZi1iYWFkNmU4ODA5NGQiLCJzdWIiOiJTOTYwMjI4In0.2HZXtCMaMqvV4UoSQ0SPfh7UTx7taTIqDxltZmvUrUf3ktT8Um7pH-MkuYaslhNJOMcWRUqjuAl2VVaSjnxJ9NctRwlaZytma8sempIWCt5FY_s_BA7ClWhOvQybyeu6yxZ13jIu_sftxpBAX_h61xzpek4WcquoYmmwOjuOR99hb3D82z0UWxbcqxXVocCEZalgjCtS5ySx3i8siY5t2hnvxd6r1v3qX3pB7oE3jL1iKWbei4tGjs1cZ8wyhXtKpq59Xa4x4T5CQONg4E3VscLJOlmao6iF1jXVYE_UEks3o9sWqhgHqU4COH17tYtP0jfOB0wnTNUbTcnR_ubciPllh4vtyBcLe-iyiYDtMo_RjSeFdzZci99DTd1iV8RAWTw0q97-XHsNRVHOSKhDj48ELNRP9WtNHCaa-58lc5KVqXxRzzBi8HmTiF0DTw4bmEofUGhc1FQX9DVNHugAbpBA4mddyCaAdkOiJR8dsFkxLztOW4Uzh1ava3l31lMLhS6dLtYNvA91eue_OriFO23Z3k0bb9QG1dt0-hUVeX0kb9HWYbBC9GHEHwpQN1bGpdxx4_lShieRuW7NT6gcisGpHak62LMQ2s0q9Qvu-prV04oKZvBKjltqaYoidsfGZY6cf30IOXIE2SnHgJJ-fueMCHRWI5i6u0a720XQ8IM"

	// Parse the JWT.
	token, err := jwt.Parse(jwtB64, jwks.KeyFuncLegacy)
	if err != nil {
		log.Fatalf("Failed to parse the JWT.\nError: %s", err.Error())
	}

	// Check if the token is valid.
	if !token.Valid {
		log.Fatalf("The token is not valid.")
	}

	log.Println("The token is valid.")
}