    // Fetch data based on path
    if (path === '/' || path === '/projects' || path === '/tasks') {
        fetchUserInfo();
    }

    if (path === '/') {
        fetchDashboardStats();
    }

    if (path === '/projects') {
        fetchProjects();
    }
...
    async function fetchUserInfo() {
        const userNameElement = document.getElementById('user-name');
        if (userNameElement) {
            userNameElement.textContent = "Production User";
        }
    }

    async function fetchDashboardStats() {
        try {
            const resp = await fetch('/api/projects/stats');
            const data = await resp.json();
            const projectCountElem = document.getElementById('total-projects-stat');
            if (projectCountElem && data.total_projects !== undefined) {
                projectCountElem.textContent = data.total_projects;
            }
        } catch (err) {
            console.error('Failed to fetch stats:', err);
        }
    }

    // Projects Logic
    async function fetchProjects() {
        const projectList = document.getElementById('project-list');
        if (!projectList) return;

        try {
            const resp = await fetch('/api/projects');
            const projects = await resp.json();
            
            projectList.innerHTML = '';
            if (projects.length === 0) {
                projectList.innerHTML = '<p class="glass glass-card" style="padding: 20px; grid-column: 1/-1; text-align: center">No projects yet. Create your first one!</p>';
                return;
            }

            projects.forEach(p => {
                const card = document.createElement('div');
                card.className = 'glass glass-card';
                card.innerHTML = `
                    <div style="display: flex; justify-content: space-between; margin-bottom: 16px">
                        <span style="background: rgba(79, 70, 229, 0.2); color: #818cf8; padding: 4px 8px; border-radius: 6px; font-size: 12px">Active</span>
                        <i class="fas fa-ellipsis-v" style="color: var(--text-muted); cursor: pointer"></i>
                    </div>
                    <h3 style="margin-bottom: 8px">${p.name}</h3>
                    <p style="color: var(--text-muted); font-size: 14px; margin-bottom: 24px">${p.description || 'No description provided.'}</p>
                    <div style="display: flex; justify-content: space-between; align-items: center">
                        <div style="display: flex;">
                             <div style="width: 28px; height: 28px; border-radius: 50%; background: #4f46e5; border: 2px solid var(--bg-dark)"></div>
                        </div>
                        <span style="font-size: 12px; color: var(--text-muted)">Project Overview</span>
                    </div>
                `;
                projectList.appendChild(card);
            });
        } catch (err) {
            console.error('Failed to fetch projects:', err);
        }
    }

    window.saveProject = async () => {
        const name = document.getElementById('project-name').value;
        const description = document.getElementById('project-desc').value;

        if (!name) {
            showToast('Project name is required', 'error');
            return;
        }

        try {
            const resp = await fetch('/api/projects', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, description })
            });

            if (resp.ok) {
                showToast('Project created successfully');
                hideModal('project-modal');
                fetchProjects();
            } else {
                showToast('Failed to create project', 'error');
            }
        } catch (err) {
            showToast('Error connecting to server', 'error');
        }
    };

    // Tasks Logic
    async function fetchTasks() {
        // Simple Kanban populate logic
        const columns = {
            'todo': document.getElementById('col-todo'),
            'in-progress': document.getElementById('col-in-progress'),
            'review': document.getElementById('col-review'),
            'done': document.getElementById('col-done')
        };
        
        if (!columns.todo) return;

        try {
            const resp = await fetch('/api/tasks');
            const tasks = await resp.json();

            // Clear columns
            Object.values(columns).forEach(col => {
                const container = col.querySelector('.task-container');
                if (container) container.innerHTML = '';
            });

            tasks.forEach(t => {
                const status = t.status || 'todo';
                const container = columns[status]?.querySelector('.task-container');
                if (container) {
                    const taskEl = document.createElement('div');
                    taskEl.className = 'glass glass-card';
                    taskEl.style.padding = '16px';
                    taskEl.style.marginBottom = '12px';
                    taskEl.innerHTML = `
                        <div style="font-size: 12px; color: var(--primary); margin-bottom: 8px; font-weight: 600">${t.priority}</div>
                        <p style="font-weight: 500; font-size: 14px">${t.title}</p>
                        <div style="margin-top: 16px; display: flex; justify-content: space-between; align-items: center">
                            <i class="fas fa-tasks" style="font-size: 12px; color: var(--text-muted)"></i>
                            <div style="width: 24px; height: 24px; border-radius: 50%; background: #4f46e5"></div>
                        </div>
                    `;
                    container.appendChild(taskEl);
                }
            });
        } catch (err) {
            console.error('Failed to fetch tasks:', err);
        }
    }

    window.saveTask = async () => {
        const title = document.getElementById('task-title').value;
        const description = document.getElementById('task-desc').value;
        const status = document.getElementById('task-status').value;
        const priority = document.getElementById('task-priority').value;

        if (!title) {
            showToast('Task title is required', 'error');
            return;
        }

        try {
            const resp = await fetch('/api/tasks', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title, description, status, priority })
            });

            if (resp.ok) {
                showToast('Task created successfully');
                hideModal('task-modal');
                fetchTasks();
            } else {
                showToast('Failed to create task', 'error');
            }
        } catch (err) {
            showToast('Error connecting to server', 'error');
        }
    };

    // Modal Helpers
    window.showModal = (id) => {
        const modal = document.getElementById(id);
        if (modal) modal.style.display = 'flex';
    };

    window.hideModal = (id) => {
        const modal = document.getElementById(id);
        if (modal) modal.style.display = 'none';
    };

    // Toast Notification System
    window.showToast = (message, type = 'success') => {
        const toast = document.createElement('div');
        toast.className = `glass glass-card toast toast-${type}`;
        toast.innerHTML = message;
        toast.style.position = 'fixed';
        toast.style.bottom = '24px';
        toast.style.right = '24px';
        toast.style.zIndex = '1000';
        toast.style.padding = '12px 24px';
        toast.style.borderLeft = `4px solid var(--${type})`;
        
        document.body.appendChild(toast);
        
        setTimeout(() => {
            toast.style.opacity = '0';
            setTimeout(() => toast.remove(), 500);
        }, 3000);
    };
});
