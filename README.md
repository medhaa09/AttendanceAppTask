# Attendance App

## Overview
The Attendance App is a web application designed to facilitate student attendance using facial recognition. The application supports multiple user roles including students and administrators. Administrators have the capability to add users, while both students and administrators can log in, mark attendance, and view attendance history. The backend is built with Go and utilizes MongoDB for storing user data and attendance records. The Face++ API is used for matching student images to mark attendance.

## Features
1. **User Registration:** Admins can add users along with their images.
2. **User Login:** Both students and admins can log in.
3. **Marking Attendance:** Uses facial recognition to match the student's face with the stored image.

## Getting Started

### Prerequisites
- Go 1.16+
- MongoDB
- Face++ API credentials

### Variables
Change the following variables in main.go and handlers.go
```
MONGODB_URL=<Your MongoDB URL>
SECRET_KEY=<Your JWT Secret Key>
```

### Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/attendance-app.git
    cd attendance-app
    ```

2. **Install dependencies:**
    ```sh
    go mod tidy
    ```

3. **Run MongoDB:**
    Ensure MongoDB is running on your machine or use a cloud instance.

4. **Run the application:**
    ```sh
    go run main.go
    ```

5. **Access the application:**
    Open your browser and go to `http://localhost:8080`


### API Endpoints

#### Authentication
- **POST** `/login`: Login for both students and admins. 
    - Request Body:
        ```json
        {
            "handle": "user123",
            "password": "password123"
        }
        ```
    - Response:
        ```json
        {
            "token": "jwt-token",
            "refresh_token": "jwt-refresh-token"
        }
        ```

#### User Management
- **POST** `/register`: Admin registers a new user.
    - Request Headers:
        ```http
        Authorization: Bearer <admin-token>
        ```
    - Request Body:
        ```json
        {
            "handle": "user123",
            "role": "student",
            "image": "<base64-encoded-image>"
        }
        ```
    - Response:
        ```json
        {
            "message": "User registered successfully"
        }
        ```

#### Attendance
- **POST** `/attendance`: Mark attendance using facial recognition.
    - Request Headers:
        ```http
        Authorization: Bearer <user-token>
        ```
    - Request Body:
        ```json
        {
            "image": "<base64-encoded-image>"
        }
        ```
    - Response:
        ```json
        {
            "message": "Attendance marked successfully"
        }
        ```



## Detailed Explanation

### Authentication
The `auth.go` file handles user authentication, including JWT generation and validation. Middleware functions ensure that only authenticated users can access certain routes.

### Database Connection
The `db.go` file in the `utils` directory handles the connection to the MongoDB database using environment variables specified in the `.env` file.

### Middleware
- **TokenAuthMiddleware:** Validates JWT tokens for secure routes.
- **IsAdmin:** Ensures that only admin users can access certain routes.

### Token Management
The `token.go` file in the `utils` directory contains functions for generating and validating JWT tokens.

### Models
The `user.go` file in the `models` directory defines the structure of user data stored in MongoDB.

### Handlers
The `auth.go` file in the `handlers` directory contains functions for handling authentication-related requests.


## Conclusion
This Attendance App provides a comprehensive solution for managing student attendance using facial recognition. The modular architecture ensures that the application is easy to maintain and extend. By following this guide, you can set up and run the application locally, and explore its various features through the provided endpoints.

For any questions or contributions, please feel free to open an issue or a pull request on the [GitHub repository](https://github.com/medhaa09/AttendanceAppTask).
