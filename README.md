go-project-demo
=====
A simple Go project demo

This project is trying to elaborate "how does a Go project look like". And further more, it will be containerised a small service with required build scripts.

What am I going to achieve in this project?
-----
Basically, this project makes a custom web server with pre-processor and post-processor implemented built-in libraries.

In this example, the server will implements a pre-processor function which sends the brief information about each request using channel. A concurrent goroutine function will receive and persist all incoming request information. And provides request per minute statistic on each HTTP response.

BDD?
-----
This project was supposed to use the BDD approach. Which is, having the feature in Gherkin format and testing scripts prepared before project implementation starts. However, I set a 5 hours time limit to finish this project. Under my estimation, I don't have enough time to complete testing scripts within this time limit.

Integration Test
-----

You may either perform the integration test on [local instance](#run-local-instance) or [docker image](#docker-image)

##### Run Local 
This project requires Go v1.12.5. Please make sure your local Go version is v1.12.5 or above.

Within the project, execute below command. The proxy will listen to port 8080 by default.

`$ make run`

##### Docker Image
You may try to containerise this service. However, the container does not persist the history data file to a mounted volume. It limits the test scope of this project.

```
make image
./run-docker-image.sh
```

The service will be listening on your port 8080.

### Integration Test Commands

Make you you have [cURL](https://curl.haxx.se/) installed on your system.

You can try below commands to test the server:

`curl -L 'http://localhost:8080/'`

You may expect a "hello" plain text.

`curl -L 'http://localhost:8080/stat'`

You may expect to get the statistic number of handled request in last 60 seconds.
