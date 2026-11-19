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
# 6. Bot FAQ Endpoint
# -----------------------------------------------------------------------------
echo ""
echo "[6/6] Testing Bot FAQ Endpoint..."

# Price query
BOT_PRICE_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from": "628123456789", "text": "Berapa biaya cuci kiloan?", "push_name": "TestUser"}')
BOT_PRICE_CAT=$(echo "$BOT_PRICE_RESP" | jq -r '.category // empty')
if [ "$BOT_PRICE_CAT" == "PRICE" ]; then
  pass "POST /bot/messages matches PRICE category"
else
  fail "POST /bot/messages PRICE query returned category: $BOT_PRICE_CAT"
fi

# Unknown query
BOT_UNKNOWN_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from": "628123456789", "text": "Apakah jual pulsa kuota?", "push_name": "TestUser"}')
BOT_UNKNOWN_CAT=$(echo "$BOT_UNKNOWN_RESP" | jq -r '.category // empty')
if [ "$BOT_UNKNOWN_CAT" == "UNKNOWN" ]; then
  pass "POST /bot/messages unrecognized query returns UNKNOWN with fallback menu"
else
  fail "POST /bot/messages UNKNOWN query returned category: $BOT_UNKNOWN_CAT"
fi

# Validation rejection
BOT_EMPTY_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{}')
if [ "$BOT_EMPTY_STATUS" -eq 400 ]; then
  pass "POST /bot/messages with empty body rejected (HTTP 400)"
else
  fail "POST /bot/messages empty body returned HTTP $BOT_EMPTY_STATUS"
fi

# -----------------------------------------------------------------------------
# 7. Interactive Multi-Turn Pickup Flow
# -----------------------------------------------------------------------------
echo ""
echo "[7/7] Testing Multi-Turn Interactive Pickup Flow..."

PICKUP_PHONE="6281777$(date +%s | tail -c 6)"

# Turn 1: Initiate pickup
TURN1_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$PICKUP_PHONE"'","text":"Kak mau pickup cucian dong","push_name":"TesterPickup"}')
if echo "$TURN1_RESP" | grep -q "REQUEST PENJEMPUTAN"; then
  pass "Turn 1: Bot initiates pickup session and prompts for address"
else
  fail "Turn 1: Failed initiating pickup session: $TURN1_RESP"
fi

# Turn 2: Provide Address
TURN2_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$PICKUP_PHONE"'","text":"Jl. Kemang Raya No. 15, Apartemen Kemang Tower B No 301"}')
if echo "$TURN2_RESP" | grep -q "pilih paket layanan"; then
  pass "Turn 2: Address recorded and active services menu presented"
else
  fail "Turn 2: Failed processing address: $TURN2_RESP"
fi

# Turn 3: Select Service Option (Cuci Reguler)
TURN3_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$PICKUP_PHONE"'","text":"Cuci Reguler"}')
if echo "$TURN3_RESP" | grep -q "Kapan kurir kami bisa menjemput"; then
  pass "Turn 3: Service selected and schedule prompt displayed"
else
  fail "Turn 3: Failed selecting service: $TURN3_RESP"
fi

# Turn 4: Provide Pickup Schedule
TURN4_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$PICKUP_PHONE"'","text":"Besok siang jam 13:00 WIB"}')
if echo "$TURN4_RESP" | grep -q "RINGKASAN REQUEST PENJEMPUTAN"; then
  pass "Turn 4: Schedule recorded and confirmation summary displayed"
else
  fail "Turn 4: Failed processing schedule: $TURN4_RESP"
fi

# Turn 5: Confirm with OK
TURN5_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$PICKUP_PHONE"'","text":"OK"}')
if echo "$TURN5_RESP" | grep -q "WAITING_PICKUP"; then
  pass "Turn 5: Order created with WAITING_PICKUP status and success dispatched"
else
  fail "Turn 5: Failed confirming pickup order: $TURN5_RESP"
fi

# Turn 6: Test Pickup Cancellation
CANCEL_PHONE="6281888$(date +%s | tail -c 6)"
curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$CANCEL_PHONE"'","text":"pickup","push_name":"CancelUser"}' > /dev/null

CANCEL_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$CANCEL_PHONE"'","text":"BATAL"}')
if echo "$CANCEL_RESP" | grep -q "dibatalkan"; then
  pass "Pickup session cancellation verified"
else
  fail "Pickup cancellation failed: $CANCEL_RESP"
fi

# Turn 7: Test Check Order Status by Phone
STATUS_RESP=$(curl -s -X POST "$BASE_URL/bot/messages" \
  -H "Content-Type: application/json" \
  -d '{"from":"'"$PICKUP_PHONE"'","text":"Cek status cucian saya dong","push_name":"TesterPickup"}')
if echo "$STATUS_RESP" | grep -q "MENUNGGU PENJEMPUTAN"; then
  pass "Turn 7: Check status by phone accurately tracks order"
else
  fail "Turn 7: Check status failed: $STATUS_RESP"
fi

# Gateway /groups endpoint test
GATEWAY_GROUPS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:3000/groups")
if [ "$GATEWAY_GROUPS_STATUS" -eq 200 ]; then
  pass "Gateway GET /groups returns HTTP 200"
else
  fail "Gateway GET /groups returned HTTP $GATEWAY_GROUPS_STATUS"
fi

# -----------------------------------------------------------------------------
# Results Summary
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
