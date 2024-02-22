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

## References

- [What Is a Load Balancer?](https://www.f5.com/glossary/load-balancer)