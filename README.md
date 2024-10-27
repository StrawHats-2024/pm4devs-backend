# API Documentation

## Table of Contents

1. [Authentication API](#authentication-api)
2. [Secrets API](#secrets-api)
3. [Group API](#group-api)
4. [User Secrets API](#user-secrets-api)
5. [Group Secrets API](#group-secrets-api)

List of all the routes present in the API:

1. `/v1/auth/register` (POST)
2. `/v1/auth/login` (POST)
3. `/v1/auth/logout` (POST)
4. `/v1/secrets` (POST, GET, PATCH, DELETE)
5. `/v1/secrets/share/user` (POST, PATCH, DELETE)
6. `/v1/secrets/share/group` (POST, PATCH, DELETE)
7. `/v1/groups` (POST, GET, PATCH, DELETE)
8. `/v1/secrets/user` (GET)
9. `/v1/secrets/group` (GET)

## Authentication API

### 1. Register User

- **Endpoint**: `/v1/auth/register`
- **Method**: POST
- **Description**: Create a new user account.
- **Request Body**:
  - `email` (string, required): User's email address
  - `password` (string, required): User's password
- **Responses**:
  - 201 Created: User successfully registered
  - 422 Unprocessable Entity: Validation errors
  - 409 Conflict: Email already registered

### 2. Login User

- **Endpoint**: `/v1/auth/login`
- **Method**: POST
- **Description**: Authenticate user and receive Auth token.
- **Request Body**:
  - `email` (string, required): User's email address
  - `password` (string, required): User's password
- **Responses**:
  - 200 OK: Successfully authenticated, returns token
  - 401 Unauthorized: Invalid credentials

### 3. Logout User

- **Endpoint**: `/v1/auth/logout`
- **Method**: POST
- **Description**: Invalidate user's authentication token.
- **Headers**:
  - `Authorization`: Bearer token
- **Responses**:
  - 200 OK: Successfully logged out
  - 401 Unauthorized: Invalid or missing token

## Secrets API

**Note**: All routes require authentication via Auth token in the Authorization header.

### 1. Create a Secret

- **Endpoint**: `/v1/secrets`
- **Method**: POST
- **Request Body**:
  - `name` (string, required): Name of the secret
  - `encrypted_data` (string, required): Encrypted value
  - `iv` (string required): Initialization Vector
- **Responses**:
  - 201 Created: Secret created successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not authenticated

### 2. Retrieve a Secret

- **Endpoint**: `/v1/secrets`
- **Method**: GET
- **Request Body**:
  - `secret_id` (integer, required): ID of the secret to retrieve
- **Responses**:
  - 200 OK: Secret retrieved successfully
  - 422 Unprocessable Entity: Invalid secret_id
  - 401 Unauthorized: User lacks permission

### 3. Update a Secret

- **Endpoint**: `/v1/secrets`
- **Method**: PATCH
- **Request Body**:
  - `secret_id` (integer, required): ID of the secret to update
  - `name` (string, required): Updated name of the secret
  - `encrypted_data` (string, required): Updated encrypted data
  - `iv` (string required): Updated initialization Vector
- **Responses**:
  - 200 OK: Secret updated successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

### 4. Delete a Secret

- **Endpoint**: `/v1/secrets`
- **Method**: DELETE
- **Request Body**:
  - `secret_id` (integer, required): ID of the secret to delete
- **Responses**:
  - 204 No Content: Secret deleted successfully
  - 422 Unprocessable Entity: Invalid secret_id
  - 401 Unauthorized: User not owner of the secret

### 5. Share Secret with User

- **Endpoint**: `/v1/secrets/share/user`
- **Method**: POST
- **Request Body**:
  - `secret_id` (integer, required): ID of the secret to share
  - `user_id` (integer, required): ID of the user to share with
  - `permission` (string, required): Either 'read-only' or 'read-write'
- **Responses**:
  - 201 Created: Secret shared successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

### 6. Update Permission for Shared Secret

- **Endpoint**: `/v1/secrets/update/user`
- **Method**: PATCH
- **Request Body**:
  - `secret_id` (integer, required): ID of the shared secret
  - `user_id` (integer, required): ID of the user to update permission for
  - `permission` (string, required): Either 'read-only' or 'read-write'
- **Responses**:
  - 200 OK: Permission updated successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

### 7. Revoke User's Permission for Shared Secret

- **Endpoint**: `/v1/secrets/revoke/user`
- **Method**: DELETE
- **Request Body**:
  - `secret_id` (integer, required): ID of the shared secret
  - `user_id` (integer, required): ID of the user to revoke access from
- **Responses**:
  - 200 OK: Permission revoked successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

### 8. Share Secret with Group

- **Endpoint**: `/v1/secrets/share/group`
- **Method**: POST
- **Description**: Share a secret with a group, granting either read-only or read-write access.
- **Request Body**:
  - `secret_id` (integer, required): ID of the secret to share
  - `group_id` (integer, required): ID of the group to share with
  - `permission` (string, required): Either 'read-only' or 'read-write'
- **Responses**:
  - 201 Created: Secret shared successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

### 9. Update Group Permission for Shared Secret

- **Endpoint**: `/v1/secrets/share/group`
- **Method**: PATCH
- **Description**: Update the permission level for a group that has access to a shared secret.
- **Request Body**:
  - `secret_id` (integer, required): ID of the shared secret
  - `group_id` (integer, required): ID of the group to update permission for
  - `permission` (string, required): Either 'read-only' or 'read-write'
- **Responses**:
  - 200 OK: Permission updated successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

### 10. Revoke Group's Permission for Shared Secret

- **Endpoint**: `/v1/secrets/share/group`
- **Method**: DELETE
- **Description**: Revoke a group's access to a shared secret.
- **Request Body**:
  - `secret_id` (integer, required): ID of the shared secret
  - `group_id` (integer, required): ID of the group to revoke access from
- **Responses**:
  - 200 OK: Permission revoked successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret

## Group API

### 1. Create New Group

- **Endpoint**: `/v1/groups`
- **Method**: POST
- **Request Body**:
  - `group_name` (string, required): Name of the group (minimum 5 characters)
- **Responses**:
  - 201 Created: Group created successfully
  - 400 Bad Request: Invalid or missing body
  - 422 Unprocessable Entity: Validation errors
  - 409 Conflict: Group name already exists

### 2. Get Group by ID

- **Endpoint**: `/v1/groups`
- **Method**: GET
- **Request Body**:
  - `group_name` (string, required): Name of the group to retrieve
- **Responses**:
  - 200 OK: Group retrieved successfully
  - 400 Bad Request: Invalid or missing body
  - 422 Unprocessable Entity: Invalid group_id
  - 404 Not Found: Group does not exist

### 3. Update Group

- **Endpoint**: `/v1/groups`
- **Method**: PATCH
- **Request Body**:
  - `group_name` (string, required): Name of the group to update
  - `new_group_name` (string, required): New name for the group (minimum 5 characters)
- **Responses**:
  - 200 OK: Group updated successfully
  - 400 Bad Request: Invalid or missing body
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the group
  - 404 Not Found: Group does not exist

### 4. Delete Group

- **Endpoint**: `/v1/groups`
- **Method**: DELETE
- **Request Body**:
  - `group_name` (string, required): Name of the group to delete
- **Responses**:
  - 204 No Content: Group deleted successfully
  - 400 Bad Request: Invalid or missing body
  - 422 Unprocessable Entity: Invalid group_id
  - 401 Unauthorized: User not creator of the group
  - 404 Not Found: Group does not exist

## User Secrets API

### Get User Secrets

- **Endpoint**: `/v1/secrets/user`
- **Method**: GET
- **Headers**:
  - `Authorization`: Bearer token
- **Responses**:
  - 200 OK: User secrets retrieved successfully
  - 401 Unauthorized: Authentication required

## Group Secrets API

### Get Group Secrets

- **Endpoint**: `/v1/secrets/group`
- **Method**: GET
- **Headers**:
  - `Authorization`: Bearer token
- **Request Body**:
  - `group_id` (integer, required): ID of the group
- **Responses**:
  - 200 OK: Group secrets retrieved successfully
  - 422 Unprocessable Entity: Invalid or missing group_id
  - 401 Unauthorized: User not a member of the group
