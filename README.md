# Go Referral System API

Welcome to the **Go Referral System API**! This project is a full-stack application where users can register, log in, and refer others using unique referral codes. Successful referrals earn points for the referrer, creating an engaging system to encourage user growth.

---

## Features

### User Authentication
- **Register**: Users can sign up with their details.
- **Login**: Users can securely log in to their accounts.

### Referral System
- **Generate Referral Codes**: Each user gets a unique referral code upon registration.
- **Referral Tracking**: Users can share their referral codes with others.
- **Reward Points**: Referrers earn points when a new user registers using their code.

### Frontend
- Built with **React** and styled using **Tailwind CSS** for a responsive and user-friendly interface.

### Backend
- Developed using **Go**, ensuring high performance and scalability.
- RESTful API design for seamless interaction between the client and server.

### Database
- Powered by **PostgreSQL** for robust and reliable data management.

---

## Technologies Used

| **Tech**      | **Purpose**                       |
|---------------|-----------------------------------|
| **Go**        | Backend API                      |
| **React**     | Frontend user interface          |
| **Tailwind CSS** | Responsive and modern UI styling |
| **PostgreSQL** | Data storage and management      |

---

## Installation and Setup

### Prerequisites
- [Go](https://golang.org/) 
- [Node.js](https://nodejs.org/) (for React frontend)
- [PostgreSQL](https://www.postgresql.org/)


### Backend Setup
1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd referral-system
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up environment variables in `.env`:
   ```plaintext
   BLUEPRINT_DB_HOST=your host name
   BLUEPRINT_DB_PORT=db port
   BLUEPRINT_DB_DATABASE=database name
   BLUEPRINT_DB_USERNAME=username 
   BLUEPRINT_DB_PASSWORD=password
   BLUEPRINT_DB_SCHEMA=db schema
   ```

### Frontend Setup
1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```
4. Run the project:
   ```bash
   make run
   ```
---

## API Endpoints

| **Method** | **Endpoint**           | **Description**                |
|------------|------------------------|--------------------------------|
| `POST`     | `/api/register`        | Register a new user            |
| `POST`     | `/api/login`           | Log in an existing user        |
| `POST`     | `/api/generateReferal` | Share referral code            |
| `GET`      | `/api/referrals`       | List all referrals for a user  |


---

## Database Schema

### Users Table
| **Column**       | **Type**    | **Description**              |
|------------------|-------------|------------------------------|
| `id`             | UUID        | Unique identifier            |
| `username`       | VARCHAR     | User's name                  |
| `email`          | VARCHAR     | User's email                 |
| `password`       | VARCHAR     | Hashed password              |
| `points`         | INT         | Total earned referral points |

### Referrals Table
| **Column**       | **Type**    | **Description**              |
|------------------|-------------|------------------------------|
| `id`             | UUID        | Unique identifier            |
| `referrer_code`  | Text        | referal code                 |
| `referred_by`    | UUID        | referrer user's ID           |
| `used`           | BOOl        | Referral code is used        |

---

## Contributing

We welcome contributions! Follow these steps:
1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add feature"
   ```
4. Push to the branch:
   ```bash
   git push origin feature-name
   ```
5. Create a pull request.

---




## Contact

For any questions or feedback, feel free to reach out via GitHub issues or email at [bhatsabiq9@.com].

