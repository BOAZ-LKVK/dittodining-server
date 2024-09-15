APP?=dittodining-mysql
DEFAULT_PASSWORD?=password

.PHONY: run-mysql
run-mysql:
	echo '[client]\npassword=${DEFAULT_PASSWORD}' > /tmp/mysqlconfig.cnf
	docker run --rm -d \
		--name ${APP} \
		-e MYSQL_ROOT_HOST=% \
		-e MYSQL_ROOT_PASSWORD=${DEFAULT_PASSWORD} \
		-v /tmp/mysqlconfig.cnf:/mysql/config.cnf \
		-p 3306:3306 \
		mysql/mysql-server:8.0.23 \
			--character-set-server=utf8mb4 \
			--explicit_defaults_for_timestamp=true
