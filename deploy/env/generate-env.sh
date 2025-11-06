#!/bin/bash

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$SCRIPT_DIR"
COMPOSE_DIR="$SCRIPT_DIR/../compose"

# Use provided ENV_SUBST variable or system envsubst
if [ -z "$ENV_SUBST" ]; then
  if ! command -v envsubst &> /dev/null; then
    echo "❌ Error: envsubst not found in system and not provided via ENV_SUBST!"
    echo "Run the script via: task env:generate"
    exit 1
  fi
  ENV_SUBST=envsubst
fi

# Load main .env file
if [ ! -f "$SCRIPT_DIR/.env" ]; then
  echo "Error: File $SCRIPT_DIR/.env not found!"
  exit 1
fi

# Export all variables from .env for use in envsubst
set -a
source "$SCRIPT_DIR/.env"
set +a

# Function to process template and create .env file
process_template() {
  local service=$1
  local template="$TEMPLATE_DIR/${service}.env.template"
  local output="$COMPOSE_DIR/${service}/.env"
  
  echo "Processing template for service $service..."
  
  if [ ! -f "$template" ]; then
    echo "⚠️ Template $template not found, skipping..."
    return 0
  fi
  
  # Create directory if it doesn't exist yet
  mkdir -p "$(dirname "$output")"
  
  # Use envsubst to replace variables in template
  $ENV_SUBST < "$template" > "$output"
  
  echo "✅ Created file $output"
}

# Determine list of services from environment variable
if [ -z "$SERVICES" ]; then
  echo "⚠️ SERVICES variable is not set. No services to process."
  exit 0
fi

# Split service list by comma
IFS=',' read -ra services <<< "$SERVICES"
echo "🔍 Processing services: ${services[*]}"

# Process templates for all specified services
success_count=0
skip_count=0
for service in "${services[@]}"; do
  process_template "$service"
  if [ -f "$TEMPLATE_DIR/${service}.env.template" ]; then
    ((success_count++))
  else
    ((skip_count++))
  fi
done

if [ $success_count -eq 0 ]; then
  echo "⚠️ No .env files created. Check the service list and template availability."
else
  echo "🎉 Generation complete: $success_count files created, $skip_count templates skipped"
fi 