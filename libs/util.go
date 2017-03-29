package gormigrate

import "time"

func Now() time.Time {
	return time.Now()
}

func UTC() time.Time {
	return time.Now().UTC()
}
