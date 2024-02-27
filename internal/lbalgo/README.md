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