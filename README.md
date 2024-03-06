# Write your own load balancer

This is the solution
for [Building your own load balancer](https://codingchallenges.fyi/challenges/challenge-load-balancer) implemented using
Go.

## Goals

A load balancer performs the following functions:

- Distributes client requests/network load efficiently across multiple servers
- Ensures high availability and reliability by sending requests only to servers that are online
- Provides the flexibility to add or subtract servers as demand dictates

Therefore, our goals for this project are to:

- Build a load balancer that can send traffic to two or more servers with different algorithms.
- Health check the servers.
- Handle a server going offline (failing a health check).
- Handle a server coming back online (passing a health check).

## Distributed computing

- Web apps are deployed on servers with finite resources.
- System capacity: The maximum number of requests a server can serve.
- Horizontal scaling adds more server to a system.

## Load Balancers

A load balancer enables distribution of network traffic dynamically across resources (on-premises or cloud) to support
an application. It distributes incoming traffic to maximize the system's capacity and minimize the time it takes to
fulfill the request.

### Load balancing
1. The client sends a request to the load balancer.
2. The load balancer selects an appropriate backend server based on its load balancing algorithm.
3. The load balancer forwards the client's request to the selected backend server.
4. The backend server processes the request and generates a response.
5. The backend server sends the response back to the load balancer.
6. The load balancer receives the response from the backend server.
7. Finally, the load balancer forwards the response back to the client that initially made the request.

### Common used Load balancers.

- HAProxy
- Nginx

### Advantages

- Application Reliability
- Application Availability
- Application Scalability
- Application Security
- Application performance

### Disadvantages

- Can be a single point of failure.
- A user with multiple requests may be served with different backend servers: How to keep session in sync for backend
  servers?
- Deploying new server versions can take longer and require more machines. How to roll traffic to new servers and drain
  requests from the old machine is an issue.

## Steps

1. -[x] Create a simple HTTP server using Rust. It should have the ability to start up on a custom port.

2. Create a basic HTTP server that can start-up, listen for incoming connections and then forward them to a single
   server.
    - [x] Allow concurrent requests.
    - [x] Create service registration endpoint.
3. Distribute incoming requests between backend servers.
    - [x] Allow user to choose algorithm when starting up.
    - [x] Copy the received request.
    - [ ] Implement Load Balancing Algorithms
        - [x] Round Robin
        - [x] Sticky Round Robin
        - [x] Weighted Round Robin
        - [x] Least connections
        - [ ] Least time
        - [ ] URL hash
        - [ ] Source IP hash
        - [ ] Consistent hashing
        - [ ] Threshold
        - [ ] Random with two choices
4. Perform periodic health check.
    - [x] Allow a health check period to be specified on the command line during initializing.
        - [x] Health check url, GET request on backend server.
    - [x] Remove unhealthy backend servers from available servers.
    - [x] Move server that came alive back to the available servers.

## References

### Videos

- [Load Balancers for System Design Interviews](https://www.youtube.com/watch?v=chyZRNT7eEo)

### Reads
- [What Is a Load Balancer?](https://www.f5.com/glossary/load-balancer)
- [Canceling blocking read from stdin](https://www.reddit.com/r/golang/comments/fsxkqr/cancelling_blocking_read_from_stdin/)
- [Round-robin load balancing](https://avinetworks.com/glossary/round-robin-load-balancing/)
- [How sticky sessions can tilt load balancers](https://medium.com/@iSooraj/how-sticky-sessions-can-tilt-load-balancers-c5dc8f50099c)
- [Kafka, Range Round-Robin 和 Sticky三種分區分配策略](https://blog.csdn.net/u010022158/article/details/106271208)
- [What is the weighted round-robin technique](https://www.educative.io/answers/what-is-the-weighted-round-robin-load-balancing-technique)
- [What is weighted round-robin](https://webhostinggeeks.com/blog/what-is-weighted-round-robin/)
- [Least connection method](https://docs.netscaler.com/en-us/citrix-adc/current-release/load-balancing/load-balancing-customizing-algorithms/leastconnection-method.html)
- [負載均衡策略之最少連接](https://mozillazg.com/2019/02/load-balancing-strategy-algorithm-weighted-least-connection.html#hidleast-connection)
- [Nginx HTTP upstream least connection module](https://github.com/nginx/nginx/blob/d8ccef021588cf79d2dae7c132a0b1225ed52c30/src/http/modules/ngx_http_upstream_least_conn_module.c)
- [Least response time load balancing](http://smallrye.io/smallrye-stork/1.1.1/load-balancer/response-time/#)
- [What is the least response time load balancing technique](https://www.educative.io/answers/what-is-the-least-response-time-load-balancing-technique)