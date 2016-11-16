# Cliware middlewares
[![Go Report Card](https://goreportcard.com/badge/github.com/delicb/cliware-middlewares)](https://goreportcard.com/report/github.com/delicb/cliware-middlewares)
[![Build Status](https://travis-ci.org/delicb/cliware-middlewares.svg?branch=master)](https://travis-ci.org/delicb/cliware-middlewares)
[![codecov](https://codecov.io/gh/delicb/cliware-middlewares/branch/master/graph/badge.svg)](https://codecov.io/gh/delicb/cliware-middlewares)
[![GoDoc](http://godoc.org/github.com/delicb/cliware-middlewares?status.png)](http://godoc.org/github.com/delicb/cliware-middlewares)

This repository contains some useful middlewares for [Cliware](https://github.com/delicb/cliware)
library.

# Install
Run `go get go.delic.rs/cliware-middlewares` in terminal.


# Content
Middlewares are separated per packages that can be useful on their own. This is done
because not everybody needs all middlewares. Currently following packages exist:

* body - handling request body (JSON, XML, string)
* cookies - handling request cookies (add, set, delete)
* headers - handling request headers (add, set, delete)
* query - handling request query parameters (add, set, delete)
* url - handling URL endpoint for request (base URL, path)

# Credits
Idea and bunch of implementation details were taken from cool GoLang HTTP client
[Gentleman](https://github.com/h2non/gentleman). Difference is that these middewares
are based on Cliware, instead of Gentleman builtin plugin mechanism.
