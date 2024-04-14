package service

import (
	"fmt"
	"github.com/MrKrisYu/koi-go-common/logger"
)

type Service struct {
	Logger *logger.Helper
	Error  error
}

func (s *Service) AddError(err error) error {
	if s.Error == nil {
		s.Error = err
	} else if err != nil {
		s.Error = fmt.Errorf("%v; %w", s.Error, err)
	}
	return s.Error
}
