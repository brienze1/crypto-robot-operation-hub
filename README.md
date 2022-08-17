<div id="top"></div>

# Crypto Data Operation Hub

1. [About the Project](#about-the-project)
    1. [Input](#input)
    2. [Output](#output)
    3. [Persistence](#persistence)
        1. [Client DB](#client-db)
            1. [Schema](#schema)
            2. [Operation](#operation)
        2. [Client Lock Cache](#client-lock-cache)
            1. [Schema](#schema-1)
            2. [Operation](#operation-1)
    4. [Rules](#rules)
    5. [Built With](#built-with)
        1. [Dependencies](#dependencies)
        2. [Compiler Dependencies](#compiler-dependencies)
        3. [Test Dependencies](#test-dependencies)
    6. [Roadmap](#roadmap)
2. [Getting Started](#getting-started)
    1. [Prerequisites](#prerequisites)
    2. [Installation](#installation)
    3. [Requirements](#requirements)
        1. [Deploying Local Infrastructure](#deploying-local-infrastructure)
    4. [Usage](#usage)
        1. [Manual Input](#manual-input)
        2. [Docker Input](#docker-input)
    5. [Testing](#testing)

## About the Project

The objective of this project is to receive a BUY or SELL event and trigger operations for active users.

### Input

The input should be received as an SNS message sent through an SQS subscription. This message will trigger the lambda
handler to perform the service. The indicators received should be used to predict marked behaviour and trigger BUY or
SELL operations.

Data could also be analysed for specific interval indicators if needed in the future (but the current behaviour does not
use that data).

Analysis summary indicators:

- STRONG_BUY
- BUY
- NEUTRAL
- SELL
- STRONG_SELL

Example of how the data received should look like:

```
{
    "summary": "BUY",
    "timestamp": "20-07-2022 02:18:10",
    "analysed_data": [
        {
          "interval": "0NE_MINUTE",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "BUY",
          "analysis": [
            {
                "metric": "SIMPLE_MOVING_AVERAGE",
                "indicator": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        },
        {
          "interval": "FIVE_MINUTES",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "STRONG_BUY",
          "analysis": [
            {
                "indicator": "SIMPLE_MOVING_AVERAGE",
                "summary": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        },
        {
          "interval": "FIFTEEN_MINUTES",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "STRONG_BUY",
          "analysis": [
            {
                "indicator": "SIMPLE_MOVING_AVERAGE",
                "summary": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        },
        {
          "interval": "THIRTY_MINUTES",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "STRONG_BUY",
          "analysis": [
            {
                "indicator": "SIMPLE_MOVING_AVERAGE",
                "summary": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        },
        {
          "interval": "ONE_HOUR",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "STRONG_BUY",
          "analysis": [
            {
                "indicator": "SIMPLE_MOVING_AVERAGE",
                "summary": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        },
        {
          "interval": "SIX_HOURS",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "STRONG_BUY",
          "analysis": [
            {
                "indicator": "SIMPLE_MOVING_AVERAGE",
                "summary": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        },
        {
          "interval": "ONE_DAY",
          "timestamp": "20-07-2022 02:18:10",
          "summary": "STRONG_BUY",
          "analysis": [
            {
                "indicator": "SIMPLE_MOVING_AVERAGE",
                "summary": "BUY",
                "score": {
                    "buy": 4,
                    "sell": 2
                }
            },
            {
                "indicator": "EXPONENTIAL_MOVING_AVERAGE",
                "summary": "NEUTRAL",
                "score": {
                    "buy": 3,
                    "sell": 3
                }
            }
          ]
        }
    ]
}
```

### Output

Since this is an async application there is no output to be returned, but operation events are generated from the data
received.
Operation events should have information needed to track client who will trigger the operation, the operation type,
crypto being traded, operation event time sent and available amount (if operation type is BUY amount should be in cash,
otherwise if operation type is SELL, amount should be used like crypto currency amount).

Example of how the line should look like:

```
{
    "client_id": "asdasd-asdasd-asdasd-asdasd",
    "operation": "BUY",
    "symbol": "BTC",
    "start_time": "2007-12-03 10:15:30",
    "available_amount": 1233.32
}
```

### Persistence

#### Client DB

Client DB is the database that contains the client information and configuration needed to trigger the operations.
For this DB, DynamoDB was chosen because of the easy implementation, schema changes are also easy to implement, speed of
operations, etc...

##### Schema

```
{
  "id": "asdasdASD",
  "active" true,
  "locked" false,
  "cash": {
    "amount": 100.00,
    "reserved": 0.00,
  },
  "crypto": [
    {
        "symbol": "BTC",
        "amount": 0.0000312
        "reserved": 0.0
    }
  ],
  "config": {
    "buy_on": "STRONG_BUY",
    "sell_on": "SELL",
    "ops_timeout_seconds": 60,
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
  }
}
```

##### Operations

This application supports the following operations to the Client DB:

- Read ops:
    - Used to read clients available for the current operation
- Write ops:
    - Used to change the lock variable (locking the client for current operation).
    - Used to change the reserved cash or crypto amount.

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

### Rules

Here are some rules that need to be implemented in this application.

Implemented:

- Data needs to be updated on the database once it is received
- The app should gather all the data saved and generate an indication
- The data should be sent via SNS topic event

Not implemented:

- Data received should be checked to see if it's newer than the one saved
- If there is no data saved on the database, the summary should be generated using only the data received

### Built With

This application is build with Node.js Typescript, code is build using a Dockerfile every deployment into the main
branch in GitHub using GitHub actions.
Local environment is created using localstack for testing purposes using
[crypto-robot-localstack](https://github.com/brienze1/crypto-robot-localstack).

#### Dependencies

- [aws-lambda](https://www.npmjs.com/package/aws-lambda): Used in Lambda Handler integration
- [aws-sdk](https://www.npmjs.com/package/aws-sdk): Used in SNS integration (Needs to be replaced for SNS specific
  dependency)
- [dynamoose](https://www.npmjs.com/package/dynamoose): Used as ORM for DynamoDB
- [winston](https://www.npmjs.com/package/winston): Used for logging purposes
- [uuid](https://www.npmjs.com/package/uuid): Used to generate uuids

#### Compiler Dependencies

- [typescript](https://www.npmjs.com/package/typescript): Used run/compile typescript code
- [eslint](https://www.npmjs.com/package/eslint): Used to enforce coding practices
- [babel](https://babeljs.io/): Used to transpile code into js on build
- [dotenv](https://www.npmjs.com/package/dotenv): Used to map .env variables

#### Test Dependencies

- [jest](https://www.npmjs.com/package/jest): Used to run unit tests
- [@cucumber/cucumber](https://www.npmjs.com/package/@cucumber/cucumber): Used to run integration tests
- [chai](https://www.npmjs.com/package/chai): Used to perform test assertions with cucumber
- [sinon](https://www.npmjs.com/package/sinon): Used to create mocks/stubs/spy's
- [aws-sdk-mock](https://www.npmjs.com/package/aws-sdk-mock): Used to create mocks for AWS integrations

### Roadmap

-   [X] Implement Behaviour tests (BDD)
-   [X] Implement Unit tests
-   [X] Implement application logic
-   [X] Create Dockerfile
-   [X] Create Docker compose for local infrastructure
-   [X] Document everything in Readme

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

### Prerequisites

- Install Node and npm

    - Windows/MacOS/Linux
        - [Manual](https://nodejs.org/)
    - macOS
        - [Homebrew](https://docs.brew.sh/Installation)
          ```bash
          brew install node
          ```
    - Linux
        - Via terminal
          ```bash
          sudo apt install nodejs
          sudo apt install npm
          ```

- Install Docker
    - [Windows/macOS/Linux/WSL](https://www.docker.com/get-started/)

### Installation

- Run the following to install dependencies and compile the project:
    - Windows/MacOS/Linux/WSL
      ```bash
      npm install && npm run build
      ```

### Requirements

To run the application locally, first a local infrastructure needs to be deployed

#### Deploying Local Infrastructure

This requires [docker](#prerequisites) to be installed. Localstack will deploy aws local integration and create the
topic used by this application to send the events.

Obs: Make sure Docker is running before.

- Start the required infrastructure via localstack using docker compose command:

    - Windows/macOS/Linux/WSL
      ```bash
      docker-compose up
      ```

- To stop localstack:
    - Windows/macOS/Linux/WSL
      ```bash
      docker-compose down
      ```

### Usage

#### Manual Input

- Start the compiled application:
    - Windows/macOS/Linux/WSL
      ```bash
      npm run dev
      ```
- To stop the application just press Ctrl+C

#### Docker Input

- In case you want to use a Docker container to run the application first you need to build the Docker image from
  Dockerfile:
    - Windows/macOS/Linux/WSL
      ```bash
      docker build -t crypto-robot-data-digest .
      ```

- And then run the new created image:
    - Windows/macOS/Linux/WSL
      ```bash
      docker run --network="host" -d -it crypto-robot-data-digest bash -c "npm install && npm run dev:docker"
      ```

### Testing

- To run the tests just type the command bellow in terminal:
    - Windows/macOS/Linux/WSL
      ```bash
      npm run test
      ```

<p align="right">(<a href="#top">back to top</a>)</p>
