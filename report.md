# 2020/04/06

## 起きている問題
- masterからのrequestがworkerまで届いていない
    - synerexでworkerからのresponseを受けれている
    - clock startも正常に動作している
    - woker-providerの出力に表示されていないだけ？

- worker-nodeidにworker-providerが登録されていない
    - registをmasterかworkerかわかる様に表示するようにする
    - worker川の登録処理をコメントアウトしたままビルドしたから？

- 再度全部buildしてから試す
    - masterでSyncError...がでてしまう
    - master,workerにproviderがうまく登録できていない？
        - &{1247099226094084097 1269267341 CLOCK_SERVICE 0xc00017a6a0 {Client:Worker} 0} SubscribeSupply Error rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp :10000: connect: connection refused"エラーが頻繁に
        - どっちにしろ、再connectする機能は必要

- 途中でworker-providerのmasterNodeId登録が切れる

- bashにsleepをいれて順番を保証したら動いた！！
    - masternodeidからworker が消えてるが、一応動いている