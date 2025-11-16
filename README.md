# DependencyChecker-DC-

## Как пользоваться утилитой?

1. Склонируйте репозиторий на свой компьютер:
```bash
git clone https://github.com/Zaharazov/DependencyChecker-DC-.git
```
2. Соберите проект командой:
```bash
go build main.go
```
3. Запустите утилиту в консоли:
```bash
./main https://github.com/rep_name
```
4. Получите результат анализа.

## Примеры вывода

* ./main https://github.com/Zaharazov/Account-Service
```
| Module Name: Account-Service
|-------------------------------------------------------
| Go Version: 1.21
|-------------------------------------------------------
| Can Be Updated: go.mongodb.org/mongo-driver ( v1.16.0 -> v1.17.6 )
| Can Be Updated: github.com/go-openapi/jsonpointer ( v0.19.5 -> v0.22.1 )
...
| Can Be Updated: golang.org/x/text ( v0.14.0 -> v0.31.0 )
```

* ./main https://github.com/Zaharazov/NewReleaseTracker-NRT-
```
| Module Name: NRT
|-------------------------------------------------------
| Go Version: 1.24.3
|-------------------------------------------------------
| There Are No Possible Updates!
```

## Используемая литература
* https://pkg.go.dev/golang.org/x/mod/semver
* https://leapcell.io/blog/handling-command-line-arguments-in-golang
* https://deepsource.com/blog/go-modules
* https://docs.github.com/en/repositories/working-with-files/using-files/viewing-and-understanding-files#viewing-the-raw-version-of-a-file
* https://pkg.go.dev/golang.org/x/mod/modfile
* https://go.dev/ref/mod#module-versions
* https://docs.github.com/ru/rest/git/trees?apiVersion=2022-11-28#get-a-tree
* https://pkg.go.dev/strings
* https://pkg.go.dev/io
* https://pkg.go.dev/net/url
* https://docs.gitlab.com/api/repositories
* https://docs.gitlab.com/development/routing/