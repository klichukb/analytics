# Analytics aggregation tool

Client/server sample application that aggregates event data.

### Communication interface
Websockets + JSON-RPC.

Pros:
* Good performance and throughput. Capable of handling real-time load.
* Websockets hold connection open, as opposed to REST for example, and not as chatty. 
* Websockets allow browser to be clients.
* Having JSON for data makes assembling messages easy for JS applications.

Cons:
* Apache Thrift can be much faster, even though is more heavy in setup.
* JSON is also a bit chatty.

### Database
MySQL.

Pros:
* Relational - can be good for complex queries involving different research of analytical data on customer's dashboard.
* Scales horizontally. Tools like MySQL Cluster and MySQL Fabric (makes as well sense to have while processing aggregated data).

Cons:
* Is said to have stability issues sometimes.

Good alternatives:
* Cassandra - highly scalable, good for real-time stuff, document store, which may be good for analytics data.
* PostgreSQL - rock-solid and tunable, but horizontal scaling (Citus) are quite new and have SPOF.

### Usage
Get a new database of MySQL, create a table for events:

    CREATE TABLE `analytics_event` (
        `id` int(10)  unsigned primary key auto_increment,
        `event_type` char(20) NOT NULL,
        `ts` timestamp NOT NULL,
        `params` json
    );


Get the app source:

    go get github.com/klichukb/analytics
    
Start server (default port is `8000`):

    MYSQL_USER=user MYSQL_PWD=pwd MYSQL_DB=db $GOPATH/bin/analytics --mode server
    
Start test clients (default port is `8000`):

    $GOPATH/bin/analytics --mode client
    
You'll see output on both sides with message exchange. 

### Data layout
The plan is keep MySQL on server side to receive all the incoming events and dump data on daily/per-event-type basis.
Currently there is a `Makefile` that holds `export` target, which dumps event data for **current day** to CSV, bzips, and uploads to an S3 bucket. Example:

    $ make export
    Dumping '/usr/local/mysql/dumps/dump_2016_11_11.csv'...
    Uploading to S3...
    upload: /usr/local/mysql/dumps/dump_2016_11_11.csv.bz2 to s3://test-bucket/dump_2016_11_11.csv.bz2
    Done.

Querying database for questions like "all events of type X during last 24 hours" are straight forward. Index on `analytics_event.ts` should be set. Optionally on `event_type` (depends on how many event types exist).

We can also slice dumps by event type, in general my assumptions regarding future use are following:
* Most likely we'll want to be able to show to customer all events their app is generating. Reports, views etc. But accessing 
* Future data access patterns are yet uknown, but I assume we'll have some kind of App ID/Customer ID, that we'll be writing into event table as well. In that case we'll want one custom to be able to access his data as quickly as possible. This assumption leads to possible sharding by App/Customer ID or range of such. Data could be dumped regularily per user and stored like that for further processing.


### Source code

Unit tests aren't covering much - still switching from Python world. Small integration test runs client/server for couple of messages.ocial Analytics aggregation.
