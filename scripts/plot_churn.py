#!/usr/bin/env python3

import argparse
import os
import pandas as pd
import matplotlib.pyplot as plt

def load_meta(path):
    meta = {}
    with open(path, "r") as f:
        for line in f:
            if not line.startswith("#"):
                continue
            parts = line[1:].strip().split(",", 1)
            if len(parts) != 2:
                continue
            key, val = parts
            meta[key.strip()] = val.strip()
    return meta

def per_node_churn_plot(df, meta, outdir):
    algo = meta.get("algo", "unknown")
    churn_op = meta.get("churn_op", "unknown")
    nodes_before = meta.get("nodes_before", "?")
    nodes_after = meta.get("nodes_after", "?")
    zipf_s = meta.get("zipf_s", "?")

    x = df["node_id"]
    before = df["count_before"]
    after = df["count_after"]

    width = 0.4
    xs = range(len(x))

    plt.figure(figsize=(12,6))
    plt.bar([i - width/2 for i in xs], before, width=width, label="before")
    plt.bar([i + width/2 for i in xs], after, width=width, label="after")

    plt.xticks(xs, x, rotation=40, ha="right")
    plt.ylabel("Count")
    plt.title(
        f"Churn per-node load ({churn_op})\n"
        f"algo={algo}, nodes {nodes_before}→{nodes_after}, zipf_s={zipf_s}"
    )
    plt.legend()
    plt.tight_layout()

    outname = f"churn_per_node_{algo}_{churn_op}_nb{nodes_before}_na{nodes_after}_zipf{zipf_s}.png"
    outpath = os.path.join(outdir, outname)
    plt.savefig(outpath, dpi=200)
    plt.close()
    print(f"[OK] Wrote {outpath}")

def summary_moved_plot(metas, outdir):
    import numpy as np

    algos = []
    moved_ratios = []
    churn_ops = set()

    for meta in metas:
        algo = meta.get("algo", "unknown")
        churn_op = meta.get("churn_op", "unknown")
        mr = float(meta.get("moved_ratio", 0.0))

        algos.append(algo)
        moved_ratios.append(mr)
        churn_ops.add(churn_op)

    df = pd.DataFrame({"algo": algos, "moved_ratio": moved_ratios})
    df = df.sort_values(by="algo")

    plt.figure(figsize=(8,5))
    plt.plot(df["algo"], df["moved_ratio"], marker="o")
    title = "Fraction of keys moved vs algorithm"
    if len(churn_ops) == 1:
        title += f" (churn_op={list(churn_ops)[0]})"
    plt.title(title)
    plt.xlabel("Algorithm")
    plt.ylabel("moved_ratio")
    plt.grid(True)
    plt.tight_layout()

    outpath = os.path.join(outdir, "summary_moved_ratio_vs_algo.png")
    plt.savefig(outpath, dpi=200)
    plt.close()
    print(f"[OK] Wrote {outpath}")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--csv",
        nargs="+",
        required=True,
        help="Churn CSV files produced by cmd/sim (mode=churn)",
    )
    parser.add_argument(
        "--outdir",
        required=True,
        help="Output directory for churn plots",
    )
    args = parser.parse_args()

    os.makedirs(args.outdir, exist_ok=True)

    metas = []

    for path in args.csv:
        meta = load_meta(path)
        metas.append(meta)

        df = pd.read_csv(path, comment="#")
        if not {"node_id", "count_before", "count_after"}.issubset(df.columns):
            raise ValueError(f"{path} does not look like a churn CSV")

        per_node_churn_plot(df, meta, args.outdir)

    summary_moved_plot(metas, args.outdir)

if __name__ == "__main__":
    main()
