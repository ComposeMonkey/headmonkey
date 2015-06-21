# headmonkey
# CROSS ORIGIN FIXED + PARSING LOGIC FIXED
server to orchestrate failure scenarios

Provide APIs to caller:
GET /proxies/index
return a list of all known proxy and their type
Sample Response =>
 [
	“src1:dst2”,
	“src3:dst4”,
]

Exampe: 

    curl -XGET localhost:8080/proxy/index

GET /proxies/?proxy=PROXY_ID
return an info object of proxy and the installable behaviours on it and the behaviours installed on it
	{
		“name”: “dummy”
	}
Example: 

    curl -XGET localhost:8080/proxies?proxy=figtest_webredis_1

PUT /proxies/?proxy=PROXY_ID JSON_OBJ
Example:

    curl -XPUT -d '{"sleep": 2, "name": "delay"}' http://localhost:8080/behavior \
        -H "Content-Type: application/json"
