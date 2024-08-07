#!/bin/bash

# Database connection details
DB_HOST="localhost"
DB_NAME="db"
DB_USER="db_user"
DB_PASSWORD="db_password"
BACKUP_PATH="/mnt/db"
TIMESTAMP=$(date +%F-%H%M)
BACKUP_FILE="$BACKUP_PATH/$DB_NAME-$TIMESTAMP.tar.gz"
GDRIVE_FOLDER_ID="GDRIVE_FOLDER_ID"
PG_DUMP_PATH="/usr/lib/postgresql/16/bin/pg_dump" # Adjust the path if needed

# Export the password environment variable
export PGPASSWORD=$DB_PASSWORD

# Get the current date and time
current_time=$(date "+%Y-%m-%d %H:%M:%S")

# Log the current date and time
echo "[$current_time] Starting backup process..."
# Perform the backup and compress it
$PG_DUMP_PATH -Fc -h $DB_HOST -U $DB_USER $DB_NAME > $BACKUP_FILE

# Upload to Google Drive
#gdrive upload --parent $GDRIVE_FOLDER_ID $BACKUP_FILE
./mytools backup $BACKUP_FILE --parents $GDRIVE_FOLDER_ID --remove-backup=true --duration 70

# Optional: Remove old backups locally (older than 7 days)
find $BACKUP_PATH -type f -name "*.gz" -mtime +7 -exec rm {} \;

# Unset the password environment variable
unset PGPASSWORD
# Log the completion time
current_time=$(date "+%Y-%m-%d %H:%M:%S")
echo "[$current_time] Backup process completed."