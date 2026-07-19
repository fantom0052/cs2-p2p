# Components

## Services

### steam-collector

Collects item price snapshots from Steam Market and writes raw snapshots to Kafka topic `raw.market.prices.steam`.

Does not own a database.

### csmoney-collector

Collects item price snapshots from CS.Money Market and writes raw snapshots to Kafka topic `raw.market.prices.csmoney`.

Does not own a database.

### lisskins-collector

Collects item price snapshots from Lis-Skins and writes raw snapshots to Kafka topic `raw.market.prices.lisskins`.

Does not own a database.

### market-proc

Reads raw market snapshots from Kafka, normalizes them, compares prices between markets, calculates spread/profit/profit percent, writes results to ClickHouse, and can publish opportunities to Kafka topic `market.opportunities`.

Uses ClickHouse for market snapshots and opportunities.

### p2p-controller

Provides the control-plane API. Stores collector configuration and operational state in Postgres, talks directly to collectors over HTTP/gRPC, and reads ClickHouse to serve market data and opportunities to the UI.

Uses Postgres for configuration/state and ClickHouse for read-only market data queries.

### p2p-ui

Web UI for configuring collectors and viewing prices/opportunities through `p2p-controller`.

Does not talk directly to databases.

## Infrastructure

### Kafka / Redpanda

Event broker for market data.

### ClickHouse

Analytical storage for market snapshots and calculated opportunities.

### Postgres

Operational storage for `p2p-controller`: collector configs, strategies, filters, and runtime state.
