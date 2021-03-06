# `monoctl` authentication refresh flow

```mermaid
sequenceDiagram
    participant U as User
    participant M as monoctl
    participant A as Ambassador
    participant G as Gateway
    participant C as CommandHandler
    U-->>+M: monoctl tenant create
    M-->>+A: calls API
    A-->>+G: authenticate request
    G-->>G: validate token
    G-->>-A: returns unauthorized
    A-->>-M: returns unauthorized
    Note right of M: Normal auth flow<br> executed by<br>monoctl.<br>See Login.
    M-->>+A: calls API
    A-->>+G: authenticate request
    G-->>G: validate token
    G-->>-A: returns authorized
    A-->>+C: calls API
    C-->>-A: returns command response
    A-->>-M: returns command response
    M-->>-U: Output command response
```