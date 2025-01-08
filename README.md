# Overview: Task-Tracker

## Brief Description

Task-Tracker is a command-line tool designed to help developers and CLI enthusiasts efficiently track their work progress. It encourages users to stay in the terminal interface, helping them manage current and future tasks without switching to a graphical interface.

---

## Target Audience

Task-Tracker is aimed at developers and programming enthusiasts who prefer working within the terminal environment.

---

# Getting Started

## Prerequisites

- You need to have Go installed on your computer to run the project.
- Alternatively, you can request a prebuilt executable by contacting [taofik.dev@gmail.com](mailto:taofik.dev@gmail.com).

## Installation

1. Clone the repository:
   ```bash
   git clone <repository_url>
   ```
2. Change to the cloned directory:
   ```bash
   cd <repository_directory>
   ```

## Running the Project

To run the application, execute the following command in the terminal:

```bash
go run . //list the commands available
```

```bash
go run . <command>
```

---

# Usage

Below are the available commands for Task-Tracker:

### General Help

To see the list of available commands:

```bash
go run . help
```

### Task Management Commands

#### Adding a Task:

```bash
task-tracker add "Task Description"
```

Example:

```bash
task-tracker add "Complete documentation"
```

#### Listing Tasks:

- To list all tasks:
  ```bash
  task-tracker list
  ```
- To list all completed tasks:
  ```bash
  task-tracker list done
  ```
- To list all pending tasks:
  ```bash
  task-tracker list todo
  ```
- To list all tasks in progress:
  ```bash
  task-tracker list inprogress
  ```

#### Deleting a Task:

```bash
task-tracker delete <task_id>
```

Example:

```bash
task-tracker delete 1
```

#### Updating a Task:

```bash
task-tracker update <task_id> "Updated Task Description"
```

Example:

```bash
task-tracker update 1 "Revise documentation section"
```

#### Changing Task Status:

- Mark a task as in progress:
  ```bash
  task-tracker mark-in-progress <task_id>
  ```
  Example:
  ```bash
  task-tracker mark-in-progress 2
  ```
- Mark a task as completed:
  ```bash
  task-tracker mark-done <task_id>
  ```
  Example:
  ```bash
  task-tracker mark-done 3
  ```
- Mark a task as pending:
  ```bash
  task-tracker mark-todo <task_id>
  ```
  Example:
  ```bash
  task-tracker mark-todo 4
  ```

---

# Contributing

Contributions are welcome! If you'd like to contribute, please fork the repository and create a pull request with your changes.

---

# License

Task-Tracker is open-source software licensed under the GNU General Public License Version 3 or any later version. You can find the license [here](https://www.gnu.org/licenses/gpl-3.0.html).

---

# Contact

For questions, feedback, or support, feel free to reach out via email:

- [bidemi64@gmail.com](mailto:bidemi64@gmail.com)
