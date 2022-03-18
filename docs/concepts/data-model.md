# Monoskope Data Model

While m8 has no classic CRUD database with tables for storing state, the transient state built based on that has.
Here is how the projected transient state is modeled:


```mermaid
erDiagram
    User {
        uuid id
        string name
        string email
    }

    UserRoleBinding {
        uuid id
        uuid user_id
        string role
        string scope
        uuid resource
    }

    Cluster {
        uuid id
        string name
        string display_name
        string api_server_address
        bytes ca_cert_bundle
    }

    Tenant {
        uuid id
        string name
        string display_name
        string prefix
    }

    TenantClusterBinding {
        uuid id
        uuid cluster_id
        uuid tenant_id
    }

    ClusterSecretStoreBinding {
        uuid id
        uuid cluster_id
        uuid secret_store_id
    }

    SecretStore {
        uuid id
        string name
        string type
        bytes config
    }

    Secret {
        uuid id
        uuid secret_upload_key_id
        uri upstream_id
        string status
        string hash
        bytes uploader_public_key
        bytes encrypted_payload
    }

    SecretUploadKey {
       uuid id
       uuid secret_store_id
       bytes public_key
       timestamp expiry
    }

    Secret ||--|| SecretUploadKey : references

    SecretUploadKey ||--|| SecretStore : references

    User ||--o{ UserRoleBinding : part_of
    Tenant ||--o{ UserRoleBinding : part_of

    Tenant ||--o{ TenantClusterBinding : part_of
    Cluster ||--o{ TenantClusterBinding : part_of

    Cluster ||--o{ ClusterSecretStoreBinding : part_of
    SecretStore ||--o{ ClusterSecretStoreBinding : part_of
```