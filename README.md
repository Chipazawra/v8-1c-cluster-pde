1С-RAS Prometheus data exporter
========
1C-RAS Экспортер метрик в Prometheus 
========
Ещё один экспортер метрик для Prometheus с 1C-RAS, не требующий установки 1C-RAC благодаря пакету https://github.com/khorevaa/ras-client.
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

### Конфигурирование экспортера выполняется средствами установки переменных окружения или параметров командной строки(имеют более высокий приоритет чем переменные окружения):
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


### Экспортируемые метрики, пример:
```
//HELP rp_hosts_active count of active rp hosts on cluster
//TYPE rp_hosts_active gauge
rp_hosts_active{cluster="Локальный кластер"} 1
//HELP rp_hosts_available_perfomance available host performance
//TYPE rp_hosts_available_perfomance gauge
rp_hosts_available_perfomance{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 100
//HELP rp_hosts_avg_back_call_time host avg back call time
//TYPE rp_hosts_avg_back_call_time gauge
rp_hosts_avg_back_call_time{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 0
//HELP rp_hosts_avg_call_time host avg call time
//TYPE rp_hosts_avg_call_time gauge
rp_hosts_avg_call_time{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 0.42917682026288584
//HELP rp_hosts_avg_db_call_time host avg db call time
//TYPE rp_hosts_avg_db_call_time gauge
rp_hosts_avg_db_call_time{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 0.0513839334863925
//HELP rp_hosts_avg_lock_call_time host avg lock call time
//TYPE rp_hosts_avg_lock_call_time gauge
rp_hosts_avg_lock_call_time{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 0
//HELP rp_hosts_avg_server_call_time host avg server call time
//TYPE rp_hosts_avg_server_call_time gauge
rp_hosts_avg_server_call_time{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 0.3705480961587519
//HELP rp_hosts_avg_threads average number of client threads
//TYPE rp_hosts_avg_threads gauge
rp_hosts_avg_threads{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 3.2406813319453436
//HELP rp_hosts_capacity host capacity
//TYPE rp_hosts_capacity gauge
rp_hosts_capacity{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 1000
//HELP rp_hosts_connections number of connections to host
//TYPE rp_hosts_connections gauge
rp_hosts_connections{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 11
//HELP rp_hosts_enable host enable
//TYPE rp_hosts_enable gauge
rp_hosts_enable{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 1
//HELP rp_hosts_memory count of active rp hosts on cluster
//TYPE rp_hosts_memory gauge
rp_hosts_memory{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 1.501164e+06
//HELP rp_hosts_memory_excess_time host memory excess time
//TYPE rp_hosts_memory_excess_time gauge
rp_hosts_memory_excess_time{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 0
//HELP rp_hosts_running host enable
//TYPE rp_hosts_running gauge
rp_hosts_running{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 1
//HELP rp_hosts_scrape_duration the time in milliseconds it took to collect the metrics
//TYPE rp_hosts_scrape_duration gauge
rp_hosts_scrape_duration 1031
//HELP rp_hosts_selection_size host selection size
//TYPE rp_hosts_selection_size gauge
rp_hosts_selection_size{cluster="Локальный кластер",host="192.168.1.1",pid="13552",port="1560",startedAt="2021-12-23 02:25:33"} 297924
```
### Расшифровка метрики:
```
rp_hosts_avg_call_time - Реакция сервера
rp_hosts_avg_server_call_time - Затрачено сервером
rp_hosts_avg_db_call_time - Затрачено СУБД
rp_hosts_avg_lock_call_time - Затрачено мененджером блокировок
rp_hosts_avg_threads - Клиентских потоков
rp_hosts_memory - Занято памяти
rp_hosts_connections - Соединений
rp_hosts_avg_back_call_time - Затрачено клиентом
rp_hosts_available_perfomance - Доступная производительность
rp_hosts_enable - Включен
rp_hosts_running - Активен
```
### Расшифровка метки:
```
сluster - Кластер
host - Хост рабочего процесса
port - Порт рабочего процесса
pid - Идентификатор рабочего процесса
startedAt - Время запуска рабочего процесса
```