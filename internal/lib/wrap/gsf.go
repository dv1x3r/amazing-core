package wrap

import "github.com/dv1x3r/amazing-core/internal/network/gsf/types/appcode"

type GSFError struct {
	err     error
	appCode appcode.AppCode
}

func (e GSFError) Error() string {
	return e.err.Error()
}

func (e GSFError) Unwrap() error {
	return e.err
}

func (e GSFError) AppCode() appcode.AppCode {
	return e.appCode
}

func WithGSFError(err error, appCode appcode.AppCode) error {
	return GSFError{
		err:     err,
		appCode: appCode,
	}
}
