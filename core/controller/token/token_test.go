package token

import (
	"testing"
)

func TestJWTHS256(t *testing.T) {
	jwt := Token{Type: Type{Typ: "JWT", Alg: "HS256"}}
	jwt.Payload = &Payload{}
	jwt.Payload.EnterpriseId = 1234
	jwt.Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRlcnByaXNlSWQiOiIxNjY4MTgyMDcyMTg3NDk4NDk4IiwiZXhwIjoxNjg5NTYzMTM3MTQyLCJyb2xlaWQiOjEsInJvbGVrZXkiOiJhZG1pbiIsInJvbGVuYW1lIjoiYWRtaW4iLCJ1c2VySWQiOiIxNjY4MTgyMDcyMTg3NDk4NDk4In0.OxCRvrLasZzZpixV6QsXyIUoZ7P-XIii7Uu4XjpQIpw`
	err := jwt.Verify("go-admin")
	if err != nil {
		t.Error(err)
	}
	jwt.Payload = nil
	p := jwt.GetPayLoad()
	println(p.EnterpriseId)
}
