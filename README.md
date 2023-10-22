# gitch

Сервис для синхронизации репозиториев. В конфиге можно указать с какой периодичностью нужно запускать синхронизацию. 

## Установка

Сейчас нужно склонировать репозиторий локально, собрать бинарник и использовать. Перед использованием переименуйте 
config.template.yml в config.yaml

Сборка сервиса

```
go build -o gitch ./cmd/gitch
```

Запуск сервиса

```
./gitch --env=config server
```

Разовая синхронизация репозиторий

```
./gitch sync --from=git@gitflic.ru:example/example.git --to=git@github.com:example/example.git
```

## Доделки

- [x] Периодический синк репозиториев
- [ ] Синхронизация по веб-хуку
- [ ] Завернуть в Docker для быстрого старта

## Ссылки 

- Весь проект работает на [go-git](https://github.com/go-git/go-git)
- Исправление [проблемы с таймаутом](https://bengsfort.github.io/articles/fixing-git-push-pull-timeout/)
