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

## Protocol

VSS uses slots to communicate, if you ever worked with microcontrollers you would get the similarities with registers.
The idea is that you can either write or read a slot. Slots cannot have more than 36 characters of data and go from slot #0 to slot #99.

All messages are plain text in order to simplify the protocol among different programming languages.

In order to read a slot you can send a read request. That would be the command ‘r’, then three digits defining the slot number .

`r000`

This will trigger a value response with the information about the slot. The response ‘v’ indicates is a value response, then three digits to determine the slot and up to 36 characters to define the value.

`v0006396A64C-1C2C-4BFC-B8F1-034758018CAC`

In this example, the slot has a UUID stored.

When a client wants to write a value on a slot, they can use the `w` command:

`w0006396A64C-1C2C-4BFC-B8F1-034758018CAC`

Same as the read command, the server will return the written value:

`v0006396A64C-1C2C-4BFC-B8F1-034758018CAC`

If there is any issue with a command, the server will return a failure with a code that can be used to identify the issue:

`f0006396A64C-1C2C-4BFC-B8F1-034758018CAC`

In some cases there are messages sent as event from the server (see broadcast slots), these kind of messages are sent at any time and use the `e` (event) response:

`e2346396A64C-1C2C-4BFC-B8F1-034758018CAC`

Same as the other examples, it would contain the `e` response, then the slot (in this case 234) and the event data (in this case a UUID).

## Configuration 

VSS has 100 configurable slots that can be used to provide different functions.
Slots are configured through configuration files, if a slot configuration changes VSS cannot enforce consistency in the data until the new configuration is propagated.
Clients must know the configuration beforehand in order to use the slots appropriately.

For example, the same VSS server can be configured to have the first 3 slots as rate limiters and the next two as multicast signal propagation slots.
This way the applications can use a single server to solve more than one problem. I mean, is already there!

### Simple memory slot

This is the most basic slot where a value can be stored. The value has a maximum of 36 characters. You can read and write on the value and there are no restrictions.
This slot has also no configuration.

### Timeout memory slot

This slot is also a memory slot but the main difference with the Simple memory slot is that it has an owner. Only the client that has last written in this slot can write again.
If the owner does not write on this slot for a certain time (timeout), it will lose the ownership and any other client can take over.

The timeout can be configured:

|Config      |Value                               |
|------------|------------------------------------|
|timeout     |Timeout value configured in seconds.|

All clients can read from this slot, but only the owner can write. If any other client tries to write it will fail. If there is no owner, the first client that writes becomes the owner.

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

Writes have no effect on this slot. Reads will return the number of tokens (or zero if there are no tokens available).

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

Writes have no effect on this slot. Reads will return 1 if the token was accepted or zero if not.

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

Writes will propagate the written value to all other clients. Reads will read the last written value.

### Multicast signal propagation

Similar to the Broadcast slot but this slot allows to send a message to a specific group of clients. This type of multicast **requires 2 consecutive slots:**
- Register/Deregister: This slot allows a client to register or deregister from the multicast. If the client is registered, it will receive the events. To register a client can write a value on this slot, to deregister it can write zero. If a client reads this slot, then a non-zero value means the client is already registered and a zero value means is not.
- Message: This will send a message the same way as the Broadcast slot but it will only send it to registered clients. Writes will propagate the written value to all other clients. Reads will read the last written value.

|Config          | Description |
|----------------|-------------|
| timeout        | Time to wait for a confirmation on the clients that the message was received. |
| dereg_tries    | Number of messages that are tried on a client until is de-registered. |


### Random signal propagation

This signal propagation slot works like the Multicast signal propagation explained before but with a major diference, the message is not sent to all registered clients, but only one. It uses a pseudo-random generator to distribute the messages among the clients.

It also **requires 2 slots**:
- Register/Deregister: This slot allows a client to register or deregister from the multicast. If the client is registered, it will receive the events. To register a client can write a value on this slot, to deregister it can write zero. If a client reads this slot, then a non-zero value means the client is already registered and a zero value means is not.
- Message: This will send a message the same way as the Broadcast slot but it will only send it to registered clients. Writes will propagate the written value to all other clients. Reads will read the last written value.

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
