- [Giới thiệu](#giới-thiệu)

### Create database mysql in docker

```bash

docker run --name payday -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=payday -p 3309:3306 -d mysql

```
