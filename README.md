This is a simple openweathermap extractor for prometheus written in
Go. Once it is compiled the usage is simple.

```
$ openweather-extractor --lat=... --lon=... --api-key=...
```

The output can be put in a prometheus prom file and picked up by the
node exporter text extractor or other prometheus tools.

