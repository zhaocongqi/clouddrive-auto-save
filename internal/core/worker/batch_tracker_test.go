package worker

import (
	"sync"
	"testing"
	"time"
)

func TestBatchTracker_RegisterAndReport(t *testing.T) {
	tracker := NewBatchTracker()
	tracker.RegisterBatch("batch_1", 2)

	var mu sync.Mutex
	var notified bool
	tracker.onComplete = func(results []BatchResult, totalDuration time.Duration) {
		mu.Lock()
		notified = true
		mu.Unlock()
	}

	tracker.ReportTask("batch_1", BatchResult{
		TaskName: "task1", Status: "success", Message: "ok", Files: []string{"a.mp4"}, Duration: 5 * time.Second,
	})

	mu.Lock()
	if notified {
		t.Fatal("should not notify after first task, batch has 2 tasks")
	}
	mu.Unlock()

	tracker.ReportTask("batch_1", BatchResult{
		TaskName: "task2", Status: "failed", Message: "err", Duration: 2 * time.Second,
	})

	mu.Lock()
	if !notified {
		t.Fatal("should notify after all tasks complete")
	}
	mu.Unlock()
}

func TestBatchTracker_SingleTaskBatch(t *testing.T) {
	tracker := NewBatchTracker()
	tracker.RegisterBatch("batch_2", 1)

	var results []BatchResult
	tracker.onComplete = func(r []BatchResult, _ time.Duration) {
		results = r
	}

	tracker.ReportTask("batch_2", BatchResult{
		TaskName: "only_task", Status: "success", Message: "done", Duration: 3 * time.Second,
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].TaskName != "only_task" {
		t.Errorf("expected task name 'only_task', got '%s'", results[0].TaskName)
	}
}

func TestBatchTracker_UnknownBatch(t *testing.T) {
	tracker := NewBatchTracker()
	// 不注册批次，直接上报——不应 panic
	tracker.ReportTask("nonexistent", BatchResult{
		TaskName: "ghost", Status: "success", Message: "", Duration: 0,
	})
}
