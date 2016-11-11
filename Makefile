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
s3_bucket=elasticbeanstalk-us-west-2-422680256038

runserver:
	@ MYSQL_USER=$(db_user) MYSQL_PWD=$(db_pwd) MYSQL_DB=$(db) $(GOPATH)bin/analytics --mode server

runclient:
	@ $(GOPATH)bin/analytics --mode client

export:
	$(eval dump_path=$(dump_dir)/dump_$(shell date "+%Y_%m_%d").csv)

	@ echo "Dumping '$(dump_path)'...";
	$(eval query=SELECT * FROM analytics_event \
		WHERE DATE(ts) = CURDATE() \
		INTO OUTFILE '$(dump_path)' \
		FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n';)
	@ MYSQL_PWD=$(db_pwd) mysql -u analytics analytics -D analytics -e "$(query)"
	@ bzip2 $(dump_path)
	@ echo "Uploading to S3...";
	@ aws s3 cp $(dump_path).bz2 s3://$(s3_bucket)/
	@ echo "Done."
