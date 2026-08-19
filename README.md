# check-dns-multi

複数のDNSサーバーへ同じ問い合わせを送り、名前解決の結果を監視するコマンドです。DNSサーバーごとの応答を表示し、監視ツールから利用できる終了コードを返します。指定したDNSサーバーへの問い合わせは並列に実行されます。

## インストール

GitHub Releasesから、実行環境に合ったパッケージをダウンロードします。現在はLinuxのamd64とarm64に対応しています。以下はLinux amd64向けの最新版をダウンロードして、`~/.local/bin` にインストールする例です。

```sh
curl -fL -o check-dns-multi.zip \
	https://github.com/monitoring-forge/check-dns-multi/releases/latest/download/check-dns-multi_linux_amd64.zip
unzip -o check-dns-multi.zip
mkdir -p "$HOME/.local/bin"
install -m 0755 check-dns-multi "$HOME/.local/bin/check-dns-multi"
```

arm64環境では、ファイル名の `linux_amd64` を `linux_arm64` に変更してください。
`~/.local/bin` にPATHが通っていない場合は、PATHに追加するか、インストール先を直接指定して実行します。

## 使い方

```text
check-dns-multi [OPTIONS]
```

### オプション

| オプション | デフォルト | 説明 |
| --- | --- | --- |
| `-v`, `--version` | - | バージョン、OS、アーキテクチャ、Goバージョン、コミットIDを表示します。 |
| `--protocol=tcp\|udp` | `udp` | DNS問い合わせに使用するプロトコルです。 |
| `-p`, `--port=PORT` | `53` | DNSサーバーのポート番号です。 |
| `-H`, `--hostname=HOST` | `127.0.0.1` | DNSサーバーのホスト名またはIPアドレスです。複数指定できます。 |
| `-Q`, `--question=NAME` | `example.com.` | 問い合わせるホスト名です。必要に応じて末尾に `.` を付けます。 |
| `-q`, `--querytype=A\|AAAA` | `A` | 問い合わせるレコード種別です。`A` または `AAAA` を指定できます。 |
| `-E`, `--expect=STRING` | なし | `A` または `AAAA` の回答値に含まれることを期待する文字列です。 |
| `--timeout=DURATION` | `5s` | 1台のDNSサーバーへのタイムアウトです。例: `500ms`, `10s`。 |
| `--all` | 無効 | すべてのDNSサーバーで名前解決が成功した場合だけOKにします。 |
| `-h`, `--help` | - | ヘルプを表示します。 |

`-H` は指定した回数だけ繰り返せます。`-E` を指定した場合、問い合わせ自体が成功しても回答に期待文字列が含まれなければ、そのサーバーは失敗扱いになります。

## 実行例

### 複数のDNSサーバーを確認する

```sh
./check-dns-multi -H 8.8.8.8 -H 1.1.1.1 -Q example.com. -q A
```

出力例:

```text
DNS OK: [8.8.8.8:53] HEADER-> ;; opcode: QUERY, status: NOERROR, id: 61175 ;; flags: qr rd ra;
[8.8.8.8:53] ANSWER-> example.com.  62  IN  A  93.184.216.34
[1.1.1.1:53] HEADER-> ;; opcode: QUERY, status: NOERROR, id: 16450 ;; flags: qr rd ra;
[1.1.1.1:53] ANSWER-> example.com.  268 IN  A  93.184.216.34
```

`HEADER->` はDNS応答のヘッダー、`ANSWER->` は応答に含まれるレコードです。実際のIPアドレス、TTL、問い合わせIDなどはDNSサーバーの応答によって変わります。また、問い合わせは並列実行されるため、DNSサーバーの表示順は一定ではありません。

### IPv6（AAAA）を確認する

```sh
check-dns-multi -H 8.8.8.8 -Q example.com. -q AAAA
```

### 応答内容を検証する

```sh
check-dns-multi -H 8.8.8.8 -H 1.1.1.1 -Q example.com. -q A -E 93.184.216.34
```

### すべてのDNSサーバーを必須にする

通常は、複数台のうち少なくとも1台が成功すればOKです。すべてのDNSサーバーが利用可能であることを確認する場合は `--all` を指定します。

```sh
check-dns-multi --all -H 8.8.8.8 -H 1.1.1.1 -Q example.com.
```

## 終了コードと判定

| 終了コード | 判定 | 条件 |
| ---: | --- | --- |
| `0` | OK | 通常は1台以上、`--all` 指定時は全台で名前解決が成功した場合 |
| `2` | CRITICAL | 通常は全台失敗、`--all` 指定時は1台以上失敗した場合 |
| `3` | UNKNOWN | オプション不正など、チェックを実行できなかった場合 |

成功時の先頭には `DNS OK:` が表示されます。失敗したサーバーには、たとえば次のようなメッセージが表示されます。

```text
[192.0.2.53:53] failed to resolve: read udp 192.0.2.53:53: i/o timeout
```

監視ツールからは、出力だけでなく終了コードも判定に使用してください。

## mackerel-agent 設定例

[mackerel-agent](https://github.com/mackerelio/mackerel-agent) で利用する場合は、`mackerel-agent.conf` の `[plugin.checks]` セクションに `command` を指定します。以下は `check-dns-multi` を `/usr/local/bin` に配置した場合の例です。

### 複数DNSサーバーのうち少なくとも1台が応答すればOK

```conf
[plugin.checks.dns_multi]
command = "/usr/local/bin/check-dns-multi -H 8.8.8.8 -H 1.1.1.1 -Q example.com. -q A"
```

### すべてのDNSサーバーが応答することを要求する

```conf
[plugin.checks.dns_multi_all]
command = "/usr/local/bin/check-dns-multi --all -H 8.8.8.8 -H 1.1.1.1 -Q example.com."
```

### 応答に特定のIPアドレスが含まれることを確認する

```conf
[plugin.checks.dns_multi_expect]
command = "/usr/local/bin/check-dns-multi -H 8.8.8.8 -H 1.1.1.1 -Q example.com. -q A -E 93.184.216.34"
```

設定後、`mackerel-agent` を再起動すると、監視結果が Mackerel に送信されます。

## ライセンス

[LICENSE](LICENSE)
