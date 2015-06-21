# headmonkey
API server to control containers running [Vaurien](https://vaurien.readthedocs.org/en/1.8/) to control failure scenarios.

The following API are supported

GET /proxies/index

return a list of all known proxy

Exampe:

    curl -XGET localhost:8080/proxy/index

Returns

    [
        “src1:dst2”,
        “src3:dst4”,
    ]


GET /proxies/?proxy=PROXY_ID

return an info object of proxy and the installable behaviours on it and the behaviours installed on it

Example:

    curl -XGET localhost:8080/proxies?proxy=figtest_webredis_1

Returns

	{
		“behavior”: “dummy”
	}

PUT /proxies/?proxy=PROXY_ID JSON_OBJ

Add new behavior to vaurien.

Example:

    curl -XPUT -d '{"sleep": 2, "name": "delay"}' http://localhost:8080/behavior \
        -H "Content-Type: application/json"

Returns 200 and an empty response.
