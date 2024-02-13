CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');
CREATE TABLE IF NOT EXISTS test_table
(
    serial                   serial PRIMARY KEY,
    varchar                  VARCHAR(100)             NOT NULL,
    integer                  INTEGER                  NOT NULL,
    numeric                  NUMERIC(10, 2)           NOT NULL,
    boolean                  BOOLEAN                  NOT NULL,
    decimal                  DECIMAL(10, 2)           NOT NULL,
    real                     REAL                     NOT NULL,
    double                   DOUBLE PRECISION         NOT NULL,
    text                     TEXT                     NOT NULL,
    bytea                    BYTEA                    NOT NULL,
    timestamp                TIMESTAMP                NOT NULL,
    timestamp_with_time_zone TIMESTAMP WITH TIME ZONE NOT NULL,
    date                     DATE                     NOT NULL,
    time                     TIME                     NOT NULL,
    time_with_time_zone      TIME WITH TIME ZONE      NOT NULL,
    interval                 INTERVAL                 NOT NULL,
    enum                     mood                     NOT NULL,
    cidr                     CIDR                     NOT NULL,
    inet                     INET                     NOT NULL,
    macaddr                  MACADDR                  NOT NULL,
    macaddr8                 MACADDR8                 NOT NULL,
    bit                      BIT(10)                  NOT NULL,
    bit_varying              BIT VARYING(10)          NOT NULL,
    uuid                     UUID                     NOT NULL,
    json                     JSON                     NOT NULL,
    jsonb                    JSONB                    NOT NULL,
    oid                      OID                      NOT NULL,
    character                CHARACTER(10)            NOT NULL,
    array_field              INTEGER[]                NOT NULL
);
CREATE TABLE IF NOT EXISTS test_table2
(
    serial      serial PRIMARY KEY,
    foreign_key serial REFERENCES test_table (serial)
)