# CLUSTER-SMI

The same as `nvidia-smi` but for multiple machines at the same time.

## install
- edit const.go
- start `cluster-smi-node` at different machines
- start `cluster-smi-server` at a specific machine
- use `cluster-smi` like `nvidia-smi`

Output should be something like

```
|  NODE     |          GPU           |       MEMORY-USAGE        | GPU-UTIL |
+-----------+------------------------+---------------------------+----------+
| machine-0 | 0: GeForce GTX 1080    | 3857MiB / 12189MiB (0%)   | 0%       |
| machine-0 | 1: GeForce GTX 1080    | 11689MiB / 12189MiB (2%)  | 5%       |
| machine-0 | 2: GeForce GTX 1080    | 10787MiB / 12189MiB (0%)  | 0%       |
| machine-0 | 3: GeForce GTX 1080    | 10965MiB / 12189MiB (2%)  | 4%       |
| machine-1 | 0: TITAN Xp            | 11667MiB / 12189MiB (8%)  | 100%     |
| machine-1 | 1: TITAN Xp            | 11667MiB / 12189MiB (7%)  | 100%     |
| machine-1 | 2: TITAN Xp            | 8497MiB / 12189MiB (21%)  | 99%      |
| machine-1 | 3: TITAN Xp            | 8499MiB / 12189MiB (32%)  | 100%     |
| machine-2 | 0: GeForce GTX 1080 Ti | 1449MiB / 11172MiB (72%)  | 99%      |
| machine-2 | 1: GeForce GTX 1080 Ti | 1453MiB / 11172MiB (55%)  | 100%     |
| machine-2 | 2: GeForce GTX 1080 Ti | 1355MiB / 11172MiB (61%)  | 98%      |
| machine-2 | 3: GeForce GTX 1080 Ti | 6812MiB / 11172MiB (46%)  | 91%      |

```