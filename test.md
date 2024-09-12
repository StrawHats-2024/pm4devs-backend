#### 14. **Update Group Name**
- **Endpoint**: `/groups`
- **Method**: `PUT`
- **Description**: Update the name of a specified group.
- **Authorization**: Admin role in the group is required (user ID taken from JWT).
- **Request Body**:
  ```json
  {
    "new_group_name": "New Group Name",
    "group_id": 123
  }
  ```
- **Response**:
  - `200 OK`: Successfully updated group name.
  ```json
  {
    "message": "Group name updated successfully."
  }
  ```

#### 15. **Delete Group**
- **Endpoint**: `/groups`
- **Method**: `DELETE`
- **Description**: Delete a specified group.
- **Authorization**: Creator of the group is required (JWT authentication).
- **Request Body**:
  ```json
  {
    "group_id": 123
  }
  ```
- **Response**:
  - `200 OK`: Successfully deleted group.
  ```json
  {
    "message": "Group deleted successfully."
  }
  ```

#### 16. **Get Group by ID**
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
