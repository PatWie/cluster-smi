# CLUSTER-SMI

The same as `nvidia-smi` but for multiple machines at the same time.

<p align="center"> <img src="./cluster-smi.jpg" width="100%"> </p>

## install
- edit `cluster-smi.yml` (is required by all parts)
- start `cluster-smi-node` at different machines
- start `cluster-smi-server` at a specific machine
- use `cluster-smi` like `nvidia-smi`

Output should be something like

```
+---------+------------------------+---------------------+----------+----------+
| Node    | Gpu                    | Memory-Usage        | Mem-Util | GPU-Util |
+---------+------------------------+---------------------+----------+----------+
| node-00 | 0: TITAN Xp            |  3857MiB / 12189MiB | 31%      | 0%       |
|         | 1: TITAN Xp            | 11689MiB / 12189MiB | 95%      | 0%       |
|         | 2: TITAN Xp            | 10787MiB / 12189MiB | 88%      | 0%       |
|         | 3: TITAN Xp            | 10965MiB / 12189MiB | 89%      | 100%     |
+---------+------------------------+---------------------+----------+----------+
| node-01 | 0: TITAN Xp            | 11667MiB / 12189MiB | 95%      | 100%     |
|         | 1: TITAN Xp            | 11667MiB / 12189MiB | 95%      | 96%      |
|         | 2: TITAN Xp            |  8497MiB / 12189MiB | 69%      | 100%     |
|         | 3: TITAN Xp            |  8499MiB / 12189MiB | 69%      | 98%      |
+---------+------------------------+---------------------+----------+----------+
| node-02 | 0: GeForce GTX 1080 Ti |  1447MiB / 11172MiB | 12%      | 8%       |
|         | 1: GeForce GTX 1080 Ti |  1453MiB / 11172MiB | 13%      | 99%      |
|         | 2: GeForce GTX 1080 Ti |  1673MiB / 11172MiB | 14%      | 0%       |
|         | 3: GeForce GTX 1080 Ti |  6812MiB / 11172MiB | 60%      | 36%      |
+---------+------------------------+---------------------+----------+----------+
```
