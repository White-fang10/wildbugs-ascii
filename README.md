# Wild Bugs ASCII Animation Server

A small Go HTTP server that streams the supplied Wild Bugs ASCII logo as a true-color ANSI terminal animation.

## Run locally

Requirements: Go 1.22+ or Docker.

```bash
go run .
```

Then open a second terminal:

```bash
curl -N http://localhost:8080/wildbugs
```

Stop playback with `Ctrl+C`. The server clears and redraws the terminal for each frame.

A one-frame diagnostic endpoint is also available:

```bash
curl -N "http://localhost:8080/wildbugs?once=1"
```

Health check:

```bash
curl http://localhost:8080/healthz
```

## Deploy to Render for free

1. Create a GitHub repository named `wildbugs-ascii`.
2. Upload these project files to the repository root.
3. Sign in to Render and choose **New + → Blueprint** (or create a new Web Service using the repository).
4. Select the GitHub repository. Render will read `render.yaml` and build the Docker image.
5. Wait for the deployment to finish. Copy the service URL, for example `https://wildbugs-ascii.onrender.com`.
6. Test it:

   ```bash
   curl -N https://YOUR-SERVICE.onrender.com/wildbugs
   ```

Replace `YOUR-SERVICE` with the actual hostname assigned by Render. No custom domain is required.

## Endpoints

- `/wildbugs` — continuously streams the animation until the client disconnects.
- `/wildbugs?once=1` — sends one frame for a quick diagnostic.
- `/healthz` — returns `ok` for Render's health check.
- `/` — short usage instructions.

## Notes

- Use a terminal with ANSI true-color support for the intended colors and animation.
- `curl -N` disables curl's output buffering. Use `Ctrl+C` to stop the stream.
- Render's free services may sleep after inactivity. The first request after sleep can take time to start.
- A long-lived animation stream consumes an open connection while playing. Disconnecting with `Ctrl+C` closes the request.
- This is a terminal renderer adapted from the artwork and scan/glitch logic in the supplied ascii.rest snippet. Browser canvas mounting has been replaced with ANSI terminal output.
