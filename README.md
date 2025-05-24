GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

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
./agent register --login test --password -test
```

Пример авторизации:
```bash
./agent login --login test --password -test
```

Пример сохранения ключа:
```bash
./agent write --path=123 --description="login&password" --value='{"username":"gh-user","password":"gh-pass"}' --max-ttl=1000
```

Пример вывода списка ключей:
```bash
./agent list
```

Пример удаления ключа:
```bash
keeper-agent delete --path my/secret/path
```

Пример флага путь до токена:
```bash
keeper-agent delete --path my/secret/path --token-file .keeper-token
```

Пример чтения ключа:
```bash
keeper-agent read --path my/secret/path
```

С флагом --out-file, чтобы сохранить значение в файл:
```bash
keeper-agent read --path my/secret/path --out-file 1.json
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
