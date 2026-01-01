

# Distributed Systems

A Distributed System is a collection of independent computers (nodes) that work together as a single system to achieve a common goal. These computers communicate via a network (like the internet or LAN) and coordinate to perform tasks such as computation, storage, or serving requests.

## Key Characteristics of Distributed Systems

- **Scalability**
- **Reliability**
- **Availability**
- **Efficiency**
- **Serviceability or Manageability**

Let's briefly review them:

### Scalability

Any distributed system that can continuously evolve in order to support the growing amount of work is considered to be scalable. A scalable system would like to achieve this scaling without performance loss.

Generally, the performance of a system, although designed (or claimed) to be scalable, declines with the system size due to the management or environment cost. For instance, network speed may become slower because machines tend to be far apart from one another. More generally, some tasks may not be distributed, either because of their inherent atomic nature or because of some flaw in the system design. At some point, such tasks would limit the speed-up obtained by distribution. A scalable architecture avoids this situation and attempts to balance the load on all the participating nodes evenly.

**Horizontal vs. Vertical Scaling:**

- **Horizontal scaling**: Means that you scale by adding more servers into your pool of resources.
- **Vertical scaling**: Means that you scale by adding more power (CPU, RAM, Storage, etc.) to an existing server.

With horizontal-scaling it is often easier to scale dynamically by adding more machines into the existing pool; Vertical-scaling is usually limited to the capacity of a single server and scaling beyond that capacity often involves downtime and comes with an upper limit.

### Reliability

A distributed system is considered reliable if it keeps delivering its services even when one or several of its software or hardware components fail. Since in such systems any failing machine can always be replaced by another healthy one, ensuring the completion of the requested task.

Obviously, redundancy has a cost and a reliable system has to pay that to achieve such resilience for services by eliminating every single point of failure.

### Availability

By definition, availability is the time a system remains operational to perform its required function in a specific period. It is a simple measure of the percentage of time that a system, service, or a machine remains operational under normal conditions.

Availability takes into account maintainability, repair time, spares availability, and other logistics considerations.

**Reliability Vs. Availability**: If a system is reliable, it is available. However, if it is available, it is not necessarily reliable.

### Efficiency

Efficiency in a distributed system measures when multiple machines work together to perform a task (like fetching, processing, or storing data). Since operations happen across many nodes, efficiency depends on how quickly and how much work is done by the system as a whole.

**Two Key Measures of Efficiency:**

- **Response Time (Latency)**: The time delay between sending a request and receiving the first result.
- **Throughput (Bandwidth)**: The rate at which the system can deliver results (e.g., how many items per second).

### Serviceability or Manageability

Serviceability or manageability is the simplicity and speed with which a system can be repaired or maintained.






