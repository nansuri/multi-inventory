#!/bin/bash

# Multi Inventory - PostgreSQL Stop Script

set -e

echo "🛑 Stopping PostgreSQL..."

docker compose -f docker-compose.postgres.yml down

echo "✅ PostgreSQL stopped successfully!"
