# Latency vs Throughput in System Design

## Overview

Latency and throughput are fundamental concepts in system design. Latency measures how long a single operation takes, while throughput measures how many operations can be completed per unit of time. These metrics often trade off against each other, and understanding them is crucial for building efficient systems.

## Definitions

- **Latency**: The time required for one request or data packet to travel from start to finish. For example, an 80 ms response time for a web page.
- **Throughput**: The amount of work done per unit time. For example, 1,000 requests per second or 100 Mbps of data transfer.

Both are key performance metrics, but latency focuses on the speed of individual operations, while throughput emphasizes overall capacity.

## Why These Concepts Exist

Engineers distinguish between latency and throughput because "system speed" can be ambiguous:

- A system might respond quickly to a single user (low latency) but fail under heavy load (low throughput).
- Conversely, a system could handle large volumes efficiently (high throughput) but with longer wait times per request (high latency).

Separating these concepts allows for:

- Designing for optimal user experience (prioritizing latency).
- Scaling to support many users or large datasets (prioritizing throughput).

## Problems Solved

Without clear metrics for latency and throughput, optimization efforts can be misguided:

- Focusing solely on throughput might improve batch processing but degrade interactive application responsiveness.
- Prioritizing latency alone could make individual requests fast but limit overall user capacity.

By measuring both, you can:

- Identify bottlenecks: High latency with available bandwidth indicates delays, not capacity issues.
- Select appropriate algorithms: Batch processing for throughput vs. immediate processing for latency.

## Where It Matters

### Web Applications
- Users expect sub-100ms response times (latency).
- Companies need servers to handle thousands of requests per second (throughput).

### Real-Time Systems (Games, Trading, VoIP)
- Extremely sensitive to latency; even minor delays are problematic.
- Often trade throughput for instant responses.

### Data-Intensive Tasks (Streaming, Backups, Big Data)
- Prioritize throughput for moving large volumes of data.
- Some initial delay is acceptable if overall transfer is efficient.

### System Design Interviews/Backend Engineering
- Trade-offs like batching database writes or data compression directly affect latency vs. throughput.

## Trade-Offs

Increasing one metric often decreases the other:

- **Batching**: Improves throughput by processing multiple requests together but increases latency due to waiting.
- **Parallelism**: Handling fewer concurrent requests reduces latency by minimizing queues but lowers total throughput.

The key design question is: "Do we optimize for low latency (interactive/real-time) or high throughput (bulk/batch)?"

## Core Design Principles

### User-Facing Paths: Optimize for Latency
When users are waiting, focus on speed:

**Must-Do:**
- Reduce network hops
- Avoid heavy computations
- Minimize slow database calls

**Techniques:**
- Caching
- Pre-computation
- Read replicas
- Simple data models
- Keep synchronous paths short

### Background/Internal Paths: Optimize for Throughput
For handling many users or large-scale operations:

**Must-Do:**
- Process requests in parallel
- Decouple systems
- Accept eventual consistency

**Techniques:**
- Queues
- Asynchronous workers
- Batching
- Horizontal scaling