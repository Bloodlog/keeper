GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

![linter workflow](https://github.com/Bloodlog/keeper/actions/workflows/.github/workflows/golangci-lint.yml/badge.svg?event=push)

![vet workflow](https://github.com/Bloodlog/keeper/actions/workflows/.github/workflows/statictest.yml/badge.svg?event=push)

![test workflow](https://github.com/Bloodlog/keeper/actions/workflows/.github/workflows/go-test.yml/badge.svg?event=push)

## ⚙️ Конфигурация сервера

Приложение настраивается с помощью **флагов командной строки**, **переменных окружения** и значений по умолчанию.  
Приоритет: `флаг > переменная окружения > значение по умолчанию`.

---

## Пример запуска сервера:
```bash
./server --dsn="postgres://user:pass@localhost:5432/keeper"
```

## Клиент

Общие флаги подключения к серверу:
* --grpc-address 127.0.0.1 
* --grpc-port 8081
* --token-file .keeper-token
* --enable-tls 
* --ca-cert=cert/public.cert

Пример регистрации:
```bash
keeper-agent register --login test --password -test
```

Пример авторизации:
```bash
keeper-agent login --login test --password -test
```

Пример сохранения ключа С JSON:
```bash
keeper-agent write --path=123 --description="login&password" --value='{"username":"gh-user","password":"gh-pass"}' --max-ttl=1000
```
Пример сохранения файла:
```bash
keeper-agent write --path secret/bar --file ./alice.jpg
```

Пример вывода списка ключей:
```bash
keeper-agent list
```

### Удаление ключей

Мягкое удаление ключа - (ключ больше не будет возвращаться, но его можно восстановить):
```bash
keeper-agent delete --path my/secret/path
```

Восстанавливает ключ после мягкого удаления.
```bash
keeper-agent delete --undelete --version 4 --path my/secret/path
```

Уничтожает саму информацию в ключе -(перезатирает скрытые данные без возможности восстановления)
```bash
keeper-agent delete --destroy --path my/secret/path
```

Удаляет весь ключ, включая все версии и всю информацию.
```bash
keeper-agent delete --metadata --path my/secret/path
```

Пример флага путь до токена:
```bash
keeper-agent delete --path my/secret/path --token-file .keeper-token
```

Пример чтения ключа:
```bash
keeper-agent read --path my/secret/path
```
response:
```bash 
Build version: N/A
Build date: N/A
Build commit: N/A
====== Metadata ======
Key              Value
---              -----
created_time     1970-01-01T00:00:00Z
deletion_time    n/a
destroyed        false
version          3

====== Data ======
Key          Value
---          -----
username     gh-user
password     gh-pass

```
response:
```bash
Build version: N/A
Build date: N/A
Build commit: N/A
====== Metadata ======
Key              Value
---              -----
created_time     1970-01-01T00:00:00Z
deletion_time    n/a
destroyed        false
version          6
Key              Value
---              -----
alice.jpg    FILE-UPLOADED
```

С флагом --out-file, чтобы сохранить значение в файл:
```bash
keeper-agent read --path my/secret/path --out-file 1.json
```

С флагом --out-file, чтобы восстановить файл значение в файл:
```bash
keeper-agent read --path secret/bar --out-file alice2.jpg
```
response:
```bash
Build version: N/A
Build date: N/A
Build commit: N/A
====== Metadata ======
Key              Value
---              -----
created_time     1970-01-01T00:00:00Z
deletion_time    n/a
destroyed        false
version          6
Key              Value
---              -----
alice.jpg    FILE-UPLOADED

✅ Secret written to file: alice2.jpg
```

## TLS:
### Сервер:
Пример генерации сертификатов tls:
```bash
./server gen-cert

```
Пример с флагом включения tls и путь до сертификата:
```bash
keeper-server --grpc-address=0.0.0.0 --grpc-port=9090 --enable-tls --cert-file=cert/public.cert --key-file=cert/private.cert
```

env:
```bash
export KEEPER_GRPC_ADDRESS=0.0.0.0
export KEEPER_GRPC_PORT=9090
export KEEPER_ENABLE_TLS=true
export KEEPER_CERT_FILE=cert/public.cert
export KEEPER_KEY_FILE=cert/private.cert
```

### Клиент

```bash
keeper-agent login --enable-tls --ca-cert=cert/public.cert
```

```bash
export KEEPER_ENABLE_TLS=true
export KEEPER_CA_CERT=cert/public.cert
```
=======
