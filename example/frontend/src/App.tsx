import React, { useState, useEffect } from "react";
import TodoForm from "./TodoForm";
import TodoList from "./TodoList";
import { Todo } from "./interfaces";



const App: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);

  useEffect(() => {
    const fetchTodos = async () => {
      const response = await fetch("/todos/");
      const data: Todo[] = await response.json();
      setTodos(data);
    };

    fetchTodos();
  }, []);

  const handleNewTodo = (todo: Todo) => {
    setTodos([...todos, todo]);
  };

  return (
    <div className="App">
      <h1>Todo App</h1>
      <TodoForm onNewTodo={handleNewTodo} />
      <TodoList todos={todos} />
    </div>
  );
}

export default App;
