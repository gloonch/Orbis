# 🌌 Orbis
**Orbis** is a real-time analytics platform at the intersection of astrology and finance decoding celestial movements to generate meaningful market signals.
It platform at the intersection of astrology and finance decoding celestial movements to generate meaningful market signals.

> *"Connecting the pulse of the cosmos to the heartbeat of the market."*

---

## ✨ Overview

Orbis integrates three live data streams:

1. **Astronomical Data** (planetary positions & aspects)
2. **Environmental Data** (extreme heat, humidity, wind, precipitation)
3. **Market Data** (volumes, prices, volatility)

…to discover and act on statistical **cause-effect pairs** like:

- “Extreme heat → ↓ stock-market volume”
- “Cold snap → ↑ energy prices”
- “Full moon square Saturn → bond yields shift”

---

## ⚙️ Key Features

-  **Astro Stream** — Real-time ephemeris via NASA/JPL APIs (Sun, Moon, planets & major aspects)
-  **Weather Stream** — Live feeds (temperature, humidity, wind, AO/HDD indices) from Tomorrow.io/OpenWeatherMap
-  **Market Stream** — Live prices, volumes, returns & volatility via IEX Cloud/Alpha Vantage
-  **Cause-Effect Rules** — User-defined “If X then Y” rule engine (e.g. heat index > 35 °C AND full moon square Mars → limit buy)
-  **Time-Series Storage** — TimescaleDB hypertables for alignment & continuous aggregates
-  **Visualization** — Grafana dashboards (multi-series charts of astro, weather & market sync’d on a common timeline)
-  **Alerts & Notifications** — gRPC/REST webhooks, Slack/Email integration when rules fire
-  **Extensible** — Modular Go services, easily add new streams, new causal pairs or machine-learning models

---

## 🧩 Tech Stack

- **Language:** Go (Golang)
- **Stream Processing:** Apache Kafka
- **Storage:** PostgreSQL (TimescaleDB) + Redis
- **Deployment:** Docker / Kubernetes
- **Astro Engine:** NASA/JPL-based ephemeris via Go libraries

---

## 🔮 Inspired By
- Ecological Economics article: [Extreme Heat & Stock Market Activity (S. Direct 2020)](https://www.sciencedirect.com/science/article/abs/pii/S092180092030015X)
- Financial Astrology 



