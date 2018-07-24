# Hesperiden 🌳

Hesperiden is a prototype service for realizing too-simple Continuous Deployment pipelines with Docker. By using a GraphQL API interface, CI/CD services can trigger the ‘recreation’ of containers with their prior settings, but using a freshly-obtained image. Access to containers is limited by using grant-specific tokens.


## Setup

The only prerequisite to run the serive is Docker. To launch the service, create a `registries.json`, derived from `registries-example.json`. This allows you to use private registries that are secured with HTTP authentication. If desired, assign an exposed port in the `docker-compose.yml` configuration. (By default, port 80 is used). Finally, launch the service via `docker-compose up -d`.

You may use the following environment variables to configure the service:

* `HTTP_PORT` specifies on which port the HTTP interface will listen, which defaults to `80`.
* `REDIS_HOST` and `REDIS_PORT` specify how to connect to the Redis server.
* `REDIS_PREFIX` specifies the prefix used to store tokens. The pattern used is `<prefix>_<token>`, and the prefix defaults to `token`.
* `PRODUCTION` will run the service in production mode when set to `true`. Using HTTP authentication is recommended when in production, see `AUTH_*` below. By default, the production environment is _not_ active.
* `AUTH_USER` and `AUTH_PASSWORD` specify the credentials used to secure the service via HTTP basic authentication.
* `RATE_LIMIT` specifies the rate limiting enforced on the HTTP interface. For the syntax, please refer [to the `limiter` package](https://github.com/ulule/limiter). The limiting defaults to 30 requests per minute, `30-M`.


## Usage

Once the service is running, you may add tokens to the Redis configuration, revoke them, and recreate containers by using tokens. Please use a GraphiQL-inspired client to learn about the GraphQL API that exposes these functions on `/graphql`, where everything is documented in detail.


## Development

Use Golang `v1.10.3` with `dep` `v0.4.1`.
