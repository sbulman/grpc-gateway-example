package v1

import (
	"context"
	"time"

	v1 "github.com/sbulman/grpc-gateway-example/pkg/api/v1"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// todo model
type todo struct {
	title       string
	description string
	reminder    time.Time
}

// ToDoService implements the v1.ToDoServiceServer proto interface
type ToDoService struct {
	v1.UnimplementedToDoServiceServer
	todos map[int64]todo
}

// NewToDoService creates and returns a new ToDoService
func NewToDoService() *ToDoService {
	return &ToDoService{
		todos: make(map[int64]todo),
	}
}

// Create a todo task
func (s *ToDoService) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	reminder, err := ptypes.Timestamp(req.Todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	id := int64(len(s.todos) + 1)
	s.todos[id] = todo{
		title:       req.Todo.Title,
		description: req.Todo.Description,
		reminder:    reminder,
	}

	return &v1.CreateResponse{Id: id}, nil
}

// Read a todo task
func (s *ToDoService) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	todo, ok := s.todos[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	reminder, err := ptypes.TimestampProto(todo.reminder)
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &v1.ReadResponse{
		Todo: &v1.ToDo{
			Id:          req.Id,
			Title:       todo.title,
			Description: todo.description,
			Reminder:    reminder,
		},
	}, nil
}

// Update a todo task
func (s *ToDoService) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	_, ok := s.todos[req.Todo.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	reminder, err := ptypes.Timestamp(req.Todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	s.todos[req.Todo.Id] = todo{
		title:       req.Todo.Title,
		description: req.Todo.Description,
		reminder:    reminder,
	}

	return &v1.UpdateResponse{Updated: 1}, nil
}

// Delete a todo task
func (s *ToDoService) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	_, ok := s.todos[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	delete(s.todos, req.Id)

	return &v1.DeleteResponse{Deleted: 1}, nil
}

// ReadAll todo tasks
func (s *ToDoService) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	todos := make([]*v1.ToDo, 0, len(s.todos))
	for k, v := range s.todos {
		reminder, err := ptypes.TimestampProto(v.reminder)
		if err != nil {
			return nil, status.Error(codes.Internal, codes.Internal.String())
		}
		todos = append(todos, &v1.ToDo{
			Id:          k,
			Title:       v.title,
			Description: v.description,
			Reminder:    reminder,
		})
	}

	return &v1.ReadAllResponse{Todos: todos}, nil
}
