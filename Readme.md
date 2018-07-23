# Hesperiden 🌳

Hesperiden is a prototype service for realizing too-simple Continuous Deployment pipelines with Docker. By using a HTTP API interface, CI/CD services can trigger the ‘recreation’ of containers with their prior settings, but using a freshly-obtained image. Access to containers is limited by using container-specific tokens.

## Usage

The only prerequisite to run the serive is Docker. To create the service:

* Create a `registries.json`, derived from `registries-example.json`. This allows you to use private registries with HTTP authentication.
* If desired, assign an exposed port in the `docker-compose.yml` configuration. (By default, port 80 is used).
* Launch the service via `docker-compose up -d`.

Once the service is running, you may add tokens to the Redis configuration — which is a workaround until we implemented a HTTP API. Tokens are expected to have the key `token_<TOKEN>` and a value of a stringified JSON object with the fields `containerId` (ID or name) and `imageTag` (`latest`, for instance).

Then, a recreation process is started by calling `/recreate?access_token=<TOKEN>`. The request is terminated once the recreation succeeded or failed.


## Development

Use Golang `v1.10.3` with `dep` `v0.4.1`.
