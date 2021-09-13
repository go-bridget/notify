# Notify

Notify is a microservice which uses Redis and a JWT user_id claim to
provide state and broadcast delivery over a websocket interface.

A JWT user_id claim was chosen to avoid database lookups. Redis is used
to retrieve an initial state object on a websocket connection, as well as
subscribing to a per-user topic.

## How it works?

The session ID JWT claims must include a `user_id` claim (string). The
user ID is then mapped to individual redis keys:

- `notify:%s` - a publish/subscribe channel for broadcasts
- `notify:%s:state` - a hash value (e.g. `[unread_messages => 1, ...]`)

Broadcasting to the Redis channels can be done from your own applications
directly. Twirp RPC API endpoints are also provided for convenience.

## Building

Notify uses Docker and Drone CI. In order to build a notify binary for
amd64 you need to do two things:

- `cd docker/build && make` - this creates a build image with needed dependencies,
- `make` - this runs the full Drone CI suite and produces a binary under `/build`

## Testing

You can use the provided `docker-compose.yml` to run a test instance of
notify, as well as a keydb redis-compatible service. In order to make
your life a bit easier, a few shorthand commands are provided:

- `make up` - runs the compose services,
- `make down` - stops/removes the compose services,
- `make logs` - tails logs for the compose services,

The testing service is exposed on port `1234`, by default the
microservice is listening to port `:3000`, and also provides a `pprof`
listener on `:6060`. Both addresses can be overriden with environment
variables `PORT` and `PPROF_PORT` respectivelly.

# Miscellaneous

It's possible to extend the microservice to:

- Provide other kinds of authorization (e.g. Pareto, SQL, ...)
- Use different kinds of queues (e.g. Kafka, zeromq, ...)
- Use a different key/value store (e.g. Memcache, SQL, ...)

These are very specific needs that aren't currently implemented. If you'd 
like to sponsor development, [send me an e-mail](mailto:black@scene-si.org).

