# leader-election

This project simulate a basic leader election process in a distributed system. Multiple server instances run concurrently, each trying to claim leadership by setting a key in a shared Redis database. The system ensures that only one server is recognized as the leader at any given time.

# Setup

Build and run the containers:

```sh
docker-compose up --build
```

Start the services using Docker Compose, this command will start Redis and three server instances.

```sh
docker-compose up
```

To stop the services, press <kbd>⌃ Control</kbd> + <kbd>C</kbd>.

# API Endpoints

Each server exposes a single HTTP endpoint that returns the current leader's ID.

Example usage:

```sh
curl http://localhost:8001/
```

Response:

```json
{ "leader": "service8001" }
```

# Data flow

```mermaid
sequenceDiagram
    participant S1 as Server 1
    participant S2 as Server 2
    participant S3 as Server 3
    participant R as Redis

    loop Every 1-10 seconds
        S1->>R: Check for leader
        S2->>R: Check for leader
        S3->>R: Check for leader

        alt No leader exists
            S1->>R: Attempt to set leader (10s expiry)
            R-->>S1: Confirm leader set
            S2->>R: Check leader (S1 is leader)
            S3->>R: Check leader (S1 is leader)
        else Leader exists
            R-->>S1: Return current leader
            R-->>S2: Return current leader
            R-->>S3: Return current leader
        end
    end

    note over S1,R: Leader key expires after 10s

    loop Next cycle
        S2->>R: Check for leader (expired)
        S2->>R: Attempt to set leader (10s expiry)
        R-->>S2: Confirm leader set
        S1->>R: Check leader (S2 is leader)
        S3->>R: Check leader (S2 is leader)
    end
```
