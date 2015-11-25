This project provides an example of how to use a multi backend adapter to aggregate content
produced by multiple backends.

To simulate multiple backend servers, we'll use [Mountebank](http://www.mbtest.org/) with two imposters, defined in imposter1.json
and imposter2.json. Once mountebank is started use curl to load the imposter definitions.

<pre>
curl -i -X POST -H 'Content-Type: application/json' -d@imposter1.json http://127.0.0.1:2525/imposters
curl -i -X POST -H 'Content-Type: application/json' -d@imposter2.json http://127.0.0.1:2525/imposters
</pre>

We want to provide an API that aggregates the data returned by the two imposters.

To do this, we need to:

* Provide an implementation of a plugin.MultiBackendHandlerFunc that performs the aggregation using multiple backends.
* Provide a factory function (type plugin.MultiBackendAdapterFactory) that can instantiate the MultiBackendHandlerFunc implementation.
* Register the factory function via plugin.RegisterMultiBackendAdapterFactory

See adapter/things.go for an implementaion of the first two items, and refer to main.go to see the factory registration.

To wire up the API, use the following commands:

<pre>
export XAVI_KVSTORE_URL=file:///`pwd`/config
./multi-backend-sample add-server -address localhost -port 5000 -name thing1svr
./multi-backend-sample add-server -address localhost -port 6000 -name thing2svr
./multi-backend-sample add-backend -name thing1 -servers thing1svr
./multi-backend-sample add-backend -name thing2 -servers thing2svr
./multi-backend-sample add-route -name things-route -backends thing1,thing2 -base-uri /things -multibackend-adapter handle-things
./multi-backend-sample add-listener -name things-listener -routes things-route
</pre>

Once the configuration is in place, boot the listener:

<pre>
./multi-backend-sample listen -ln things-listener -address 0.0.0.0:8080
</pre>

