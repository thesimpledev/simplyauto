# Lock to Window Feature

## Overview
Add the ability to lock the auto clicker to a specific application window. If the cursor moves off that window, clicking pauses. When the cursor returns, clicking resumes.

## UX Flow
1. User enables "Lock to Window"
2. User clicks on the target app (similar to the existing fixed-position picker)
3. App captures the window under the cursor using `WindowFromPoint()`
4. Walk up to the top-level parent with `GetAncestor()` to avoid locking to a child control
5. Store the HWND and display the window title (via `GetWindowText()`) in the UI (e.g. "Locked to: Roblox")

## Click Loop Behavior
- Before each click, get cursor position with `GetCursorPos()`
- Call `WindowFromPoint()` and walk up to top-level parent
- Compare against the stored HWND
  - **Match** — click fires normally
  - **No match** — skip the click (effectively pauses)
- Clicks resume instantly when cursor moves back over the locked window

## Edge Cases
- If the locked window is closed/crashes, check with `IsWindow()` and unlock/stop gracefully
- HWND becomes invalid — don't click into whatever replaced it

## Notes
- All APIs needed (`WindowFromPoint`, `GetAncestor`, `GetWindowText`, `IsWindow`) are in `user32.dll` which is already loaded
- Overhead is negligible — these calls are nanosecond-fast
- This pairs well with the macro system too (play recording into a specific window)
