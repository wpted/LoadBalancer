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
an application. It distributes incoming traffic to maximize the system's capacity and minimize the time it takes to fulfill the request.

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
- A user with multiple requests may be served with different backend servers: How to keep session in sync for backend servers?
- Deploying new server versions can take longer and require more machines. How to roll traffic to new servers and drain requests from the old machine is an issue.

## Steps

1. -[x] Create a simple HTTP server using Rust. It should have the ability to start up on a custom port.

2. Create a basic HTTP server that can start-up, listen for incoming connections and then forward them to a single
   server.
    - [x] Allow concurrent requests.
    - [x] Create service registration endpoint.
    - [ ] Forward should match the incoming HTTP request methods and paths.
3. Distribute incoming requests between backend servers.
    - [ ] Allow user to choose algorithm when starting up.
    - [ ] Implement Load Balancing Algorithms
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
4. Perform periodic health check.
    - [ ] Allow a health check period to be specified on the command line.
      - [x] Health check url, GET request on backend server.
    - [x] Remove unhealthy backend servers from available servers.
    - [x] Move server that came alive back to the available servers.

## Unsolved problems
1. A repl on another routine is fine, but it requires input for the for loop to proceed, which blocks the receiving from l.ReplDone(). I'm disabling the repl for now.
![blockedForLoop.png](static%2FblockedForLoop.png)

## References

### Videos
- [Load Balancers for System Design Interviews](https://www.youtube.com/watch?v=chyZRNT7eEo)

### Reads
- [What Is a Load Balancer?](https://www.f5.com/glossary/load-balancer)
- [Canceling blocking read from stdin](https://www.reddit.com/r/golang/comments/fsxkqr/cancelling_blocking_read_from_stdin/)