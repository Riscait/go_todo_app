package store

import (
	"errors"

	"github.com/Riscait/go_todo_app/entity"
)

var (
	Tasks = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}

	ErrNotFound = errors.New("not found")
)

type TaskStore struct {
	// 動作確認用の仮実装のためあえてexportする。
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

// Add はタスクを追加する。
func (s *TaskStore) Add(t *entity.Task) (entity.TaskID, error) {
	s.LastID++
	t.ID = s.LastID
	s.Tasks[t.ID] = t
	return t.ID, nil
}

func (s *TaskStore) All() entity.Tasks {
	tasks := make(entity.Tasks, len(s.Tasks))
	for i, t := range s.Tasks {
		tasks[i-1] = t
	}
	return tasks
}

func (ts *TaskStore) Get(id entity.TaskID) (*entity.Task, error) {
	if ts, ok := ts.Tasks[id]; ok {
		return ts, nil
	}
	return nil, ErrNotFound
}
