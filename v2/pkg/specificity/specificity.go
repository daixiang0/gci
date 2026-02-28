package specificity

type MatchSpecificity interface {
	IsMoreSpecific(other MatchSpecificity) bool
	Equal(other MatchSpecificity) bool
}

type MisMatch struct{}

func (m MisMatch) IsMoreSpecific(other MatchSpecificity) bool {
	// MisMatch is the least specific, it's never more specific than anything
	return false
}

func (m MisMatch) Equal(other MatchSpecificity) bool {
	_, ok := other.(MisMatch)
	return ok
}

type StandardMatch struct{}

func (s StandardMatch) IsMoreSpecific(other MatchSpecificity) bool {
	_, isMisMatch := other.(MisMatch)
	_, isDefault := other.(DefaultMatch)
	return isMisMatch || isDefault
}

func (s StandardMatch) Equal(other MatchSpecificity) bool {
	_, ok := other.(StandardMatch)
	return ok
}

type DefaultMatch struct{}

func (d DefaultMatch) IsMoreSpecific(other MatchSpecificity) bool {
	// Default is only more specific than MisMatch
	_, isMisMatch := other.(MisMatch)
	return isMisMatch
}

func (d DefaultMatch) Equal(other MatchSpecificity) bool {
	_, ok := other.(DefaultMatch)
	return ok
}

type Match struct {
	Length int
}

func (m Match) IsMoreSpecific(other MatchSpecificity) bool {
	if _, ok := other.(MisMatch); ok {
		return true
	}
	if _, ok := other.(DefaultMatch); ok {
		return true
	}
	if _, ok := other.(StandardMatch); ok {
		return true
	}
	if otherMatch, ok := other.(Match); ok {
		return m.Length > otherMatch.Length
	}
	return false
}

func (m Match) Equal(other MatchSpecificity) bool {
	if otherMatch, ok := other.(Match); ok {
		return m.Length == otherMatch.Length
	}
	return false
}

type NameMatch struct{}

func (n NameMatch) IsMoreSpecific(other MatchSpecificity) bool {
	if _, isMisMatch := other.(MisMatch); isMisMatch {
		return true
	}
	if _, isDefault := other.(DefaultMatch); isDefault {
		return true
	}
	if _, isStandard := other.(StandardMatch); isStandard {
		return true
	}
	if _, isMatch := other.(Match); isMatch {
		return true
	}
	return false
}

func (n NameMatch) Equal(other MatchSpecificity) bool {
	_, ok := other.(NameMatch)
	return ok
}

type LocalModule struct{}

func (l LocalModule) IsMoreSpecific(other MatchSpecificity) bool {
	_, isMisMatch := other.(MisMatch)
	_, isDefault := other.(DefaultMatch)
	_, isStandard := other.(StandardMatch)
	_, isMatch := other.(Match)
	_, isName := other.(NameMatch)
	return isMisMatch || isDefault || isStandard || isMatch || isName
}

func (l LocalModule) Equal(other MatchSpecificity) bool {
	_, ok := other.(LocalModule)
	return ok
}
