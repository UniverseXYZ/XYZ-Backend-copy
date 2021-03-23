package core

import (
	"github.com/barnbridge/barnbridge-backend/eth/bestblock"
	"github.com/barnbridge/barnbridge-backend/processor"
	"github.com/barnbridge/barnbridge-backend/scraper"
	"github.com/barnbridge/barnbridge-backend/taskmanager"
)

type Features struct {
	Backfill    bool
	Lag         FeatureLag
	Automigrate bool
	Uncles      bool
}

type FeatureLag struct {
	Enabled bool
	Value   int64
}

type Config struct {
	BestBlockTracker         bestblock.Config
	TaskManager              taskmanager.Config
	Scraper                  scraper.Config
	PostgresConnectionString string
	Features                 Features
	AbiPath                  string
	Processor                processor.Config
}
