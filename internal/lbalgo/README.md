# Load Balancing Algorithms

### Round-robin

Client requests are routed to available servers on a cyclical basis, which worked bests when servers have roughly
identical computing capabilities and storage capacity.

### Sticky Round-robin

Sticky sessions, is a load balancing technique where a load balancer routes a users' subsequent requests to the same
backend server they initially connected to.

This is achieved by storing session information in th form of cookies or other mechanisms. The idea behind sticky
sessions is to maintain user state and prevent issues with sessions being lost or disrupted when users interact with a
dynamic web application.

Sticky Round Robin direct specific clients to the same server, sacrificing load balancing fairness.

Pros:

- Seamless user experience
- Session persistence
- Improve performance by simplifying load balancer operations.
- Stateful application
- Less data consistency challenges
- Trouble shoot for a user would be easier ( All issues are isolated to a single server )

Cons:

- Uneven traffic across backend servers.
- Failed server causes interrupted users session ( Fail-over mechanism in need )

### Weighted Round-robin

Load balancing by assigning numeric weights to all the servers. The weights can be assigned based on factors such as
servers processing power or total bandwidth.

> Suppose we have three servers —ServerA, ServerB, ServerC— with weights (5, 2, 1) that are waiting to serve incoming
> requests behind the load balancer.
>
> The load balancer will forward the first five requests to ServerA, the next two requests to ServerB, and then one
> request to ServerC.
>
> If any of the other incoming requests arrive, the load balancer will forward those requests back to ServerA again for
> the next five incoming requests, then ServerB will get its turn, and after that the requests will be forwarded to
> ServerC. The cycle will continue on this way.

By using the Weighted Round Robin algorithm, network administrators can ensure a more balanced and efficient use of
resources, leading to improved performance and user experience.

Pros:

1. Ensuring all servers are used according to their capacity
2. Reducing the risk of server overload
3. Provides fault tolerance by redirecting requests in case of server failure

Cons:

1. How to accurately assign weight? Accurate assignments of weights can be a complex task.
2. Frequent changes of server capacity may lead to inaccurate weights. How to update weights correctly?

### Least connection

Least connection is a load balancing algorithm used in computer networking to distribute incoming network traffic
among multiple servers or resources. The basic principle behind the least connection algorithm is to direct new requests
to the server that is currently serving the fewest number of active connections. This helps to ensure that the load is
evenly distributed across all available servers, thus optimizing performance and preventing any single server from
becoming overwhelmed.

1. When a new request comes in, the load balancer examines the current number of active connections on each server.
2. The load balancer then selects the server with the fewest active connections.
3. The new request is forwarded to the selected server.
4. The load balancer keeps track of the number of connections to each server and updates this information as requests
   are
   processed and connections are closed.

Pros:

1. Traffic is distributed dynamically based on the current load on each server

### Least Response Time Load Balancing

The least-response-time load balancing strategy collects response times of the calls made with service instances and
picks an instance based on this information.

Erroneous responses are treated as responses with a long response time, by default 60 seconds. This can be controlled
with the error-penalty attribute.

The algorithm for service instance selection is as follows:

- if there is a service instance that wasn't used before - use it, otherwise:
- if there are any service instances with collected response times - select the one for which score is the lowest,
  otherwise:
- select a random instance

The score for an instance decreases in time if an instance is not used. This way we ensure that instances that haven't
been used in a long time, are retried.

