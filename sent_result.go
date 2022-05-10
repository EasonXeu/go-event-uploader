package main

import "fmt"

type ResultDetail struct {
	Success      bool
	RequestId    string
	HttpCode     string
	ErrorCode    string
	ErrorMessage string
	TimeStampMs  int64
	CostMs       int64
}

type SentResult struct {
	reservedResultDetails []*ResultDetail
	success               bool
}

func createResultDetail(success bool,
	requestId string, httpCode string, errorCode string, errorMessage string,
	timeStampMs int64, costMs int64) *ResultDetail {
	return &ResultDetail{
		Success:      success,
		RequestId:    requestId,
		HttpCode:     httpCode,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		TimeStampMs:  timeStampMs,
		CostMs:       costMs,
	}
}

func createResult() *SentResult {
	return &SentResult{
		reservedResultDetails: []*ResultDetail{},
		success:               false,
	}
}

func (result *SentResult) addResultDetail(resultDetail *ResultDetail) {
	result.reservedResultDetails = append(result.reservedResultDetails, resultDetail)
}

func (result *SentResult) successful() bool {
	return result.success
}

func (result *SentResult) GetReservedResultDetail() []*ResultDetail {
	return result.reservedResultDetails
}

func (result *SentResult) GetHttpCode() string {
	if len(result.reservedResultDetails) == 0 {
		return ""
	}
	cursor := len(result.reservedResultDetails) - 1
	return result.reservedResultDetails[cursor].HttpCode
}

func (result *SentResult) GetErrorCode() string {
	if len(result.reservedResultDetails) == 0 {
		return ""
	}
	cursor := len(result.reservedResultDetails) - 1
	return result.reservedResultDetails[cursor].ErrorCode
}

func (result *SentResult) GetErrorMessage() string {
	if len(result.reservedResultDetails) == 0 {
		return ""
	}
	cursor := len(result.reservedResultDetails) - 1
	return result.reservedResultDetails[cursor].ErrorMessage
}

func (result *SentResult) GetRequestId() string {
	if len(result.reservedResultDetails) == 0 {
		return ""
	}
	cursor := len(result.reservedResultDetails) - 1
	return result.reservedResultDetails[cursor].RequestId
}

func (result *SentResult) GetTimeStampMs() int64 {
	if len(result.reservedResultDetails) == 0 {
		return 0
	}
	cursor := len(result.reservedResultDetails) - 1
	return result.reservedResultDetails[cursor].TimeStampMs
}

func (result *SentResult) GetCostMs() int64 {
	if len(result.reservedResultDetails) == 0 {
		return 0
	}
	cursor := len(result.reservedResultDetails) - 1
	return result.reservedResultDetails[cursor].CostMs
}

func (result *SentResult) toString() string {
	return fmt.Sprintf("%#v", result)
}
