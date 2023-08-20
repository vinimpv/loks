import React from "react";
import { Todo } from "./interfaces";

interface TodoListProps {
    todos: Todo[];
}


const TodoList: React.FC<TodoListProps> = ({ todos }) => {
  return (
    <ul>
      {todos.map((todo) => (
        <li key={todo.id}>
          {todo.title} - {todo.description}
        </li>
      ))}
    </ul>
  );
}

export default TodoList;
