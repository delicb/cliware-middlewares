# Cliware middlewares
[![Go Report Card](https://goreportcard.com/badge/github.com/delicb/cliware-middlewares)](https://goreportcard.com/report/github.com/delicb/cliware-middlewares)
[![Build Status](https://travis-ci.org/delicb/cliware-middlewares.svg?branch=master)](https://travis-ci.org/delicb/cliware-middlewares)
[![codecov](https://codecov.io/gh/delicb/cliware-middlewares/branch/master/graph/badge.svg)](https://codecov.io/gh/delicb/cliware-middlewares)
![status](https://img.shields.io/badge/status-beta-red.svg)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/delicb/cliware-middlewares)

This repository contains some useful middlewares for [Cliware](https://github.com/delicb/cliware)
library.

# Install
Run `go get go.delic.rs/cliware-middlewares` in terminal.

# Usage
Cliware middlewares are set of tools that can be used to create useful HTTP clients.
On there own they do not do much, but to see how they they can be used to put
together cool HTTP client, take a loko at [GWC](https://github.com/delicb/gwc) (HTTP
client built on top of Cliware and Cliware-middlewares).

# Content
Middlewares are separated per packages that can be useful on their own. This is done
to avoid dependencies between middlewars and create cleaner naming schema. 
Currently following packages exist:

* auth - authentication via header support
* body - handling request body, support setting JSON, XML, string and from io.Reader
* cookies - handling request cookies (add, set, delete)
* errors - handling HTTP error status codes and converting them to GoLang errors
* headers - handling request headers (add, set, delete)
* query - handling request query parameters (add, set, delete)
* responsebody - managing respones body, get json, string or write raw content to own writer
* retry - request retry mechanism based on custom classifier and with custom backoff
* url - handling URL endpoint for request (base URL, path)

# Credits
Idea and bunch of implementation details were taken from cool GoLang HTTP client
[Gentleman](https://github.com/h2non/gentleman). Difference is that these middewares
are based on Cliware, instead of Gentleman builtin plugin mechanism.

Some ideas for retry middlewares are from [go-resiliency](https://github.com/eapache/go-resiliency)
project.

# Licence
Cliware-middlewares is released under MIT licence.

