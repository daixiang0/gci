package specificity

type Module struct{}

func (m Module) IsMoreSpecific(than MatchSpecificity) bool {
	return isMoreSpecific(m, than)
}

func (m Module) Equal(to MatchSpecificity) bool {
	return equalSpecificity(m, to)
}

func (Module) class() specificityClass {
	return ModuleClass
}

func (Module) String() string {
	return "Module"
}
