package tasks

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/welovemedia/ffmate/e2e/tests"
	"github.com/welovemedia/ffmate/internal/dto"
)

var task dto.Task
var tasks []dto.Task

func RunTests() error {
	webhookServer, webhookChan := tests.SetupWebhookServer(4)

	tests.RegisterWebhook(webhookServer, dto.TASK_CREATED)
	tests.RegisterWebhook(webhookServer, dto.TASK_UPDATED)
	tests.RegisterWebhook(webhookServer, dto.TASK_DELETED)

	if err := testTaskCreation(); err != nil {
		return err
	}
	tests.WaitForWebhook(webhookChan, dto.TASK_CREATED)
	tests.WaitForWebhook(webhookChan, dto.TASK_UPDATED) // queue -> running
	tests.WaitForWebhook(webhookChan, dto.TASK_UPDATED) // running - failed

	if err := testTaskList(); err != nil {
		return err
	}

	if err := testTaskDeletion(); err != nil {
		return err
	}
	tests.WaitForWebhook(webhookChan, dto.TASK_DELETED)

	return nil
}

func testTaskCreation() error {
	resp, err := resty.New().R().
		SetBody(&dto.NewTask{
			Name:    "Test Task",
			Command: "-i {input} -c copy {output}",
		}).
		SetResult(&task).
		Post("http://localhost:3000/api/v1/tasks")

	if err != nil || resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to create task: %v", resp)
	} else {
		fmt.Printf("Created '%s' (uuid: %s)\n", task.Name, task.Uuid)
	}

	return nil
}

func testTaskList() error {
	resp, err := resty.New().R().
		SetResult(&tasks).
		Get("http://localhost:3000/api/v1/tasks")

	if err != nil || resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to list tasks: %v", resp)
	} else {
		fmt.Printf("Listed '%d' tasks\n", len(tasks))
	}

	return nil
}

func testTaskDeletion() error {
	resp, err := resty.New().R().
		Delete("http://localhost:3000/api/v1/tasks/" + task.Uuid)

	if err != nil || resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to delete task: %v", resp)
	} else {
		fmt.Printf("Deleted '%s' (uuid: %s)\n", task.Name, task.Uuid)
	}

	return nil
}
