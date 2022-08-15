# crypto-robot-operation-hub

Verify and triggers active clients operations individually

build :

```bash
  go build -o bin/operation-hub cmd/operation-hub/main.go
```

clean :

```bash
  rm -r bin
```

run :

```bash
  go run cmd/operation-hub/main.go
```

WIP rules:

- Client must be active
- Client must not be locked
- Client must have enough cash to buy minimum allowed amount of crypto
- Client must have the coin symbol selected inside `config.symbols` variable
- Operations should not be triggered after monthly sell cap has been reached
    - Operations should also check if the amount won't go over when sell operation is triggered, for example if monthly
      total amount is 20.000,00 and the cap is 25.000,00 the maximum operation value triggered should be of 2.500,00,
      because when the operation is completed and the crypto is sold the expectation is that the value should be equal
      or close to the bought amount (witch mas of 2500) totalizing 25.000,00 monthly sell value.
- Buy operations should be triggered when the summary received is equal or less restricting than the `config.buy_on`
  value.
    - For example if the config value is equal to `BUY` and a `STRONG_BUY` analysis was received, the operation should
      be allowed, and the opposite should be denied.
- Sell operations should be triggered when the summary received is equal or less restricting than the `config.sell_on`
  value.
    - For example if the config value is equal to `SELL` and a `STRONG_SELL` analysis was received, the operation should
      be allowed, and the opposite should be denied.
- Operations should not be triggered if `daily_summary.proffit` has a negative value of more than or equal to
  the `config.day_stop_loss` value.
    - `daily_summary.day` value should be checked to see if current day has changed, in this case, the values
      should be updated to start a new day.
- Operations should not be triggered if `monthly_summary.proffit` has a negative value of more than or equal to
  the `config.month_stop_loss` value.
    - `monthly_summary.month` value should be checked to see if current month has changed, in this case, the values
      should be updated to start a new month.

[//]: # (TODO check for more rules)


WIP client db example:

```
{
  "id": "asdasdASD",
  "active" true,
  "locked" false,
  "cash": 100.00,
  "crypto": [
    {
        "symbol": "BTC",
        "amount": 0.0000312
    }
  ],
  "config": {
    "buy_on": "STRONG_BUY",
    "sell_on": "SELL",
    "operation_stop_loss": 50.00,
    "day_stop_loss": 500.00,
    "month_stop_loss": 500.00,
    "month_sell_cap": 25000.00,
    "symbols": [
        "BTC",
        "SOL"
    ]
  },
  "monthly_summary": {
    "month": "08/2022",
    "amount_sold": 23000.42,
    "amount_bought": 37123.42,
    "profit": 1032.32,
    "crypto": [
        {
            "symbol": "BTC",
            "average_buy_value": 230020.42,
            "average_sell_value": 235020.42,
            "amount_sold": 0.00231,
            "amount_bought": 0.00431,
        }
    ]
  },
  "daily_summary": {
    "day": "14/08/2022"
    "amount_sold": 23000.42,
    "amount_bought": 37123.42,
    "profit": -53.00,
     "crypto": [
        {
            "symbol": "BTC",
            "average_buy_value": 230020.42,
            "average_sell_value": 235020.42,
            "amount_sold": 0.00231,
            "amount_bought": 0.00431,
        }
    ]
  },
  "active_operations": [
    {
        "status": "AWAITING_SELL",
        "symbol": "BTC",
        "start_time": "2007-12-03 10:15:30",
        "buy_time": "2007-12-03 10:15:30",
        "crypto_amount_bought": 0.001023,
        "buy_quote": 234223.23,
        "stop_loss_quote": 233223.23,
        "available_amount": 1233.32,
        "operation_amount": 1233.32
    },
    {
        "status": "AWAITING_BUY",
        "symbol": "BTC",
        "start_time": "2007-12-03 10:15:30",
        "buy_time": "",
        "crypto_amount_bought": 0.0,
        "buy_quote": 0.0,
        "stop_loss_quote": 0.0,
        "available_amount": 1233.32,
        "operation_amount": 0.0
    }
  ]
}
```
