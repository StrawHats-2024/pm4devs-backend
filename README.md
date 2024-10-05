# Secrets API Documentation

**Base URL:** `/v1/secrets`

**Authentication:** All routes require authentication via a valid JWT token in the Authorization header (e.g., `Authorization: Bearer <token>`), unless stated otherwise.

## Endpoints

### 1. Create a Secret

**POST /v1/secrets**

Creates a new secret entry for the authenticated user.

#### Request Body
```json
{
  "name": "string",          // Required, name of the secret (e.g., "Bank Account")
  "encrypted_data": "string" // Required, the encrypted value (e.g., an encrypted password)
}
```

#### Responses

- **201 Created**
  ```json
  {
    "message": "Success! Your secret has been created.",
    "secret_id": "integer"   // The ID of the newly created secret
  }
  ```

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "name": "must be provided",
      "encrypted_data": "must be provided"
    }
  }
  ```

- **401 Unauthorized**
  - If the user is not authenticated or token is invalid.

### 2. Retrieve a Secret

**GET /v1/secrets**

Retrieves a secret by its ID for the authenticated user.

#### Request Body
```json
{
  "secret_id": 1  // Required, ID of the secret to retrieve
}
```

#### Responses

- **200 OK**
  ```json
  {
    "message": "Success!",
    "data": {
      "secret_id": "integer",
      "name": "string",
      "encrypted_data": "string",
      "owner_id": "integer"
    }
  }
  ```

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "secret_id": "must be provided"
    }
  }
  ```

- **401 Unauthorized**
  - If the user does not have permission to access the secret.

### 3. Update a Secret

**PATCH /v1/secrets**

Updates an existing secret for the authenticated user.

#### Request Body
```json
{
  "secret_id": 1,             // Required, ID of the secret to update
  "name": "newname",          // Required, updated name of the secret
  "encrypted_data": "newdata" // Required, updated encrypted data
}
```

#### Responses

- **200 OK**
  ```json
  {
    "message": "Success!"
  }
  ```

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "secret_id": "must be provided",
      "name": "must be provided",
      "encrypted_data": "must be provided"
    }
  }
  ```

- **401 Unauthorized**
  ```json
  {
    "message": "Only owner can update a secret"
  }
  ```

### 4. Delete a Secret

**DELETE /v1/secrets**

Deletes an existing secret for the authenticated user.

#### Request Body
```json
{
  "secret_id": 1  // Required, ID of the secret to delete
}
```

#### Responses

- **204 No Content**
  - No response body.

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "secret_id": "must be provided"
    }
  }
  ```

- **401 Unauthorized**
  ```json
  {
    "message": "Only owner can delete a secret"
  }
  ```

### 5. Share Secret with User

**POST /v1/secrets/share/user**

Allows a secret owner to share a secret with another user, granting either "read-only" or "read-write" access.

#### Request Body
```json
{
  "secret_id": 1,           // Required, must be a valid secret ID owned by the current user
  "user_id": 2,             // Required, the ID of the user to share the secret with
  "permission": "read-only" // Required, must be either 'read-only' or 'read-write'
}
```

#### Responses

- **201 Created**
  ```json
  {
    "message": "Secret shared successfully with the user."
  }
  ```

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "secret_id": "must be provided",
      "user_id": "must be provided",
      "permission": "must be 'read-only' or 'read-write'"
    }
  }
  ```

- **401 Unauthorized**
  ```json
  {
    "message": "Only secret owner can manage access"
  }
  ```

### 6. Update Permission for Shared Secret

**PATCH /v1/secrets/update/user**

Allows a secret owner to update the permission level of another user who has access to the secret.

#### Request Body
```json
{
  "secret_id": 1,            // Required, must be a valid secret ID owned by the current user
  "user_id": 2,              // Required, the ID of the user whose permission is being updated
  "permission": "read-write" // Required, must be either 'read-only' or 'read-write'
}
```

#### Responses

- **200 OK**
  ```json
  {
    "message": "Permission updated successfully for the user."
  }
  ```

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "secret_id": "must be provided",
      "user_id": "must be provided",
      "permission": "must be 'read-only' or 'read-write'"
    }
  }
  ```

- **401 Unauthorized**
  ```json
  {
    "message": "Only secret owner can manage access"
  }
  ```

### 7. Revoke User's Permission for Shared Secret

**DELETE /v1/secrets/revoke/user**

Allows a secret owner to revoke access to a shared secret from a specific user.

#### Request Body
```json
{
  "secret_id": 1, // Required, must be a valid secret ID owned by the current user
  "user_id": 2    // Required, the ID of the user whose permission is being revoked
}
```

#### Responses

- **200 OK**
  ```json
  {
    "message": "Permission revoked successfully for the user."
  }
  ```

- **422 Unprocessable Entity**
  ```json
  {
    "error": {
      "secret_id": "must be provided",
      "user_id": "must be provided"
    }
  }
  ```

- **401 Unauthorized**
  ```json
  {
    "message": "Only secret owner can manage access"
  }
  ```

## Group API Documentation

### 1. Create New Group

**POST /v1/groups**

Creates a new group.

#### Request Body
```json
{
  "group_name": "string" // Must be a string of at least 5 characters
}
```

#### Responses

- **201 Created**
  ```json
  {
    "Message": "Success!",
    "data": {
      "group_name": "string",
      "id": "int64"
    }
  }
  ```

- **400 Bad Request**
  - Invalid or missing body (e.g., "group_name" key missing).

- **422 Unprocessable Entity**
  - Group name too short (less than 5 characters).
  - Validation error for the request body format.

- **409 Conflict**
  - If the group name already exists.

### 2. Get Group by ID

**GET /v1/groups**

Retrieves a group by its ID.

#### Request Body
```json
{
  "group_id": "int64" // Must be a positive integer
}
```

#### Responses

- **200 OK**
  ```json
  {
    "Message": "Success!",
    "Data": {
      "id": "int64",
      "name": "string",
      "creator_id": "int64",
      "created_at": "timestamp"
    }
  }
  ```

- **400 Bad Request**
  - Missing body or incorrect format.

- **422 Unprocessable Entity**
  - If group_id is invalid or zero.

- **404 Not Found**
  - If the group does not exist.

### 3. Update Group

**PATCH /v1/groups**

Updates the name of an existing group. Only the group creator can update the group name.

#### Request Body
```json
{
  "new_group_name": "string", // Must be a string of at least 5 characters
  "group_id": "int64"         // Must be a positive integer
}
```

#### Responses

- **200 OK**
  ```json
  {
    "Message": "Success!"
  }
  ```

- **400 Bad Request**
  - Missing or invalid body.

- **422 Unprocessable Entity**
  - Group name is too short, or group ID is invalid.

- **401 Unauthorized**
  - If the user is not the owner of the group.

- **404 Not Found**
  - If the group does not exist.

### 4. Delete Group

**DELETE /v1/groups**

Deletes a group by its ID. Only the creator of the group can delete it.

#### Request Body
```json
{
  "group_id": "int64" // Must be a positive integer
}
```

#### Responses

- **204 No Content**
  - Success with no body content returned.

- **400 Bad Request**
  - Invalid or missing body.

- **422 Unprocessable Entity**
  - If group_id is invalid or zero.

- **401 Unauthorized**
  - If the user is not the creator of the group.

- **404 Not Found**
  - If the group does not exist.
