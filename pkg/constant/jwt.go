package constant

import "os"

var SecretKey = []byte(os.Getenv("JWT_SECRET"))
