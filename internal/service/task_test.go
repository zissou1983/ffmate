package service

import (
	"testing"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTaskTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Task{}, &model.Webhook{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	s := sev.New("test", "", "", 3000)
	s.SetDB(db)
	Init(s)

	return db, s
}

func TestTaskService(t *testing.T) {
	db, _ := setupTaskTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	t.Run("Create and find task", func(t *testing.T) {
		newTask := &dto.NewTask{
			InputFile:  "/test/input.mp4",
			OutputFile: "/test/output.mp4",
			Command:    "ffmpeg -i ${INPUT_FILE} ${OUTPUT_FILE}",
			Priority:   1,
		}

		task, err := TaskService().NewTask(newTask, "", "test")
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		found, err := TaskService().GetTaskByUuid(task.Uuid)
		if err != nil {
			t.Fatalf("Failed to find task: %v", err)
		}

		if found.Command.Raw != newTask.Command {
			t.Errorf("Expected command %s, got %s", newTask.Command, found.Command)
		}
	})

	t.Run("List tasks", func(t *testing.T) {
		tasks, total, err := TaskService().ListTasks(0, 10, "")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if total < 1 {
			t.Error("Expected at least one task")
		}

		if len(*tasks) < 1 {
			t.Error("Expected at least one task in list")
		}
	})

	t.Run("Cancel task", func(t *testing.T) {
		newTask := &dto.NewTask{
			InputFile:  "/test/input.mp4",
			OutputFile: "/test/output.mp4",
			Command:    "test",
		}

		task, err := TaskService().NewTask(newTask, "", "test")
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		cancelled, err := TaskService().CancelTask(task.Uuid)
		if err != nil {
			t.Fatalf("Failed to cancel task: %v", err)
		}

		if cancelled.Status != dto.DONE_CANCELED {
			t.Errorf("Expected status %s, got %s", dto.DONE_CANCELED, cancelled.Status)
		}
	})

	t.Run("Delete task", func(t *testing.T) {
		newTask := &dto.NewTask{
			InputFile:  "/test/input.mp4",
			OutputFile: "/test/output.mp4",
			Command:    "test",
		}

		task, err := TaskService().NewTask(newTask, "", "test")
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		err = TaskService().DeleteTask(task.Uuid)
		if err != nil {
			t.Fatalf("Failed to delete task: %v", err)
		}

		_, err = TaskService().GetTaskByUuid(task.Uuid)
		if err == nil {
			t.Error("Expected error when finding deleted task")
		}
	})
}
