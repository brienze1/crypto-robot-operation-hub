package dto

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
)

// TODO validade if variables need to be exported
type AnalysisDto struct {
	Summary      enum.Summary `json:"summary"`
	Timestamp    string       `json:"timestamp"`
	AnalyzedData interface{}  `json:"analyzed_data"`
}

func (analysisDto AnalysisDto) ToAnalysis() model.Analysis {
	analysis := model.Analysis{
		Summary:      analysisDto.Summary,
		Timestamp:    analysisDto.Timestamp,
		AnalyzedData: "",
	}

	return analysis
}
