#!/bin/bash

# HTTP Mode Test Script for Planka MCP Server
# This script tests all HTTP endpoints using curl

# Don't use set -e - we want to run all tests even if some fail

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
MCP_ENDPOINT="${BASE_URL}/mcp"
HEALTH_ENDPOINT="${BASE_URL}/health"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# E2E test resource tracking (for cleanup)
E2E_PROJECT_ID=""
E2E_BOARD_ID=""
E2E_LIST_ID=""
E2E_CARD_IDS=()
E2E_TASK_IDS=()

# Helper function to print test header
print_test() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Test: $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Helper function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

# Helper function to print failure
print_failure() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

# Helper function to make JSON-RPC request
jsonrpc_request() {
    local method=$1
    local params=$2
    local id=${3:-1}
    
    local request
    if [ -z "$params" ]; then
        request=$(cat <<EOF
{
  "jsonrpc": "2.0",
  "method": "$method",
  "id": $id
}
EOF
        )
    else
        request=$(cat <<EOF
{
  "jsonrpc": "2.0",
  "method": "$method",
  "params": $params,
  "id": $id
}
EOF
        )
    fi
    
    echo "$request" | curl -s -X POST \
        -H "Content-Type: application/json" \
        -d @- \
        "$MCP_ENDPOINT"
}

# Helper function to check JSON-RPC response
check_jsonrpc_response() {
    local response=$1
    local expected_method=$2
    
    # Check if response contains jsonrpc: "2.0"
    if echo "$response" | jq -e '.jsonrpc == "2.0"' > /dev/null 2>&1; then
        # Check for error
        if echo "$response" | jq -e 'has("error") and .error != null' > /dev/null 2>&1; then
            local error_code=$(echo "$response" | jq -r '.error.code // "unknown"')
            local error_message=$(echo "$response" | jq -r '.error.message // "unknown error"')
            print_failure "JSON-RPC error: $error_code - $error_message"
            echo "$response" | jq '.'
            return 1
        fi
        # Check for result (result can be null for notifications, which is valid)
        if echo "$response" | jq -e 'has("result")' > /dev/null 2>&1; then
            return 0
        fi
    fi
    
    print_failure "Invalid JSON-RPC response"
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
    return 1
}

# Helper function to extract JSON from tool call response
extract_tool_result() {
    local response=$1
    # Try to extract as JSON string first, then parse
    local text=$(echo "$response" | jq -r '.result.content[0].text // empty' 2>/dev/null)
    if [ -z "$text" ] || [ "$text" = "null" ] || [ "$text" = "" ]; then
        # If no text field, return empty JSON
        echo "{}"
        return 0
    fi
    # Check if text is already a JSON object/array (starts with { or [)
    if echo "$text" | jq -e . > /dev/null 2>&1; then
        # Already valid JSON, return as-is
        echo "$text"
    else
        # Try to parse as JSON string (handles escaped JSON strings)
        # Use try-catch in jq to handle parse errors gracefully
        local parsed=$(echo "$text" | jq -r 'try fromjson catch .' 2>/dev/null)
        if [ -n "$parsed" ] && [ "$parsed" != "null" ]; then
            echo "$parsed"
        else
            # If parsing fails, return empty object
            echo "{}"
        fi
    fi
}

# Helper function to call a tool and get the result
call_tool() {
    local tool_name=$1
    local arguments=$2
    local request_id=${3:-100}
    
    local params
    if [ -z "$arguments" ] || [ "$arguments" = "{}" ]; then
        params="{\"name\": \"$tool_name\", \"arguments\": {}}"
    else
        params="{\"name\": \"$tool_name\", \"arguments\": $arguments}"
    fi
    
    local response=$(jsonrpc_request "tools/call" "$params" $request_id)
    if check_jsonrpc_response "$response" "tools/call"; then
        extract_tool_result "$response"
    else
        return 1
    fi
}

# Test 1: Health Check
test_health() {
    print_test "Health Check (GET /health)"
    
    response=$(curl -s -X GET "$HEALTH_ENDPOINT")
    
    if echo "$response" | jq -e '.status == "ok"' > /dev/null 2>&1; then
        print_success "Health check passed"
        echo "$response" | jq '.'
    else
        print_failure "Health check failed"
        echo "$response"
        return 1
    fi
}

# Test 2: Initialize
test_initialize() {
    print_test "Initialize (POST /mcp - initialize)"
    
    response=$(jsonrpc_request "initialize" '{"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}' 1)
    
    if check_jsonrpc_response "$response" "initialize"; then
        print_success "Initialize successful"
        echo "$response" | jq '.result.serverInfo'
    else
        return 1
    fi
}

# Test 3: Initialized Notification
test_initialized() {
    print_test "Initialized Notification (POST /mcp - notifications/initialized)"
    
    response=$(jsonrpc_request "notifications/initialized" "" 2)
    
    if check_jsonrpc_response "$response" "notifications/initialized"; then
        print_success "Initialized notification accepted"
    else
        return 1
    fi
}

# Test 4: List Tools
test_list_tools() {
    print_test "List Tools (POST /mcp - tools/list)"
    
    response=$(jsonrpc_request "tools/list" "" 3)
    
    if check_jsonrpc_response "$response" "tools/list"; then
        tool_count=$(echo "$response" | jq '.result.tools | length')
        print_success "List tools successful (found $tool_count tools)"
        echo "$response" | jq '.result.tools[] | {name: .name, description: .description}' | head -20
        if [ "$tool_count" -gt 20 ]; then
            echo "... and $((tool_count - 20)) more tools"
        fi
    else
        return 1
    fi
}

# Test 5: Get Projects
test_get_projects() {
    print_test "Get Projects (POST /mcp - tools/call - get_projects)"
    
    params='{"name": "get_projects", "arguments": {}}'
    response=$(jsonrpc_request "tools/call" "$params" 4)
    
    if check_jsonrpc_response "$response" "tools/call"; then
        project_count=$(echo "$response" | jq '.result.content[0].text | fromjson | length')
        print_success "Get projects successful (found $project_count projects)"
        echo "$response" | jq -r '.result.content[0].text | fromjson | .[] | "  - \(.name) (ID: \(.id))"' | head -10
    else
        return 1
    fi
}

# Test 6: Get Project (if projects exist)
test_get_project() {
    print_test "Get Project (POST /mcp - tools/call - get_project)"
    
    # First get projects to find an ID
    params='{"name": "get_projects", "arguments": {}}'
    projects_response=$(jsonrpc_request "tools/call" "$params" 5)
    
    if ! check_jsonrpc_response "$projects_response" "tools/call"; then
        print_failure "Cannot get project - failed to fetch projects list"
        return 1
    fi
    
    project_id=$(echo "$projects_response" | jq -r '.result.content[0].text | fromjson | .[0].id // empty')
    
    if [ -z "$project_id" ] || [ "$project_id" = "null" ]; then
        echo -e "${YELLOW}⚠ No projects found, skipping get_project test${NC}"
        return 0
    fi
    
    params="{\"name\": \"get_project\", \"arguments\": {\"projectId\": \"$project_id\"}}"
    response=$(jsonrpc_request "tools/call" "$params" 6)
    
    if check_jsonrpc_response "$response" "tools/call"; then
        project_name=$(echo "$response" | jq -r '.result.content[0].text | fromjson | .name')
        print_success "Get project successful: $project_name"
    else
        return 1
    fi
}

# Test 7: Get Boards (if projects exist)
test_get_boards() {
    print_test "Get Boards (POST /mcp - tools/call - get_boards)"
    
    # First get projects to find an ID
    params='{"name": "get_projects", "arguments": {}}'
    projects_response=$(jsonrpc_request "tools/call" "$params" 7)
    
    if ! check_jsonrpc_response "$projects_response" "tools/call"; then
        print_failure "Cannot get boards - failed to fetch projects list"
        return 1
    fi
    
    project_id=$(echo "$projects_response" | jq -r '.result.content[0].text | fromjson | .[0].id // empty')
    
    if [ -z "$project_id" ] || [ "$project_id" = "null" ]; then
        echo -e "${YELLOW}⚠ No projects found, skipping get_boards test${NC}"
        return 0
    fi
    
    params="{\"name\": \"get_boards\", \"arguments\": {\"projectId\": \"$project_id\"}}"
    response=$(jsonrpc_request "tools/call" "$params" 8)
    
    if check_jsonrpc_response "$response" "tools/call"; then
        board_count=$(echo "$response" | jq -r '.result.content[0].text | fromjson | length')
        print_success "Get boards successful (found $board_count boards)"
        echo "$response" | jq -r '.result.content[0].text | fromjson | .[] | "  - \(.name) (ID: \(.id))"' | head -10
    else
        return 1
    fi
}

# Test 8: CORS Preflight
test_cors() {
    print_test "CORS Preflight (OPTIONS /mcp)"
    
    response=$(curl -s -X OPTIONS \
        -H "Origin: http://localhost:3000" \
        -H "Access-Control-Request-Method: POST" \
        -H "Access-Control-Request-Headers: Content-Type" \
        -v "$MCP_ENDPOINT" 2>&1)
    
    # Check for CORS headers in verbose output
    if echo "$response" | grep -q "Access-Control-Allow-Origin"; then
        print_success "CORS headers present"
        echo "$response" | grep -i "access-control" || true
    else
        print_failure "CORS headers missing"
        return 1
    fi
}

# Test 9: Invalid Method
test_invalid_method() {
    print_test "Invalid Method (POST /mcp - invalid_method)"
    
    response=$(jsonrpc_request "invalid_method" "" 9)
    
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        error_code=$(echo "$response" | jq -r '.error.code')
        # Accept any JSON-RPC error code as long as there's an error
        # Common codes: -32601 (Method not found), -32600 (Invalid Request), -32603 (Internal error)
        if [[ "$error_code" =~ ^-326[0-9][0-9]$ ]] || [[ "$error_code" =~ ^-327[0-9][0-9]$ ]]; then
            print_success "Invalid method correctly rejected (error code: $error_code)"
        else
            print_failure "Unexpected error code: $error_code"
            echo "$response" | jq '.'
            return 1
        fi
    else
        print_failure "Expected error for invalid method, but got success"
        return 1
    fi
}

# Test 10: Invalid JSON
test_invalid_json() {
    print_test "Invalid JSON (POST /mcp - malformed JSON)"
    
    response=$(echo "invalid json" | curl -s -X POST \
        -H "Content-Type: application/json" \
        -d @- \
        "$MCP_ENDPOINT")
    
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        print_success "Invalid JSON correctly rejected"
        echo "$response" | jq '.error'
    else
        print_failure "Expected error for invalid JSON"
        echo "$response"
        return 1
    fi
}

# E2E Test: Create Project (idempotent)
e2e_create_project() {
    local project_name="MCP Test Project - $(date +%s)"
    
    # Check if test project already exists
    local projects=$(call_tool "get_projects" "{}" 200)
    if [ $? -eq 0 ]; then
        local existing_id=$(echo "$projects" | jq -r ".[] | select(.name | startswith(\"MCP Test Project\")) | .id" | head -1)
        if [ -n "$existing_id" ] && [ "$existing_id" != "null" ]; then
            E2E_PROJECT_ID="$existing_id"
            echo -e "${YELLOW}⚠ Using existing test project: $existing_id${NC}"
            return 0
        fi
    fi
    
    # Create new project
    local result=$(call_tool "create_project" "{\"name\": \"$project_name\"}" 201)
    if [ $? -eq 0 ]; then
        E2E_PROJECT_ID=$(echo "$result" | jq -r '.id')
        print_success "Created project: $project_name (ID: $E2E_PROJECT_ID)"
        return 0
    else
        print_failure "Failed to create project"
        return 1
    fi
}

# E2E Test: Create Board (idempotent)
e2e_create_board() {
    if [ -z "$E2E_PROJECT_ID" ]; then
        print_failure "Cannot create board: no project ID"
        return 1
    fi
    
    # Check if test board already exists
    local boards=$(call_tool "get_boards" "{\"projectId\": \"$E2E_PROJECT_ID\"}" 202)
    if [ $? -eq 0 ]; then
        local existing_id=$(echo "$boards" | jq -r ".[] | select(.name == \"MCP Test Board\") | .id" | head -1)
        if [ -n "$existing_id" ] && [ "$existing_id" != "null" ]; then
            E2E_BOARD_ID="$existing_id"
            echo -e "${YELLOW}⚠ Using existing test board: $existing_id${NC}"
            return 0
        fi
    fi
    
    # Create new board
    local result=$(call_tool "create_board" "{\"projectId\": \"$E2E_PROJECT_ID\", \"name\": \"MCP Test Board\"}" 203)
    if [ $? -eq 0 ]; then
        E2E_BOARD_ID=$(echo "$result" | jq -r '.id')
        print_success "Created board: MCP Test Board (ID: $E2E_BOARD_ID)"
        return 0
    else
        print_failure "Failed to create board"
        return 1
    fi
}

# E2E Test: Create List (idempotent)
e2e_create_list() {
    if [ -z "$E2E_BOARD_ID" ]; then
        print_failure "Cannot create list: no board ID"
        return 1
    fi
    
    # Check if test list already exists
    local lists=$(call_tool "get_lists" "{\"boardId\": \"$E2E_BOARD_ID\"}" 204)
    if [ $? -eq 0 ]; then
        local existing_id=$(echo "$lists" | jq -r ".[] | select(.name == \"MCP Test List\") | .id" | head -1)
        if [ -n "$existing_id" ] && [ "$existing_id" != "null" ]; then
            E2E_LIST_ID="$existing_id"
            echo -e "${YELLOW}⚠ Using existing test list: $existing_id${NC}"
            return 0
        fi
    fi
    
    # Create new list
    local result=$(call_tool "create_list" "{\"boardId\": \"$E2E_BOARD_ID\", \"name\": \"MCP Test List\", \"position\": 1}" 205)
    if [ $? -eq 0 ]; then
        E2E_LIST_ID=$(echo "$result" | jq -r '.id')
        print_success "Created list: MCP Test List (ID: $E2E_LIST_ID)"
        return 0
    else
        print_failure "Failed to create list"
        return 1
    fi
}

# E2E Test: Create Cards
e2e_create_cards() {
    if [ -z "$E2E_LIST_ID" ]; then
        print_failure "Cannot create cards: no list ID"
        return 1
    fi
    
    local card_names=("Test Card 1" "Test Card 2" "Test Card 3")
    local position=1
    
    for card_name in "${card_names[@]}"; do
        local result=$(call_tool "create_card" "{\"listId\": \"$E2E_LIST_ID\", \"name\": \"$card_name\", \"position\": $position}" $((206 + position)))
        if [ $? -eq 0 ]; then
            local card_id=$(echo "$result" | jq -r '.id')
            E2E_CARD_IDS+=("$card_id")
            print_success "Created card: $card_name (ID: $card_id)"
            ((position++))
        else
            print_failure "Failed to create card: $card_name"
        fi
    done
    
    if [ ${#E2E_CARD_IDS[@]} -gt 0 ]; then
        return 0
    else
        return 1
    fi
}

# E2E Test: Create Tasks for Cards
e2e_create_tasks() {
    if [ ${#E2E_CARD_IDS[@]} -eq 0 ]; then
        echo -e "${YELLOW}⚠ No cards available to add tasks${NC}"
        return 0
    fi
    
    local task_names=("Task 1" "Task 2" "Task 3")
    local request_id=300
    
    for card_id in "${E2E_CARD_IDS[@]}"; do
        local position=1
        for task_name in "${task_names[@]}"; do
            local result=$(call_tool "create_task" "{\"cardId\": \"$card_id\", \"name\": \"$task_name\", \"position\": $position}" $request_id)
            if [ $? -eq 0 ] && [ -n "$result" ]; then
                local task_id=$(echo "$result" | jq -r '.id // empty' 2>/dev/null)
                if [ -n "$task_id" ] && [ "$task_id" != "null" ] && [ "$task_id" != "" ]; then
                    E2E_TASK_IDS+=("$task_id")
                    print_success "Created task: $task_name for card $card_id (Task ID: $task_id)"
                else
                    print_failure "Failed to extract task ID from response: $result"
                fi
                ((position++))
                ((request_id++))
            else
                print_failure "Failed to create task: $task_name for card $card_id"
            fi
        done
    done
    
    if [ ${#E2E_TASK_IDS[@]} -gt 0 ]; then
        return 0
    else
        return 1
    fi
}

# E2E Test: Cleanup - Delete all created resources
e2e_cleanup() {
    print_test "E2E Cleanup - Deleting created resources"
    
    local cleanup_failed=0
    
    # Delete tasks
    for task_id in "${E2E_TASK_IDS[@]}"; do
        if [ -n "$task_id" ] && [ "$task_id" != "null" ]; then
            local result=$(call_tool "delete_task" "{\"taskId\": \"$task_id\"}" 400)
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}  ✓ Deleted task: $task_id${NC}"
            else
                echo -e "${YELLOW}  ⚠ Failed to delete task: $task_id${NC}"
                cleanup_failed=1
            fi
        fi
    done
    
    # Delete cards
    for card_id in "${E2E_CARD_IDS[@]}"; do
        if [ -n "$card_id" ] && [ "$card_id" != "null" ]; then
            local result=$(call_tool "delete_card" "{\"cardId\": \"$card_id\"}" 401)
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}  ✓ Deleted card: $card_id${NC}"
            else
                echo -e "${YELLOW}  ⚠ Failed to delete card: $card_id${NC}"
                cleanup_failed=1
            fi
        fi
    done
    
    # Delete list (if delete_list tool exists)
    if [ -n "$E2E_LIST_ID" ] && [ "$E2E_LIST_ID" != "null" ]; then
        local result=$(call_tool "delete_list" "{\"listId\": \"$E2E_LIST_ID\"}" 402 2>/dev/null)
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}  ✓ Deleted list: $E2E_LIST_ID${NC}"
        else
            echo -e "${YELLOW}  ⚠ delete_list tool not available, skipping${NC}"
        fi
    fi
    
    # Delete board (if delete_board tool exists)
    if [ -n "$E2E_BOARD_ID" ] && [ "$E2E_BOARD_ID" != "null" ]; then
        local result=$(call_tool "delete_board" "{\"boardId\": \"$E2E_BOARD_ID\"}" 403 2>/dev/null)
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}  ✓ Deleted board: $E2E_BOARD_ID${NC}"
        else
            echo -e "${YELLOW}  ⚠ delete_board tool not available, skipping${NC}"
        fi
    fi
    
    # Delete project (if delete_project tool exists)
    if [ -n "$E2E_PROJECT_ID" ] && [ "$E2E_PROJECT_ID" != "null" ]; then
        local result=$(call_tool "delete_project" "{\"projectId\": \"$E2E_PROJECT_ID\"}" 404 2>/dev/null)
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}  ✓ Deleted project: $E2E_PROJECT_ID${NC}"
        else
            echo -e "${YELLOW}  ⚠ delete_project tool not available, skipping${NC}"
            echo -e "${YELLOW}  ⚠ Note: Test project '$E2E_PROJECT_ID' may need manual cleanup${NC}"
        fi
    fi
    
    if [ $cleanup_failed -eq 0 ]; then
        print_success "Cleanup completed successfully"
        return 0
    else
        print_failure "Some cleanup operations failed"
        return 1
    fi
}

# E2E Test: Full end-to-end test
test_e2e_full() {
    print_test "E2E Test: Full Create/Delete Workflow"
    
    echo "Creating resources..."
    e2e_create_project
    e2e_create_board
    e2e_create_list
    e2e_create_cards
    e2e_create_tasks
    
    echo ""
    echo -e "${BLUE}Verifying created resources...${NC}"
    
    # Verify project
    if [ -n "$E2E_PROJECT_ID" ]; then
        local project=$(call_tool "get_project" "{\"projectId\": \"$E2E_PROJECT_ID\"}" 250)
        if [ $? -eq 0 ]; then
            local project_name=$(echo "$project" | jq -r '.name')
            print_success "Verified project: $project_name"
        fi
    fi
    
    # Verify board
    if [ -n "$E2E_BOARD_ID" ]; then
        local board=$(call_tool "get_board" "{\"boardId\": \"$E2E_BOARD_ID\"}" 251)
        if [ $? -eq 0 ]; then
            local board_name=$(echo "$board" | jq -r '.name')
            print_success "Verified board: $board_name"
        fi
    fi
    
    # Verify list
    if [ -n "$E2E_LIST_ID" ]; then
        local list=$(call_tool "get_list" "{\"listId\": \"$E2E_LIST_ID\"}" 252)
        if [ $? -eq 0 ]; then
            local list_name=$(echo "$list" | jq -r '.name')
            print_success "Verified list: $list_name"
        fi
    fi
    
    # Verify cards
    if [ ${#E2E_CARD_IDS[@]} -gt 0 ]; then
        local cards=$(call_tool "get_cards" "{\"listId\": \"$E2E_LIST_ID\"}" 253)
        if [ $? -eq 0 ]; then
            local card_count=$(echo "$cards" | jq 'length')
            print_success "Verified cards: $card_count cards found"
        fi
    fi
    
    # Verify tasks
    if [ ${#E2E_CARD_IDS[@]} -gt 0 ] && [ ${#E2E_TASK_IDS[@]} -gt 0 ]; then
        local first_card_id="${E2E_CARD_IDS[0]}"
        local tasks=$(call_tool "get_tasks" "{\"cardId\": \"$first_card_id\"}" 254)
        if [ $? -eq 0 ]; then
            local task_count=$(echo "$tasks" | jq 'length')
            print_success "Verified tasks: $task_count tasks found for first card"
        fi
    fi
    
    echo ""
    echo -e "${BLUE}Cleaning up resources...${NC}"
    e2e_cleanup
    
    print_success "E2E test completed"
}

# Main execution
main() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    Planka MCP HTTP Mode Test Suite                          ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo "Testing server at: $BASE_URL"
    echo "MCP endpoint: $MCP_ENDPOINT"
    echo ""
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}Error: jq is required but not installed.${NC}"
        echo "Install it with: sudo apt-get install jq (or brew install jq on macOS)"
        exit 1
    fi
    
    # Check if server is reachable
    if ! curl -s -f "$HEALTH_ENDPOINT" > /dev/null 2>&1; then
        echo -e "${RED}Error: Cannot connect to server at $BASE_URL${NC}"
        echo "Make sure the server is running with: ./mcp-planka --http --http-port 8080"
        exit 1
    fi
    
    # Run all tests (continue even if one fails)
    test_health
    test_initialize
    test_initialized
    test_list_tools
    test_get_projects
    test_get_project
    test_get_boards
    test_cors
    test_invalid_method
    test_invalid_json
    
    # E2E Test: Full create/delete workflow
    test_e2e_full
    
    # Print summary
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Test Summary${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    
    total=$((TESTS_PASSED + TESTS_FAILED))
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}✗ Some tests failed${NC}"
        exit 1
    fi
}

# Run main function
main "$@"

