package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func serve(path string, x []string) error {
	if len(x) != 1 {
		return fmt.Errorf("serve requires 1 argument: <host:port>")
	}

	hostPort := x[0]

	// Validate that the journal exists
	if err := validate(path); err != nil {
		return err
	}

	server := &todoServer{
		path: path,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", server.serveHTML)
	mux.HandleFunc("GET /api/todos", server.getTodos)
	mux.HandleFunc("POST /api/todos", server.addTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", server.removeTodo)
	mux.HandleFunc("POST /api/sync", server.doSync)

	fmt.Printf("Starting todo server at http://%s\n", hostPort)
	return http.ListenAndServe(hostPort, mux)
}

type todoServer struct {
	path     string
	todoPath string
}

func (ts *todoServer) serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, htmlTemplate)
}

func (ts *todoServer) getTodos(w http.ResponseWriter, r *http.Request) {
	t, err := todoLoad(ts.path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t.t)
}

func (ts *todoServer) addTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		http.Error(w, "todo text cannot be empty", http.StatusBadRequest)
		return
	}

	if err := todoAdd(ts.path, []string{text}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (ts *todoServer) removeTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := todoComplete(ts.path, []string{idStr}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (ts *todoServer) doSync(w http.ResponseWriter, r *http.Request) {
	if err := doGitPull(ts.path); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err := doGitPush(ts.path); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "synced"})
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Todo List</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
			background-color: #f5f5f5;
			padding: 20px;
		}

		.container {
			max-width: 600px;
			margin: 0 auto;
			background: white;
			border-radius: 8px;
			box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			padding: 20px;
		}

		h1 {
			color: #333;
			margin-bottom: 20px;
			font-size: 24px;
		}

		.input-group {
			display: flex;
			gap: 10px;
			margin-bottom: 20px;
		}

		input[type="text"] {
			flex: 1;
			padding: 10px;
			border: 1px solid #ddd;
			border-radius: 4px;
			font-size: 14px;
		}

		input[type="text"]:focus {
			outline: none;
			border-color: #4CAF50;
		}

		button {
			padding: 10px 20px;
			background-color: #4CAF50;
			color: white;
			border: none;
			border-radius: 4px;
			cursor: pointer;
			font-size: 14px;
		}

		button:hover {
			background-color: #45a049;
		}

		button.sync {
			background-color: #2196F3;
			flex: 0 0 auto;
		}

		button.sync:hover {
			background-color: #0b7dda;
		}

		.todo-list {
			list-style: none;
		}

		.todo-item {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 12px;
			border-bottom: 1px solid #eee;
		}

		.todo-item:last-child {
			border-bottom: none;
		}

		.todo-text {
			flex: 1;
			color: #333;
		}

		.todo-delete {
			padding: 6px 12px;
			background-color: #f44336;
			font-size: 12px;
		}

		.todo-delete:hover {
			background-color: #da190b;
		}

		.empty {
			text-align: center;
			color: #999;
			padding: 40px 20px;
		}

		.status {
			text-align: center;
			padding: 10px;
			margin-bottom: 10px;
			border-radius: 4px;
			display: none;
		}

		.status.show {
			display: block;
		}

		.status.success {
			background-color: #d4edda;
			color: #155724;
			border: 1px solid #c3e6cb;
		}

		.status.error {
			background-color: #f8d7da;
			color: #721c24;
			border: 1px solid #f5c6cb;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Todo List</h1>
		<div id="status" class="status"></div>
		<div class="input-group">
			<input type="text" id="todoInput" placeholder="Add a new todo..." />
			<button onclick="addTodo()">Add</button>
			<button class="sync" onclick="syncTodos()">Sync</button>
		</div>
		<ul id="todoList" class="todo-list">
			<li class="empty">Loading...</li>
		</ul>
	</div>

	<script>
		async function loadTodos() {
			try {
				const res = await fetch('/api/todos');
				if (!res.ok) throw new Error('Failed to load todos');
				const todos = await res.json();
				renderTodos(todos);
			} catch (err) {
				showStatus('Failed to load todos', 'error');
			}
		}

		function renderTodos(todos) {
			const list = document.getElementById('todoList');
			if (!todos || todos.length === 0) {
				list.innerHTML = '<li class="empty">No todos yet</li>';
				return;
			}
			list.innerHTML = todos.map((todo, idx) => {
				return '<li class="todo-item"><span class="todo-text">' + 
					escapeHtml(todo) + 
					'</span><button class="todo-delete" onclick="removeTodo(' + idx + ')">Remove</button></li>';
			}).join('');
		}

		async function addTodo() {
			const input = document.getElementById('todoInput');
			const text = input.value.trim();
			if (!text) return;

			try {
				const res = await fetch('/api/todos', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ text })
				});
				if (!res.ok) throw new Error('Failed to add todo');
				input.value = '';
				showStatus('Todo added', 'success');
				loadTodos();
			} catch (err) {
				showStatus('Failed to add todo', 'error');
			}
		}

		async function removeTodo(id) {
			try {
				const res = await fetch('/api/todos/' + id, {
					method: 'DELETE'
				});
				if (!res.ok) throw new Error('Failed to remove todo');
				showStatus('Todo removed', 'success');
				loadTodos();
			} catch (err) {
				showStatus('Failed to remove todo', 'error');
			}
		}

		async function syncTodos() {
			try {
				const res = await fetch('/api/sync', {
					method: 'POST'
				});
				if (!res.ok) throw new Error('Failed to sync');
				showStatus('Synced successfully', 'success');
				loadTodos();
			} catch (err) {
				showStatus('Failed to sync', 'error');
			}
		}

		function showStatus(msg, type) {
			const status = document.getElementById('status');
			status.textContent = msg;
			status.className = 'status show ' + type;
			setTimeout(() => status.classList.remove('show'), 3000);
		}

		function escapeHtml(text) {
			const div = document.createElement('div');
			div.textContent = text;
			return div.innerHTML;
		}

		document.getElementById('todoInput').addEventListener('keypress', (e) => {
			if (e.key === 'Enter') addTodo();
		});

		loadTodos();
	</script>
</body>
</html>
`
