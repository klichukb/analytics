# # # #
#
# Usage:
#
# 	 make server - start listening to websocket connections on :8000
# 	 make clients - launch mock clients.
# 	 make export - dump database CSV of analytics event for current day.
#
# Run this script daily at the time of minimum load to upload dialy dump
# to S3 storage. Modify optionally to clean up local database using some
# strategy, for example by removing all records older than 2 days, and thus
# allowing to resolve any poissible inconsistencies.
#
# # # #

dump_dir=/usr/local/mysql/dumps
db_user=analytics
db_pwd=analytics
db=analytics
s3_bucket=my-bucket

server:
	@ MYSQL_USER=$(db_user) MYSQL_PWD=$(db_pwd) MYSQL_DB=$(db) \
	  go run wrapper.go models.go api.go db.go server.go

clients:
	@ go run wrapper.go models.go client.go

export:
	$(eval dump_path=$(dump_dir)/dump_$(shell date "+%Y_%m_%d").csv)

	@ echo "Dumping '$(dump_path)'...";
	$(eval query=SELECT * FROM analytics_event \
		WHERE DATE(ts) = CURDATE() \
		INTO OUTFILE '$(dump_path)' \
		FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n';)
	@ MYSQL_PWD=$(db_pwd) mysql -u analytics analytics -D analytics -e "$(query)"
	@ echo "Uploading to S3...";
	@ aws s3 cp $(dump_path) s3://$(s3_bucket)/
