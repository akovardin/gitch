# Gitch

Сервис для синхронизации репозиториев.

## Установка

Запуск в докере

```sh
docker run -v /opt/gitch/data:/data  -p 8080:8080 -d registry.gitflic.ru/project/kovardin/gitch/gitch:latest --dir /data --dev  --http :8080 serve 
```

`opt/gitch/data` - директория с данными 

## Ссылки 

- Весь проект работает на [go-git](https://github.com/go-git/go-git)
- Исправление [проблемы с таймаутом](https://bengsfort.github.io/articles/fixing-git-push-pull-timeout/)
