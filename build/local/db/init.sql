CREATE DATABASE crypto_robot;

CREATE USER operation_hub WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE crypto_robot TO operation_hub;

\c crypto_robot

CREATE TABLE crypto
(
    id     INT PRIMARY KEY,
    symbol character varying(15) NOT NULL
);

CREATE TABLE clients
(
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    active              BOOLEAN          default false             NOT NULL,
    locked              BOOLEAN          default false             NOT NULL,
    locked_until        TIMESTAMP        default CURRENT_TIMESTAMP NOT NULL,
    cash_amount         DOUBLE PRECISION default 0                 NOT NULL,
    cash_reserved       DOUBLE PRECISION default 0                 NOT NULL,
    crypto_amount       DOUBLE PRECISION default 0                 NOT NULL,
    crypto_reserved     DOUBLE PRECISION default 0                 NOT NULL,
    buy_on              INT              default 0                 NOT NULL,
    sell_on             INT              default 0                 NOT NULL,
    ops_timeout_seconds INT              default 0                 NOT NULL,
    operation_stop_loss DOUBLE PRECISION                           NOT NULL,
    day_stop_loss       DOUBLE PRECISION                           NOT NULL,
    month_stop_loss     DOUBLE PRECISION                           NOT NULL
);

CREATE TABLE client_symbols
(
    crypto_id INT  NOT NULL,
    client_id UUID NOT NULL,
    PRIMARY KEY (crypto_id, client_id),
    CONSTRAINT fk_client foreign key (client_id) references clients (id) on delete CASCADE,
    CONSTRAINT fk_symbol foreign key (crypto_id) references crypto (id) on delete CASCADE
);

CREATE TABLE clients_summary
(
    id            SERIAL PRIMARY KEY,
    client_id     UUID                       NOT NULL,
    type          character varying(15)      NOT NULL,
    day           INT,
    month         INT,
    year          INT,
    amount_sold   DOUBLE PRECISION default 0 NOT NULL,
    amount_bought DOUBLE PRECISION default 0 NOT NULL,
    profit        DOUBLE PRECISION default 0 NOT NULL,
    CONSTRAINT fk_client foreign key (client_id) references clients (id) on delete CASCADE
);

CREATE TABLE clients_crypto_summary
(
    id                 SERIAL PRIMARY KEY,
    client_id          UUID                       NOT NULL,
    crypto_id          INT                        NOT NULL,
    summary_id         INT                        NOT NULL,
    average_buy_value  DOUBLE PRECISION default 0 NOT NULL,
    average_sell_value DOUBLE PRECISION default 0 NOT NULL,
    amount_bought      DOUBLE PRECISION default 0 NOT NULL,
    profit             DOUBLE PRECISION default 0 NOT NULL,
    CONSTRAINT fk_client foreign key (client_id) references clients (id) on delete CASCADE,
    CONSTRAINT fk_symbol foreign key (crypto_id) references crypto (id) on delete CASCADE,
    CONSTRAINT fk_summary foreign key (summary_id) references clients_summary (id) on delete CASCADE
);


-- functions
CREATE OR REPLACE FUNCTION create_day_summary()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
BEGIN
    IF NOT EXISTS(SELECT 1
                  FROM clients_summary
                  WHERE NEW.id = client_id
                    AND type = 'DAY'
                    AND day = date_part('day', (SELECT current_timestamp))
                    AND month = date_part('month', (SELECT current_timestamp))
                    AND year = date_part('year', (SELECT current_timestamp))) THEN
        INSERT INTO clients_summary(client_id, type, day, month, year)
        VALUES (NEW.id, 'DAY', date_part('day', (SELECT current_timestamp)),
                date_part('month', (SELECT current_timestamp)),
                date_part('year', (SELECT current_timestamp)));
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION create_month_summary()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
BEGIN
    IF NOT EXISTS(SELECT 1
                  FROM clients_summary
                  WHERE NEW.id = client_id
                    AND type = 'MONTH'
                    AND month = date_part('month', (SELECT current_timestamp))
                    AND year = date_part('year', (SELECT current_timestamp))) THEN
        INSERT INTO clients_summary(client_id, type, month, year)
        VALUES (NEW.id, 'MONTH', date_part('month', (SELECT current_timestamp)),
                date_part('year', (SELECT current_timestamp)));
    END IF;
    RETURN NEW;
END;
$$;

-- triggers
CREATE OR REPLACE TRIGGER trg_create_day_summary_insert_client
    AFTER INSERT
    ON clients
    FOR EACH ROW
EXECUTE FUNCTION create_day_summary();

CREATE OR REPLACE TRIGGER trg_create_day_summary_update_client
    AFTER UPDATE
    ON clients
    FOR EACH ROW
EXECUTE FUNCTION create_day_summary();

CREATE OR REPLACE TRIGGER trg_create_month_summary_insert_client
    AFTER INSERT
    ON clients
    FOR EACH ROW
EXECUTE FUNCTION create_month_summary();

CREATE OR REPLACE TRIGGER trg_create_month_summary_update_client
    AFTER UPDATE
    ON clients
    FOR EACH ROW
EXECUTE FUNCTION create_month_summary();

-- Inserts test data
INSERT INTO crypto (id, symbol)
VALUES (1, 'BTC'),
       (2, 'SOL');

INSERT INTO clients (id, active, locked, cash_amount, crypto_amount, buy_on, sell_on, operation_stop_loss,
                                  day_stop_loss, month_stop_loss)
VALUES (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000),
       (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000),
       (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000),
       (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000),
       (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000),
       (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000),
       (DEFAULT, true, false, 10000, 0.00001, 2, 2, 300, 1000, 10000);

INSERT INTO client_symbols (client_id, crypto_id)
SELECT id, 1
FROM clients;

INSERT INTO client_symbols (client_id, crypto_id)
SELECT id, 2
FROM clients;


