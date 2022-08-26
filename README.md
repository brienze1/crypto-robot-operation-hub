<div id="top"></div>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/brienze1/crypto-robot-operation-hub/blob/main/LICENSE)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/brienze1/crypto-robot-operation-hub)
![Build](https://img.shields.io/github/workflow/status/brienze1/crypto-robot-operation-hub/Build?label=Build)
[![Coverage Status](https://coveralls.io/repos/github/brienze1/crypto-robot-operation-hub/badge.svg?branch=main)](https://coveralls.io/github/brienze1/crypto-robot-operation-hub?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/brienze1/crypto-robot-operation-hub)](https://goreportcard.com/report/github.com/brienze1/crypto-robot-operation-hub)
[![Golang](https://img.shields.io/github/go-mod/go-version/brienze1/crypto-robot-operation-hub)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/brienze1/crypto-robot-operation-hub.svg)](https://pkg.go.dev/github.com/brienze1/crypto-robot-operation-hub)

# Crypto Data Operation Hub

1. [About the Project](#about-the-project)
    1. [Input](#input)
    2. [Output](#output)
    3. [Persistence](#persistence)
        1. [Client DB](#client-db)
            1. [Schema](#schema)
            2. [Operation](#operation)
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
3. [About Me](#about-me)

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
otherwise if operation type is SELL, amount should be used like cryptocurrency amount).

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
    "locked_until" "20-07-2022 02:18:10",
    "locked" false,
    "cash_amount": 100,
    "cash_reserved": 0.00,
    "crypto_amount": 0.0000312,
    "crypto_symbol": "BTC",
    "crypto_reserved": 0.0,
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
    ],
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

##### Operation

This application supports the following operations to the Client DB:

- Read ops:
    - Used to read clients available for the current operation

### Rules

Here are some rules that need to be implemented in this application.

Implemented:

- None

Not implemented:

[//]: # (- Client must have the coin symbol selected inside `config.symbols` variable to operate it)

- Client must be active
- Client must not be locked
- Current date must be greater than locked_until value
- Client must have enough cash to buy minimum allowed amount of crypto
- Client must have enough crypto to sell minimum allowed amount
- Buy operations should be triggered when the summary received is equal or less restricting than the `config.buy_on`
  value.
    - For example if the config value is equal to `BUY` and a `STRONG_BUY` analysis was received, the operation should
      be allowed, and the opposite should be denied.
- Sell operations should be triggered when the summary received is equal or less restricting than the `config.sell_on`
  value.
    - For example if the config value is equal to `SELL` and a `STRONG_SELL` analysis was received, the operation should
      be allowed, and the opposite should be denied.

- Operations should not be triggered after monthly sell cap has been reached
    - Operations should also check if the amount won't go over when sell operation is triggered, for example if monthly
      total amount is 20.000,00 and the cap is 25.000,00 the maximum operation value triggered should be of 2.500,00,
      because when the operation is completed and the crypto is sold the expectation is that the value should be equal
      or close to the bought amount (witch mas of 2500) totalizing 25.000,00 monthly sell value.
    - If monthly cap is 0 it can be ignored.
- Operations should not be triggered if `daily_summary.proffit` has a negative value of more than or equal to
  the `config.day_stop_loss` value.
    - `daily_summary.day` value should be checked to see if current day has changed, in this case, the values
      should be updated to start a new day.
- Operations should not be triggered if `monthly_summary.proffit` has a negative value of more than or equal to
  the `config.month_stop_loss` value.
    - `monthly_summary.month` value should be checked to see if current month has changed, in this case, the values
      should be updated to start a new month.

### Built With

This application is build with Golang, code is build using a Dockerfile every deployment into the main branch in GitHub
using GitHub actions. Local environment is created using localstack for testing purposes using
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

- [] Implement Behaviour tests (BDD)
- [] Implement Unit tests
- [] Implement application logic
- [] Create Dockerfile
- [] Create Docker compose for local infrastructure
- [] Document everything in Readme

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

### Prerequisites

- Install Golang

    - Windows/MacOS/Linux
        - [Manual](https://go.dev/dl/)
    - macOS
        - [Homebrew](https://docs.brew.sh/Installation)
          ```bash
          brew install go
          ```
    - Linux
        - Via terminal
          ```bash
          sudo add-apt-repository ppa:longsleep/golang-backports
          sudo apt update
          sudo apt install golang-go
          ```

- Install Docker
    - [Windows/macOS/Linux/WSL](https://www.docker.com/get-started/)

### Installation

- Run the following to install project dependencies:
    - Windows/MacOS/Linux/WSL
      ```bash
      go mod download
      ```

- Run the following to compile the project and generate executable:
    - Windows/MacOS/Linux/WSL
      ```bash
      go build -o bin/operation-hub cmd/operation-hub/main.go
      ```

Note: the binary generated will be available at `./bin` folder.

### Requirements

To run the application locally, first a local infrastructure needs to be deployed

#### Deploying Local Infrastructure

This requires [docker](#prerequisites) to be installed. Localstack will deploy aws local integration and create the
topic used by this application to send the events.

Obs: Make sure Docker is running before.

- Start the required infrastructure via localstack using docker compose command:

    - Windows/macOS/Linux/WSL
      ```bash
      docker-compose -f ./build/local/docker-compose.yml up
      ```

- To stop localstack:
    - Windows/macOS/Linux/WSL
      ```bash
      docker-compose -f ./build/local/docker-compose.yml down
      ```

### Usage

#### Manual Input

- Start the compiled application:
    - Windows/macOS/Linux/WSL
      ```bash
      go run cmd/operation-hub/main.go
      ```
- To stop the application just press Ctrl+C

#### Docker Input

- In case you want to use a Docker container to run the application first you need to build the Docker image from
  Dockerfile:
    - Windows/macOS/Linux/WSL
      ```bash
      docker build -t crypto-robot-operation-hub .
      ```

- And then run the new created image:
    - Windows/macOS/Linux/WSL
      ```sh
      docker run --network="host" -d -it crypto-robot-operation-hub bash \
      -c "OPERATION_HUB_ENV=localstack go run ./cmd/local/main_local.go"
      ```

### Testing

- To run the tests just type the command bellow in terminal:
    - Windows/macOS/Linux/WSL
      ```bash
      go test ./...
      ```

<p align="right">(<a href="#top">back to top</a>)</p>

## About me

Hello! :)

My name is Luis Brienze, and I'm a Software Engineer.

I focus primarily on software development, but I'm also good at system architecture, mentoring other developers,
etc... I've been in the IT industry for 4+ years, during this time I worked for companies like Itau, Dock, Imagine
Learning and
EPAM.

I graduated from UNESP studying Automation and Control Engineering in 2022, and I also took multiple courses on Udemy
and Alura.

My main stack is Java, but I'm also pretty good working with Kotlin and Typescript (both server side).
I have quite a good knowledge of AWS Cloud, and I'm also very conformable working with Docker.

During my career, while working with QA's, I've also gained a lot of valuable experience with testing applications in
general from unit/integrated
testing using TDD and BDD, to performance testing apps with JMeter for example.

If you want to talk to me, please fell free to reach me anytime at [LinkedIn](https://www.linkedin.com/in/luisbrienze/)
or [e-mail](mailto:lfbrienze@gmail.com?subject=GitHUB Crypto robot operation hub).

<p align="right">(<a href="#top">back to top</a>)</p>