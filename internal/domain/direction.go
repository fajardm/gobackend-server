package domain

import (
	"encoding/json"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type Direction int

const (
	DirectionAscending Direction = iota + 1
	DirectionDescending
)

var stringToDirection = map[string]Direction{
	"ASC":  DirectionAscending,
	"DESC": DirectionDescending,
}

var directionToString = map[Direction]string{
	DirectionAscending:  "ASC",
	DirectionDescending: "DESC",
}

func DirectionFromString(s string) (Direction, error) {
	d, ok := stringToDirection[s]
	if !ok {
		return 0, errors.New(errors.InvalidJSON, "invalid direction")
	}
	return d, nil
}

func (d Direction) String() string {
	return directionToString[d]
}

func (d Direction) MarshalText() ([]byte, error) {
	s := d.String()
	if s == "" {
		return nil, errors.New(errors.InvalidJSON, "invalid direction")
	}
	return []byte(s), nil
}

func (d *Direction) UnmarshalJSON(data []byte) error {
	var direction string
	if err := json.Unmarshal(data, &direction); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	di, err := DirectionFromString(direction)
	if err != nil {
		return err
	}
	*d = di
	return nil
}
