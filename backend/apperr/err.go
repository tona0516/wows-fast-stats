package apperr

import (
	"fmt"
	"reflect"
)

type AppError struct {
	Code string
	Name string
	Raw  error
}

func (d AppError) WithRaw(err error) AppError {
	d.Raw = err

	return d
}

func (d AppError) Error() string {
	if d.Raw != nil {
		return fmt.Sprintf("%s %s %s", d.Code, d.Name, d.Raw.Error())
	}

	return fmt.Sprintf("%s %s", d.Code, d.Name)
}

type wg struct {
	AccountInfo       AppError
	AccountList       AppError
	ClansAccountInfo  AppError
	ClansInfo         AppError
	EncyclopediaShips AppError
	ShipsStats        AppError
	EncyclopediaInfo  AppError
	BattleArenas      AppError
	BattleTypes       AppError
}

type ns struct {
	Req   AppError
	Parse AppError
}

type cache struct {
	Serialize   AppError
	Deserialize AppError
}

type cfg struct {
	Read   AppError
	Update AppError
}

type ss struct {
	Save AppError
}

type tai struct {
	Get  AppError
	Save AppError
}

type unreg struct {
	Warship AppError
}

type ac struct {
	Parse     AppError
	Retry     AppError
	Read      AppError
	Unmarshal AppError
}

type srvcfg struct {
	InvalidInstallPath AppError
	InvalidAppID       AppError
	InvalidFontSize    AppError
	NoUserConfig       AppError
	InvalidPlayerName  AppError
	InvalidPattern     AppError
	MessageTooLong     AppError
}

type srvprep struct {
	DeleteCache AppError
}

type srvrw struct {
	NewWatcher  AppError
	WatcherAdd  AppError
	WatcherChan AppError
}

type srvss struct {
	SaveDialog AppError
	Canceled   AppError
}

type app struct {
	OpenDir AppError
}

func newDetailStruct[T any](codePrefix string, codeStart int) T {
	detailType := reflect.TypeOf(AppError{})

	var t T
	ps := reflect.ValueOf(&t)
	s := ps.Elem()

	for i := 0; i < s.NumField(); i++ {
		varName := s.Type().Field(i).Name
		f := s.FieldByName(varName)
		if f.Type() == detailType {
			detail := AppError{Code: fmt.Sprintf("%s%d", codePrefix, i+codeStart), Name: varName}
			f.Set(reflect.ValueOf(detail))
		}
	}

	return t
}

var (
	ErrNoTimeKey          = fmt.Errorf("no time key")
	ErrNoDataKey          = fmt.Errorf("no data key")
	ErrInvalidInstallPath = fmt.Errorf("選択したフォルダに「WorldOfWarships.exe」が存在しません。")
	ErrInvalidAppID       = fmt.Errorf("WG APIと通信できません。AppIDが間違っている可能性があります。")
	ErrInvalidFontSize    = fmt.Errorf("不正な文字サイズです。")
	ErrNoUserConfig       = fmt.Errorf("未設定の状態のため登録できません。「設定」から入力してください。")
	ErrInvalidPlayerName  = fmt.Errorf("存在しないプレイヤー名です。")
	ErrInvalidPattern     = fmt.Errorf("不正なパターンです。")
	ErrMessageTooLong     = fmt.Errorf("メモの文字数が140文字を超えています。")
)

//nolint:gochecknoglobals
var (
	Wg      = newDetailStruct[wg]("I", 100)
	Ns      = newDetailStruct[ns]("I", 200)
	Cache   = newDetailStruct[cache]("I", 300)
	Cfg     = newDetailStruct[cfg]("I", 400)
	Ss      = newDetailStruct[ss]("I", 500)
	Tai     = newDetailStruct[tai]("I", 600)
	Unreg   = newDetailStruct[unreg]("I", 700)
	Ac      = newDetailStruct[ac]("I", 800)
	SrvCfg  = newDetailStruct[srvcfg]("S", 100)
	SrvPrep = newDetailStruct[srvprep]("S", 200)
	SrvRw   = newDetailStruct[srvrw]("S", 300)
	SrvSs   = newDetailStruct[srvss]("S", 400)
	App     = newDetailStruct[app]("A", 100)
)
