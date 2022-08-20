package model

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum"

type Analysis struct {
	Summary      enum.SummaryEnum
	Timestamp    string
	AnalyzedData string
}
