import { useState, useEffect } from "react";
import { Events, Window } from "@wailsio/runtime";

import { TodoItem } from "./todo_item";
import { TodoList } from "./todo_list";

import type { Todo } from "../bindings/todo/internal/model";
import type { Priority } from "./types";

import { TodoService, SettingService, WindowService } from "../bindings/todo";



import "./App.css";

export default function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [content, setContent] = useState<string>("");

  const [editingID, setEditingID] = useState<string | null>(null);
  const [draft, setDraft] = useState<string>("");

  const [showCompleted, setShowCompleted] = useState("false");
  const [pinned, setPinned] = useState<boolean | null>(null);
  const [isSavingPinned, setIsSavingPinned] = useState(false);
  const [locked, setLocked] = useState<boolean | null>(null);
  const [isSavingLocked, setIsSavingLocked] = useState(false);
  const [autostartEnabled, setAutostartEnabled] = useState<boolean | null>(null);
  const [isSavingAutostart, setIsSavingAutostart] = useState(false);

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
    fetchWindowState().catch(console.error);
    fetchAutostartState().catch(console.error);
  }, []);

  useEffect(() => {
    const unsubscribe = Events.On("autostart-changed", (event) => {
      if (typeof event.data === "boolean") {
        setAutostartEnabled(event.data);
      }
    });

    return unsubscribe;
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

  async function toggleTodo(todo: Todo) {
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

  function startEditing(todo: Todo) {
    setEditingID(todo.id);
    setDraft(todo.content);
  }

  function cancelEditing() {
    setEditingID(null);
    setDraft("");
  }

  async function saveEditing(todo: Todo) {
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

  async function fetchWindowState() {
    const state = await WindowService.LoadWindowSettings();
    setPinned(state?.pinned ?? false);
    setLocked(state?.locked ?? false);
  }

  async function togglePinned() {
    if (pinned === null || isSavingPinned) return;

    const nextPinned = !pinned;
    setIsSavingPinned(true);
    try {
      await WindowService.SetPinned(nextPinned);
      setPinned(nextPinned);
    } catch (error) {
      console.error(error);
      alert("更新置顶状态失败，请重试");
    } finally {
      setIsSavingPinned(false);
    }
  }

  async function toggleLocked() {
    if (locked === null || isSavingLocked) return;

    const nextLocked = !locked;
    setIsSavingLocked(true);
    try {
      await WindowService.SetLocked(nextLocked);
      setLocked(nextLocked);
    } catch (error) {
      console.error(error);
      alert("更新位置锁定状态失败，请重试");
    } finally {
      setIsSavingLocked(false);
    }
  }

  async function fetchAutostartState() {
    const enabled = await WindowService.IsAutostartEnabled();
    setAutostartEnabled(enabled);
  }

  async function toggleAutostart() {
    if (autostartEnabled === null || isSavingAutostart) return;

    const nextEnabled = !autostartEnabled;
    setIsSavingAutostart(true);

    try {
      await WindowService.SetAutostartEnabled(nextEnabled);
      setAutostartEnabled(nextEnabled);
    } catch (error) {
      console.error("更新开机启动状态失败", error);
      alert("更新开机启动状态失败，请重试");
    } finally {
      setIsSavingAutostart(false);
    }
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

  async function minimiseWindow() {
    try {
      await Window.Minimise();
    } catch (error) {
      console.error("最小化窗口失败", error);
    }
  }

  async function toggleFullscreen() {
    try {
      await Window.ToggleFullscreen();
    } catch (error) {
      console.error("切换全屏失败", error);
    }
  }

  async function closeWindow() {
    try {
      // 后端的 WindowClosing hook 会将此操作转换为隐藏到系统托盘。
      await Window.Close();
    } catch (error) {
      console.error("关闭窗口失败", error);
    }
  }

  return (
    <main className="sticky-note">
      <header className={`note-header${locked ? " is-position-locked" : ""}`}>
        <div className="note-heading">
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

        <div className="window-controls" aria-label="窗口控制">
          <button
            className="window-control-button"
            type="button"
            aria-label="最小化"
            title="最小化"
            onClick={minimiseWindow}
          >
            <span aria-hidden="true">—</span>
          </button>
          <button
            className="window-control-button"
            type="button"
            aria-label="切换全屏"
            title="切换全屏"
            onClick={toggleFullscreen}
          >
            <span className="fullscreen-icon" aria-hidden="true" />
          </button>
          <button
            className="window-control-button window-close-button"
            type="button"
            aria-label="关闭并隐藏到系统托盘"
            title="关闭并隐藏到系统托盘"
            onClick={closeWindow}
          >
            <span aria-hidden="true">×</span>
          </button>
        </div>

        <div className="header-actions">
          <button
            className="toggle-completed-button"
            type="button"
            onClick={toggleShowCompleted}
          >
            {showCompleted === "true" ? "隐藏已完成" : "显示已完成"}
          </button>
          <button
            className="toggle-pinned-button"
            type="button"
            disabled={pinned === null || isSavingPinned}
            aria-pressed={pinned ?? undefined}
            onClick={togglePinned}
          >
            {pinned ? "取消置顶" : "置顶显示"}
          </button>
          <button
            className="toggle-locked-button"
            type="button"
            disabled={locked === null || isSavingLocked}
            aria-pressed={locked ?? undefined}
            onClick={toggleLocked}
          >
            {locked ? "解除锁定" : "锁定位置"}
          </button>
          <button
            className="toggle-autostart-button"
            type="button"
            disabled={autostartEnabled === null || isSavingAutostart}
            aria-pressed={autostartEnabled ?? undefined}
            onClick={toggleAutostart}
          >
            {autostartEnabled ? "关闭开机自启" : "开机自启"}
          </button>
        </div>
      </header>

      <TodoList
        noteID={todos[0]?.noteId ?? "1"}
        todos={visibleTodos}
        setTodos={setTodos}
        sortingEnabled={showCompleted === "true"}
        renderTodo={(todo) => (
          <TodoItem
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
