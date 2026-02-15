# Nginx

Use the **host-installed** nginx (e.g. Ubuntu `apt install nginx`), not the Docker one.

1. **API (api.efoyetastore.com)**  
   Copy and enable the host config:
   ```bash
   sudo cp nginx/host-api.efoyetastore.com.conf /etc/nginx/sites-available/api.efoyetastore.com
   sudo ln -sf /etc/nginx/sites-available/api.efoyetastore.com /etc/nginx/sites-enabled/
   ```
   Edit the SSL cert paths if yours differ, then:
   ```bash
   sudo nginx -t && sudo systemctl reload nginx
   ```

2. **Backend**  
   Ensure the backend is reachable on the host (e.g. `docker compose` publishes port 8090). The host config proxies `/api/` to `http://127.0.0.1:8090/`.

The other files in this folder (`default.conf`, `entrypoint.sh`) are for reference or for running nginx in Docker if you choose to later.
