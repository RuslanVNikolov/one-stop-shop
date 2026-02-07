# Architecture Documentation

## System Overview

One Stop Shop is a microservices-based platform built with Go and React, designed to help car owners track their vehicle history and connect with mechanics.

### Architecture Diagram

```
                           ┌─────────────────┐
                           │   Next.js Web   │
                           │    Frontend     │
                           │  (TypeScript)   │
                           └────────┬────────┘
                                    │ HTTPS/REST
                                    ↓
                    ┌───────────────────────────┐
                    │    API Gateway (Future)   │
                    │      nginx/traefik        │
                    └──────────┬────────────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ↓                      ↓                      ↓
┌───────────────┐      ┌───────────────┐     ┌───────────────┐
│ Auth Service  │      │  Car Service  │     │Mechanic Service│
│   Port 8001   │      │   Port 8002   │     │   Port 8003   │
│               │      │               │     │               │
│ • Registration│      │ • Car CRUD    │     │ • Profiles    │
│ • Login/JWT   │      │ • Service Rec.│     │ • Reviews     │
│ • User Mgmt   │      │ • History     │     │ • Ratings     │
└───────┬───────┘      └───────┬───────┘     └───────┬───────┘
        │                      │                     │
        └──────────────────────┴─────────────────────┘
                               │
                    ┌──────────┴──────────┐
                    │                     │
                    ↓                     ↓
            ┌──────────────┐      ┌──────────────┐
            │  PostgreSQL  │      │   MongoDB    │
            │              │      │              │
            │ • Users      │      │ • Audit Logs │
            │ • Cars       │      │ • History    │
            │ • Mechanics  │      │ • Events     │
            │ • Reviews    │      │              │
            └──────────────┘      └──────────────┘
                    │
                    ↓
            ┌──────────────┐
            │    Kafka     │      ┌──────────────────┐
            │              │◄─────│ Notification Svc │
            │ • Events     │      │   Port 8004      │
            │ • Analytics  │      │                  │
            └──────────────┘      │ • Email (future) │
                                  │ • Push (future)  │
                                  └──────────────────┘
```

---

## Core Entities

### User
```
- id: UUID (primary key)
- email: string (unique)
- password_hash: string
- full_name: string
- role: enum(car_owner, mechanic, admin)
- created_at: timestamp
- updated_at: timestamp
```

### Car
```
- id: UUID (primary key)
- owner_id: UUID (foreign key → User)
- vin: string (unique, optional)
- make: string
- model: string
- year: integer
- license_plate: string (optional)
- current_mileage: integer
- created_at: timestamp
- updated_at: timestamp
```

### ServiceRecord
```
- id: UUID (primary key)
- car_id: UUID (foreign key → Car)
- mechanic_id: UUID (foreign key → Mechanic, nullable)
- service_date: date
- mileage: integer
- service_type: string
- description: text
- cost: decimal
- parts_used: jsonb
- notes: text
- receipt_urls: string[]
- created_at: timestamp
- updated_at: timestamp
```

### Mechanic
```
- id: UUID (primary key)
- user_id: UUID (foreign key → User, nullable)
- business_name: string
- license_number: string (optional)
- address: text
- phone: string
- specialties: string[]
- years_experience: integer
- claimed: boolean
- average_rating: decimal(3,2)
- total_reviews: integer
- created_at: timestamp
- updated_at: timestamp
```

### Review
```
- id: UUID (primary key)
- mechanic_id: UUID (foreign key → Mechanic)
- reviewer_id: UUID (foreign key → User)
- service_record_id: UUID (foreign key → ServiceRecord, nullable)
- rating: integer (1-5)
- title: string
- comment: text
- response: text (mechanic reply)
- helpful_count: integer
- created_at: timestamp
- updated_at: timestamp
```

---

## Service Responsibilities

### Auth Service
**Technology**: Go, JWT, bcrypt  
**Port**: 8001  
**Database**: PostgreSQL

**Responsibilities**:
- User registration with email validation
- Login with JWT token generation
- Token refresh mechanism
- Password reset (future)
- Role-based access control

**Key Endpoints**:
```
POST   /api/auth/register      - Create new user
POST   /api/auth/login         - Authenticate user
POST   /api/auth/refresh       - Refresh access token
GET    /api/auth/me            - Get current user info
POST   /api/auth/logout        - Invalidate refresh token
```

**Key Concepts to Learn**:
- JWT token generation and validation
- Password hashing with bcrypt
- HTTP middleware for authentication
- Error handling in Go

---

### Car Service
**Technology**: Go, PostgreSQL, MongoDB, Kafka  
**Port**: 8002  
**Databases**: PostgreSQL + MongoDB

**Responsibilities**:
- CRUD operations for cars
- Service record management
- VIN validation and decoding
- History timeline generation
- Publish service events to Kafka

**Key Endpoints**:
```
POST   /api/cars                      - Add new car
GET    /api/cars                      - List user's cars
GET    /api/cars/:id                  - Get car details
PUT    /api/cars/:id                  - Update car info
DELETE /api/cars/:id                  - Remove car
POST   /api/cars/:id/services         - Add service record
GET    /api/cars/:id/services         - List service history
GET    /api/cars/:id/services/:sid    - Get service details
PUT    /api/cars/:id/services/:sid    - Update service
DELETE /api/cars/:id/services/:sid    - Delete service
```

**Key Concepts to Learn**:
- Database transactions
- Foreign key relationships
- JSON handling in Go
- Kafka producer implementation
- MongoDB integration for audit logs

---

### Mechanic Service
**Technology**: Go, PostgreSQL, Kafka  
**Port**: 8003  
**Database**: PostgreSQL

**Responsibilities**:
- Mechanic profile CRUD
- Review submission and management
- Rating aggregation
- Search and filtering
- Mechanic claiming process

**Key Endpoints**:
```
POST   /api/mechanics                 - Create mechanic profile
GET    /api/mechanics                 - Search mechanics
GET    /api/mechanics/:id             - Get mechanic details
PUT    /api/mechanics/:id             - Update profile
POST   /api/mechanics/:id/claim       - Claim profile
GET    /api/mechanics/:id/reviews     - List reviews
POST   /api/mechanics/:id/reviews     - Submit review
PUT    /api/reviews/:id               - Update review
DELETE /api/reviews/:id               - Delete review
POST   /api/reviews/:id/helpful       - Mark review helpful
```

**Key Concepts to Learn**:
- Complex SQL queries (joins, aggregations)
- Full-text search
- Database triggers for rating updates
- Input validation and sanitization

---

### Notification Service
**Technology**: Go, Kafka  
**Port**: 8004  
**Dependencies**: Kafka

**Responsibilities**:
- Consume events from Kafka
- Send notifications (email, push - future)
- Process async tasks
- Generate analytics (future)

**Kafka Topics Consumed**:
```
- service.created       - New service record added
- review.created        - New review submitted
- mechanic.claimed      - Mechanic claimed profile
- user.registered       - New user signup
```

**Key Concepts to Learn**:
- Kafka consumer groups
- Go concurrency (goroutines, channels)
- Worker pools
- Error handling in async systems

---

## Data Flow Examples

### Example 1: Creating a Service Record

```
1. User submits service record via Frontend
   ↓
2. Car Service validates request
   ↓
3. Car Service saves to PostgreSQL (service_records table)
   ↓
4. Car Service creates audit log in MongoDB
   ↓
5. Car Service publishes "service.created" event to Kafka
   ↓
6. Notification Service consumes event
   ↓
7. (Future) Send confirmation email to user
```

### Example 2: Submitting a Review

```
1. User writes review for mechanic
   ↓
2. Mechanic Service validates user owned the service
   ↓
3. Mechanic Service saves review to PostgreSQL
   ↓
4. Database trigger updates mechanic's average_rating
   ↓
5. Mechanic Service publishes "review.created" to Kafka
   ↓
6. Notification Service notifies mechanic
```

---

## Database Design Decisions

### PostgreSQL (Primary)
Used for all relational data where ACID compliance matters:
- User accounts and authentication
- Car ownership
- Mechanic profiles
- Reviews and ratings
- Service records

**Why PostgreSQL?**
- Strong consistency guarantees
- Complex queries and joins
- Foreign key constraints
- Triggers for automated rating updates
- JSON support for flexible fields (parts_used)

### MongoDB (Secondary)
Used for append-only audit logs and history:
- Service change history
- Event logs
- Future analytics data

**Why MongoDB?**
- Flexible schema for varied log types
- Excellent write performance
- Easy to query time-series data
- No schema migrations for logs

---

## Technology Choices & Learning Goals

### Why Go?
- **Concurrency**: Learn goroutines and channels
- **Performance**: Fast execution, low memory footprint
- **Microservices**: Excellent for building APIs
- **Industry Standard**: Used by Docker, Kubernetes, many startups

**Go Patterns to Practice**:
- Middleware chains
- Dependency injection
- Error handling (no exceptions)
- Interface-based design
- Worker pools

### Why Kafka?
- **Event-Driven Architecture**: Decouple services
- **Scalability**: Handle high throughput
- **Reliability**: Message persistence and replay
- **Learning**: Industry-standard messaging

**Kafka Patterns to Practice**:
- Producer/Consumer implementation
- Consumer groups
- Offset management
- Error handling and retries

### Why Next.js?
- **SEO**: Server-side rendering
- **Performance**: Built-in optimizations
- **Developer Experience**: Hot reload, TypeScript support
- **Full-Stack**: API routes for BFF pattern if needed

---

## Security Considerations

### Authentication
- Passwords hashed with bcrypt (cost factor 10)
- JWT access tokens (15 min expiry)
- Refresh tokens (7 day expiry)
- HTTPS only in production

### Authorization
- Role-based access control
- Resource ownership validation
- JWT claims for user identity

### Input Validation
- Email format validation
- VIN format validation
- SQL injection prevention (parameterized queries)
- XSS prevention (sanitize inputs)

### Data Privacy
- Password never logged or returned in responses
- Sensitive data encrypted at rest (future)
- GDPR compliance considerations (future)

---

## Deployment Strategy

### Phase 1: Local Development
- Docker Compose for all infrastructure
- Run services individually during development

### Phase 2: Single VPS
- Deploy to DigitalOcean/Hetzner ($5-10/month)
- Docker Compose for orchestration
- Nginx as reverse proxy
- Let's Encrypt for SSL

### Phase 3: Cloud (Future)
- Consider Fly.io, Railway, or Render
- Managed databases
- Auto-scaling
- CDN for static assets

---

## Performance Considerations

### Caching Strategy (Future)
- Redis for session storage
- Cache mechanic search results
- Cache aggregated ratings
- Cache user car lists

### Database Optimization
- Indexes on foreign keys
- Compound indexes for common queries
- Pagination for large result sets
- Connection pooling

### API Optimization
- Response compression (gzip)
- Request rate limiting
- Query optimization
- Batch operations where possible

---

## Monitoring & Observability (Future)

### Logging
- Structured logging with JSON
- Correlation IDs across services
- Log aggregation (ELK stack or similar)

### Metrics
- Prometheus for metrics collection
- Grafana for dashboards
- Track: request latency, error rates, throughput

### Alerting
- Alert on error rate spikes
- Alert on service downtime
- Alert on database connection issues

---

## Development Workflow

### Git Workflow
```
main (production)
  ↑
develop (integration)
  ↑
feature/* (individual features)
```

### Service Development Process
1. Design API endpoints
2. Write tests (TDD approach)
3. Implement handlers
4. Add middleware
5. Test integration
6. Deploy to dev environment

### Database Changes
1. Write migration (up/down)
2. Test locally
3. Review schema changes
4. Apply to dev environment
5. Document changes

---

## Next Steps

### Immediate (Week 1-2)
1. Implement Auth Service
2. Create frontend login/register
3. Test JWT flow end-to-end

### Short Term (Week 3-4)
1. Build Car Service
2. Add service record management
3. Create car dashboard frontend

### Medium Term (Month 2)
1. Implement Mechanic Service
2. Add review system
3. Integrate Kafka for events

### Long Term (Month 3+)
1. Add MongoDB audit logging
2. Build notification system
3. Add search and filtering
4. Deploy to production

---

## Resources & References

### Go Learning
- [Go by Example](https://gobyexample.com)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Web Examples](https://gowebexamples.com)

### Kafka
- [Kafka Go Client](https://github.com/segmentio/kafka-go)
- [Kafka Documentation](https://kafka.apache.org/documentation/)

### PostgreSQL
- [pgx driver](https://github.com/jackc/pgx)
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com)

### General Architecture
- [Microservices Patterns](https://microservices.io/patterns)
- [The Twelve-Factor App](https://12factor.net)
