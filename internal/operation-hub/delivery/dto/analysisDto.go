package dto

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
)

type AnalysisDto struct {
	Summary   enum.SummaryEnum `json:"summary"`
	Timestamp string           `json:"timestamp"`
}

func (analysisDto AnalysisDto) ToAnalysis() model.Analysis {
	analysis := model.Analysis{
		Summary:   analysisDto.Summary,
		Timestamp: analysisDto.Timestamp,
	}

	return analysis
}
