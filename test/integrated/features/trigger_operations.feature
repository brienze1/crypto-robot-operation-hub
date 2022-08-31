Feature: Trigger operations
  In order to trigger operations
  There must be at least one client available
  I must receive a summary analysis via sqs

  Scenario: Trigger operation for one client
    Given dynamoDb is "up"
    And there is 1 client available in dynamodb
    And binance api is "up"
    And sns service is "up"
    When I receive message with summary equals "STRONG_BUY"
    Then there should be 1 messages sent via sns
    And process should exit with 0