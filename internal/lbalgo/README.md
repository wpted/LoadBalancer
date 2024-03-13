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

### Source IP Hashing Load Balancing

Hashes are shorter and easier to use than the information that they are based on, while retaining enough information to
ensure that no two different pieces of information generate the same hash and are therefore confused with one another.

Most of the hashing algorithms calculate two hash values:

- A hash of the service’s IP address and port.
- A hash of the incoming URL, the domain name, the source IP address, the destination IP address, or the source and
  destination IP addresses, depending on the configured hash method.

Step:

1. Computes hash value
2. Select a service
3. Is the service up? No, select a new service.
4. Forward the request to the selected service.

### Power of Two Choices

> It All Falls Apart with Multiple Guides
> So far, we’ve had one guide (that is, load balancer) with a complete view of the queues and response time in the
> arrivals hall. That guide tries to make the best choice for each traveler based on the information he knows.
> Now consider what happens if we have several guides, each directing travelers independently. The guides have
> independent
> views of the queue lengths and queue wait times – they only consider the travelers that they send to each queue.
> This scenario is prone to an undesirable behavior, where all the guides notice that one queue is momentarily shorter
> and faster, and all send travelers to that queue. Simulations show that this “herd behavior” distributes travelers
> in a way
> that is unbalanced and unfair. In the same way, several independent load balancers can overload some upstream
> servers,
> no matter which “best choice” algorithm you use.


> Instead of making the absolute best choice using incomplete data, with “power of two choices” you pick two queues at
> random and chose the better option of the two, avoiding the worse choice. “Power of two choices” is efficient to
> implement. You don’t have to compare all queues to choose the best option each time; instead, you only need to compare
> two. And, perhaps unintuitively, it works better at scale than the best‑choice algorithms. It avoids the undesired herd
> behavior by the simple approach of avoiding the worst queue and distributing traffic with a degree of randomness.