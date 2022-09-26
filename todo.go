package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	IsDone      bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) AddTask(task string) {
	todo := item{
		Task:        task,
		IsDone:      false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) CompleteTask(index int) error {
	todosList := *t

	if index < 0 || index > len(todosList) {
		return errors.New("invalid index")
	}

	todosList[index-1].CompletedAt = time.Now()
	todosList[index-1].IsDone = true

	return nil
}

func (t *Todos) DeleteTask(index int) error {
	todosList := *t

	if index < 0 || index > len(todosList) {
		return errors.New("invalid index")
	}

	*t = append(todosList[:index-1], todosList[index:]...)

	return nil
}

func (t *Todos) LoadFile(filename string) error {
	file, err := ioutil.ReadFile(filename) // Read file using ioutil(deprecated)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) { // Provided filename/filepath does not exist
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) StoreFile(filename string) error {

	data, err := json.Marshal(t)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644) // 0644 is the write permission

}

func (t *Todos) PrintTodos() {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: gray("#")},
			{Align: simpletable.AlignCenter, Text: gray("Task")},
			{Align: simpletable.AlignCenter, Text: gray("isDone")},
			{Align: simpletable.AlignCenter, Text: gray("CreatedAt")},
			{Align: simpletable.AlignCenter, Text: gray("CompletedAt")},
		},
	}

	var cells [][]*simpletable.Cell

	for index, item := range *t {
		index++
		task := blue(item.Task)
		isDone := blue("No")
		completed := ""

		if item.IsDone {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			isDone = green("Yes")
			completed = item.CompletedAt.Format(time.RFC822)
		}

		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index)},
			{Text: task},
			{Text: isDone},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: completed},
		})
	}
	table.Body = &simpletable.Body{
		Cells: cells,
	}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos", t.countPendingTask()))},
	}}

	fmt.Println(table.String())
}

func (t *Todos) countPendingTask() int {
	cnt := 0
	for _, item := range *t {
		if !item.IsDone {
			cnt++
		}
	}

	return cnt
}
