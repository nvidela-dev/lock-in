I'm not finding a TUI that matches exactly what I wanna do

A todo/trello app to manage my todos, should be a TUI

It will function through commands

At the bottom, I'll be able to see my projects, much like I can see my open windows by session in tmux

_____________________________________

# 1 | Task 1 | Ready

# 2 | Task 2 | Done

# 3 | Task 3 | In Progress

   |
   |_ #3.1 | Subtask 1 | In Progress

[1: Project 1 |2: Project 2 |3: Project 3]

Hey claude, in project 1, find task 3.1 and open a change request for this and mark it as progress

-------------------------------------

- Should include a man + command manual with all the commands

The commands should be of 3 types, normal or input commands,
Command List:
    Normal:
    a - Add Task to project
    s - Add Subtask to Task
    k or arrow up - Go up 1 Task in the List
    j or arrow down - Go down 1 task in the List
    h or arrow left - collapse task
    l or arrow right - expand task subtasks
    d - mark item as done
    r - mark item as ready
    p - mark item as in progress

  Requires Input: // Text inside  [] will signify Input ///
    G [3] - Go to 3
    D [3] - Mark 3 as done
    S [3] - Add Subtask to 3

This would be kinda how the TUI looks.
