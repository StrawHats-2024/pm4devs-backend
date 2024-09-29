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

### 2. Login

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

- **Method**: `POST`
- **Endpoint**: `/auth/verify-token`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    ```json
    {
      "token": "JWT_TOKEN"
    }
    ```

---

### 4. Refresh Token

**Description**: Refreshes the authentication token.

- **Method**: `POST`
- **Endpoint**: `/auth/refresh`
- **Request Headers**:
    - `Authorization: Bearer <token>` (required)

- **Response**:

    **Success (200)**:
    ```json
    {
      "token": "JWT_TOKEN"
    }
    ```
