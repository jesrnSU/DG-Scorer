# Disc Golf Watch Scorekeeper

A work-in-progress personal project for recording disc golf scores on a smartwatch (currently WearOS) and reviewing your played rounds in a web dashboard.

I'm building this as there is no free alternative for WearOS, as the only choice is paying for UDisc premium.

This project involves building a backend from the ground up using Go, that communicates to the frontend using gin and REST. The web frontend is created using Astro and TypeScript while the watch frontend is written in Kotlin and Compose. The database uses Postgresql running in Docker Compose. The GO api will eventually also run in docker Compose. 

## Project vision

- Record hole-by-hole scores directly on a smartwatch.
- Use location functionalities of the smartwatch for automatic course selection and optional throw location tracking. 
- Support scoring for up to 4 players. 
- Keep scores available during offline play and synchronize them when connected.
- Review previous rounds and performance trends in a web dashboard. 

## Technology and design

| Area | Technologies |
| --- | --- |
| Backend API | Go, Gin |
| Database | PostgreSQL, pgx |
| Web frontend | Astro, TypeScript |
| Containers and development infrastructure | Docker, Docker Compose |
| Automated checks | GitHub Actions |

The project also provides a practical setting for learning containers and virtualisation, starting with running PostgreSQL in a Docker container and managing its configuration and persistent storage through Docker Compose.

The backend separates HTTP handling, scoring and validation rules, and database access. A storage interface connects the application service to PostgreSQL, keeping application logic independent of the database implementation. Transactional writes keep each round and its hole scores together.

CI checks Go formatting, builds the backend and frontend, and runs Go static analysis. Automated tests are not currently included.
