package group

import "github.com/rs/zerolog/log"

var logger = log.With().Str("module", "group").Logger()
