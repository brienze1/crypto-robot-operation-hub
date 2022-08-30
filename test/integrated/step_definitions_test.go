package integrated

import (
	"errors"
	"github.com/cucumber/godog"
)

func iEat(arg1 int) error {
	return godog.ErrPending
}

func thereAreGodogs(arg1 int) error {
	return errors.New("test")
}

func thereShouldBeRemaining(arg1 int) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I eat (\d+)$`, iEat)
	ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
