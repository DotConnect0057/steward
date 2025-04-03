# steward

steward is a tool designed to simplify the management of Linux systems, whether standalone or in a cluster. It provides a unified solution for configuration, package management, status monitoring, and more. With a declarative approach at its core, Steward ensures consistency and reliability, while also supporting procedure-based management for tasks that require step-by-step execution.

steward is intended to simulate common package management like npm, go mod into linux management.

## How to use

**Build go package**
```
git clone https://github.com/DotConnect0057/steward.git
cd steward
go build -o st
```

### Initialize config
This command will generate config.yaml and config.json files
```
st init
```
### Apply configuration
This command will execute package installation according steward config file.
```
st apply
```

## Features Todo

- **Declarative Configuration Management**:
  - Define the desired state of your Linux systems, and Steward will ensure they match.
  - Simplifies complex configurations with a single source of truth.

- **Procedure-Based Management**:
  - Supports step-by-step execution for tasks that require procedural workflows.
  - Ideal for custom or one-off operations.

- **Cluster Management**:
  - Manage multiple Linux systems as a cluster with ease.
  - Apply configurations and monitor statuses across all nodes.

- **Package Management**:
  - Install, update, and remove packages seamlessly.
  - Supports popular Linux package managers like `apt`, `yum`, and `dnf`.