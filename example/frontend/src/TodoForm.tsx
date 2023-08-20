import React, { useState } from "react";
import { Todo } from "./interfaces";


interface TodoFormProps {
  onNewTodo: (todo: Todo) => void;
}

const TodoForm: React.FC<TodoFormProps> = ({ onNewTodo }) => {
  const [title, setTitle] = useState<string>("");
  const [description, setDescription] = useState<string>("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const response = await fetch("/todos/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, description }),
    });

    const data: Todo = await response.json();
    onNewTodo(data);
    setTitle("");
    setDescription("");
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Title"
      />
      <input
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="Description"
      />
      <button type="submit">Add</button>
    </form>
  );
}

export default TodoForm;
