package wrap

type GSFError struct {
	err        error
	resultCode int32
	appCode    int32
}

func (e GSFError) Error() string {
	return e.err.Error()
}

func (e GSFError) Unwrap() error {
	return e.err
}

func (e GSFError) ResultCode() int32 {
	return e.resultCode
}

func (e GSFError) AppCode() int32 {
	return e.appCode
}

func WithGSFError(err error, resultCode int32, appCode int32) error {
	return GSFError{
		err:        err,
		resultCode: resultCode,
		appCode:    appCode,
	}
}
