# # # #
#
# Usage:
#
# 	 make server - start listening to websocket connections on :8000
# 	 make clients - launch mock clients.
# 	 make export - dump database CSV of analytics event for current day.
#
# # # #

dump_dir=/usr/local/mysql/dumps
db_user=analytics
db_pwd=analytics
db=analytics

server:
	@ MYSQL_USER=$(db_user) MYSQL_PWD=$(db_pwd) MYSQL_DB=$(db) \
	  go run wrapper.go models.go api.go db.go server.go

clients:
	@ go run wrapper.go models.go client.go

export:
	$(eval dump_path=$(dump_dir)/dump_$(shell date "+%Y_%m_%d").csv)
	@ echo "Dumping '$(dump_path)'...";
	$(eval query=SELECT * FROM analytics_event INTO OUTFILE '$(dump_path)' FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n';)
	@ MYSQL_PWD=$(db_pwd) mysql -u analytics analytics -D analytics -e "$(query)"
