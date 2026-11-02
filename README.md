# Disc Golf Watch Scorekeeper

A work-in-progress personal project for recording disc golf scores on a smartwatch and reviewing rounds in a web dashboard. The goal is to make scorekeeping convenient on the course and help players track their progress over time.

I’m building this project to develop practical experience with Go, relational databases, web development, wearable applications, containers, and virtualisation. It connects these areas through a complete workflow: recording scores, validating and storing rounds, and presenting useful statistics.

## Project vision

- Record hole-by-hole scores directly on a smartwatch.
- Keep scores available during offline play and synchronize them when connected.
- Review previous rounds and performance trends in a web dashboard.

## Current progress

The backend currently validates rounds, calculates total strokes and scores relative to par, and saves rounds and hole scores to PostgreSQL in a single transaction. The Astro dashboard has an initial browser-to-API health check.

The smartwatch application, offline synchronization, round-history interface, and analytics are planned. The watch platform has not been selected yet.

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

Automated tests cover scoring, validation, save-service behavior, and HTTP handling. CI checks Go formatting, runs backend tests and static analysis, and builds the frontend.