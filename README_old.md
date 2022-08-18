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
            "profit": -53.00
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
            "profit": -53.00
        }
    ]
  },
  "active_operations": [
    {
        "id": "asdasd-asdasd-asdasd-asdasd",
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
        "id": "asdasd-asdasd-asdasd-asdasd",
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


#### Client Lock Cache

Client Lock Cache is the database that contains the client lock information needed to avoid multiple trigger the
operations to the same client.
For this DB, Redis was chosen because of the easy implementation, really fast operation speed and liability.

##### Schema

Locks are being save into the key `CRYPTO_ROBOT_OPERATION_HUB_CLIENT_LOCK_<client_id>` with it´s value being
the `client_id` itself.

##### Operations

This application supports the following operations to the Client Lock Cache:

- Read ops:
    - Used to find locks active for clients.
- Write ops:
    - Used to persist locks for clients (Uses ttl)
