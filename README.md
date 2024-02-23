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

## Load Balancers

A load balancer enables distribution of network traffic dynamically across resources (on-premises or cloud) to support
an application.

### Advantages

- Application Availability
- Application Scalability
- Application Security
- Application performance

### Load Balancing Algorithms

- [ ] Round Robin
- [ ] Sticky Round Robin
- [ ] Weighted Round Robin
- [ ] Least connections
- [ ] Least time
- [ ] URL hash
- [ ] Source IP hash
- [ ] Consistent hashing
- [ ] Threshold
- [ ] Random with two choices

## Steps

1. -[x] Create a simple HTTP server using Rust. It should have the ability to start up on a custom port.

2. Create a basic HTTP server that can start-up, listen for incoming connections and then forward them to a single
   server.
    - [x] Allow concurrent requests.
    - [ ] Create service registration endpoint.
    - [ ] Forward should match the incoming HTTP request methods and paths.
3. Distribute incoming requests between backend servers.
    - [ ] Allow user to choose algorithm when starting up.
4. Perform periodic health check.
    - [ ] Allow a health check period to be specified on the command line.
      - [ ] Health check url, GET request on backend server.
    - [ ] Remove unhealthy backend servers from available servers.
    - [ ] Move server that came alive back to the available servers.

## References

- [What Is a Load Balancer?](https://www.f5.com/glossary/load-balancer)