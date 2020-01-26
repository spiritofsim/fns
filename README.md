# Библиотека для проверки чеков в ФНС

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/spiritofsim/fns)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/spiritofsim/fns/Go)
![Codecov](https://img.shields.io/codecov/c/github/SpiritOfSim/fns)
[![Go Report Card](https://goreportcard.com/badge/github.com/spiritofsim/fns)](https://goreportcard.com/report/github.com/spiritofsim/fns) 

Библиотека позволяет проверить чек в ФНС, а так же получить по нему полную информацию, включая список покупок.

Для получения пароля можно воспользоваться функцией `Register`, либо установить приложение "Проверка касcового чека" и зарегистрироваться через него.

* [IOS](https://appsto.re/ru/TKUSfb.i)
* [Android](https://play.google.com/store/apps/details?id=ru.fns.billchecker)

В обоих случаях, после успешной регистрации пароль придет в смс на указанный номер.

### Примеры 

* `Register(context.Background(), "<email>", "<name>", "<phone>")` - регистрация пользователя
* `CheckReceipt(context.Background(), "9251440300046840", 1, 29414, 1250830908, time.Date(2020, 1, 15, 21, 10, 0, 0, time.UTC), 1030)` - проверка чека
* `GetReceipt(context.Background(), "<phone>", "<password>", "9251440300046840", 29414, 1250830908)` - получение полной информации о чеке
* `ParseQrStr("t=20200115T2110&s=1030.00&fn=9251440300046840&i=29414&fp=1250830908&n=1")` - конвертация данных с QR кода на чеке

Перед получением полной информации о чеке необходимо хотя бы раз выполнить проверку существования чека с помощью `CheckReceipt` 

