# cs2-p2p

CS2 P2P market monitoring system.

## Services

```text
services/
  steam-collector     collects Steam Market snapshots
  csmoney-collector   collects CS.Money Market snapshots
  lisskins-collector  collects Lis-Skins snapshots
  market-proc         processes market snapshots and calculates opportunities
  p2p-controller      API/control plane for collectors and UI
  p2p-ui              web UI
```

## Infrastructure

```text
deployments/docker-compose
```

Local infrastructure currently contains Redpanda as a Kafka-compatible broker and Redpanda Console.

## Architecture

```text
External Markets -> Collectors -> Kafka -> market-proc -> ClickHouse

p2p-ui -> p2p-controller
p2p-controller <-> Collectors
p2p-controller <-> Postgres
p2p-controller -> ClickHouse
```

See [docs/architecture/components.md](docs/architecture/components.md).
