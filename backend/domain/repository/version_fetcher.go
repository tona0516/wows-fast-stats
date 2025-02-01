package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type VersionFetcher interface {
	Fetch(currentSemver string) (model.LatestRelease, error)
}
