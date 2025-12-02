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

---

# 🎨 Interactive Visualization Tool (New Feature)

This project now includes an **interactive web-based visualization tool** for demonstrating consistent hashing algorithms in real time.

Using a Go backend + React (TypeScript) + D3.js frontend, the visualizer provides:

* A **live animated consistent-hash ring**
* Visualization of **nodes** and **keys** on the ring
* **Real-time rebalancing** when nodes are added or removed
* Smooth **key-movement animations**
* Support for **all four algorithms**:

  * Plain Consistent Hash (RingCH)
  * Jump Consistent Hashing
  * Maglev
  * CH-BL (Consistent Hashing with Bounded Loads)
* Hover interactions (highlight keys or nodes)
* Fully interactive controls:

  * Add Node
  * Remove Node
  * Regenerate Keys
  * Select Algorithm
  * Adjust key count

All of this runs locally and integrates directly with the Go sharding implementations.

---

## 🚀 Visualizer Setup

### 1. Run the Go backend

From project root:

```bash
go run ./cmd/visualizer
```

This starts an HTTP server at:

```
http://localhost:8080
```

### 2. Run the frontend

```bash
cd web
npm install
npm run dev
```

Open the UI at:

```
http://localhost:5173
```

---

## 🖥 What You Can Do in the Visualizer

### ✨ Add or Remove Nodes

* Watch the ring update instantly
* Keys animate smoothly to new owners

### 🔁 Switch Algorithms Dynamically

Compare behaviors visually:

* `ring`
* `jump`
* `maglev`
* `chbl`

### 🎞 Real-time Key Reassignment

* See exactly **which keys** moved and **where** they went
* Understand load balancing behavior intuitively

### 🎛 Control Panel

* Algorithm dropdown
* Number of keys
* Buttons for node operations
* Regenerate keys instantly

---

## 📊 Why This Matters

This visualization makes it easy to see:

### **How Jump achieves minimal churn**

* Only ~1/(N+1) keys move when adding a node.

### **How Maglev ensures uniform distribution**

* Nodes fill evenly around the ring.

### **How CH-BL enforces load bounds**

* Nodes never exceed `c × average load`, even with skew.

### **How plain CH behaves as baseline**

* Provides a simple reference point for evaluating improvements.

---

## 🎥 Perfect for Live Demos

In your 20-minute presentation, you can:

* Show the algorithms on static plots (from Python scripts)
* Then switch to the visualizer and:

  * Add a node live → keys animate
  * Remove a node → keys reassign
  * Switch algos → completely different behaviors

Your professor will understand the differences instantly.

---

## 🧩 Architecture of the Visualizer

```
Go Backend (cmd/visualizer)
│
├── Exposes JSON APIs:
│     GET /state
│     POST /add-node
│     POST /remove-node
│     POST /regenerate-keys
│
└── Uses algorithms from:
      pkg/router/ringch
      pkg/router/jump
      pkg/router/maglev
      pkg/router/chbl

React + TypeScript + D3.js Frontend (web/)
│
├── Components:
│     Ring.tsx (SVG ring + nodes)
│     Key.tsx (SVG keys)
│     ControlPanel.tsx (UI controls)
│     App.tsx (main shell)
│
└── Animations:
      D3 transitions, easing functions
```

---

## 📚 How This Integrates with the Simulator

The visualizer is complementary to your CLI simulator:

* Simulator → produces CSVs for academic plots
* Visualizer → gives live intuition for real-time rebalancing

Both use the exact same Go algorithm implementations.

---

## 📦 Recommended Workflow

1. **Run simulations** (dist + churn)
2. **Generate plots** (Python)
3. **Run the visualizer**
4. **Show interactive demo in presentation**

This gives you:

* Hard data → CV, Max/Avg, churn ratios
* Visual intuition → live ring animations

---

## 🏁 Conclusion

This visualizer elevates the project from a basic simulation to a **full interactive distributed-systems demo**.
It clearly showcases how different consistent hashing strategies behave:

* Balance
* Churn
* Key movement
* Skew handling
* Capacity bounds

Perfect for impressing your professor, your classmates, and even future interviewe
