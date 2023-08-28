package services

import "github.com/rs/zerolog/log"

var logger = log.With().Str("module", "services").Logger()
