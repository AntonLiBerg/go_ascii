package world

type MenuChoice struct {
	Text   string
	Update func(*World)
}

type MenuChoices struct {
	ShouldShow bool
	IsOpen     bool
	Header     string
	Choices    []MenuChoice
}

func cloneMenuChoices(menuChoices MenuChoices) MenuChoices {
	clone := MenuChoices{
		ShouldShow: menuChoices.ShouldShow,
		IsOpen:     menuChoices.IsOpen,
		Header:     menuChoices.Header,
		Choices:    append([]MenuChoice(nil), menuChoices.Choices...),
	}
	if clone.Choices == nil {
		clone.Choices = []MenuChoice{}
	}
	return clone
}
