# Sprint 05 - Web Dashboard & POS (Point of Sale)

## Goal

Provide a modern, responsive Web Dashboard & POS for FAN Laundry cashier and owner with seamless 2-way WhatsApp synchronization.

## Business Problem

Counter walk-in orders, petty cash expenses, and order status transitions were disconnected from WhatsApp notifications and reporting.

## Scope

- **Milestone 1**: Database Schema & Backend REST APIs (Multi-item POS, Expenses, Master Data, Audit Trail)
- **Milestone 2**: Frontend Setup, Fresh Layout & Halaman Home (KPI Cards, Quick Action Buttons, Multi-Outlet Switcher)
- **Milestone 3**: Walk-in POS Kasir & Order Lifecycle (5-tab queue, 1-click status progression, payments, 3-type receipt printing, photo documentation)
- **Milestone 4**: Manajemen Pengeluaran Kas Kecil (Petty Cash categories, history, and profit-loss deduction)
- **Milestone 5**: Laporan Keuangan, Operasional & CRM (Daily revenue line chart, unpaid receivables, service volume breakdown, top 10 customers)
- **Milestone 6**: Master Settings & RBAC Pegawai (Role-based privacy toggle for cashier, thermal receipt customizer, perfume & discount master)

## Milestones Status

- [x] **Milestone 1**: Database Schema & Backend REST APIs
  - Migration `000003_pos_and_expenses.up.sql`
  - SQLC Queries: `pos.sql`, `expense.sql`, `report.sql`, `master.sql`, `order_item.sql`, `order_photo.sql`, `order_history.sql`
  - Repositories: `POSRepository`, `ExpenseRepository`, `ReportRepository`, `MasterRepository`
  - Services: `POSService`, `ExpenseService`, `ReportService`, `MasterService`, `NotificationService`
  - Handlers: `/pos/...`, `/expenses/...`, `/reports/...`, `/master/...` with CORS middleware
  - 112/112 Go unit tests passing
- [ ] **Milestone 2**: Frontend Setup, Fresh Modern Layout & Halaman Home
- [ ] **Milestone 3**: Walk-in POS Kasir & Order Lifecycle Management
- [ ] **Milestone 4**: Manajemen Pengeluaran Kas Kecil (Petty Cash)
- [ ] **Milestone 5**: Laporan Keuangan, Operasional & CRM
- [ ] **Milestone 6**: Master Settings & Role-Based Access Control
