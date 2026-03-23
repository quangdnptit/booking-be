# 🤖 AI Assistant Guidelines

## 📌 Project Overview
This is a backend microservices system built with:
- Go and Java (Spring Boot)
- PostgreSQL, Redis
- Kafka for async communication

---

## 🧱 Architecture Rules
- Follow Clean Architecture
- Separate layers:
    - handler / controller
    - service / usecase
    - repository
- No business logic in controllers

---

## 🧠 Coding Principles
- Follow SOLID principles
- Prefer composition over inheritance
- Write small, testable functions

---

## 🔄 Concurrency (Go)
- Use goroutines safely
- Avoid shared memory when possible
- Prefer channels for communication

---

## 🗄️ Database Rules
- Use transactions where needed
- Avoid N+1 queries
- Optimize indexes

---

## 📡 API Design
- RESTful conventions
- Return proper HTTP status codes
- Always return JSON

---

## 📊 Logging & Tracing
- Use structured logging (JSON)
- Always include:
    - traceId
    - requestId
- Log at key business steps

---

## ⚡ Performance
- Avoid unnecessary allocations
- Use caching (Redis) where needed
- Measure before optimizing

---

## ❌ Avoid
- God classes
- Tight coupling
- Blocking operations in async flows