# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).


## Вспомогательный материал

Скомпилировать для автотестов
```
cd ~/go/src/self_dev/sprints/shortener/ & go build -o ./cmd/shortener/shortener ./cmd/shortener/*.go &chmod +x ./cmd/shortener/shortener 
```

Запустить автотесты
```
~/go/external_untrusted_bin/shortenertest-darwin-arm64 -test.v -test.run='^TestIteration1$' -binary-path=./cmd/shortener/shortener
```

```
~/go/external_untrusted_bin/shortenertest-darwin-arm64 -test.v -test.run='^TestIteration2$' -binary-path=./cmd/shortener/shortener -source-path=./
```

Пример запуска с указанием флагов
```
go run ./cmd/shortener/main.go -a=localhost:8080 -b=http://localhost:8080
```

```
go run ./cmd/shortener/main.go -l "Debug" -d "host=localhost user=urls password=jf6y5Sfnxsu
R sslmode=disable port=5432"

```
go run ./cmd/shortener/main.go -l "Debug" -f "./test.json"
```

go run ./cmd/shortener/main.go -l "Debug" -d "host=localhost user=urls password=jf6y5Sfnxsu
R sslmode=disable port=5432"
```

Пример запуска локальных тестов
```
go test ./... -cover -v
```

Пример формирования моков
```
create mocks
mockgen -destination=internal/app/mocks/mock_db.go -package=mocks github.com/kripsy/shortener/internal/app/handlers Repository
```

Пример запуска автотестов из CI CD
 ~/go/src/self_dev/go-autotests/bin/shortenertest -test.v -test.run='^TestIteration12$' -binary-path=./cmd/shortener/shortener -source-path=./ --database-dsn='host=localhost user=postgres  sslmode=disable port=5432 '