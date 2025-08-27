#!/bin/sh
# Aggressively clean up Docker resources after parallel Hive tests

# Stop all running containers (if any are stuck)
docker stop $(docker ps -aq) || true

# Remove all containers (including exited and dead)
docker rm -f $(docker ps -aq) || true

# Remove all dangling images
docker rmi -f $(docker images -qf dangling=true) || true

# Remove all unused volumes
docker volume rm $(docker volume ls -qf dangling=true) || true

# Remove all unused networks
docker network prune -f || true

# System-wide cleanup (use with caution)
docker system prune -af --volumes || true

echo "Docker cleanup completed."