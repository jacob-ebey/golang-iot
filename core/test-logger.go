package core

type TestLogger struct {
	LoggedErrors        []error
	LoggedErrorsChannel chan error
}

func (logger *TestLogger) LogError(err error) {
	logger.LoggedErrors = append(logger.LoggedErrors, err)
	if logger.LoggedErrorsChannel != nil {
		logger.LoggedErrorsChannel <- err
	}
}
