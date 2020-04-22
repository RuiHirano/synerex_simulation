# データ交換基盤の分散によるオーバーヘッド抑制の調査

## 自分のパソコンでの結果
### 1. データ交換基板一つにおけるパフォーマンスを評価
Node:1
data-plt:1
worker-provider:20

agentprovder20: 
agents0: 30ms
800: 50ms
1600: 60ms
2200: 97ms
2900: 100ms
7500: 190ms
10500: 280ms
13000: 440ms
14800: 540ms
15800: 700ms
17500: 980ms

17500 * 20 = 350000agents


### 2. データ交換基板三つによるパフォーマンスを評価
Macパソコン
Node:1
data-plt:20
worker-provider:20

agents0: 51ms
800: 76ms
1600: 92ms
2400: 105ms
3900: 150ms
6150: 185ms
9500: 340ms
12500: 460ms
15800: 630ms
17400: 750ms
18000: 900ms

### 考察
どちらもそこまで変わらない結果となった。
通信がある分、2のが不利な感じ
1のagentproviderがcontainerなのが有利かも（local通信なので）

## 仮想サーバ4で検証
### 1. データ交換基板一つにおけるパフォーマンスを評価
Node:4
data-plt:1
worker-provider:20

agentprovder20: 
agents0: 30ms
800: 50ms
1600: 60ms
2200: 97ms
2900: 100ms
7500: 190ms
10500: 280ms
13000: 440ms
14800: 540ms
15800: 700ms
17500: 980ms

17500 * 20 = 350000agents


### 2. データ交換基板三つによるパフォーマンスを評価
Node:4
data-plt:20
worker-provider:20

agents0: 26ms
829: 55ms
1650: 86ms
3200: 160ms
4500: 260ms
5800: 280ms
7000: 390ms
9500: 610ms
11000: 690ms
13000: 840ms
14000: 980ms


### 考察




