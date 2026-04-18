# Fix Failed Task Persistence in Monitor Plan

**Goal:** Ensure that when a user dismisses (clicks "X") a failed task in the "Task Execution Monitor", it stays hidden and doesn't reappear on the next 5-second poll.

**Changes:**

### 1. Backend API (`internal/api/router.go`)
- **Add Endpoint**: `POST /api/tasks/:id/dismiss`.
  - Logic: Update the task's `stage` field to `"Dismissed"` in the database.
- **Update Query**: In `getDashboardStats`, update the `runningTasksList` query.
  - Filter: Exclude tasks where `status = "failed" AND stage = "Dismissed"`.
  - This ensures that only "new" or "unacknowledged" failures are shown.

### 2. Frontend API (`web/src/api/task.js`)
- Add `dismissTask(taskId)` function to call the new backend endpoint.

### 3. Frontend UI (`web/src/views/Dashboard.vue`)
- Update the `dismissTask` function:
  - It should now be `async`.
  - Call the `dismissTask` API before removing the task from the local `runningTasks` array.
  - This persists the "hidden" state to the database.

**Re-appearance Logic:** 
- No changes needed to `worker.go`. When a task is re-run, the worker naturally calls `updateProgress` which sets `stage = "Started"`, automatically clearing the `"Dismissed"` state and making it visible again if it fails or runs.
