![UnitTest](https://github.com/vinyguedess/pismo-test/actions/workflows/unit_test.yml/badge.svg)

# Pismo: Test
Development test for Pismo anti-fraud team.

## Getting started
First, create a `.env` file using `.env.default` as example.

## Executing container
`make up` command is configured to lift up containers and ssh into app container.

## Running tests
Once you're into running container, you can run `make test` command. It will execute tests and print coverage information in bash but it's also possible to check deep coverage data in `coverage/index.html` generated in the end of the tests.

## Running API
`make run` command is configured to run API in development mode. It will run API in port 8000 mapped to port 6073 and you can access it in `http://localhost:6073`.
Swagger documentation is accessible through `http://localhost:6073/swagger`. There you can test all endpoints developed.
