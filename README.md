## A prototype application (for education purpose) that provides loadbalancing and service discovery built only with Go standard
### Features:
- Enable loadbalancing according to info in `config.json`;
- Middleware platform (not great but its ok);
- Discover which services are onlineand gather data about then;
- Save current info about servers when interrupted or exited.
- Expose an endpoint(/info) with real time loadbalancing info
```
{
	"host": "localhost",
	"port": "8080",
	"services": [
		{
			"id": "localhost:8000SOMETHING",
			"url": "localhost:8000",
			"alive": false,
			"healthcheck": {
				"delay": 5,
				"endpoint": "/",
				"duration": 0.283,
				"got_error": true
			},
			"secret_key": "SOMETHING",
			"timeout": 3,
			"endpoints": {
				"/": {
					"allowed_methods": [
						"GET"
					],
					"allowed_hosts": [
						"localhost"
					]
				}
			}
		},
		{
			"id": "localhost:7000SOMETHING",
			"url": "localhost:7000",
			"alive": false,
			"healthcheck": {
				"delay": 5,
				"endpoint": "/",
				"duration": 0.416,
				"got_error": true
			},
			"secret_key": "SOMETHING",
			"timeout": 3,
			"endpoints": {
				"/": {
					"allowed_methods": [
						"GET"
					],
					"allowed_hosts": [
						"localhost"
					]
				}
			}
		}
	]
}
```


