# API Documentation

## Authentication API

**Base URL**: `/api`

### 1. Register ✅

**Description**: Registers a new user.

- **Method**: `POST`
- **Endpoint**: `/auth/register`
- **Request Body**:
    - `name` (string, optional): User's full name.
    - `email` (string, required): User's email.
    - `password` (string, required): User's password.
    - `passwordConfirm` (string, required): Password confirmation.

- **Response**:

    **Success (200)**:
    ```json
    {
      "id": "RECORD_ID",
      "collectionId": "_pb_users_auth_",
      "collectionName": "users",
      "username": "username123",
      "verified": false,
      "emailVisibility": true,
      "email": "test@example.com",
      "created": "2022-01-01 01:00:00.123Z",
      "updated": "2022-01-01 23:59:59.456Z",
      "name": "test",
      "avatar": "filename.jpg"
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Failed to create record.",
      "data": {
        "name": {
          "code": "validation_required",
          "message": "Missing required value."
        }
      }
    }
    ```

---

### 2. Login ✅

**Description**: Authenticates a user and returns a JWT token.

- **Method**: `POST`
- **Endpoint**: `/auth/login`
- **Request Body**:
    - `email` (string, required): User's email.
    - `password` (string, required): User's password.

- **Response**:

    **Success (200)**:
    ```json
    {
      "token": "JWT_TOKEN",
      "record": {
        "id": "RECORD_ID",
        "collectionId": "_pb_users_auth_",
        "collectionName": "users",
        "username": "username123",
        "verified": false,
        "emailVisibility": true,
        "email": "test@example.com",
        "created": "2022-01-01 01:00:00.123Z",
        "updated": "2022-01-01 23:59:59.456Z",
        "name": "test",
        "avatar": "filename.jpg"
      }
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Failed to authenticate.",
      "data": {
        "identity": {
          "code": "validation_required",
          "message": "Missing required value."
        }
      }
    }
    ```

---

### 3. Verify Token

**Description**: Verifies the provided authentication token.

- **Method**: `GET`
- **Endpoint**: `/auth/verify-token`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    **Fail (403)**:
    
---

### 4. Refresh Token  

**Description**: Refreshes the authentication token.

- **Method**: `POST`
- **Endpoint**: `/auth/refresh-token`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    ```json
    {
      "token": "JWT_TOKEN"
    }
    ```


## Secrets API

### 1. Get All Secrets of Current User

**Description**: Fetches all secrets belonging to the authenticated user.

- **Method**: `GET`
- **Endpoint**: `/secrets/user`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    ```json
    {
      "page": 1,
      "perPage": 30,
      "totalPages": 1,
      "totalItems": 2,
      "items": [
        {
          "id": "RECORD_ID",
          "collectionId": "gn4sv0yna2iqmf6",
          "collectionName": "secrets",
          "created": "2022-01-01 01:00:00.123Z",
          "updated": "2022-01-01 23:59:59.456Z",
          "name": "test",
          "encrypted_data": "test",
          "owner": "RELATION_RECORD_ID"
        }
      ]
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Something went wrong while processing your request.",
      "data": {}
    }
    ```

---

### 2. Get All Secrets of a Group

**Description**: Fetches all secrets of a specific group.

- **Method**: `GET`
- **Endpoint**: `/secrets/group/{groupId}`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    ```json
    {
      "page": 1,
      "perPage": 30,
      "totalPages": 1,
      "totalItems": 2,
      "items": [
        {
          "id": "RECORD_ID",
          "collectionId": "gn4sv0yna2iqmf6",
          "collectionName": "secrets",
          "created": "2022-01-01 01:00:00.123Z",
          "updated": "2022-01-01 23:59:59.456Z",
          "name": "test",
          "encrypted_data": "test",
          "owner": "RELATION_RECORD_ID"
        }
      ]
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Something went wrong while processing your request. Invalid filter.",
      "data": {}
    }
    ```

    **Error (403)**:
    ```json
    {
      "code": 403,
      "message": "You are not allowed to perform this request.",
      "data": {}
    }
    ```

---

### 3. Get All Shared Secrets of Current User

**Description**: Fetches all secrets shared with the authenticated user.

- **Method**: `GET`
- **Endpoint**: `/secrets/shared`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    ```json
    {
      "page": 1,
      "perPage": 30,
      "totalPages": 1,
      "totalItems": 2,
      "items": [
        {
          "id": "RECORD_ID",
          "collectionId": "gn4sv0yna2iqmf6",
          "collectionName": "secrets",
          "created": "2022-01-01 01:00:00.123Z",
          "updated": "2022-01-01 23:59:59.456Z",
          "name": "test",
          "encrypted_data": "test",
          "owner": "RELATION_RECORD_ID"
        }
      ]
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Something went wrong while processing your request. Invalid filter.",
      "data": {}
    }
    ```

    **Error (403)**:
    ```json
    {
      "code": 403,
      "message": "You are not allowed to perform this request.",
      "data": {}
    }
    ```

---

### 4. Create Secret

**Description**: Creates a new secret for the authenticated user.

- **Method**: `POST`
- **Endpoint**: `/secrets`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)
- **Request Body**:
    - `name` (string, required): Name of the secret.
    - `encrypted_data` (string, required): Encrypted data of the secret.

- **Response**:

    **Success (200)**:
    ```json
    {
      "id": "RECORD_ID",
      "collectionId": "gn4sv0yna2iqmf6",
      "collectionName": "secrets",
      "created": "2022-01-01 01:00:00.123Z",
      "updated": "2022-01-01 23:59:59.456Z",
      "name": "test",
      "encrypted_data": "test",
      "owner": "RELATION_RECORD_ID"
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Failed to create record.",
      "data": {
        "name": {
          "code": "validation_required",
          "message": "Missing required value."
        }
      }
    }
    ```

    **Error (403)**:
    ```json
    {
      "code": 403,
      "message": "You are not allowed to perform this request.",
      "data": {}
    }
    ```

---

### 5. Update Secret

**Description**: Updates an existing secret belonging to the authenticated user.

- **Method**: `PUT`
- **Endpoint**: `/secrets/{secretId}`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)
- **Request Body**:
    - `name` (string, required): Name of the secret.
    - `encrypted_data` (string, required): Encrypted data of the secret.

- **Response**:

    **Success (200)**:
    ```json
    {
      "id": "RECORD_ID",
      "collectionId": "gn4sv0yna2iqmf6",
      "collectionName": "secrets",
      "created": "2022-01-01 01:00:00.123Z",
      "updated": "2022-01-01 23:59:59.456Z",
      "name": "test",
      "encrypted_data": "test",
      "owner": "RELATION_RECORD_ID"
    }
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Failed to update record.",
      "data": {
        "name": {
          "code": "validation_required",
          "message": "Missing required value."
        }
      }
    }
    ```

    **Error (403)**:
    ```json
    {
      "code": 403,
      "message": "You are not allowed to perform this request.",
      "data": {}
    }
    ```

    **Error (404)**:
    ```json
    {
      "code": 404,
      "message": "The requested resource wasn't found.",
      "data": {}
    }
    ```

---

### 6. Delete Secret

**Description**: Deletes a secret belonging to the authenticated user.

- **Method**: `DELETE`
- **Endpoint**: `/secrets/{secretId}`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (204)**:
    ```json
    null
    ```

    **Error (400)**:
    ```json
    {
      "code": 400,
      "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
      "data": {}
    }
    ```

    **Error (404)**:
    ```json
    {
      "code": 404,
      "message": "The requested resource wasn't found.",
      "data": {}
    }
    ```
