#!/bin/bash

# Multi Inventory - PostgreSQL Setup Script

set -e

echo "🐘 Starting PostgreSQL for Multi Inventory..."

# Start PostgreSQL
docker compose -f docker-compose.postgres.yml up -d

echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Wait for PostgreSQL to be healthy
until docker exec multi_inventory_postgres pg_isready -U postgres > /dev/null 2>&1; do
  echo "Waiting for PostgreSQL to start..."
  sleep 2
done

echo "✅ PostgreSQL is ready!"
echo ""
echo "Database connection details:"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: multi_inventory"
echo "  User: postgres"
echo "  Password: postgres"
echo ""
echo "To connect to the database:"
echo "  psql -h localhost -U postgres -d multi_inventory"
echo ""
echo "To stop PostgreSQL:"
echo "  docker compose -f docker-compose.postgres.yml down"
echo ""
echo "To stop and remove data:"
echo "  docker compose -f docker-compose.postgres.yml down -v"
