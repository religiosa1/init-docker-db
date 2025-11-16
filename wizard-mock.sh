#!/bin/bash

# Mock wizard using gum to test bubbles TUI aesthetics
# This script demonstrates the UI/UX without any actual functionality

set -e

# Check if gum is installed
if ! command -v gum &> /dev/null; then
    echo "Error: gum is not installed. Please install it first:"
    echo "  https://github.com/charmbracelet/gum"
    exit 1
fi

# Display header
gum style \
    --foreground 212 \
    --border-foreground 212 \
    --border double \
    --align center \
    --width 50 \
    --margin "1 2" \
    --padding "1 4" \
    "Init Docker DB" "Interactive Wizard (Mock)"

echo ""

# Step 1: Database Type Selection
gum style --foreground 14 "Select database type:"
DB_TYPE=$(gum choose --selected="postgres" "postgres" "mysql" "mssql" "mongo" "redis")

echo ""
gum style --foreground 10 "✓ Selected: $DB_TYPE"
echo ""

# Set defaults based on database type
case $DB_TYPE in
    postgres)
        DEFAULT_USER="postgres"
        DEFAULT_PASSWORD="postgres"
        SUPPORTS_DATABASE=true
        SUPPORTS_AUTH=true
        ;;
    mysql)
        DEFAULT_USER="mysql"
        DEFAULT_PASSWORD=""
        SUPPORTS_DATABASE=true
        SUPPORTS_AUTH=true
        ;;
    mssql)
        DEFAULT_USER="mssql"
        DEFAULT_PASSWORD="Password12"
        SUPPORTS_DATABASE=true
        SUPPORTS_AUTH=true
        ;;
    mongo)
        DEFAULT_USER="mongo"
        DEFAULT_PASSWORD=""
        SUPPORTS_DATABASE=true
        SUPPORTS_AUTH=true
        ;;
    redis)
        DEFAULT_USER=""
        DEFAULT_PASSWORD=""
        SUPPORTS_DATABASE=false
        SUPPORTS_AUTH=false
        ;;
esac

# Step 2: Database Name (conditional)
if [ "$SUPPORTS_DATABASE" = true ]; then
    gum style --foreground 14 "Database name:"
    DB_NAME=$(gum input --placeholder "db")
    if [ -z "$DB_NAME" ]; then
        DB_NAME="db"
    fi
    echo ""
    gum style --foreground 10 "✓ Database name: $DB_NAME"
    echo ""
fi

# Step 3: Database User (conditional)
if [ "$SUPPORTS_AUTH" = true ]; then
    gum style --foreground 14 "Database user:"
    DB_USER=$(gum input --placeholder "$DEFAULT_USER")
    if [ -z "$DB_USER" ]; then
        DB_USER="$DEFAULT_USER"
    fi
    echo ""
    gum style --foreground 10 "✓ Database user: $DB_USER"
    echo ""
fi

# Step 4: Database Password (conditional)
if [ "$SUPPORTS_AUTH" = true ]; then
    gum style --foreground 14 "Database password:"
    if [ -n "$DEFAULT_PASSWORD" ]; then
        DB_PASSWORD=$(gum input --password --placeholder "$DEFAULT_PASSWORD")
    else
        DB_PASSWORD=$(gum input --password --placeholder "(none)")
    fi
    if [ -z "$DB_PASSWORD" ]; then
        DB_PASSWORD="$DEFAULT_PASSWORD"
    fi

    # Mock MSSQL password validation
    if [ "$DB_TYPE" = "mssql" ]; then
        # In real implementation, this would validate password complexity
        # For mock, we just show the validation would happen
        gum style --foreground 3 --italic "  (Password complexity validation would occur here for MSSQL)"
    fi

    echo ""
    gum style --foreground 10 "✓ Password set"
    echo ""
fi

# Step 5: Container Name (always asked)
# Generate a random-looking name as placeholder
RANDOM_NAME="$(shuf -n1 -e bold brave calm cool daring eager fancy gentle happy jolly kind lively merry noble proud quiet smart swift vital witty)-$(shuf -n1 -e ant bat cat dog elk fox hen jay owl pig ram yak)"

gum style --foreground 14 "Docker container name:"
CONTAINER_NAME=$(gum input --placeholder "$RANDOM_NAME")
if [ -z "$CONTAINER_NAME" ]; then
    CONTAINER_NAME="$RANDOM_NAME"
fi

echo ""
gum style --foreground 10 "✓ Container name: $CONTAINER_NAME"
echo ""

# Display summary
echo ""
gum style \
    --foreground 212 \
    --border-foreground 212 \
    --border rounded \
    --align left \
    --width 50 \
    --margin "1 2" \
    --padding "1 2" \
    "Configuration Summary" \
    "" \
    "Database Type:    $DB_TYPE" \
    "$([ "$SUPPORTS_DATABASE" = true ] && echo "Database Name:    $DB_NAME")" \
    "$([ "$SUPPORTS_AUTH" = true ] && echo "Database User:    $DB_USER")" \
    "$([ "$SUPPORTS_AUTH" = true ] && echo "Password:         ********")" \
    "Container Name:   $CONTAINER_NAME"

echo ""
gum style --foreground 2 --bold "✓ Mock complete! This is how the wizard would look with gum/bubbles TUI."
echo ""

# Show what would happen next
gum style --foreground 8 --italic "In the real implementation, the container would be created now..."
