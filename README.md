This project provides an example of how to use a multi backend adapter to aggregate content
produced by multiple backends.

To simulate multiple backend servers, we'll use [Mountebank](http://www.mbtest.org/) with two imposters, defined in imposter1.json
and imposter2.json. Once mountebank is started use curl to load the imposter definitions.

<pre>
curl -i -X POST -H 'Content-Type: application/json' -d@imposter1.json http://127.0.0.1:2525/imposters
curl -i -X POST -H 'Content-Type: application/json' -d@imposter2.json http://127.0.0.1:2525/imposters
</pre>

