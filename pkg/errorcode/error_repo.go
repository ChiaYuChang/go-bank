package errorcode

type errorRepo map[ErrorCode]*Error

func (er errorRepo) RegisterErr(code ErrorCode, httpStatusCode int, message string) error {
	_, ok := er.GetErr(code)
	if ok {
		return ErrErrorCodeHasBeenUsed
	}
	er[code] = NewError(code, httpStatusCode, message)
	return nil
}

func (er errorRepo) RegisterErrFromErr(err error, code ErrorCode, httpStatusCode int) error {
	return er.RegisterErr(code, httpStatusCode, err.Error())
}

func (er errorRepo) GetErr(code ErrorCode) (*Error, bool) {
	e, ok := er[code]
	if !ok {
		return nil, false
	}
	return e.Clone(), true
}

func (er errorRepo) MustGetErr(code ErrorCode) *Error {
	e, _ := er.GetErr(code)
	return e
}
