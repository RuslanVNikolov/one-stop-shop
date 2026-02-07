# One Stop Shop

**One place for everything about your car.**

A comprehensive platform for car owners to track vehicle maintenance history and connect with trusted mechanics through an integrated review system.

---

## Vision

One Stop Shop aims to be the single source of truth for car owners, providing:
- Complete vehicle service history tracking
- Mechanic discovery and reviews
- Maintenance reminders and analytics
- Integration with service providers

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Architecture**: Microservices
- **Databases**:
    - PostgreSQL (relational data)
    - MongoDB (audit logs & history)
- **Message Queue**: Apache Kafka
- **Cache**: Redis

### Frontend
- **Framework**: Next.js 14+ (React)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Query

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Development**: Local / Docker
- **Production**: TBD (VPS or Cloud)

---

## Project Structure

```
one-stop-shop/
├── backend/
│   ├── auth-service/       # Authentication & user management
│   ├── car-service/        # Vehicle & service record management
│   ├── mechanic-service/   # Mechanic profiles & reviews
│   ├── notification-service/ # Event processing & notifications
│   └── shared/             # Shared utilities, models, middleware
├── frontend/               # Next.js application
├── infrastructure/         # Docker configs, deployment scripts
├── docs/                   # Architecture & planning docs
└── scripts/                # Development & deployment utilities
```

---

## Core Features

### For Car Owners
- Track complete vehicle service history
- Store receipts and service documentation
- Find and review mechanics
- Get maintenance reminders
- Export service history for resale

### For Mechanics
- Create and manage business profiles
- Receive and respond to reviews
- Showcase specialties and certifications
- Build reputation through verified service records

---

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- Docker & Docker Compose
- GoLand (or any Go IDE)

### Quick Start

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd one-stop-shop
```

2. **Start infrastructure services**
```bash
docker-compose up -d
```

3. **Run backend services**
```bash
cd backend/auth-service
go run main.go
```

4. **Run frontend**
```bash
cd frontend
npm install
npm run dev
```

See [docs/QUICK_START.md](docs/QUICK_START.md) for detailed setup instructions.

---

## Development

### Service Ports
- Auth Service: `8001`
- Car Service: `8002`
- Mechanic Service: `8003`
- Notification Service: `8004`
- Frontend: `3000`

### Database Access
- PostgreSQL: `localhost:5432`
- MongoDB: `localhost:27017`
- Redis: `localhost:6379`
- Kafka: `localhost:9092`

---

## Documentation

- [Architecture Overview](docs/ARCHITECTURE.md) - System design and technical decisions
- [Quick Start Guide](docs/QUICK_START.md) - Development setup and workflow
- [API Documentation](docs/API.md) - Endpoint specifications (coming soon)
- [Database Schema](infrastructure/postgres/init/01-schema.sql) - PostgreSQL schema

---

## Roadmap

### Phase 1: MVP (Current)
- [x] Project setup and architecture
- [ ] Authentication service
- [ ] Basic car & service tracking
- [ ] Mechanic profiles & reviews
- [ ] Simple web interface

### Phase 2: Enhanced Features
- [ ] MongoDB audit logging
- [ ] Kafka event streaming
- [ ] Advanced search & filtering
- [ ] Email notifications
- [ ] Mobile-responsive design

### Phase 3: Advanced
- [ ] Maintenance reminders
- [ ] Analytics dashboard
- [ ] Service provider integrations
- [ ] Mobile app
- [ ] Multi-language support

---

## Contributing

This is currently a personal learning project. Contributions and suggestions are welcome once the MVP is complete.

---

## License

MIT License - See [LICENSE](LICENSE) for details

---

## Contact

For questions or suggestions, please open an issue.

---

**Built with Go, React, and a passion for making car ownership easier.**
