package model

import "github.com/rs/zerolog/log"

var logger = log.With().Str("module", "model").Logger()
