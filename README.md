go-project-demo
=====
A simple Go project demo

This project is trying to elaborate "how does a Go project look like". And further more, it will be containerised a small service with required build scripts.

What am I going to achieve in this project?
-----
Basically, this project makes a custom web server with pre-processor and post-processor implemented built-in libraries.

In this example, the server will implements a pre-processor function which sends the brief information about each request using channel. A concurrent goroutine function will receive and persist all incoming request information. And provides request per minute statistic on foreground UI.

BDD?
-----
This project was supposed to use the BDD approach. Which is, having the feature in Gherkin format and testing scripts prepared before project implementation starts. However, I set a 5 hours time limit to finish this project. Under my estimation, I don't have enough time to complete testing scripts within this time limit.
