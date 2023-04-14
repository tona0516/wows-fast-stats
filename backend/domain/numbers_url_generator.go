package domain

import (
	"strconv"
	"strings"
)

const baseURL = "https://asia.wows-numbers.com/"

type NumbersURLGenerator struct {
}

func (n *NumbersURLGenerator) PlayerPage(accountID int, accountName string) string {
    return baseURL + "player/" + strconv.Itoa(accountID) + "," + accountName
}

func (n *NumbersURLGenerator) ShipPage(shipID int, shipName string) string {
    return baseURL + "ship/" + strconv.Itoa(shipID) + "," + strings.ReplaceAll(shipName, " ", "-")
}
