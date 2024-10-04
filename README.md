# Secrets API Documentation

**Base URL:** `/v1/secrets`

**Authentication:** All routes require authentication via a valid JWT token in the Authorization header (e.g., `Authorization: Bearer <token>`), except when stated otherwise.

## Endpoints:

### 1. Create a Secret

**Endpoint:** `POST /v1/secrets`

**Description:** Creates a new secret entry for the authenticated user.

**Request Body:**
```json
{
  "name": "string",                 // Required, name of the secret (e.g., "Bank Account")
  "encrypted_data": "string"        // Required, the encrypted value (e.g., an encrypted password)
}
```

**Responses:**

- **201 Created**
```json
{
  "message": "Success! Your secret has been created.",
  "secret_id": "integer"        // The ID of the newly created secret
}
```

- **422 Unprocessable Entity**
```json
{
  "error": {
    "name": "must be provided",                 // If 'name' is missing
    "encrypted_data": "must be provided"        // If 'encrypted_data' is missing
  }
}
```

- **401 Unauthorized**
  - If the user is not authenticated or token is invalid.

### 2. Retrieve a Secret

**Endpoint:** `GET /v1/secrets`

**Description:** Retrieves a secret by its ID for the authenticated user.

**Request Body:**
```json
{
  "secret_id": 1        // Required, ID of the secret to retrieve
}
```

**Responses:**

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
    "secret_id": "must be provided"    // If 'secret_id' is missing or invalid
  }
}
```

- **401 Unauthorized**
  - If the user does not have permission to access the secret.

### 3. Update a Secret

**Endpoint:** `PATCH /v1/secrets`

**Description:** Updates an existing secret for the authenticated user.

**Request Body:**
```json
{
  "secret_id": 1,                // Required, ID of the secret to update
  "name": "newname",             // Required, updated name of the secret
  "encrypted_data": "newdata"    // Required, updated encrypted data
}
```

**Responses:**

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
    "secret_id": "must be provided",        // If 'secret_id' is invalid or missing
    "name": "must be provided",             // If 'name' is missing
    "encrypted_data": "must be provided"    // If 'encrypted_data' is missing
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

**Endpoint:** `DELETE /v1/secrets`

**Description:** Deletes an existing secret for the authenticated user.

**Request Body:**
```json
{
  "secret_id": 1     // Required, ID of the secret to delete
}
```

**Responses:**

- **204 No Content**
  - No response body.

- **422 Unprocessable Entity**
```json
{
  "error": {
    "secret_id": "must be provided"    // If 'secret_id' is missing or invalid
  }
}
```

- **401 Unauthorized**
```json
{
  "message": "Only owner can delete a secret"
}
```

# Group API Documentation

## 1. Create New Group

- **Endpoint:** `/v1/groups`
- **Method:** POST
- **Authentication:** Required
- **Description:** Creates a new group.

### Request Body:

```json
{
  "group_name": "string"
}
```

- `group_name`: Must be a string of at least 5 characters.

### Responses:

#### 201 Created:

```json
{
  "Message": "Success!",
  "data": {
    "group_name": "string",
    "id": "int64"
  }
}
```

#### 400 Bad Request:
- Invalid or missing body (e.g., "group_name" key missing).

#### 422 Unprocessable Entity:
- Group name too short (less than 5 characters).
- Validation error for the request body format.

#### 409 Conflict:
- If the group name already exists.

## 2. Get Group by ID

- **Endpoint:** `/v1/groups`
- **Method:** GET
- **Authentication:** Required
- **Description:** Retrieves a group by its ID.

### Request Body:

```json
{
  "group_id": "int64"
}
```

- `group_id`: Must be a positive integer.

### Responses:

#### 200 OK:

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

#### 400 Bad Request:
- Missing body or incorrect format.

#### 422 Unprocessable Entity:
- If group_id is invalid or zero.

#### 404 Not Found:
- If the group does not exist.

## 3. Update Group

- **Endpoint:** `/v1/groups`
- **Method:** PATCH
- **Authentication:** Required
- **Description:** Updates the name of an existing group. Only the group creator can update the group name.

### Request Body:

```json
{
  "new_group_name": "string",
  "group_id": "int64"
}
```

- `new_group_name`: Must be a string of at least 5 characters.
- `group_id`: Must be a positive integer.

### Responses:

#### 200 OK:

```json
{
  "Message": "Success!"
}
```

#### 400 Bad Request:
- Missing or invalid body.

#### 422 Unprocessable Entity:
- Group name is too short, or group ID is invalid.

#### 401 Unauthorized:
- If the user is not the owner of the group.

#### 404 Not Found:
- If the group does not exist.

## 4. Delete Group

- **Endpoint:** `/v1/groups`
- **Method:** DELETE
- **Authentication:** Required
- **Description:** Deletes a group by its ID. Only the creator of the group can delete it.

### Request Body:

```json
{
  "group_id": "int64"
}
```

- `group_id`: Must be a positive integer.

### Responses:

#### 204 No Content:
- Success with no body content returned.

#### 400 Bad Request:
- Invalid or missing body.

#### 422 Unprocessable Entity:
- If group_id is invalid or zero.

#### 401 Unauthorized:
- If the user is not the creator of the group.

#### 404 Not Found:
- If the group does not exist.
