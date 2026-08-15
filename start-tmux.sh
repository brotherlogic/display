#!/bin/bash

# Ensure the 'display' session exists
if ! tmux has-session -t display 2>/dev/null; then
  # Create a new session named 'display', detached
  cd /workspaces/display
  tmux new-session -d -s display
 fi
