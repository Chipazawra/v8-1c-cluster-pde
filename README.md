1C-RAS Prometheus data exporter
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

### Результатом проделанных действий на `http://<host-ip>:3000` или есть вы делали это на своей машине `http://localhost:3000` будет доступна Grafan c demo дабордом:

![image](https://user-images.githubusercontent.com/18016416/147658207-6f03553a-354c-43dd-9ac5-a95accd9bbe9.png)

