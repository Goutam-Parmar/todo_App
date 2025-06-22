# ✅ TodoApp - A Simple Todo API using Go (Golang)

Hello! This is a backend app that helps users **register**, **login**, and **manage their todo tasks**.

---

## 🔧 What This App Can Do

- A user can **register** with name, email, and password
- Login with email and password
- After login, user gets a **token** (like an entry pass)
- With that token, user can:
  - ✅ Add a new todo
  - 📃 View all their todos
  - ✏️ Update their todo
  - ❌ Delete their todo
  - ✔️ Mark it as completed

---

## 🛠️ Tech Used

| Technology | Use                     |
|------------|-------------------------|
| Go (Golang) | Backend logic           |
| PostgreSQL | Database for storing data |
| bcrypt     | Password encryption     |
| Gorilla Mux | Routing APIs           |

---

## 📁 Folder Structure

TodoApp/
│
├── model/ # All request and response formats
├── handler/ # Code for API endpoints
├── db/migration/ # SQL scripts for creating tables
├── cmd/main.go # Starting point of the app


Server running at http://localhost:8081
