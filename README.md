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
./xavi-multi-backend-sample add-server -address localhost -port 5000 -name thing1svr
./xavi-multi-backend-sample add-server -address localhost -port 6000 -name thing2svr
./xavi-multi-backend-sample add-backend -name thing1 -servers thing1svr
./xavi-multi-backend-sample add-backend -name thing2 -servers thing2svr
./xavi-multi-backend-sample add-route -name things-route -backends thing1,thing2 -base-uri /things -plugins SessionId,Timing,Recovery -multibackend-adapter handle-things
./xavi-multi-backend-sample add-listener -name things-listener -routes things-route
</pre>

Once the configuration is in place, boot the listener:

<pre>
./xavi-multi-backend-sample listen -ln things-listener -address 0.0.0.0:8080
</pre>

And curl away...

<pre>
curl localhost:8080/things
</pre>


For using a multi-backend server route with an HTTPs backend, generate a key and cert to use, first by building the
generate cert program that come with golang:

<pre>
cd $GOROOT
cd src/crypto/tls
go build generate_cert.go
</pre>

Then using it to generate your cert:

<pre>
$GOROOT/src/crypto/tls/generate_cert -ca -host `hostname`
</pre>

Note the hostname used for TLS will be verified by golang.

Once the cert and key are in place, update imposter2-https.json with the key and cert.

<pre>
curl -i -X POST -H 'Content-Type: application/json' -d@imposter1.json http://127.0.0.1:2525/imposters
curl -i -X POST -H 'Content-Type: application/json' -d@imposter2-https.json http://127.0.0.1:2525/imposters
</pre>

Then apply the following configuration:

<pre>
export XAVI_KVSTORE_URL=file:///`pwd`/config
./xavi-multi-backend-sample add-server -address localhost -port 5000 -name thing1svr
./xavi-multi-backend-sample add-server -address `hostname` -port 6000 -name thing2svr
./xavi-multi-backend-sample add-backend -name thing1 -servers thing1svr
./xavi-multi-backend-sample add-backend -name thing2 -servers thing2svr -cacert-path ./cert.pem -tls-only=true
./xavi-multi-backend-sample add-route -name things-route -backends thing1,thing2 -base-uri /things -plugins SessionId,Timing,Recovery -multibackend-adapter handle-things
./xavi-multi-backend-sample add-listener -name things-listener -routes things-route
</pre>

Once the configuration is in place, boot the listener:

<pre>
./xavi-multi-backend-sample listen -ln things-listener -address 0.0.0.0:8080
</pre>

And curl away...

<pre>
curl localhost:8080/things
</pre>

Note the above can also be run with a health check. To health check the https endpoint, modify
the server command for thing2svr:

<pre>
./xavi-multi-backend-sample add-server -address `hostname` -port 6000 -name thing2svr -health-check http-get
</pre>