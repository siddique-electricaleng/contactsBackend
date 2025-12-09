package config

import "contacts/internal/env"

const APIVersion = "v1"

// fix the error below
var JWT_SECRET []byte = []byte(env.GetStringNoFallback("JWT_SECRET"))
