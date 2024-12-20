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
  - `user_email` (string, required): Email of the user to share with
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
  - `user_email` (string, required): Email of the user to share with
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
  - `user_email` (string, required): Email of the user to share with
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
  - `group_name` (string, required): Name of the group to share with
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
  - `group_name` (string, required): Name of the group to share with
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
  - `group_name` (string, required): Name of the group to share with
- **Responses**:
  - 200 OK: Permission revoked successfully
  - 422 Unprocessable Entity: Validation errors
  - 401 Unauthorized: User not owner of the secret


### 11. Get Secrets Shared By User
- **Endpoint**: `/v1/secrets/sharedby/user`
- **Method**: GET
- **Description**: Retrieves all secrets that the authenticated user has shared with other users.
- **Request Body**: None
- **Response Body**:
  ```json
  {
    "message": "Success!",
    "data": [
      {
        // Array of shared secret objects
      }
    ]
  }
  ```
- **Responses**:
  - 200 OK: Successfully retrieved shared secrets
  - 405 Method Not Allowed: Invalid HTTP method
  - 401 Unauthorized: User not authenticated
  - 500 Internal Server Error: Server-side error occurred

### 12. Get Secrets Shared To User's Groups
- **Endpoint**: `/v1/secrets/sharedto/group`
- **Method**: GET 
- **Description**: Retrieves all secrets that have been shared with groups that the authenticated user belongs to.
- **Request Body**: None
- **Response Body**:
  ```json
  {
    "message": "Success!",
    "data": [
      {
        // Array of shared secret objects
      }
    ]
  }
  ```
- **Responses**:
  - 200 OK: Successfully retrieved secrets shared to user's groups
  - 405 Method Not Allowed: Invalid HTTP method
  - 401 Unauthorized: User not authenticated
  - 500 Internal Server Error: Server-side error occurred

### 13. Get Secrets Shared To User
- **Endpoint**: `/v1/secrets/sharedto/user`
- **Method**: GET
- **Description**: Retrieves all secrets that have been directly shared with the authenticated user by other users.
- **Request Body**: None
- **Response Body**:
  ```json
  {
    "message": "Success!",
    "data": [
      {
        // Array of secret objects shared with the user
      }
    ]
  }
  ```
- **Responses**:
  - 200 OK: Successfully retrieved secrets shared to user
  - 405 Method Not Allowed: Invalid HTTP method
  - 401 Unauthorized: User not authenticated
  - 500 Internal Server Error: Server-side error occurred



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

### 2. Get Group by Name

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

### 5. List user groups

- **Endpoint**: `/v1/groups/user`
- **Method**: GET
- **Responses**:
- 200 OK: Group updated successfully
- 401 Unauthorized: User not owner of the group

### 6. Add User to Group

- **Endpoint**: `/v1/groups/add_user`
- **Method**: POST
- **Request Body**:
  - `group_name` (string, required): Name of the group to which the user will be added.
  - `user_email` (string, required): Email of the user to add to the group.
- **Responses**:
  - **200 OK**: User added successfully.
  - **400 Bad Request**: Invalid or missing `group_name` or `user_email`.
  - **401 Unauthorized**: Only the group owner can add members to the group.
  - **404 Not Found**: Group or user not found.

### 7. Remove User from Group

- **Endpoint**: `/v1/groups/remove_user`
- **Method**: POST
- **Request Body**:
  - `group_name` (string, required): Name of the group from which the user will be removed.
  - `user_email` (string, required): Email of the user to remove from the group.
- **Responses**:
  - **200 OK**: User removed successfully.
  - **400 Bad Request**: Invalid or missing `group_name` or `user_email`, or if attempting to remove the group creator.
  - **401 Unauthorized**: Only the group owner can remove members from the group.
  - **404 Not Found**: Group or user not found.

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
