**Project Overview:**
The goal of this project is to create a password manager for developers, enabling users to manage and share sensitive data like passwords, SSH keys, and API keys securely. The system will support user authentication, end-to-end encryption, sharing features, groups, and have a web interface, TUI, and CLI.

### Core Features:
1. **Authentication**: JWT tokens with email/password-based authentication.
2. **Sharing**: Users can share passwords, SSH keys, and API keys with other users or groups.
3. **Groups**: Group management for sharing secrets.
4. **End-to-End Encryption**: Secrets are encrypted on the client-side before being saved to the database.
5. **Accessibility**: Web interface, TUI (Text User Interface), and CLI (for scripting).

### Database Design

#### Entities:
1. **User**: Stores information about the users.
2. **Secret**: Stores secrets (passwords, SSH keys, API keys).
3. **Group**: Represents groups that users can create for sharing secrets.
4. **UserGroup**: Junction table that links users and groups (many-to-many relationship).
5. **SharedSecret**: Represents secrets shared between users or groups.
6. **AuditLog**: Tracks actions performed by users like creating, updating, or sharing secrets.

#### Attributes:
1. **User**:
   - `user_id`: Primary Key (PK)
   - `email`: Email address (unique)
   - `password_hash`: Hashed password
   - `created_at`: Account creation timestamp
   - `last_login`: Last login timestamp

2. **Secret**:
   - `secret_id`: Primary Key (PK)
   - `user_id`: Foreign Key (FK) referencing `User`
   - `secret_type`: Enum ('password', 'ssh_key', 'api_key')
   - `encrypted_data`: Encrypted secret (end-to-end encryption applied)
   - `description`: User-defined description of the secret
   - `created_at`: Timestamp of when the secret was created
   - `updated_at`: Timestamp of the last update

3. **Group**:
   - `group_id`: Primary Key (PK)
   - `group_name`: Name of the group
   - `created_by`: FK referencing `User` (who created the group)
   - `created_at`: Timestamp when the group was created

4. **UserGroup**:
   - `user_group_id`: Primary Key (PK)
   - `user_id`: FK referencing `User`
   - `group_id`: FK referencing `Group`
   - `role`: Enum ('admin', 'member') (Permissions in the group)

5. **SharedSecret**:
   - `shared_secret_id`: Primary Key (PK)
   - `secret_id`: FK referencing `Secret`
   - `shared_with_user`: FK referencing `User` (optional, for sharing with individual users)
   - `shared_with_group`: FK referencing `Group` (optional, for sharing with groups)
   - `permissions`: Enum ('read', 'write') (Permissions for the shared secret)
   - `shared_at`: Timestamp of sharing

6. **AuditLog**:
   - `log_id`: Primary Key (PK)
   - `user_id`: FK referencing `User`
   - `action`: Action performed (e.g., 'create_secret', 'update_secret', 'share_secret')
   - `details`: Additional information about the action
   - `timestamp`: Timestamp when the action occurred


### Table Definitions

1. **User Table**:

```sql
CREATE TABLE User (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);
```

2. **Secret Table**:

```sql
CREATE TABLE Secret (
    secret_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES User(user_id) ON DELETE CASCADE,
    secret_type VARCHAR(50) CHECK (secret_type IN ('password', 'ssh_key', 'api_key')),
    encrypted_data TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

3. **Group Table**:

```sql
CREATE TABLE Group (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    created_by INT REFERENCES User(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

4. **UserGroup Table**:

```sql
CREATE TABLE UserGroup (
    user_group_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES User(user_id) ON DELETE CASCADE,
    group_id INT REFERENCES Group(group_id) ON DELETE CASCADE,
    role VARCHAR(50) CHECK (role IN ('admin', 'member'))
);
```

5. **SharedSecret Table**:

```sql
CREATE TABLE SharedSecret (
    shared_secret_id SERIAL PRIMARY KEY,
    secret_id INT REFERENCES Secret(secret_id) ON DELETE CASCADE,
    shared_with_user INT REFERENCES User(user_id) ON DELETE SET NULL,
    shared_with_group INT REFERENCES Group(group_id) ON DELETE SET NULL,
    permissions VARCHAR(50) CHECK (permissions IN ('read', 'write')),
    shared_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

6. **AuditLog Table**:

```sql
CREATE TABLE AuditLog (
    log_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES User(user_id),
    action VARCHAR(255) NOT NULL,
    details TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Relationships:
- **User ↔ Secret**: A one-to-many relationship, where a user can own multiple secrets.
- **User ↔ Group**: A many-to-many relationship via the `UserGroup` junction table.
- **Group ↔ SharedSecret**: A group can be associated with multiple shared secrets.
- **User ↔ SharedSecret**: A user can receive secrets shared with them.
- **AuditLog ↔ User**: A user can perform multiple actions, which are tracked in the audit log.

### Key Functional Requirements Mapped to the Schema:
1. **Authentication**: The `User` table stores the `email` and `password_hash`. JWT tokens are managed externally during authentication.
2. **End-to-End Encryption**: Secrets are encrypted on the client-side and stored as `encrypted_data` in the `Secret` table.
3. **Sharing**: The `SharedSecret` table links secrets with users or groups, with optional read/write permissions.
4. **Group Management**: The `Group` and `UserGroup` tables handle groups and their membership roles (admin/member).
5. **Audit Logs**: User actions like secret creation, updates, and sharing are tracked in the `AuditLog` table.

### Conclusion:
This schema supports the secure management, sharing, and access control of sensitive developer data with proper relationships to enable features like authentication, encryption, and auditing.


