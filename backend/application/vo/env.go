package vo

type Env struct {
	// Note: set by ldflags in Makefile
	Str string
}

func (e *Env) IsProduction() bool {
	return e.Str == "production"
}
