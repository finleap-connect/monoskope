# `monoctl` authentication flow

```mermaid
sequenceDiagram
    participant U as User
    participant M as monoctl
    participant A as Ambassador
    participant G as Gateway
    participant I as Identity Provider
    U-->>+M: monoctl auth login
    M-->>+G: calls GetAuthInformation
    G-->>-M: returns AuthInformation
    M-->>M: opens browser
    M-->>+I: navigate browser to AuthCodeURL from AuthInformation
    I-->>I: determines user’s identity
    I-->>-M: redirect with auth_code
    M-->>+G: calls ExchangeAuthCode
    G-->>+I: exchange auth_code
    I-->>-G: returns id_token
    G-->>G: verifies id_token and claims
    G-->>G: issues access_token signed by M8
    G-->>-M: returns access_token, expiry
    M-->>-U: output login success
```