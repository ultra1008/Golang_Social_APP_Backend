# Highload social network
[![Build](https://github.com/niklod/highload-social-network/workflows/Build/badge.svg)](https://github.com/niklod/highload-social-network/actions)
[![Test](https://github.com/niklod/highload-social-network/workflows/Test/badge.svg)](https://github.com/niklod/highload-social-network/actions)  

Учебный проект для курса OTUS Highload Achitect

## Стек

Golang 1.15 + MySQL 8.0

## Лог бенчмарков

Для нагрузочного тестирования использовалась утилита wrk

1. Поиск пользователей по имени и фамилии.
Без оптимизации
Запросы идут напрямую в http сервер Go. В таблицах БД отсутствуют индексы, кроме PK на id сущностей.

