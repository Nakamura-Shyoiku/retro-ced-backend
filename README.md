# Development 

## Pre-requisites 

* Golang version go1.16.6 linux/amd64
* docker-compose `docker-compose version 1.29.2+`
* migrate 4.12.1 (for db migrations)


## Setup db and migrations

1. docker-compose up 
2. make migrate-up
3. make run 

# Quickstart

The project now uses go modules, so in order to use it locally, clone the project to a location outside of your GOPATH and run `go build` at the root of the project.

**Note:**

By default, the codebase connects to cloud mysql instance. Do not continue to use this database for further development work if possible.
Instead, create a `mysql` docker container locally and import the data from the cloud database into the local database.

## API Testing

This is a temporary solution: To check that the following endpoints still work

```
curl -X GET http://localhost:8080/v1/product/categories/bags/0
curl -X GET http://localhost:8080/v1/product/brands/hermes/0
curl -X GET http://localhost:8080/v1/product/brand
curl -X GET http://localhost:8080/v1/product/search/chan/0

# Note: This API needs to be altered slightly, a person who is not logged in is supposed to received error 401/403 Unauthorized access
curl -X GET http://localhost:8080/v1/product/favourites/0 


```


