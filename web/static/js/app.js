document.addEventListener('DOMContentLoaded', () => {
    // Current user state
    let currentUser = null;

    // Highlight active nav link
    const path = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav-link');
    navLinks.forEach(link => {
        if (link.getAttribute('href') === path) {
            link.classList.add('active');
        }
    });

    // Logout handler
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            await fetch('/api/auth/logout', { method: 'POST' });
            window.location.href = '/login';
        });
    }

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

    if (path === '/tasks') {
        fetchTasks();
    }

    async function fetchUserInfo() {
        const userNameElement = document.getElementById('user-name');
        if (userNameElement) {
            userNameElement.textContent = "Production User";
        }
    }

    async function fetchDashboardStats() {
        try {
            const [projResp, taskResp] = await Promise.all([
                fetch('/api/projects/stats'),
                fetch('/api/tasks/stats')
            ]);
            
            const projStats = await projResp.json();
            const taskStats = await taskResp.json();

            document.getElementById('total-projects-stat').textContent = projStats.total_projects;
            document.getElementById('total-tasks-stat').textContent = taskStats.total_tasks;

            // Initialize Chart if on dashboard
            const chartCtx = document.getElementById('taskChart');
            if (chartCtx) {
                const statuses = taskStats.statuses || {};
                new Chart(chartCtx, {
                    type: 'doughnut',
                    data: {
                        labels: ['Todo', 'In Progress', 'Review', 'Done'],
                        datasets: [{
                            data: [
                                statuses['todo'] || 0,
                                statuses['in-progress'] || 0,
                                statuses['review'] || 0,
                                statuses['done'] || 0
                            ],
                            backgroundColor: ['#94a3b8', '#3b82f6', '#f59e0b', '#10b981'],
                            borderWidth: 0,
                            borderRadius: 4
                        }]
                    },
                    options: {
                        plugins: { legend: { position: 'bottom', labels: { color: '#94a3b8' } } },
                        cutout: '70%'
                    }
                });
            }
            // Fetch Activities
            fetchActivities();
        } catch (err) {
            console.error('Failed to fetch dashboard stats:', err);
        }
    }

    async function fetchActivities() {
        const container = document.getElementById('activity-feed');
        if (!container) return;

        try {
            const resp = await fetch('/api/activities');
            const data = await resp.json();
            
            container.innerHTML = '';
            if (!data || data.length === 0) {
                container.innerHTML = '<p style="color: var(--text-muted); padding: 20px; text-align: center">No recent activity.</p>';
                return;
            }

            data.forEach(act => {
                const item = document.createElement('div');
                item.style.padding = '12px';
                item.style.borderBottom = '1px solid var(--glass-border)';
                item.style.display = 'flex';
                item.style.gap = '12px';
                item.style.alignItems = 'center';

                const iconMap = { 'project': 'fa-rocket', 'task': 'fa-check-circle' };
                const colorMap = { 'create': 'var(--success)', 'update': 'var(--primary)', 'delete': 'var(--danger)' };

                item.innerHTML = `
                    <div style="width: 32px; height: 32px; border-radius: 8px; background: rgba(255,255,255,0.03); display: flex; align-items: center; justify-content: center; color: ${colorMap[act.action] || 'white'}">
                        <i class="fas ${iconMap[act.entity_type] || 'fa-info-circle'}"></i>
                    </div>
                    <div>
                        <div style="font-size: 13px; font-weight: 500">${act.details}</div>
                        <div style="font-size: 11px; color: var(--text-muted)">${new Date(act.created_at).toLocaleString()}</div>
                    </div>
                `;
                container.appendChild(item);
            });
        } catch (err) {
            console.error('Failed to fetch activities:', err);
        }
    }

    // Projects Logic
    async function fetchProjects() {
        const projectList = document.getElementById('project-list');
        if (!projectList) return;

        try {
            const resp = await fetch('/api/projects');
            const data = await resp.json();
            
            projectList.innerHTML = '';
            if (!data || data.length === 0) {
                projectList.innerHTML = '<div style="grid-column: 1/-1; text-align: center; padding: 40px; color: var(--text-muted)">No projects yet. Create one to get started!</div>';
                return;
            }

            data.forEach(p => {
                const card = document.createElement('div');
                card.className = 'glass glass-card project-card';
                card.innerHTML = `
                    <div style="display: flex; justify-content: space-between; align-items: start; margin-bottom: 16px">
                        <div style="width: 48px; height: 48px; border-radius: 12px; background: rgba(79, 70, 229, 0.1); border: 1px solid var(--primary); display: flex; align-items: center; justify-content: center; color: var(--primary)">
                            <i class="fas fa-rocket"></i>
                        </div>
                        <button onclick="deleteProject('${p.id}')" class="btn btn-ghost" style="color: var(--danger); padding: 4px">
                            <i class="fas fa-trash-alt"></i>
                        </button>
                    </div>
                    <h3 style="margin-bottom: 8px">${p.name}</h3>
                    <p style="color: var(--text-muted); font-size: 14px; margin-bottom: 24px; min-height: 40px">${p.description || 'No description provided.'}</p>
                    <div style="display: flex; justify-content: space-between; align-items: center">
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 65%"></div>
                        </div>
                        <span style="font-size: 12px; font-weight: 600">65%</span>
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
    window.updateTaskStatus = async (id, newStatus) => {
        try {
            const resp = await fetch(`/api/tasks/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: newStatus })
            });
            if (resp.ok) {
                showToast(`Status updated to ${newStatus}`);
                fetchTasks();
                fetchDashboardStats();
            } else {
                showToast('Failed to update status', 'error');
            }
        } catch (err) { showToast('Error connecting to server', 'error'); }
    };

    window.deleteProject = async (id) => {
        if (!confirm('Are you sure you want to delete this project?')) return;
        await fetch(`/api/projects/${id}`, { method: 'DELETE' });
        fetchProjects();
        showToast('Project deleted');
    };

    window.deleteTask = async (id) => {
        if (!confirm('Are you sure you want to delete this task?')) return;
        await fetch(`/api/tasks/${id}`, { method: 'DELETE' });
        fetchTasks();
        showToast('Task deleted');
    };

    async function fetchTasks() {
        const columns = {
            'todo': document.getElementById('col-todo'),
            'in-progress': document.getElementById('col-in-progress'),
            'review': document.getElementById('col-review'),
            'done': document.getElementById('col-done')
        };
        
        if (!columns.todo) return;

        // Populate dropdowns
        const projDropdown = document.getElementById('task-project-id');
        if (projDropdown && projDropdown.options.length <= 1) {
            const resp = await fetch('/api/projects');
            const data = await resp.json();
            data.forEach(p => {
                const opt = document.createElement('option');
                opt.value = p.id;
                opt.textContent = p.name;
                projDropdown.appendChild(opt);
            });
        }

        const userDropdown = document.getElementById('task-assigned-to');
        if (userDropdown && userDropdown.options.length <= 1) {
            const resp = await fetch('/api/users');
            const data = await resp.json();
            data.forEach(u => {
                const opt = document.createElement('option');
                opt.value = u.id;
                opt.textContent = u.full_name || u.username;
                userDropdown.appendChild(opt);
            });
        }

        try {
            const resp = await fetch('/api/tasks');
            const tasks = await resp.json();

            Object.values(columns).forEach(col => {
                const container = col.querySelector('.task-container');
                if (container) container.innerHTML = '';
            });

            if (!tasks || tasks.length === 0) return;

            tasks.forEach(t => {
                const status = t.status || 'todo';
                const container = columns[status]?.querySelector('.task-container');
                if (container) {
                    const nextStatusMap = { 'todo': 'in-progress', 'in-progress': 'review', 'review': 'done' };
                    const nextStatus = nextStatusMap[status];

                    const taskEl = document.createElement('div');
                    taskEl.className = 'glass glass-card';
                    taskEl.style.padding = '16px';
                    taskEl.style.marginBottom = '12px';
                    taskEl.style.position = 'relative';
                    taskEl.innerHTML = `
                        <button onclick="deleteTask('${t.id}')" style="position: absolute; top: 8px; right: 8px; background: none; border: none; color: var(--danger); cursor: pointer; padding: 4px; font-size: 12px; opacity: 0.6; transition: opacity 0.2s">
                            <i class="fas fa-trash-alt"></i>
                        </button>
                        <div style="font-size: 10px; color: var(--primary); margin-bottom: 4px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em">${t.priority}</div>
                        <p style="font-weight: 500; font-size: 14px; margin-right: 20px; line-height: 1.4">${t.title}</p>
                        <div style="margin-top: 16px; display: flex; justify-content: space-between; align-items: center">
                            <div style="display: flex; gap: 8px">
                                ${nextStatus ? `
                                    <button onclick="updateTaskStatus('${t.id}', '${nextStatus}')" class="btn btn-ghost" style="font-size: 11px; padding: 4px 8px; background: rgba(255,255,255,0.03)">
                                        Next <i class="fas fa-arrow-right" style="font-size: 9px; margin-left: 4px"></i>
                                    </button>
                                ` : ''}
                            </div>
                            <div title="Assigned To" style="width: 24px; height: 24px; border-radius: 50%; background: var(--primary); display: flex; align-items: center; justify-content: center; font-size: 10px; color: white; font-weight: bold; box-shadow: 0 0 10px rgba(79, 70, 229, 0.3)">
                                ${t.assigned_to ? 'U' : '?'}
                            </div>
                        </div>
                    `;
                    taskEl.onmouseover = () => { taskEl.querySelector('button').style.opacity = '1'; };
                    taskEl.onmouseout = () => { taskEl.querySelector('button').style.opacity = '0.6'; };
                    container.appendChild(taskEl);
                }
            });
        } catch (err) {
            console.error('Failed to fetch tasks:', err);
        }
    }

    window.saveTask = async () => {
        const project_id = document.getElementById('task-project-id').value;
        const assigned_to = document.getElementById('task-assigned-to').value;
        const title = document.getElementById('task-title').value;
        const description = document.getElementById('task-desc').value;
        const status = document.getElementById('task-status').value;
        const priority = document.getElementById('task-priority').value;

        if (!project_id) {
            showToast('Please select a project', 'error');
            return;
        }
        if (!title) {
            showToast('Task title is required', 'error');
            return;
        }

        try {
            const resp = await fetch('/api/tasks', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ project_id, assigned_to, title, description, status, priority })
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
