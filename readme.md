# Tricky

##### HTTP/HTTPS transparent proxy with content substitution possibility.

Current status: under development, pre-alpha

Version: 0.1.0

### Features

Tricky allows you substitute content of page with content of other page with help of substitution rules.
Substitution rule is a structure of key with associated link.
For example:

Test -> https://google.com

Bingo -> https://github.com

That rule means that if you request a url http://example.vo/?test=true and setup a Tricky as a proxy in Browser/Curl
then at response you get content of http://google.com.
Tricky search keys in request query-params.
You can setup new rules with REST API.

### Dev notes

Built with idea of a [clean architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) in mind