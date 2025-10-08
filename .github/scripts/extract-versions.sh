#!/bin/bash

# Script for extracting tool versions from Taskfile.yml
# Usage: .github/scripts/extract-versions.sh

# Path to Taskfile.yml
TASKFILE="Taskfile.yml"

# Check file existence
if [ ! -f "$TASKFILE" ]; then
  echo "Error: File $TASKFILE not found" >&2
  exit 1
fi

# Extract all variables from vars section
echo "Extracting variables from Taskfile.yml:"

# Determine the start and end of vars section
VARS_START=$(grep -n "^vars:" "$TASKFILE" | cut -d: -f1)
if [ -z "$VARS_START" ]; then
  echo "Error: vars section not found in $TASKFILE" >&2
  exit 1
fi

VARS_START=$((VARS_START + 1))

# Find the next section after vars or end of file
NEXT_SECTION=$(tail -n +$VARS_START "$TASKFILE" | grep -n "^[a-z]" | head -1 | cut -d: -f1)
if [ -n "$NEXT_SECTION" ]; then
  VARS_END=$((VARS_START + NEXT_SECTION - 2))
else
  VARS_END=$(wc -l < "$TASKFILE")
fi

# Extract all lines from vars section
VARS_SECTION=$(sed -n "${VARS_START},${VARS_END}p" "$TASKFILE")

# Initialize associative array to store variables
declare -A VARS

# Extract name and value of each variable
while IFS= read -r line; do
  # Skip empty lines and comments
  if [[ "$line" =~ ^[[:space:]]*$ || "$line" =~ ^[[:space:]]*# ]]; then
    continue
  fi
  
  # Extract name and value
  if [[ "$line" =~ ^[[:space:]]*([A-Z_0-9]+):\ *\'([^\']*)\' ]]; then
    var_name=${BASH_REMATCH[1]}
    var_value=${BASH_REMATCH[2]}
    VARS["$var_name"]="$var_value"
    echo "- $var_name: ${VARS[$var_name]}"
  elif [[ "$line" =~ ^[[:space:]]*([A-Z_0-9]+):\ *\"([^\"]*)\" ]]; then
    var_name=${BASH_REMATCH[1]}
    var_value=${BASH_REMATCH[2]}
    VARS["$var_name"]="$var_value"
    echo "- $var_name: ${VARS[$var_name]}"
  elif [[ "$line" =~ ^[[:space:]]*([A-Z_0-9]+):\ *(.*) ]]; then
    var_name=${BASH_REMATCH[1]}
    var_value=${BASH_REMATCH[2]}
    VARS["$var_name"]="$var_value"
    echo "- $var_name: ${VARS[$var_name]}"
  fi
done <<< "$VARS_SECTION"

# Find modules list
if [ -n "${VARS[MODULES]}" ]; then
  MODULES="${VARS[MODULES]}"
  echo "- found modules: $MODULES"
else
  # If not found in vars, try to find in another place (for backward compatibility)
  MODULES=$(sed -n 's/.*MODULES: \(.*\)/\1/p' "$TASKFILE" | head -1)
  echo "- modules (from old format): $MODULES"
fi

# Set GitHub Actions variables
if [ -n "$GITHUB_ENV" ]; then
  echo "Setting variables in GITHUB_ENV:"
  # Export all variables
  for var_name in "${!VARS[@]}"; do
    echo "$var_name=${VARS[$var_name]}" >> $GITHUB_ENV
    echo "  $var_name -> GITHUB_ENV"
  done
  # For compatibility, add MODULES separately if it's not in vars
  if [ -z "${VARS[MODULES]}" ] && [ -n "$MODULES" ]; then
    echo "MODULES=$MODULES" >> $GITHUB_ENV
    echo "  MODULES -> GITHUB_ENV"
  fi
fi

if [ -n "$GITHUB_OUTPUT" ]; then
  echo "Setting variables in GITHUB_OUTPUT:"
  # Export all variables
  for var_name in "${!VARS[@]}"; do
    echo "$var_name=${VARS[$var_name]}" >> $GITHUB_OUTPUT
    echo "  $var_name -> GITHUB_OUTPUT"
  done
  # For compatibility, add MODULES separately if it's not in vars
  if [ -z "${VARS[MODULES]}" ] && [ -n "$MODULES" ]; then
    echo "MODULES=$MODULES" >> $GITHUB_OUTPUT
    echo "  MODULES -> GITHUB_OUTPUT"
  fi
fi 