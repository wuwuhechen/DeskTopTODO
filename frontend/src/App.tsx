import { useState, useEffect } from "react";

import { TodoItem as TodoRow } from "./todo_item";
import { TodoList } from "./todo_list";

import type { Priority, TodoItem } from "./types";

import { TodoService, SettingService } from "../bindings/todo";



import "./App.css";

export default function App() {
  const [todos, setTodos] = useState<TodoItem[]>([]);
  const [content, setContent] = useState<string>("");

  const [editingID, setEditingID] = useState<string | null>(null);
  const [draft, setDraft] = useState<string>("");

  const [showCompleted, setShowCompleted] = useState("false");

  // 完成状态只影响展示分组，不修改保存的拖拽排序值。
  // 未完成事项保持既有顺序，完成事项同样保留彼此间的既有顺序。
  const activeTodos = todos.filter((todo) => !todo.completed);
  const completedTodos = todos.filter((todo) => todo.completed);
  const visibleTodos =
    showCompleted === "true"
      ? [...activeTodos, ...completedTodos]
      : activeTodos;

  useEffect(() => {
    TodoService.List()
      .then((items) => setTodos(items ?? []))
      .catch(console.error);

    fetchShowCompleted().catch(console.error);
  }, []);

  async function addTodo() {
    const text = content.trim();
    if (!text) return;

    try {
      const todo = await TodoService.Create(text);
      if (todo === null) {
        throw new Error("创建待办未返回数据");
      }
      setTodos((current) => [...current, todo]);
      setContent("");
    } catch (error) {
      console.error(error);
    }
  }

  async function toggleTodo(todo: TodoItem) {
    try {
      await TodoService.Toggle(todo.id, !todo.completed);

      setTodos((current) =>
        current.map((item) =>
          item.id === todo.id ? { ...item, completed: !item.completed } : item,
        ),
      );
    } catch (error) {
      console.error(error);
    }
  }

  async function deleteTodo(id: string) {
    try {
      await TodoService.Delete(id);
      setTodos((current) => current.filter((todo) => todo.id !== id));
    } catch (error) {
      console.error(error);
    }
  }

  function startEditing(todo: TodoItem) {
    setEditingID(todo.id);
    setDraft(todo.content);
  }

  function cancelEditing() {
    setEditingID(null);
    setDraft("");
  }

  async function saveEditing(todo: TodoItem) {
    const text = draft.trim();
    if (!text) return;

    await TodoService.UpdateContent(todo.id, text);

    setTodos((current) =>
      current.map((item) =>
        item.id === todo.id ? { ...item, content: text } : item,
      ),
    );

    cancelEditing();
  }

  async function fetchShowCompleted() {
    try {
      const value = await SettingService.GetShowCompleted();
      setShowCompleted(value);
    } catch (error) {
      console.error(error);
    }
  }

  async function toggleShowCompleted() {
    const newValue = showCompleted === "true" ? "false" : "true";
    setShowCompleted(newValue);
    await SettingService.SetShowCompleted(newValue);
  }

  async function updatePriority(id: string, priority: Priority) {
    const previousTodos = todos;

    setTodos((current) =>
      current.map((todo) => (todo.id === id ? { ...todo, priority } : todo)),
    );
    
    try {
      await TodoService.UpdatePriority(id, priority);
    } catch (error) {
      setTodos(previousTodos);
      alert("更新优先级失败，请重试");
      return;
    }

    try {
      const latestTodos = await TodoService.List();
      setTodos(latestTodos ?? []);
    } catch (error) {
      // 优先级已保存；刷新失败时保留乐观更新，等待下次读取恢复同步。
      console.error("刷新待办列表失败", error);
    }
  }

  return (
    <main className="sticky-note">
      <header className="note-header">
        <div>
          <h1>今天要做什么</h1>
          <span>
            已完成 {todos.filter((todo) => todo.completed).length}/{todos.length}
          </span>
          {showCompleted !== "true" && (
            <span className="sorting-hint" role="status">
              显示已完成事项后可拖拽排序
            </span>
          )}
        </div>

        <button
          className="toggle-completed-button"
          type="button"
          onClick={toggleShowCompleted}
        >
          {showCompleted === "true" ? "隐藏已完成" : "显示已完成"}
        </button>
      </header>

      <TodoList
        noteID={todos[0]?.noteId ?? "1"}
        todos={visibleTodos}
        setTodos={setTodos}
        sortingEnabled={showCompleted === "true"}
        renderTodo={(todo) => (
          <TodoRow
            key={todo.id}
            todo={todo}
            isEditing={editingID === todo.id}
            draft={draft}
            onToggle={toggleTodo}
            onStartEditing={startEditing}
            onDraftChange={setDraft}
            onSave={saveEditing}
            onCancel={cancelEditing}
            onDelete={deleteTodo}
            onPriorityChange={updatePriority}
          />
        )}
      />

      <footer>
        <input
          value={content}
          placeholder="输入待办事项"
          onChange={(e) => setContent(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              addTodo();
              setContent("");
            }
          }}
        />
      </footer>
    </main>
  );
}
