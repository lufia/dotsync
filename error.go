package main

type alreadyShownError struct {
	err error
}

func (e *alreadyShownError) Error() string {
	return e.err.Error()
}

func (e *alreadyShownError) HasAlreadyShown() bool {
	return true
}
