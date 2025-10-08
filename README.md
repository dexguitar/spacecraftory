# Spacecraftory

To call commands from Taskfile, you need to install Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

The project uses GitHub Actions for continuous integration and delivery. Main workflows:

- **CI** (`.github/workflows/ci.yml`) - checks code on every push and pull request
  - Code linting
  - Security check
  - Automatic version extraction from Taskfile.yml
