# VerySimpleService

VSS is an extremely fast and simple service that helps distributed systems by centralizing some information.
Distributed systems are complicated, there are too many moving parts and sometimes a simple task becomes really complicated.

This is why we created VSS, because sometimes the problem can be easily resolved by removing all the distributed parts of it. But having a centralized solution, even for a small part of the problem usually doesn‚Äôt come for free, it just generates a single point of failure.

This is why VSS is created with a very Unix approach, let's do one thing and do it right!

There are not many things that can be done with VSS, but these things, like the classic Unix CLI tools are the building blocks for bigger things.

VSS is created with the following requirements:
- Is fast, all requests must have single digit latency.
- Is focused on throughput, it can handle tens of thousands of clients and support thousands of requests per second.
- It is resilient, a VSS cluster is designed to have zero downtime and high availability.
- Chaos is in its core, It is created and enforced to fail continuously, this way when there is a real failure you won‚Äôt notice. 

VSS does not store data.

A VSS cluster allows to maintain availability but does not enforce data persistence. VSS servers are used to keep track or propagate what is happening in the moment but should not be used to store information.

## Configuration 

VSS has 100 configurable slots that can be used to provide different functions.
Slots are configured through configuration files, if a slot configuration changes VSS cannot enforce consistency in the data until the new configuration is propagated.
Clients must know the configuration beforehand in order to use the slots appropriately.

For example, the same VSS server can be configured to have the first 3 slots as rate limiters and the next two as multicast signal propagation slots.
This way the applications can use a single server to solve more than one problem. I mean, is already there!

### Token bucket limiter

This limiter uses the classic token bucket approach to control the rate of events. Applications can request tokens from the limiter and the limiter will return the number of tokens assigned.
This can be used for example by a distributed fleet of API servers, allowing them to centralize the rate limit for the calls.

The token bucket approach adds a certain number of tokens per period (for example a second), to a bucket. Applications can take tokens from the bucket. After the tokens are depleted, it won't return any more tokens until the next period is reached.

For example, let's say we have an API that provides 100 requests per second, and we also want to allow the application to allow brief spikes in traffic up to 2x the maximum amount.
In this case we can configure our slot with the following:

|Config      |Value |
|------------|------|
|bucket_size |200   |
|period      |second|
|refresh_rate|100   |

The complete configuration options are:

|Config          | Description |
|----------------|-------------|
| bucket_size	 | Max amount of tokens that can be accumulated. |
| period	     | The refresh period for the tokens, it can be: second, minute or hour |
| refresh_rate	 | The number of tokens added on every refresh period. Default: 1 |
| tokens_per_req | This is the number of tokens that are assigned on every request. This is used to reduce the number of calls to the server, applications can have more tokens available to be used, when those are depleted it can ask for more. If the number is not available, the available number will be returned. Default: 1 |

### Leaky bucket limiter

This limiter works by defining an imaginary bucket that has a leak on it. The idea is that the leak is the rate how those tokens get delivered at a constant rate.

The bucket has a limited capacity, then it would remove tokens at a constant rate (defined by config). Every time we do a request, we put a token in the bucket, if there is enough room in the bucket, then the request will be approved, if there is not enough room the request will be denied. When a bucket is full, it will return 0 to all the requests until a token is leaked from it, then it will have room to receive a new token and so on.

This allows applications to have a burst of requests but after some time the bucket will fill up and the requests will start at a constant rate.

|Config          | Description |
|----------------|-------------|
| bucket_size	 | Max amount of tokens that can be accumulated. |
| period	     | The refresh period for the tokens, it can be: second, minute or hour |
| refresh_rate	 | The number of tokens leaked on every refresh period. Default: 1 |

As it can be interpreted from this description, this algorithm is more network intensive than the previous one because we need to constantly query the bucket to figure out if there is room when is full.

### Broadcast signal propagation

Anything sent to this slot is propagated as a message to all the other clients. Any client connected to VSS at this point will receive the event at least once.
This means that the message could be received more than once.

The message to be sent has a maximum of 36 characters, this allows to send an ID or a UUID to all the hosts.

This kind of slot is used to notify other clients about a new event or to propagate a signal.

The only configuration for this slot is the following:

|Config          | Description |
|----------------|-------------|
| timeout        | Time to wait for a confirmation on the clients that the message was received. |

This slot will only acknowledge the command when all the messages are sent, so take into account that the more clients connected or the hardest those clients are to reach, it will delay the confirmation.

### Multicast signal propagation

Similar to the Broadcast slot but this slot allows to send a message to a specific group of clients. This slot has two type of commands:
- Register/Deregister: This command allows a client to register or deregister from the multicast. If the client is registered, it will receive the events.
- Message: This will send a message the same way as the Broadcast slot but it will only send it to registered clients.

|Config          | Description |
|----------------|-------------|
| timeout        | Time to wait for a confirmation on the clients that the message was received. |
| dereg_tries    | Number of messages that are tried on a client until is de-registered. |

### Random signal propagation

This signal propagation slot works like the Multicast signal propagation explained before but with a major diference, the message is not sent to all registered clients, but only one. It uses a pseudo-random generator to distribute the messages among the clients.

It has the same configuration as the previous slot:

|Config          | Description |
|----------------|-------------|
| timeout        | Time to wait for a confirmation on the clients that the message was received. |
| dereg_tries    | Number of messages that are tried on a client until is de-registered. |

## Cluster configuration 

Usually VSS clusters have no more than 3 instances. The clients will connect simultaneously to all the instances, this way they can recover instantly if there is an issue.
The instances will automatically connect to each other in a complete mesh configuration each one of them will have N-1 connections.
After connecting to the other instances they will automatically elect a leader, this leader will be in charge to define the cluster configuration and routing but will not be the main data instance. 
It will also define which instance is the fallback leader and which one is the main data instance.
