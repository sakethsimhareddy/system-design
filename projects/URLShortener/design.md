# URL Shortener System Design

## Requirements

### Functional Requirements
1. Given a URL, our service should generate a shorter and unique alias of it. This is called a short link.
2. When users access a short link, our service should redirect them to the original link.
3. Provide simple analytics for URLs.

### Non-Functional Requirements
1. The system should be highly available. This is required because, if our service is down, all the URL redirections will start failing.
2. URL redirection should happen in real-time with minimal latency.
3. Shortened links should not be guessable (not predictable).

## Capacity Estimation and Constraints

Given constraint is the system should be able to handle 200 URLs/s.

Since this is a read-heavy system, let's consider a 1:10 write-to-read ratio (new URLs vs. redirections).

- New URL registrations: 1/10 * 200 = 20 URLs/s
- URL redirections: 9/10 * 200 = 180 URLs/s

Total URLs in DB per month: 20 * 60 * 60 * 24 * 30 = 50 million approximately.

Assuming each record is 300 bytes, total storage needed: 300 bytes * 50M = 15 GB.

Assuming 20% of URLs account for 80% of traffic, then 20/100 * 180 * 60 * 60 * 24 * 30 * 300 bytes ≈ 9 GB of memory needed for cache.

### Summary
- New URLs: 20/s
- URL redirections: 180/s
- Storage for 1 month: 15 GB
- Memory for cache: 9 GB

## System APIs

We can have REST APIs to expose the functionality of our service. Following could be the definitions of the APIs for creating and accessing URLs:

### 1. Create URL
```
POST /api/v1/create
```
**Parameters:**
- `api_dev_key` (string): The API developer key of a registered account. This will be used to, among other things, throttle users based on their allocated quota.
- `original_url` (string): Original URL to be shortened.

**Returns:** (string)
A successful insertion returns the shortened URL; otherwise, it returns an error code.

### 2. Redirect URL
```
GET /api/v1/{tiny_url}
```
**Parameters:**
- `tiny_url` (string): The short URL which users call for redirection to the long URL.

**Returns:** Redirect to the original URL.

### Abuse Prevention
How do we detect and prevent abuse? A malicious user can put us out of business by consuming all URL keys in the current design. To prevent abuse, we can limit users via their `api_dev_key`. Each `api_dev_key` can be limited to a certain number of URL creations per some time period (which may be set to a different duration per developer key).

## Database Design

A few observations about the nature of the data we will store:
1. We need to store billions of records.
2. Each object we store is small (less than 1K).
3. There are no relationships between records—other than storing which user created a URL.
4. Our service is read-heavy.

### Database Schema
We would need two tables: one for storing information about the URL mappings, and one for the user's data who created the short link.

#### URL Table
- `original_url` (string)
- `tiny_url` (string)
- `created_by` (int)
- `request_count` (int)

#### User Table
- `user_id` (int)
- `username` (string)
- `password_hashed` (string)

Since there are no heavy relationships, I am thinking of using NoSQL for horizontal scaling. Since we're using NoSQL like MongoDB, we can use an index on `short_url` for fast retrieval.

## High-Level Design (HLD)

```
Client -> DNS
    ↓
Load Balancer ----------------------------------> Write API
    ↓
Read API
    ↓
Analytics Service (async call) <- API Servers (stateless)
    |
    ↓
Load Balancer
    ↓
Cache (Redis)
    ↓         ↑
    ↓         |
Database (Master) <------------------------------ Update
```

## Redirect Flow

When a user requests a short URL:

1. API checks cache first
2. If found → redirect immediately
3. If not → fetch from DB, update cache
4. Before returning, make an async call to analytics service (for low latency and high throughput)

## Create URL Flow

When a user posts a long URL:

1. Check the DB if URL already exists under this user ID
2. If it exists, return duplicate
3. If not, create short URL
4. Update the DB
5. Update the cache

## Basic System Design and Algorithm

### Encoding Actual URL
We can compute a unique hash (e.g., MD5 or SHA256, etc.) of the given URL. The hash can then be encoded for displaying. This encoding could be base36 ([a-z, 0-9]) or base62 ([A-Z, a-z, 0-9]) and if we add '-' and '.' we can use base64 encoding. We can use encode username + long URL.

## Failure Handling & Scaling

### Failures
1. **Cache down**: Fall back to DB
2. **DB down**: Use DB replication
3. **High traffic**: Horizontal scaling by adding new instances via load balancer
