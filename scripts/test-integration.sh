#!/usr/bin/env bash
# =============================================================================
# LaundryOS — End-to-End Integration & Regression Test Suite
# =============================================================================
set -e

BASE_URL="${API_URL:-http://localhost:8080}"
PASSED=0
FAILED=0

pass() {
  echo "  ✅ PASS: $1"
  PASSED=$((PASSED + 1))
}

fail() {
  echo "  ❌ FAIL: $1"
  FAILED=$((FAILED + 1))
}

echo "=========================================================="
echo " Starting LaundryOS Integration & Regression Test Suite"
echo " Target API: $BASE_URL"
echo "=========================================================="

# -----------------------------------------------------------------------------
# 1. Health & Liveness Checks
# -----------------------------------------------------------------------------
echo ""
echo "[1/6] Testing Health & Swagger Endpoints..."

HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
if [ "$HEALTH_STATUS" -eq 200 ]; then
  pass "GET /health returns HTTP 200"
else
  fail "GET /health returned HTTP $HEALTH_STATUS"
fi

SWAGGER_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/swagger/doc.json")
if [ "$SWAGGER_STATUS" -eq 200 ]; then
  pass "GET /swagger/doc.json returns HTTP 200"
else
  fail "GET /swagger/doc.json returned HTTP $SWAGGER_STATUS"
fi

# -----------------------------------------------------------------------------
# 2. Validation & Error Handling (Regression Testing)
# -----------------------------------------------------------------------------
echo ""
echo "[2/6] Testing Request Validation & Rejection..."

VAL_CUST_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/customers" \
  -H "Content-Type: application/json" -d '{}')
if [ "$VAL_CUST_STATUS" -eq 400 ]; then
  pass "POST /customers with empty body rejected (HTTP 400)"
else
  fail "POST /customers with empty body returned HTTP $VAL_CUST_STATUS"
fi

VAL_ORDER_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" -d '{"customer_id": "not-a-uuid", "service_id": "not-a-uuid"}')
if [ "$VAL_ORDER_STATUS" -eq 400 ]; then
  pass "POST /orders with invalid UUID rejected (HTTP 400)"
else
  fail "POST /orders with invalid UUID returned HTTP $VAL_ORDER_STATUS"
fi

# -----------------------------------------------------------------------------
# 3. Customer Full Lifecycle
# -----------------------------------------------------------------------------
echo ""
echo "[3/6] Testing Customer Lifecycle..."

RANDOM_PHONE="08$(date +%s%N | cut -b1-9)"
CUST_CREATE_RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/customers" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test Integration Customer",
    "phone": "'"$RANDOM_PHONE"'",
    "email": "test@laundryos.internal",
    "address": "Jl. Testing No. 123"
  }')

CUST_STATUS=$(echo "$CUST_CREATE_RESP" | grep "HTTP_STATUS" | cut -d: -f2)
CUST_BODY=$(echo "$CUST_CREATE_RESP" | grep -v "HTTP_STATUS")
CUSTOMER_ID=$(echo "$CUST_BODY" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ "$CUST_STATUS" -eq 201 ] && [ -n "$CUSTOMER_ID" ]; then
  pass "POST /customers creates new customer (ID: $CUSTOMER_ID)"
else
  fail "POST /customers failed (HTTP $CUST_STATUS)"
fi

# Duplicate phone rejection
DUP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/customers" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Duplicate Tester",
    "phone": "'"$RANDOM_PHONE"'"
  }')
if [ "$DUP_STATUS" -eq 400 ]; then
  pass "POST /customers with duplicate phone rejected (HTTP 400)"
else
  fail "POST /customers duplicate phone returned HTTP $DUP_STATUS"
fi

# Get detail
CUST_GET_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/customers/$CUSTOMER_ID")
if [ "$CUST_GET_STATUS" -eq 200 ]; then
  pass "GET /customers/$CUSTOMER_ID returns HTTP 200"
else
  fail "GET /customers/$CUSTOMER_ID returned HTTP $CUST_GET_STATUS"
fi

# Update
CUST_UPD_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$BASE_URL/customers/$CUSTOMER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test Integration Customer Updated",
    "phone": "'"$RANDOM_PHONE"'",
    "email": "updated@laundryos.internal"
  }')
if [ "$CUST_UPD_STATUS" -eq 200 ]; then
  pass "PUT /customers/$CUSTOMER_ID updates successfully"
else
  fail "PUT /customers/$CUSTOMER_ID returned HTTP $CUST_UPD_STATUS"
fi

# Status transition
CUST_PATCH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH "$BASE_URL/customers/$CUSTOMER_ID/status" \
  -H "Content-Type: application/json" -d '{"status": "INACTIVE"}')
if [ "$CUST_PATCH_STATUS" -eq 200 ]; then
  pass "PATCH /customers/$CUSTOMER_ID/status transitions status to INACTIVE"
else
  fail "PATCH /customers/$CUSTOMER_ID/status returned HTTP $CUST_PATCH_STATUS"
fi

# -----------------------------------------------------------------------------
# 4. Order Full Lifecycle
# -----------------------------------------------------------------------------
echo ""
echo "[4/6] Testing Order Lifecycle..."

# Reactivate customer for order placement
curl -s -X PATCH "$BASE_URL/customers/$CUSTOMER_ID/status" \
  -H "Content-Type: application/json" -d '{"status": "ACTIVE"}' > /dev/null

SERVICE_ID="a0000000-0000-0000-0000-000000000001"

ORDER_CREATE_RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "'"$CUSTOMER_ID"'",
    "service_id": "'"$SERVICE_ID"'",
    "total_weight": 3.5,
    "total_price": 28000.0,
    "pickup_required": true,
    "notes": "Regression test order"
  }')

ORD_STATUS=$(echo "$ORDER_CREATE_RESP" | grep "HTTP_STATUS" | cut -d: -f2)
ORD_BODY=$(echo "$ORDER_CREATE_RESP" | grep -v "HTTP_STATUS")
ORDER_ID=$(echo "$ORD_BODY" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
ORDER_NUM=$(echo "$ORD_BODY" | grep -o '"order_number":"[^"]*' | cut -d'"' -f4)

if [ "$ORD_STATUS" -eq 201 ] && [ -n "$ORDER_ID" ]; then
  pass "POST /orders creates order (ID: $ORDER_ID, Number: $ORDER_NUM)"
else
  fail "POST /orders failed (HTTP $ORD_STATUS)"
fi

# Get order detail
ORD_GET_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/orders/$ORDER_ID")
if [ "$ORD_GET_STATUS" -eq 200 ]; then
  pass "GET /orders/$ORDER_ID returns HTTP 200"
else
  fail "GET /orders/$ORDER_ID returned HTTP $ORD_GET_STATUS"
fi

# Update order
ORD_UPD_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$BASE_URL/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "'"$CUSTOMER_ID"'",
    "service_id": "'"$SERVICE_ID"'",
    "total_weight": 5.0,
    "total_price": 40000.0,
    "pickup_required": false,
    "notes": "Updated regression notes"
  }')
if [ "$ORD_UPD_STATUS" -eq 200 ]; then
  pass "PUT /orders/$ORDER_ID updates order fields"
else
  fail "PUT /orders/$ORDER_ID returned HTTP $ORD_UPD_STATUS"
fi

# List orders
ORD_LIST_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/orders")
if [ "$ORD_LIST_STATUS" -eq 200 ]; then
  pass "GET /orders lists orders successfully"
else
  fail "GET /orders returned HTTP $ORD_LIST_STATUS"
fi

# Delete order
ORD_DEL_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/orders/$ORDER_ID")
if [ "$ORD_DEL_STATUS" -eq 204 ]; then
  pass "DELETE /orders/$ORDER_ID deletes order (HTTP 204)"
else
  fail "DELETE /orders/$ORDER_ID returned HTTP $ORD_DEL_STATUS"
fi

ORD_DEL_CHECK=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/orders/$ORDER_ID")
if [ "$ORD_DEL_CHECK" -eq 404 ]; then
  pass "GET deleted order properly returns HTTP 404"
else
  fail "GET deleted order returned HTTP $ORD_DEL_CHECK"
fi

# -----------------------------------------------------------------------------
# 5. Customer Soft Delete & Cleanup
# -----------------------------------------------------------------------------
echo ""
echo "[5/6] Testing Soft Delete & Cleanup..."

CUST_DEL_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/customers/$CUSTOMER_ID")
if [ "$CUST_DEL_STATUS" -eq 204 ]; then
  pass "DELETE /customers/$CUSTOMER_ID soft-deletes customer (HTTP 204)"
else
  fail "DELETE /customers/$CUSTOMER_ID returned HTTP $CUST_DEL_STATUS"
fi

CUST_DEL_CHECK=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/customers/$CUSTOMER_ID")
if [ "$CUST_DEL_CHECK" -eq 404 ]; then
  pass "GET soft-deleted customer returns HTTP 404"
else
  fail "GET soft-deleted customer returned HTTP $CUST_DEL_CHECK"
fi

# -----------------------------------------------------------------------------
# 6. Results Summary
# -----------------------------------------------------------------------------
echo ""
echo "=========================================================="
echo " Integration & Regression Test Results:"
echo " Passed: $PASSED"
echo " Failed: $FAILED"
echo "=========================================================="

if [ "$FAILED" -gt 0 ]; then
  exit 1
fi
