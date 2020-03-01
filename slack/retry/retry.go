package retry

import (
	"math/rand"
	"strings"
	"time"
)

type Errors struct {
	Errors []error
}

func NewRetryErrors() *Errors {
	return &Errors{Errors: []error{}}
}

func (e *Errors) Error() string {
	errs := []string{}
	for _, err := range e.Errors {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, "\n")
}

func (e *Errors) append(err error) {
	e.Errors = append(e.Errors, err)
}

// Interval returns RetryBackOff
func Interval(trial uint, function func() error, interval time.Duration) error {
	return BackOff(trial, interval, 0, function)
}

func BackOff(trial uint, meanInterval time.Duration, randFactor float64, function func() error) error {
	errors := NewRetryErrors()
	for trial > 0 {
		trial--
		if err := function(); err == nil {
			return nil
		} else {
			errors.append(err)
		}

		if trial <= 0 {
			break
		} else if randFactor <= 0 || meanInterval <= 0 {
			time.Sleep(meanInterval)
		} else {
			interval := randInterval(meanInterval, randFactor)
			time.Sleep(interval)
		}
	}

	if len(errors.Errors) > 0 {
		return errors
	}

	return nil
}

func randInterval(duration time.Duration, factor float64) time.Duration {
	interval := float64(duration)
	delta := factor * interval
	max := interval + delta
	min := interval - delta

	return time.Duration(min + (rand.Float64() * (max - min + 1)))
}
