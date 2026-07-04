# Start Dev (OpenBSD)

Start or stop the backend and frontend dev servers on the OpenBSD server.

## Usage

`/start-dev-openbsd [--stop]`

- No flag (default): stop any running instances, then start both services in the background.
- `--stop`: stop both services only; do not start them.

## Steps

1. Read `$ARGUMENTS` to determine the mode. If `--stop` is present, set mode to **stop-only**; otherwise set mode to **restart** (stop then start).

2. Read `CLAUDE.local.md` to obtain all connection details:
   - `SSH_USER` — SSH login username
   - `SSH_HOST` — SSH hostname/IP (the `.145` interface)
   - `SSH_PASS` — SSH password
   - `SSHPASS_BIN` — full path to the `sshpass` binary
   - `SSH_BIN` — full path to the `ssh` binary
   - `REPO_PATH` — absolute path to the repo root on the server
   - `BACKEND_ADDRESS` — backend listen address (e.g. `192.168.8.144:8090`)
   - `FRONTEND_ADDRESS` — frontend dev server address (e.g. `192.168.8.144:5173`)

3. Stop the backend:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST 'pkill keybook; echo backend_stopped'
   ```

4. Stop the frontend (kill the node/vite process running the dev server):
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST 'pkill -f "vite"; echo frontend_stopped'
   ```
   A non-zero exit here just means nothing was running — treat it as success.

5. **If mode is stop-only**: report both services stopped and exit.

6. Start the backend in the background:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/backend/build &&
     nohup ./keybook serve --dev --http $BACKEND_ADDRESS > keybook.log 2>&1 &
     echo $!
   '
   ```
   Report the printed PID.

7. Start the frontend dev server in the background:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/frontend &&
     nohup npm run dev > dev.log 2>&1 &
     echo $!
   '
   ```
   Report the printed PID.

8. Report both services as running with their addresses:
   - Backend: `http://$BACKEND_ADDRESS`
   - Frontend: `http://$FRONTEND_ADDRESS`

## Notes

All connection details (SSH host, user, password, binary paths, repo path, addresses) are in `CLAUDE.local.md`, which is gitignored. Do not hardcode any of those values here.

- The backend binary must already exist at `$REPO_PATH/backend/build/keybook`. If it does not, build it first with `go build -o build/keybook ./cmd/` from `$REPO_PATH/backend`.
- Both processes are started with `nohup` so they survive SSH session close.
- Logs: backend → `$REPO_PATH/backend/build/keybook.log`, frontend → `$REPO_PATH/frontend/dev.log`.
