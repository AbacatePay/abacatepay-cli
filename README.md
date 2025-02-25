# Abacatepay-cli
AbacatePay CLI for you to run your webhooks locally

# Tech
It was made with Node.js and TypeScript.

# How To Build

After cloning the repo, run `npm install` to install all the dependencies.
Use `npm run dev` to run the app in development mode.
Use `npm run build` to build the app for production. It will be generated in the `build` folder and will follow your system's architecture.

# How To Run
Use the Generated Executable to run the app. It has the following arguments:

- `--help`: Show the help menu.
- `-t` or `--target`: URL of the local server to forward the requests to.
- `-l` or `--logger`: Use the logger.
- `-lp` or `--log_prefix`: Use the logger prefix.
- `-lt` or `--log_time`: Use the logger time.
- `--version`: Show the version.

# Example Run
```bash
abacate -t "http://localhost:3000"
```