# HAYVN - assesment task

## Batch API handler

I enjoyed coding up this task, as it seemed to be more of a real world example, compared to some of the assignments I've received in the past. I appreciated the ability to write the task in any language or framework of my choice, and because of this, I chose GO.

Golang has quickly become my favourite backend language. The built in concurrency, scalability and robust standard library make it a pleasure to code with. My previous position was that of a blockchain developer, and majority of blockchain nodes are run with [Geth](https://geth.ethereum.org/) making golang the obvious choice for interfacing with the blockchain. 

I found this task challenging at the initial look. I wasn't quite sure what the goal was at first, but after breaking the problem down, I realised that I was building an API batcher for routing messages in what could easily be a micro service environment. 

I chose to implement _go channels_ and _workers_ to effectively make this completely asyncronous and allow for complete scalability. I also elected to make all configurations editable via environment variables.

The interfaces were chosen to allow for _hot-swappable_ modules. For example, the current implementation pushes the resultant aggregated messages to a REST endpoint, but this could easily be traded for a GRPC implementation or PUB/SUB. 

I chose to use the go-chi router as a base plate for the API layer, and included the swagger library to create "API docs". This allows for easy testing and integrations.


## Variables:
	HOST_NAME           The host on which to serve the API                          default: http://localhost
    HOST_PORT           The port on which to serve the API                          default: 8081
	RESPONSE_URL        The Base URL on which to post the aggregated messages       default: http://locahost
	RESPONSE_PORT       The Port on which to post the aggregated messages           default: 8082
	RATE_LIMIT_SECONDS  The time in seconds to batch the messages received      default: 10

## Steps for running:

Build the docker image:

```docker build -t hayvn-test .```

Run the container, specifying env variables:

```
docker run -e HOST_NAME='http://localhost' -e RESPONSE_URL='http://localhost' -e HOST_PORT=8081 -e RESPONSE_PORT=8082 -e RATE_LIMIT_SECONDS=10 -p 8081:8081 hayvn-test
```

The docs will be available at http://localhost:8081/docs where you can send a health-check REST call, as well as post messages to the service. The service will log out the resultant aggregated messages POST body.


## Example:

Messages sent:
```
{
    "destination": "operations-channel",
    "text": "An important event has occurred",
    "timestamp": "2022-05-19T10:55:06.479Z"
}

{
    "destination": "operations-channel",
    "text": "A second important event has occurred",
    "timestamp": "2022-05-19T10:55:06.479Z"
}

{
    "destination": "compliance",
    "text": "An important event has occurred",
    "timestamp": "2022-05-19T10:55:06.479Z"
}
```

Resulting Message:
```
{
    "batches": [
        {
            "destination": "compliance",
            "messages": [
                {
                    "text": "An important event has occurred",
                    "timestamp": "2022-05-19T10:55:06.479Z"
                }
            ]
        },
        {
            "destination": "operations-channel",
            "messages": [
                {
                    "text": "An important event has occurred",
                    "timestamp": "2022-05-19T10:55:06.479Z"
                },
                {
                    "text": "A second important event has occurred",
                    "timestamp": "2022-05-19T10:55:06.479Z"
                }
            ]
        }
    ]
}
```