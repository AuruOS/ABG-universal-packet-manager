<div align="center">
  <img src="abg-logo.svg" height="120">
  <h1 align="center">abg</h1>

  <p align="center">abg is a universal and secure package manager designed for AuruOS. It integrates with subsystems, supports Android apps, source builds, and multiple package managers—all inside a secure environment.</p>
  <small>We doing secure enviroment by Distrobox but in feature we doing our</small>
</div>

<br/>

## Help

```bash
abg is a flexible and secure package manager with subsystem support.

Usage:
  abg [enter your command]

Available Commands:
  [subsystem]    Work within a specific subsystem with its own environment
  android        Install or manage Android applications
  build          Build and install software from source
  pkgmanagers    Use a specific package manager (e.g. apt, dnf, pacman, etc.)
  subsystems     Manage isolated environments (create, enter, delete)
  completion     Generate autocompletion script for your shell

Flags:
  -v, --version  Print the version of abg
