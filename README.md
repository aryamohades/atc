Live link (use the Reset button to start the simulation!): https://atc-demo-dr9b7.ondigitalocean.app/

<img width="1099" alt="Screenshot 2024-10-04 at 10 38 50 AM" src="https://github.com/user-attachments/assets/8ef08d6a-53aa-4ccc-91a7-02327db7956d">  

**Description**: A web application and backend for a proof of concept "air traffic control" system to view status of encounters ingested from some EHR system. The simulation behavior is to periodically create a new random encounter (each with a specific type and severity) until 100 encounters are created. Encounters that are in a 'waiting' status for too long are assigned an alert level which increases over time (depending on encounter type severity). By default, the encounters that have the highest severity and alert level are pushed to the top for visibility. Once an encounter is successfully triaged, the alert level is reset. The system tracks how long an encounter sits in each status type ('waiting', 'triaged') as well as when the encounter starts and when it is completed for historical purposes e.g. analytics. Use the "Reset" button to start a fresh simulation run.

**Stack**:
- Go (https://go.dev/)
- Postgres (https://www.postgresql.org/)
- HTMX (https://htmx.org/)
- Alpine.js (https://alpinejs.dev/)
- DaisyUI (https://daisyui.com/)
- tailwindcss (https://tailwindcss.com/)
- esbuild (https://esbuild.github.io/)
- Docker (https://www.docker.com/)

**Infrastructure**:
- Server: DigitalOcean (https://www.digitalocean.com/)
- Database: Neon (https://neon.tech/)

**Notes**:
- The application is built into a single binary (including static files like JavaScript and CSS thanks to Go's embed feature) and containerized with Docker for orchestration and maximum portability
- Database migrations handled by a pre-deploy job in DigitalOcean. Migrations managed by goose (https://pressly.github.io/goose/)

**Improvements**
- Filters e.g. Type, Status, Severity, etc
- Switch live-reload mechanism to push instead of pull using server sent events for more "realtime" updates and less data-transfer
- Use a rules engine to enable more sophisticated/customizable logic around alerting
