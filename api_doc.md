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
