# Sprint 02 - Backend Foundation

## Goal

Build the Backend foundation and core business APIs.

## Business Problem

Business data must be centralized inside the Backend instead of the WhatsApp Gateway so every module uses the same source of truth.

---

# Milestone 1 - Backend Initialization

## Status

Planned

## Goal

Prepare the Go Backend foundation and application lifecycle.

## Business Problem

The Backend must provide a stable application foundation before implementing any business features.

## Scope

- Chi HTTP Server
- Configuration loading
- Health endpoint
- Database initialization
- Graceful shutdown

## Proposed Changes

### Backend

#### [NEW] HTTP Server

- Initialize Chi.
- Configure application startup.
- Register health endpoint.

#### [NEW] Configuration

- Load application configuration.
- Load database configuration.

#### [NEW] Database

- Initialize database connection.
- Verify connection during startup.

## Risks

- Incorrect configuration may prevent application startup.
- Database connection failure should stop the application gracefully.

## Deliverables

- Backend application
- Health endpoint
- Database connection

## Out of Scope

- Customer module
- Order module
- Swagger
- Authentication
- Database migration

## Acceptance Criteria

- Backend starts successfully.
- GET /health returns HTTP 200.
- Database connection established.
- Graceful shutdown works.
- Lint passes.
- Docker build passes.

## Verification Plan

### Automated Tests

- Run lint.
- Build Docker image.

### Manual Verification

1. Start the Backend.
2. Verify GET /health returns HTTP 200.
3. Verify database connection is established.
4. Stop the application gracefully.

---

# Milestone 2 - Database Migration

## Status

Planned

## Goal

Initialize the project database schema.

## Business Problem

The Backend requires a versioned database schema to support future business modules.

## Scope

- Migration tool configuration
- Initial migration
- Customer table
- Order table

## Proposed Changes

### Database

#### [NEW] Migration

- Configure migration tool.
- Create initial migration.

#### [NEW] Schema

- Create customers table.
- Create orders table.

## Risks

- Incorrect migration order may break future migrations.
- Migration rollback must remain functional.

## Deliverables

- Initial database schema

## Out of Scope

- CRUD API
- Swagger
- Business validation

## Acceptance Criteria

- Migration up works.
- Migration down works.
- Customers table created.
- Orders table created.
- Lint passes.
- Docker build passes.

## Verification Plan

### Automated Tests

- Run migration up.
- Run migration down.
- Run lint.

### Manual Verification

1. Execute migration up.
2. Verify tables exist.
3. Execute migration down.
4. Verify tables are removed.

---

# Milestone 3 - Customer Module

## Status

Planned

## Goal

Implement Customer CRUD API.

## Business Problem

Customer information must be centralized so every business process references the same customer data.

## Scope

- Customer model
- Repository
- Service
- HTTP handler
- REST API

## Proposed Changes

### Customer Module

#### [NEW] Model

- Create Customer model.

#### [NEW] Repository

- Create Customer repository.

#### [NEW] Service

- Create Customer service.

#### [NEW] HTTP Handler

- Create Customer REST endpoints.

#### [MODIFY] Router

- Register Customer routes.

## Risks

- Customer validation rules may evolve in future milestones.
- Keep business logic inside the service layer only.

## Deliverables

- Customer CRUD API

## Out of Scope

- Order module
- Authentication
- AI
- Dashboard

## Acceptance Criteria

- Create Customer.
- Update Customer.
- Delete Customer.
- Get Customer Detail.
- List Customers.
- CRUD tested manually.
- Lint passes.
- Docker build passes.

## Verification Plan

### Automated Tests

- Run lint.

### Manual Verification

1. Create Customer.
2. Update Customer.
3. Delete Customer.
4. Get Customer Detail.
5. List Customers.

---

# Milestone 4 - Order Module

## Status

Planned

## Goal

Implement Order CRUD API.

## Business Problem

Orders must be stored centrally to support future payment, laundry workflow, and reporting.

## Scope

- Order model
- Repository
- Service
- HTTP handler
- REST API

## Proposed Changes

### Order Module

#### [NEW] Model

- Create Order model.

#### [NEW] Repository

- Create Order repository.

#### [NEW] Service

- Create Order service.

#### [NEW] HTTP Handler

- Create Order REST endpoints.

#### [MODIFY] Router

- Register Order routes.

## Risks

- Future order workflow may introduce additional statuses.
- Keep the initial implementation focused on CRUD only.

## Deliverables

- Order CRUD API

## Out of Scope

- Payment
- WhatsApp Notification
- Dashboard
- Pickup Workflow

## Acceptance Criteria

- Create Order.
- Update Order.
- Delete Order.
- Get Order Detail.
- List Orders.
- CRUD tested manually.
- Lint passes.
- Docker build passes.

## Verification Plan

### Automated Tests

- Run lint.

### Manual Verification

1. Create Order.
2. Update Order.
3. Delete Order.
4. Get Order Detail.
5. List Orders.

---

# Milestone 5 - Swagger Documentation

## Status

Planned

## Goal

Generate API documentation.

## Business Problem

API consumers require standardized documentation for integration.

## Scope

- Swagger configuration
- Customer API documentation
- Order API documentation

## Proposed Changes

### Swagger

#### [NEW] Swagger Configuration

- Configure Swagger generation.

#### [MODIFY] API Documentation

- Document Customer endpoints.
- Document Order endpoints.

## Risks

- Documentation may become outdated if APIs change.
- Keep documentation synchronized with implementation.

## Deliverables

- Swagger UI

## Out of Scope

- SDK generation
- API versioning

## Acceptance Criteria

- Swagger generated.
- Customer endpoints documented.
- Order endpoints documented.
- Swagger accessible.
- Lint passes.
- Docker build passes.

## Verification Plan

### Automated Tests

- Generate Swagger.
- Run lint.

### Manual Verification

1. Open Swagger UI.
2. Verify Customer endpoints.
3. Verify Order endpoints.

---

# Milestone 6 - API Standardization

## Status

Planned

## Goal

Standardize the Backend API before integrating with external services.

## Business Problem

A consistent API contract simplifies future integration with the WhatsApp Gateway and Dashboard.

## Scope

- Standard response format
- Error handling
- Request validation
- Logging
- README update

## Proposed Changes

### API

#### [MODIFY] Response

- Standardize success responses.
- Standardize error responses.

#### [MODIFY] Validation

- Validate incoming requests.

#### [MODIFY] Logging

- Improve request and application logging.

#### [MODIFY] Documentation

- Update README.

## Risks

- Changing response formats may impact API consumers.
- Logging should remain simple and readable.

## Deliverables

- Standardized Backend API

## Out of Scope

- Gateway integration
- Dashboard integration
- AI integration

## Acceptance Criteria

- Response format standardized.
- Error responses standardized.
- Validation implemented.
- Logging improved.
- README updated.
- Lint passes.
- Docker build passes.

## Verification Plan

### Automated Tests

- Run lint.
- Build Docker image.

### Manual Verification

1. Verify health endpoint.
2. Verify successful responses.
3. Verify validation errors.
4. Verify logging output.
5. Verify README examples.

---

# Sprint Deliverables

- Backend Foundation
- Database Migration
- Customer CRUD API
- Order CRUD API
- Swagger Documentation
- Standardized Backend API

---

# Definition of Done

- All milestone acceptance criteria completed.
- Manual verification completed.
- Lint passes.
- Docker build passes.
- CURRENT.md updated.
- RELEASE_NOTES.md updated.
- Git tag created.