#!/usr/bin/env python3

import argparse
import pandas as pd
import matplotlib.pyplot as plt
import os
import re

"""
This script reads simulator CSV output and produces:
1. Per-node bar plot
2. Summary comparison plots (CV vs algo, Max/Avg vs algo)

Assumes CSV output of the form:

node_id,count
node-0,6212
node-1,6401
...
#algo,chbl
#nodes,16
#keys,100000
#zipf_s,1.2
#mean,6250.0
#max,7813
#cv,0.04321
"""

def load_sim_csv(path):
    data = []
    meta = {}

    with open(path, "r") as f:
        for line in f:
            line = line.strip()
            if line.startswith("#"):
                # parse meta rows, like "#algo,chbl"
                parts = line[1:].split(",")
                if len(parts) == 2:
                    key, val = parts
                    meta[key.strip()] = val.strip()
                continue
            # regular CSV row
            parts = line.split(",")
            if len(parts) == 2 and parts[0] != "node_id":
                data.append({
                    "node_id": parts[0],
                    "count": int(parts[1])
                })

    df = pd.DataFrame(data)
    return df, meta


def per_node_bar_plot(df, meta, outdir):
    algo = meta.get("algo", "unknown")
    nodes = meta.get("nodes", "?")
    keys = meta.get("keys", "?")
    zipf_s = meta.get("zipf_s", "?")

    plt.figure(figsize=(12,6))
    plt.bar(df["node_id"], df["count"], color="steelblue")
    plt.xticks(rotation=40, ha='right')
    plt.title(f"Per-node load distribution\nalgo={algo}, nodes={nodes}, keys={keys}, zipf_s={zipf_s}")
    plt.ylabel("Count")
    plt.tight_layout()

    outpath = os.path.join(outdir, f"per_node_{algo}_nodes{nodes}_zipf{zipf_s}.png")
    plt.savefig(outpath, dpi=200)
    plt.close()
    print(f"[OK] Wrote {outpath}")


def summary_line_plot(csv_paths, stat_key, outpath):
    """
    stat_key: 'cv' or 'max'
    """
    records = []

    for path in csv_paths:
        df, meta = load_sim_csv(path)

        algo = meta.get("algo")
        cv = float(meta.get("cv", 0))
        max_val = float(meta.get("max", 0))
        mean_val = float(meta.get("mean", 1))
        max_over_avg = max_val / mean_val

        if stat_key == "cv":
            y = cv
        elif stat_key == "max_over_avg":
            y = max_over_avg
        else:
            raise ValueError("stat_key must be 'cv' or 'max_over_avg'")

        records.append({
            "algo": algo,
            "y": y
        })

    rec_df = pd.DataFrame(records).sort_values(by="algo")

    plt.figure(figsize=(8,5))
    plt.plot(rec_df["algo"], rec_df["y"], marker="o")
    plt.title(f"{stat_key} vs algorithm")
    plt.xlabel("Algorithm")
    plt.ylabel(stat_key)
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(outpath, dpi=200)
    plt.close()
    print(f"[OK] Wrote {outpath}")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--csv", nargs="+", required=True,
                        help="CSV files from simulator")
    parser.add_argument("--outdir", required=True,
                        help="Output directory for plots")
    args = parser.parse_args()

    os.makedirs(args.outdir, exist_ok=True)

    # For each CSV → make per-node bar chart
    for path in args.csv:
        df, meta = load_sim_csv(path)
        per_node_bar_plot(df, meta, args.outdir)

    # Summary plots for CV and Max/Avg
    summary_line_plot(args.csv, "cv", os.path.join(args.outdir, "summary_cv_vs_algo.png"))
    summary_line_plot(args.csv, "max_over_avg", os.path.join(args.outdir, "summary_maxoveravg_vs_algo.png"))


if __name__ == "__main__":
    main()
