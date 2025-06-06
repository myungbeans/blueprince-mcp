#!/bin/bash

# Test script for Blue Prince MCP Tools CLI
# Make sure the MCP server is running before executing these commands

echo "=== Blue Prince MCP Tools Test Examples ==="
echo "Note: Make sure the MCP server is running with: go run ./cmd/server/main.go"
echo ""

# Test list command
echo "1. Testing list command..."
./bin/blueprince-tools list --verbose
echo ""

# Test create command
echo "2. Testing create command..."
./bin/blueprince-tools create people/test_character.md \
  --title "Test Character" \
  --category "people" \
  --primary-subject "test_character" \
  --tags "people,test,character" \
  --confidence "low" \
  --status "theory" \
  --content "# Test Character\n\nThis is a test character for CLI testing." \
  --verbose
echo ""

# Test read command
echo "3. Testing read command..."
./bin/blueprince-tools read people/test_character.md --verbose
echo ""

# Test update command
echo "4. Testing update command..."
./bin/blueprince-tools update people/test_character.md \
  --title "Updated Test Character" \
  --category "people" \
  --primary-subject "test_character" \
  --tags "people,test,character,updated" \
  --confidence "medium" \
  --status "confirmed" \
  --content "# Updated Test Character\n\nThis is an updated test character.\n\nNew information:\n- Updated via CLI\n- Status confirmed" \
  --verbose
echo ""

# Test read updated content
echo "5. Testing read of updated content..."
./bin/blueprince-tools read people/test_character.md --verbose
echo ""

echo "=== Test Complete ==="
echo "To clean up the test file, you can delete: people/test_character.md from your vault"