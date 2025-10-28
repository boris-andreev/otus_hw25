#!/bin/sh
echo "Database is ready! Running migrations..."
goose up -env=.env 
echo "Starting server..."
./server