package usecase

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
	"wfs/backend2/common"
	"wfs/backend2/domain"
	"wfs/backend2/infra/clans"
	"wfs/backend2/infra/numbers"
	"wfs/backend2/infra/persistence"
	"wfs/backend2/infra/wargaming"

	"maps"

	"github.com/abadojack/whatlanggo"
	"github.com/samber/do"
	"golang.org/x/sync/errgroup"
)

type GetBattle struct {
	wargaming wargaming.API
	clans     clans.API
	numbers   numbers.API
}

func NewGetBattle(i *do.Injector) (*GetBattle, error) {
	return &GetBattle{
		wargaming: do.MustInvoke[wargaming.API](i),
		clans:     do.MustInvoke[clans.API](i),
		numbers:   do.MustInvoke[numbers.API](i),
	}, nil
}

func (g *GetBattle) Invoke(tempArenaInfo persistence.TempArenaInfo) (domain.Battle, error) {
	start := time.Now()

	warshipsChan := make(chan WarshipsResult)
	go g.getWarships(warshipsChan)

	battleTypesChan := make(chan common.Result[wargaming.BattleTypes])
	go func() {
		resp, err := g.wargaming.BattleTypes()
		battleTypesChan <- common.Result[wargaming.BattleTypes]{Value: resp, Error: err}
	}()

	battleArenasChan := make(chan common.Result[wargaming.BattleArenas])
	go func() {
		resp, err := g.wargaming.BattleArenas()
		battleArenasChan <- common.Result[wargaming.BattleArenas]{Value: resp, Error: err}
	}()

	accountNames := tempArenaInfo.AccountNames()

	accountList, err := g.wargaming.AccountList(accountNames)
	if err != nil {
		return domain.Battle{}, err
	}

	accountIDs := make([]int, 0, len(accountList))
	for _, v := range accountList {
		accountIDs = append(accountIDs, v.AccountID)
	}

	accountInfoChan := make(chan common.Result[wargaming.AccountInfo])
	go func() {
		resp, err := g.wargaming.AccountInfo(accountIDs)
		accountInfoChan <- common.Result[wargaming.AccountInfo]{Value: resp, Error: err}
	}()

	shipsStatsChan := make(chan ShipsStatsResult)
	go g.getShipsStats(accountIDs, shipsStatsChan)

	clansChan := make(chan ClansResult)
	go g.getClans(accountIDs, clansChan)

	warships := <-warshipsChan
	fmt.Println("warships", warships)

	battleTypes := <-battleTypesChan
	fmt.Println("battleTypes", battleTypes)

	battleArenas := <-battleArenasChan
	fmt.Println("battleArenas", battleArenas)

	accountInfo := <-accountInfoChan
	fmt.Println("accountInfo", common.PrettyJSON(accountInfo))

	shipsStats := <-shipsStatsChan
	fmt.Println("shipsStats", common.PrettyJSON(shipsStats))

	clans := <-clansChan
	fmt.Println("clans", clans)

	elapsed := time.Since(start)
	log.Println("Binomial took", elapsed)

	// TODO: Battleを生成

	return domain.Battle{}, nil
}

type EncycShipsMap map[int]wargaming.EncycShipsData
type EncycShipsResult common.Result[EncycShipsMap]

func (g *GetBattle) getEncycShips(resultChan chan EncycShipsResult) {
	firstResp, err := g.wargaming.EncycShips(1)
	if err != nil {
		resultChan <- EncycShipsResult{Error: err}

		return
	}

	result := make(EncycShipsMap)
	mu := sync.Mutex{}

	eg, _ := errgroup.WithContext(context.Background())
	for i := 2; i < firstResp.Meta.PageTotal+1; i++ {
		eg.Go(func() error {
			resp, err := g.wargaming.EncycShips(i)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			maps.Copy(result, resp.Data)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		resultChan <- EncycShipsResult{Error: err}

		return
	}

	resultChan <- EncycShipsResult{Value: result}
}

type WarshipsResult common.Result[domain.Warships]

func (g *GetBattle) getWarships(result chan WarshipsResult) {
	encycShipsChan := make(chan EncycShipsResult)
	go g.getEncycShips(encycShipsChan)

	numbersExpectedChan := make(chan common.Result[numbers.Expected])
	go func() {
		resp, err := g.numbers.Fetch()
		numbersExpectedChan <- common.Result[numbers.Expected]{Value: resp, Error: err}
	}()

	encycShipsResult := <-encycShipsChan
	if encycShipsResult.Error != nil {
		result <- WarshipsResult{Error: encycShipsResult.Error}

		return
	}

	expectedStats := <-numbersExpectedChan
	if expectedStats.Error != nil {
		result <- WarshipsResult{Error: expectedStats.Error}

		return
	}

	warships := make(domain.Warships)

	for shipID, ship := range encycShipsResult.Value {
		expected := expectedStats.Value.Data[shipID]

		warships[shipID] = domain.Warship{
			ID:            shipID,
			Name:          ship.Name,
			Tier:          ship.Tier,
			Type:          domain.NewShipType(ship.Type),
			Nation:        domain.Nation(ship.Nation),
			IsPremium:     ship.IsPremium,
			AverageDamage: expected.AverageDamageDealt,
			AverageFrags:  expected.AverageFrags,
			WinRate:       expected.WinRate,
		}
	}

	result <- WarshipsResult{Value: warships}
}

type ShipsStatsMap map[int]map[int]wargaming.ShipsStatsData
type ShipsStatsResult common.Result[ShipsStatsMap]

func (g *GetBattle) getShipsStats(accountIDs []int, resultChan chan ShipsStatsResult) {
	result := make(ShipsStatsMap)
	mu := sync.Mutex{}

	eg, _ := errgroup.WithContext(context.Background())
	for _, accountID := range accountIDs {
		eg.Go(func() error {
			resp, err := g.wargaming.ShipsStats(accountID)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			for _, stats := range resp {
				if _, ok := result[accountID]; !ok {
					result[accountID] = make(map[int]wargaming.ShipsStatsData)
				}

				for shipID, data := range stats {
					result[accountID][shipID] = data
				}
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		resultChan <- ShipsStatsResult{Error: err}

		return
	}

	resultChan <- ShipsStatsResult{Value: result, Error: nil}
}

type ClansResult common.Result[domain.Clans]

func (g *GetBattle) getClans(accountIDs []int, resultChan chan ClansResult) {
	clansAccountInfo, err := g.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		resultChan <- ClansResult{Error: err}

		return
	}

	clanIDs := make([]int, 0)

	for _, clan := range clansAccountInfo {
		if clan.ClanID == 0 {
			continue
		}

		clanIDs = append(clanIDs, clan.ClanID)
	}

	clansInfo, err := g.wargaming.ClansInfo(clanIDs)
	if err != nil {
		resultChan <- ClansResult{Error: err}

		return
	}

	clanTags := make([]string, 0, len(clansInfo))
	for _, clan := range clansInfo {
		if slices.Contains(clanTags, clan.Tag) {
			continue
		}

		clanTags = append(clanTags, clan.Tag)
	}

	autoCompletes := make(map[string]clans.Autocomplete)
	mu := sync.Mutex{}

	eg, _ := errgroup.WithContext(context.Background())
	for _, clanTag := range clanTags {
		eg.Go(func() error {
			resp, err := g.clans.FetchAutoComplete(clanTag)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			autoCompletes[clanTag] = resp

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		resultChan <- ClansResult{Error: err}

		return
	}

	result := make(domain.Clans)

	for _, accountID := range accountIDs {
		clanID := clansAccountInfo[accountID].ClanID
		if clanID == 0 {
			continue
		}

		clanInfo := clansInfo[clanID]
		clanTag := clanInfo.Tag
		hexColor := autoCompletes[clanTag].HexColor(clanTag)

		var urlRegExp = regexp.MustCompile(`https?://[^\s]+`)
		description := urlRegExp.ReplaceAllString(clanInfo.Description, "")
		description = strings.ReplaceAll(description, "\n", "")

		var lang string

		if len(description) > 0 {
			options := whatlanggo.Options{
				Whitelist: map[whatlanggo.Lang]bool{
					whatlanggo.Jpn: true,
					whatlanggo.Kor: true,
					whatlanggo.Cmn: true,
				},
			}

			info := whatlanggo.DetectWithOptions(description, options)
			lang = info.Lang.Iso6391()
		}

		result[accountID] = domain.Clan{
			ID:       clanID,
			Tag:      clanTag,
			HexColor: hexColor,
			Lang:     lang,
		}
	}

	resultChan <- ClansResult{Value: result}
}
