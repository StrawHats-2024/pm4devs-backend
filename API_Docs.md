## API Design 

### **Base URL**: 
`https://api.passwordmanager.dev/v1`

---

### **Authentication & Authorization**

#### 1. **User Registration**
- **Endpoint**: `/auth/register`
- **Method**: `POST`
- **Description**: Registers a new user in the system.
- **Request**:
  ```json
  {
    "email": "user@example.com",
    "username": "parikshith",
    "password": "user_password"
  }
  ```
- **Response**:
  - `201 Created`: User registered successfully.
  ```json
  {
      "token": "jwt_token",
    "message": "User registered successfully",
      "user_id": 1
  }
  ```
  - `400 Bad Request`: Invalid email or password format.

#### 2. **User Login**
- **Endpoint**: `/auth/login`
- **Method**: `POST`
- **Description**: Logs in the user and returns a JWT token.
- **Request**:
  ```json
  {
    "email": "user@example.com",
    "password": "user_password"
  }
  ```
- **Response**:
  - `200 OK`: User authenticated successfully.
  ```json
  {
    "token": "jwt_token",
    "user_id": "1234"
  }
  ```
  - `401 Unauthorized`: Invalid credentials. 
  - `400 Bad Request`: Invalid email or password format.

#### 3 **Token Verification**
- **Endpoint**: `/auth/verify-token`
- **Method**: `POST`
- **Description**: Verifies a given JWT token and returns its validity along with associated user information.

##### **Request**:
- **Body**: 
  ```json
  {
    "token": "jwt_token"
  }
  ```
- **Fields**:
  - `token` (string, required): The JWT token to be verified.

##### **Response**:
- **200 OK**: Token is valid.
  ```json
  {
    "valid": true,
    "user_id": 1234,
    "message": "Token is valid"
  }
  ```
- **401 Unauthorized**: Token is invalid or expired.
  ```json
  {
    "error": "Invalid token"
  }
  ```
- **400 Bad Request**: Missing or improperly formatted token.
  ```json
  {
    "error": "Bad request"
  }
  ```

##### **Errors**:
- `400 Bad Request`: The request body is missing or the token is not in the correct format.
- `401 Unauthorized`: The token is either invalid or has expired.


#### 3. **Refresh JWT Token**
- **Endpoint**: `/auth/refresh`
- **Method**: `POST`
- **Description**: Refreshes the JWT token.
- **Request**:
  ```json
  {
    "refresh_token": "old_refresh_token"
  }
  ```
- **Response**:
  - `200 OK`: New JWT token.
  ```json
  {
    "token": "new_jwt_token"
  }
  ```

---

### **Secrets Management**

#### 4. **Create Secret**
- **Endpoint**: `/secrets`
- **Method**: `POST`
- **Description**: Create a new secret (password, SSH key, or API key).
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "user_id": 1,
    "secret_type": "password", // or 'ssh_key', 'api_key'
    "encrypted_data": "encrypted_secret_here",
    "description": "My database password"
  }
  ```
- **Response**:
  - `201 Created`: Secret created successfully.
  ```json
  {
    "secret_id": "5678",
    "message": "Secret created successfully"
  }
  ```
  - `400 Bad Request`: Invalid data or missing fields.

#### 5. **Get All Secrets**
- **Endpoint**: `/secrets`
- **Method**: `GET`
- **Description**: Retrieve all secrets for the authenticated user.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: List of secrets.
  ```json
  [
    {
      "secret_id": "5678",
      "secret_type": "password",
      "description": "My database password",
      "created_at": "2024-09-07T12:00:00Z"
    },
    {
      "secret_id": "5679",
      "secret_type": "api_key",
      "description": "GitHub API Key",
      "created_at": "2024-09-07T12:30:00Z"
    }
  ]
  ```

#### 6. **Get Single Secret**
- **Endpoint**: `/secrets/{secret_id}`
- **Method**: `GET`
- **Description**: Retrieve a specific secret by its ID.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: Secret retrieved successfully.
  ```json
  {
    "secret_id": "5678",
    "secret_type": "password",
    "encrypted_data": "encrypted_secret_here",
    "description": "My database password"
  }
  ```
  - `404 Not Found`: Secret not found or user does not have access.

#### 7. **Update Secret**
- **Endpoint**: `/secrets/{secret_id}`
- **Method**: `PUT`
- **Description**: Update an existing secret's details.
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "encrypted_data": "new_encrypted_secret_here",
    "description": "Updated database password"
  }
  ```
- **Response**:
  - `200 OK`: Secret updated successfully.
  ```json
  {
    "message": "Secret updated successfully"
  }
  ```

#### 8. **Delete Secret**
- **Endpoint**: `/secrets/{secret_id}`
- **Method**: `DELETE`
- **Description**: Delete an existing secret.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: Secret deleted successfully.
  ```json
  {
    "message": "Secret deleted successfully"
  }
  ```

---

### **Sharing Secrets**

#### 9. **Share Secret with User**
- **Endpoint**: `/secrets/{secret_id}/share`
- **Method**: `POST`
- **Description**: Share a secret with another user.
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "user_email": "otheruser@example.com",
    "permissions": "read" // or 'write'
  }
  ```
- **Response**:
  - `200 OK`: Secret shared successfully.
  ```json
  {
    "message": "Secret shared successfully"
  }
  ```
  - `404 Not Found`: User not found.
  - `403 Forbidden`: You don't have permission

#### 10. **Share Secret with Group**
- **Endpoint**: `/secrets/{secret_id}/share/group`
- **Method**: `POST`
- **Description**: Share a secret with a group.
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "group_id": "123",
    "permissions": "read" // or 'write'
  }
  ```
- **Response**:
  - `200 OK`: Secret shared with group successfully.
  ```json
  {
    "message": "Secret shared with group successfully"
  }
  ```
  - `404 Not Found`: Group not found or user not part of the group.

#### 11. **Get Shared Secrets**
- **Endpoint**: `/secrets/shared`
- **Method**: `GET`
- **Description**: Retrieve all secrets shared with the user.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: List of shared secrets.
  ```json
  [
    {
      "secret_id": "5678",
      "shared_by": "user@example.com",
      "permissions": "read",
      "description": "Shared secret description"
    }
  ]
  ```

  #### 11. **Revoke Access from User**
- **Endpoint**: `/secrets/{secret_id}/revoke`
- **Method**: `POST`
- **Description**: Revoke access to a secret from a specified user.
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "user_email": "otheruser@example.com"
  }
  ```
- **Response**:
  - `200 OK`: Access revoked successfully.
  ```json
  {
    "message": "Access revoked successfully."
  }
  ```
  - `404 Not Found`: User not found or does not have access to the secret.

#### 12. **Revoke Access from Group**
- **Endpoint**: `/secrets/{secret_id}/revoke/group`
- **Method**: `POST`
- **Description**: Revoke access to a secret from a specified group.
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "group_id": "123"
  }
  ```
- **Response**:
  - `200 OK`: Access revoked from group successfully.
  ```json
  {
    "message": "Access revoked from group successfully."
  }
  ```
  - `404 Not Found`: Group not found or user not part of the group.



---

### **Group Management**

#### 13. **Create Group**
- **Endpoint**: `/groups`
- **Method**: `POST`
- **Description**: Create a new group.
- **Authorization**: JWT token required.
- **Request**:
  ```json
  {
    "group_name": "DevOps Team"
  }
  ```
- **Response**:
  - `201 Created`: Group created successfully.
  ```json
  {
    "group_id": "123",
    "message": "Group created successfully"
  }
  ```

#### 14. **Get User Groups**
- **Endpoint**: `/groups`
- **Method**: `GET`
- **Description**: Retrieve all groups the authenticated user is part of.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: List of user groups.
  ```json
  [
    {
      "group_id": "123",
      "group_name": "DevOps Team",
      "role": "admin"
    }
  ]
  ```
  #### 15. **Update Group Name**
- **Endpoint**: `/groups{group_id}`
- **Method**: `PUT`
- **Description**: Update the name of a specified group.
- **Authorization**: Admin role in the group is required (user ID taken from JWT).
- **Request Body**:
  ```json
  {
    "new_group_name": "New Group Name",
  }
  ```
- **Response**:
  - `200 OK`: Successfully updated group name.
  ```json
  {
    "message": "Group name updated successfully."
  }
  ```

#### 16. **Delete Group**
- **Endpoint**: `/groups/{group_id}`
- **Method**: `DELETE`
- **Description**: Delete a specified group.
- **Authorization**: Creator of the group is required (JWT authentication).
- **Response**:
  - `200 OK`: Successfully deleted group.
  ```json
  {
    "message": "Group deleted successfully."
  }
  ```

#### 17. **Get Group by ID**
- **Endpoint**: `/groups/{group_id}`
- **Method**: `GET`
- **Description**: Retrieve details of a specific group by its ID.
- **Authorization**: User must be a member of the group.
- **Response**:
  - `200 OK`: Successfully retrieved group details.
  ```json
  {
    "group": {
      "group_id": 4,
      "group_name": "testing rename",
      "created_by": 4,
      "created_at": "0001-01-01T00:00:00Z"
    },
    "users": [
      {
        "user_id": 1,
        "email": "user1@example.com",
        "username": "user1",
        "role": "member"
      },
      {
        "user_id": 4,
        "email": "palegar.parikshith4@gmail.com",
        "username": "paroi",
        "role": "admin"
      }
    ]
  }
  ```

#### 18. **Add User to Group**
- **Endpoint**: `/groups/{group_id}/add_user`
- **Method**: `POST`
- **Description**: Add a user to a group.
- **Authorization**: JWT token required (admin permissions needed).
- **Request**:
  ```json
  {
    "user_email": "newuser@example.com",
    "role": "member" // or 'admin'
  }
  ```
- **Response**:
  - `200 OK`: User added to the group.
  ```json
  {
    "message": "User added to the group"
  }
  ```
  - `404 Not Found`: User or group not found.

#### 19. **Remove User from Group**
- **Endpoint**: `/groups/{group_id}/remove_user`
- **Method**: `DELETE`
- **Description**: Remove a user from a group.
- **Authorization**: JWT token required (admin permissions needed).
- **Request**:
  ```json
  {
    "user_email": "newuser@example.com"
  }
  ```
- **Response**:
  - `200 OK`: User removed from the group.
  ```json
  {
    "message": "User removed from the group"
  }
  ```

---

### **Audit Logging**

#### 20. **Get Audit Logs**
- **Endpoint**: `/audit/logs`
- **Method**: `GET`
- **Description**: Retrieve all actions performed by the user for audit purposes, including secret access, updates, and group management activities.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: List of audit logs.
  ```json
  [
    {
      "log_id": "123",
      "action": "created_secret",
      "description": "Created a new secret 'My database password'",
      "timestamp": "2024-09-07T12:00:00Z"
    },
    {
      "log_id": "124",
      "action": "shared_secret",
      "description": "Shared secret '5678' with user@example.com",
      "timestamp": "2024-09-07T12:30:00Z"
    }
  ]
  ```

#### 21. **Get Audit Logs for a Secret**
- **Endpoint**: `/audit/logs/secret/{secret_id}`
- **Method**: `GET`
- **Description**: Retrieve all actions performed on a specific secret.
- **Authorization**: JWT token required.
- **Response**:
  - `200 OK`: List of audit logs for the specific secret.
  ```json
  [
    {
      "log_id": "125",
      "action": "updated_secret",
      "description": "Updated secret 'My database password'",
      "timestamp": "2024-09-07T13:00:00Z"
    }
  ]
  ```

---

### **Command-Line Interface (CLI)**

The CLI can be built to interact with the API. Below are examples of commands that align with the API:

#### 22. **Get Secret via CLI**
- **Command**: 
  ```bash
  passwordmanager-cli get secret --secret_id=5678
  ```
- **Description**: Fetch a specific secret using the CLI.
- **Response**:
  - Prints the decrypted secret in the terminal after user authentication via JWT token.
  
#### 23. **Set Secret via CLI**
- **Command**:
  ```bash
  passwordmanager-cli set secret --secret_type=password --description="My DB password"
  ```
- **Description**: Creates a new secret via CLI.
- **Response**:
  - Returns the secret ID and a success message.

#### 24. **Search Secret via CLI**
- **Command**:
  ```bash
  passwordmanager-cli search secret --query="database"
  ```
- **Description**: Searches for secrets that match the query.
- **Response**:
  - Lists matching secrets along with their descriptions.

---

### **Security Considerations**
- **End-to-End Encryption**: Secrets are encrypted on the client side before being transmitted to the server, ensuring the server only stores encrypted data. Encryption keys are never sent to the server.
- **JWT Authentication**: The system uses JWT for stateless authentication. Access tokens expire after a certain period and can be refreshed using refresh tokens.
- **Role-Based Access Control**: For group functionalities, permissions can be set at a granular level (read/write access).
- **Audit Logs**: All sensitive actions (creating, updating, deleting secrets) are logged for auditing purposes.

---

### **Error Handling**
- **401 Unauthorized**: Returned when the JWT token is invalid or expired.
- **403 Forbidden**: Returned when a user attempts to access a resource they do not have permission for.
- **404 Not Found**: Returned when the requested resource (secret, group, user) is not found.
- **500 Internal Server Error**: Returned in the case of a server-side issue.

---

### **Rate Limiting & Throttling**
To protect the API from abuse, rate-limiting should be implemented. For instance:
- **User Auth API**: Max 5 requests per minute for login attempts.
- **Secret Management API**: Max 100 requests per minute.

---

