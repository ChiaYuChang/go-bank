package errorcode

type errorRepo map[ErrorCode]*Error

func NewErrorRepo(opts ...option) (errorRepo, error) {
	repo := make(errorRepo)
	for _, opt := range opts {
		err := opt(repo)
		if err != nil {
			return repo, err
		}
	}
	return repo, nil
}

func NewErrorRepoWithDefaultErrors() errorRepo {
	repo, _ := NewErrorRepo(
		WithSuccess(),
		WithAuthErr(),
		WithClientErr(),
		WithServerErr(),
	)
	return repo
}

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
