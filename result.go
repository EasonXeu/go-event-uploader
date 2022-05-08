package main

type Attempt struct {
	Success      bool
	RequestId    string
	HttpCode     string
	ErrorCode    string
	ErrorMessage string
	TimeStampMs  int64
	CostMs       int64
}

func createAttempt(success bool,
	requestId string, httpCode string, errorCode string, errorMessage string,
	timeStampMs int64, costMs int64) *Attempt {
	return &Attempt{
		Success:      success,
		RequestId:    requestId,
		HttpCode:     httpCode,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		TimeStampMs:  timeStampMs,
		CostMs:       costMs,
	}
}

type Result struct {
	attemptList []*Attempt
	success     bool
}

func (result *Result) successful() bool {
	return result.success
}

func (result *Result) GetReservedAttempts() []*Attempt {
	return result.attemptList
}

func (result *Result) GetHttpCode() string {
	if len(result.attemptList) == 0 {
		return ""
	}
	cursor := len(result.attemptList) - 1
	return result.attemptList[cursor].HttpCode
}

func (result *Result) GetErrorCode() string {
	if len(result.attemptList) == 0 {
		return ""
	}
	cursor := len(result.attemptList) - 1
	return result.attemptList[cursor].ErrorCode
}

func (result *Result) GetErrorMessage() string {
	if len(result.attemptList) == 0 {
		return ""
	}
	cursor := len(result.attemptList) - 1
	return result.attemptList[cursor].ErrorMessage
}

func (result *Result) GetRequestId() string {
	if len(result.attemptList) == 0 {
		return ""
	}
	cursor := len(result.attemptList) - 1
	return result.attemptList[cursor].RequestId
}

func (result *Result) GetTimeStampMs() int64 {
	if len(result.attemptList) == 0 {
		return 0
	}
	cursor := len(result.attemptList) - 1
	return result.attemptList[cursor].TimeStampMs
}

func (result *Result) GetCostMs() int64 {
	if len(result.attemptList) == 0 {
		return 0
	}
	cursor := len(result.attemptList) - 1
	return result.attemptList[cursor].CostMs
}

func initResult() *Result {
	return &Result{
		attemptList: []*Attempt{},
		success:     false,
	}
}
