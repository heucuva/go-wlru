# go-wlru
Thread-safe LRU cache with permanency and context-based expiration

## Operational Complexity (Time)

| Operation | Best        | Average     | Worst       |
|-----------|-------------|-------------|-------------|
| Access    | Θ(1)        | Θ(1)        | O(1)        |
| Search    | Θ(1)        | Θ(1)        | O(n)        |
| Insertion | Θ(1)        | Θ(1)        | O(n)        |
| Deletion  | Θ(1)        | Θ(1)        | O(n)        |
| Snapshot  | Θ(n)        | Θ(n)        | Θ(n)        |

## Operation Complexity (Space)

| Complexity | Value           |
|------------|-----------------|
| Best       | Ω(2n)           |
| Average    | Ω(2n)           |
| Worst      | Ω(n + n log(n)) |

## Usage
This is a simple example LRU cache structure made with API request lookup caching in mind. If you decide to use this, do so at your own peril.

## Thread Safety
It should be thread-safe on all operations.
