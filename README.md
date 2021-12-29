1С-RAS Prometheus data exporter
========
1C-RAS Экспортер метрик в Prometheus 
========
Ещё один экспортер метрик для Prometheus с 1C-RAS, не требующий установки 1C-RAC.
Экспортер работает в двух режимах push и pull.
Для работы в режиме push требуется pushgateway в таргетах Prometheus. 
На текущий момент экспортируются показатели запущеных rpHosts.

Demo
========
### Склонируйте репозиторий, отредактируйте файл .env:
```
git clone https://github.com/Chipazawra/v8-1c-cluster-pde
```
### Содержимое ./.env:
```
RAS_HOST=<ras host>
RAS_PORT=<ras port>
CLS_USER=<ras user - если есть>
CLS_PASS=<ras pass - если есть>
```
### Запустите docker-compose:
```
docker-compose up
```

### Результатом проделанных действий на `http://<host-ip>:3000` или если вы делали это на своей машине `http://localhost:3000` будет доступна Grafana c demo дашбордом:

![image](https://user-images.githubusercontent.com/18016416/147658562-322a2f01-61d7-496a-a256-57d11ae6beae.png)

### Конфигурирование экспортера выполняется средствами установки переменных окружения или параметров командной строки(имеют более высокий приоритет чем переменный окружения):
```
RAS_HOST - хост где запущен 1С-RAS, при запуске через терминал --ras-host
RAS_PORT - порт где запущен 1С-RAS, при запуске через терминал --ras-port
CLS_USER - пользователь 1С-RAS
CLS_PASS - пароль пользователя 1С-RAS
MODE - режим работы экспортера принимает 2 значения push/pull 
PULL_EXPOSE - порт хоста на котором запущен экспортер `http://<host>:<PULL_EXPOSE>/metrics`, при запуске через терминал --pull-expose. Имеет смысл только в режиме pull

Имеют смысл только в режиме push: 

PUSH_INTERVAL - интервал в миллисекундах, с которым экпортер отправляет метрики в pushgateway, при запуске через терминал --push-interval. 
PUSH_HOST - хост pushgateway, с которым экпортер отправляет метрики в pushgateway, при запуске через терминал --push-host
PUSH_PORT - порт pushgateway, с которым экпортер отправляет метрики в pushgateway, при запуске через терминал --push-port
```
