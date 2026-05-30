# Frontend OpenBSD

Build, lint, check, or manage dependencies for the frontend SvelteKit project on the OpenBSD server via SSH.

The frontend lives in `frontend/`. Available operations:

- **build** — run `npm run build` on the server; static output lands in `frontend/build/`
- **lint** — run `npm run lint` on the server (Prettier + ESLint check)
- **check** — run `npm run check` on the server (svelte-check + TypeScript)
- **dev** — run `npm run dev` on the server (dev server with hot-reload)
- **install** — run `npm install` on the server
- **reinstall** — remove `node_modules/` then run `npm install` (clean install)

## Usage

`/frontend-openbsd [build|lint|check|dev|install|reinstall]`

If no argument is given, default to **build**.

## Steps

1. Read `$ARGUMENTS` to determine the operation (build / lint / check / dev / install / reinstall). Default: build.

2. Read `CLAUDE.local.md` to obtain all connection details:
   - `SSH_USER` — SSH login username
   - `SSH_HOST` — SSH hostname/IP (the `.145` interface)
   - `SSH_PASS` — SSH password
   - `SSHPASS_BIN` — full path to the `sshpass` binary
   - `SSH_BIN` — full path to the `ssh` binary
   - `REPO_PATH` — absolute path to the repo root on the server
   - `DEV_ADDRESS` — frontend dev server address (for reporting after `dev`)

3. Push the current branch to the remote so the server has the latest code:
   ```powershell
   git push
   ```

4. Pull on the server — get the current branch name first with `git rev-parse --abbrev-ref HEAD`, then substitute into:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH &&
     git fetch --all &&
     git checkout <branch-name> &&
     git pull origin <branch-name>
   '
   ```

5. **lint** — run:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/frontend &&
     npm run lint
   '
   ```

6. **check** — run:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/frontend &&
     npm run check
   '
   ```

7. **build** — run:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/frontend &&
     npm run build
   '
   ```
   Report the contents of `frontend/build/` on success so the user can confirm the output.

8. **dev** — run (foreground; user must Ctrl-C to stop):
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/frontend &&
     npm run dev
   '
   ```
   After starting, report the dev server address from `$DEV_ADDRESS`.

9. **install** — run:
   ```sh
   $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
     cd $REPO_PATH/frontend &&
     npm install
   '
   ```

10. **reinstall** — remove node_modules then install:
    ```sh
    $SSHPASS_BIN -p '$SSH_PASS' $SSH_BIN -o StrictHostKeyChecking=no $SSH_USER@$SSH_HOST '
      cd $REPO_PATH/frontend &&
      rm -rf node_modules &&
      npm install
    '
    ```

11. Report success or failure clearly, including any stderr from the remote command.

## Notes

All connection details (SSH host, user, password, binary paths, repo path) are in `CLAUDE.local.md`, which is gitignored. Do not hardcode any of those values here.

- Build output: `frontend/build/` (SvelteKit static adapter)
- The frontend is served independently — not through PocketBase `pb_public`
