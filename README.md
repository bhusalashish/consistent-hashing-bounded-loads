# consistent-hashing-bounded-loads

### A Comparative Study of Jump, Maglev, and Bounded-Load Consistent Hashing (CH-BL)

This project implements three modern sharding / load-distribution algorithms:

* **Jump Consistent Hashing**
* **Maglev Load Balancing (Google, NSDI’16)**
* **Consistent Hashing with Bounded Loads (CH-BL)** with:

  * vnode ring
  * per-node capacity enforcement
  * two-choice fallback

It includes a flexible **Go-based simulator** for generating uniform and Zipf-skewed workloads and a **Python plotting pipeline** to analyze:

* Per-node distribution
* Coefficient of Variation (CV)
* Max/Avg imbalance
* Key movement under churn (future extension)

Suitable for Distributed Systems coursework (CMPE 273), infra engineers, and anyone studying scalable routing or load balancing.

---

## 📦 Project Structure

```
consistent-hashing-bounded-loads/
├── cmd/
│   └── sim/              # Simulator CLI (Go)
├── internal/
│   └── ring/             # Vnode consistent hash ring for CH-BL
├── pkg/
│   ├── hash/             # xxhash64 hashing utilities
│   ├── metrics/          # CV, StdDev, Max/Avg helpers
│   ├── router/           # Algorithm routers (jump, maglev, chbl)
│   └── routercore/       # Shared interfaces + router options
├── scripts/
│   └── plot_results.py   # Python plotting script
├── results/              # CSV outputs from simulator
├── plots/                # Generated charts
└── README.md
```

---

## 🚀 Features

### Jump Consistent Hashing

* O(1) lookup
* Minimal remapping when nodes change
* Used in Google Bigtable, Cloud Pub/Sub

### Maglev Load Balancing

* Permutation-based lookup table
* Near-perfect distribution
* O(1) lookup
* Low churn
* Used in Google frontend load balancers (NSDI’16)

### Consistent Hashing with Bounded Loads (CH-BL)

* Guarantees **max load ≤ c × average load**
* Uses consistent hashing + forwarding
* Includes:

  * Vnode ring
  * Per-node capacity
  * ExpectedKeys-based bound
  * Two-choice fallback to reduce walk lengths

### Simulator

* Uniform & Zipf workloads
* Configurable params
* Outputs CSV for analysis

### Plotting

* Bar charts for per-node distribution
* CV vs algorithm
* Max/Avg vs algorithm
* High-resolution PNG output

---

## 🛠 Quickstart

### Install Go & Python dependencies

```bash
go mod tidy
pip install matplotlib pandas
```

---

## 🎛 Run the Simulator

### Jump (uniform workload)

```bash
go run ./cmd/sim \
  -algo jump -nodes 16 -keys 100000 \
  -zipf-s 0.0 \
  -out results/jump_uniform.csv
```

### Maglev

```bash
go run ./cmd/sim \
  -algo maglev -nodes 16 -keys 100000 \
  -table-size 65537 \
  -out results/maglev_uniform.csv
```

### CH-BL

```bash
go run ./cmd/sim \
  -algo chbl -nodes 16 -keys 100000 \
  -load-factor 1.25 \
  -vnodes 100 \
  -walk-threshold 8 \
  -out results/chbl_uniform.csv
```

---

## 📊 Generate Plots

```bash
python3 scripts/plot_results.py \
  --csv results/jump_uniform.csv \
       results/maglev_uniform.csv \
       results/chbl_uniform.csv \
  --outdir plots
```

This produces:

* `per_node_jump_nodes16_zipf0.png`
* `per_node_maglev_nodes16_zipf0.png`
* `per_node_chbl_nodes16_zipf0.png`
* `summary_cv_vs_algo.png`
* `summary_maxoveravg_vs_algo.png`

---

## 📈 Example Interpretation

* **Jump**: good uniformity but vulnerable under heavy Zipf skew.
* **Maglev**: extremely uniform distribution; very low churn.
* **CH-BL**: guarantees strict per-node load cap (`c × avg`) even under skew; ideal for cache & storage backends.
* **Two-choice fallback**: reduces CH-BL walk lengths at high load.

---

## 🔧 Parameters Reference

| Algorithm | Parameter       | Description                         |
| --------- | --------------- | ----------------------------------- |
| Jump      | `HashSeed`      | Hash seed for determinism           |
| Maglev    | `TableSize`     | Size of permutation table           |
| CH-BL     | `LoadFactor`    | `c` factor for calculating capacity |
| CH-BL     | `Vnodes`        | Virtual nodes per physical node     |
| CH-BL     | `WalkThreshold` | Steps before two-choice fallback    |
| CH-BL     | `ExpectedKeys`  | Used to compute capacity            |

---

## 🧪 Testing

```bash
go test ./...
```

Includes:

* Hash determinism tests
* Jump determinism + minimal remap tests
* Maglev determinism tests
* Basic CH-BL correctness tests
* Metrics tests

---

## 📚 References

* Lamping & Veach — *Jump Consistent Hashing*
* Eisenbud et al. — *Maglev: A Fast and Reliable Software Network Load Balancer* (NSDI’16)
* Mirrokni, Thorup, Zadimoghaddam — *Consistent Hashing with Bounded Loads*

---

## 👤 Author

Ashish Bhusal
San José State University
CMPE 273 — Distributed Systems (Fall 2025)

---

## 📝 License

MIT License
