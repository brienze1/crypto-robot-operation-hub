package integrated

import "github.com/cucumber/godog"

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^dynamoDb is "([^"]*)"$`, dynamoDbIs)
	ctx.Step(`^there is (\d+) client available in dynamodb$`, thereIsClientAvailableInDynamodb)
	ctx.Step(`^binance api is "([^"]*)"$`, binanceApiIs)
	ctx.Step(`^I receive message with summary equals "([^"]*)"$`, iReceiveMessageWithSummaryEquals)
	ctx.Step(`^there should be (\d+) message sent via sns$`, thereShouldBeMessageSentViaSns)
	ctx.Step(`^process should exit with (\d+)$`, processShouldExitWith)
}

func dynamoDbIs(status string) error {
	return godog.ErrPending
}

func thereIsClientAvailableInDynamodb(numberOfClients int) error {
	return godog.ErrPending
}

func binanceApiIs(status string) error {
	return godog.ErrPending
}

func iReceiveMessageWithSummaryEquals(summary string) error {
	return godog.ErrPending
}

func thereShouldBeMessageSentViaSns(arg1 int) error {
	return godog.ErrPending
}

func processShouldExitWith(status int) error {
	return godog.ErrPending
}
