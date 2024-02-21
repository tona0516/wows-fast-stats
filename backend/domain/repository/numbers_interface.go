package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type NumbersInterface interface {
	ExpectedStats() (model.ExpectedStats, error)
}
