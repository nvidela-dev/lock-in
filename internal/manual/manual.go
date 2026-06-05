package manual

const Text = `Lock-In Command Manual

Normal commands
  a          Add task to active project
  s          Add subtask to selected task
  j, down    Move down one visible task
  k, up      Move up one visible task
  [          Switch to previous project
  ]          Switch to next project
  h, left    Collapse selected task
  l, right   Expand selected task
  d          Mark selected item done; confirms cascade when it has children
  r          Mark selected item ready
  p          Mark selected item in progress
  X          Delete selected task or subtask after y/n confirmation
  c          Create project
  e          Rename active project
  x          Delete active project after y/n confirmation
  ?          Open this manual
  esc        Cancel prompt or close manual
  q, ctrl+c  Quit

Input commands
  G          Prompt for task number, then jump to it
  g          Prompt for project number, then switch to it
  D          Prompt for task number, then mark it done; confirms cascade when it has children
  R          Prompt for task number, then mark it ready
  P          Prompt for task number, then mark it in progress
  S          Prompt for task number, then prompt for a subtask title

Task numbers
  Numbers are display positions inside the active project: 1, 2, 3.1, 3.2.
  Hidden subtasks can still be targeted by number. Jumping to one expands its ancestors.

Project numbers
  Project numbers are display positions in the bottom bar.
  Use [ and ] for adjacent project navigation, or g to jump by number.
`

const Usage = `Usage:
  lock-in        Start the TUI
  lock-in man    Print the command manual
  lock-in --help Print this help

Data:
  Uses os.UserConfigDir()/lock-in/state.json by default.
  Set LOCK_IN_DATA=/path/state.json to use a custom state file.
`
