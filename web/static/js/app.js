document.addEventListener('DOMContentLoaded', () => {
    // Current user state
    let currentUser = null;
    let currentTaskId = null;

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

    if (path === '/profile') {
        fetchProfile();
    }

    // Global Search
    const searchInput = document.getElementById('global-search');
    if (searchInput) {
        let debounceTimer;
        searchInput.addEventListener('input', () => {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                const query = searchInput.value.trim();
                if (!query) return;
                
                if (path === '/projects') fetchProjects(query);
                else if (path === '/tasks') fetchTasks(query);
                else {
                    window.location.href = `/tasks?q=${encodeURIComponent(query)}`;
                }
            }, 500);
        });
        
        // Handle query param on page load
        const urlParams = new URLSearchParams(window.location.search);
        const initialQuery = urlParams.get('q');
        if (initialQuery) {
            searchInput.value = initialQuery;
        }
    }

    async function fetchUserInfo() {
        const userNameElement = document.getElementById('user-name');
        if (!userNameElement) return;

        try {
            const resp = await fetch('/api/users/profile');
            const user = await resp.json();
            userNameElement.textContent = user.full_name || user.username;
        } catch (err) {
            userNameElement.textContent = "User";
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
                container.innerHTML = `
                    <div style="text-align: center; padding: 40px; color: var(--text-muted)">
                        <i class="fas fa-stream" style="font-size: 24px; margin-bottom: 12px; opacity: 0.3"></i>
                        <p style="font-size: 13px">No recent activity to show.</p>
                    </div>`;
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
                    <div style="flex: 1">
                        <div style="font-size: 13px; font-weight: 500">
                            <span style="color: var(--primary); font-weight: 600">${act.user_name}</span> ${act.details.toLowerCase().includes(act.user_name.toLowerCase()) ? act.details.substring(act.details.indexOf(':') + 1).trim() : act.details}
                        </div>
                        <div style="font-size: 11px; color: var(--text-muted); display: flex; align-items: center; gap: 4px; margin-top: 8px">
                        <i class="far fa-calendar-alt"></i> ${window.formatDate(act.created_at)}
                    </div>
                    </div>
                `;
                container.appendChild(item);
            });
        } catch (err) {
            console.error('Failed to fetch activities:', err);
        }
    }

    // Projects Logic
    async function fetchProjects(query = '') {
        const projectList = document.getElementById('project-list');
        if (!projectList) return;

        try {
            const url = query ? `/api/projects?q=${encodeURIComponent(query)}` : '/api/projects';
            const resp = await fetch(url);
            const data = await resp.json();
            
            projectList.innerHTML = '';
            if (!data || data.length === 0) {
                projectList.innerHTML = `
                    <div style="grid-column: 1/-1; text-align: center; padding: 60px; color: var(--text-muted)">
                        <i class="fas fa-rocket" style="font-size: 48px; margin-bottom: 20px; opacity: 0.1"></i>
                        <h3>No projects found</h3>
                        <p style="margin-top: 8px">Create your first project to start managing tasks with your team.</p>
                        <button class="btn btn-primary" onclick="showModal('project-modal')" style="margin-top: 24px">
                            <i class="fas fa-plus"></i> Create Project
                        </button>
                    </div>`;
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
                        <div style="display: flex; gap: 8px; align-items: center">
                            <div style="background: rgba(16, 185, 129, 0.1); color: var(--success); padding: 4px 12px; border-radius: 20px; font-size: 11px; font-weight: 600; border: 1px solid rgba(16, 185, 129, 0.2)">
                                ${p.task_count || 0} Tasks
                            </div>
                            <button onclick="deleteProject('${p.id}')" class="btn btn-ghost" title="Delete Project" style="color: var(--danger); padding: 4px; opacity: 0.6; transition: opacity 0.2s; background: transparent; border: none; cursor: pointer">
                                <i class="fas fa-trash-alt"></i>
                            </button>
                        </div>
                    </div>
                    <h3 style="margin-bottom: 8px">${p.name}</h3>
                    <p style="color: var(--text-muted); font-size: 14px; margin-bottom: 20px; line-height: 1.5; height: 42px; overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical;">
                        ${p.description || 'No description provided.'}
                    </p>
                    <div style="display: flex; justify-content: space-between; align-items: center; border-top: 1px solid var(--glass-border); pt: 16px; margin-top: auto">
                        <span style="font-size: 12px; color: var(--text-muted)">Created ${window.formatDate(p.created_at)}</span>
                        <a href="/tasks?project_id=${p.id}" class="btn btn-primary" style="padding: 6px 14px; font-size: 12px">View Board</a>
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
        const btnId = 'save-project-btn';

        if (!name) {
            showToast('Project name is required', 'error');
            return;
        }

        setLoading(btnId, true);
        try {
            const resp = await fetch('/api/projects', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, description })
            });

            if (resp.ok) {
                showToast('Project created successfully');
                hideModal('project-modal');
                // Reset form
                document.getElementById('project-name').value = '';
                document.getElementById('project-desc').value = '';
                fetchProjects();
            } else {
                const errText = await resp.text();
                showToast(`Failed: ${errText}`, 'error');
            }
        } catch (err) {
            showToast('Server connection failed', 'error');
        } finally {
            setLoading(btnId, false);
        }
    };
    // Tasks Logic
    window.updateTaskStatus = async (id, newStatus) => {
        const btn = event?.target?.closest('button');
        if (btn) {
            btn.disabled = true;
            btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';
        }

        try {
            const resp = await fetch(`/api/tasks/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: newStatus })
            });

            if (resp.ok) {
                showToast(`Status updated to ${newStatus.replace('-', ' ')}`, 'success');
                await Promise.all([fetchTasks(), fetchDashboardStats()]);
            } else {
                const error = await resp.text();
                showToast(`Update failed: ${error}`, 'error');
            }
        } catch (err) { 
            showToast('Connection error. Please try again.', 'error'); 
        } finally {
            if (btn) {
                btn.disabled = false;
                btn.innerHTML = 'Next <i class="fas fa-arrow-right" style="font-size: 9px; margin-left: 4px"></i>';
            }
        }
    };

    async function deleteAccount() {
        if (!confirm("Are you SURE you want to delete your account? This action is irreversible.")) return;

        try {
            const resp = await fetch('/api/users/profile', { method: 'DELETE' });
            if (resp.ok) {
                showToast("Account deleted successfully.", "success");
                setTimeout(() => logout(), 1500);
            } else {
                throw new Error("Failed to delete account");
            }
        } catch (err) {
            showToast(err.message, "error");
        }
    }

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

    window.openTaskDetails = async (id) => {
        currentTaskId = id;
        try {
            const [taskResp, commentResp] = await Promise.all([
                fetch(`/api/tasks/${id}`),
                fetch(`/api/tasks/${id}/comments`)
            ]);
            
            const task = await taskResp.json();
            const comments = await commentResp.json();

            document.getElementById('detail-title').innerText = task.title;
            document.getElementById('detail-desc').innerText = task.description || 'No description provided.';
            document.getElementById('detail-project').innerText = task.project_name || 'Individual Task';
            document.getElementById('detail-assignee').textContent = task.assigned_to_name || 'Unassigned';
            const priorityEl = document.getElementById('detail-priority');
            priorityEl.textContent = task.priority.toUpperCase();
            priorityEl.className = `status-pill priority-${task.priority}`;

            renderComments(comments);
            showModal('task-details-modal');
        } catch (err) {
            showToast('Failed to load task details', 'error');
        }
    };

    function renderComments(comments) {
        const container = document.getElementById('comments-container');
        container.innerHTML = '';
        
        if (!comments || comments.length === 0) {
            container.innerHTML = '<p style="color: var(--text-muted); font-size: 13px; text-align: center; padding: 10px">No comments yet.</p>';
            return;
        }

        comments.forEach(c => {
            const div = document.createElement('div');
            div.style.marginBottom = '16px';
            div.style.padding = '12px';
            div.style.background = 'rgba(255,255,255,0.02)';
            div.style.borderRadius = '8px';
            div.innerHTML = `
                <div style="display: flex; justify-content: space-between; margin-bottom: 4px; font-size: 11px">
                    <span style="font-weight: 600; color: var(--primary)">${c.user_name}</span>
                    <span style="color: var(--text-muted)">${window.formatDate(c.created_at)}</span>
                </div>
                <div style="font-size: 13px; line-height: 1.4">${c.content}</div>
            `;
            container.appendChild(div);
        });
        container.scrollTop = container.scrollHeight;
    }

    window.submitComment = async () => {
        const input = document.getElementById('comment-input');
        const content = input.value.trim();
        if (!content || !currentTaskId) return;

        try {
            const resp = await fetch(`/api/tasks/${currentTaskId}/comments`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content })
            });

            if (resp.ok) {
                input.value = '';
                // Reload comments
                const commentResp = await fetch(`/api/tasks/${currentTaskId}/comments`);
                const comments = await commentResp.json();
                renderComments(comments);
            } else {
                showToast('Failed to post comment', 'error');
            }
        } catch (err) {
            showToast('Error connecting to server', 'error');
        }
    };

    async function fetchTasks(query = '') {
        const columns = {
            'todo': document.getElementById('col-todo'),
            'in-progress': document.getElementById('col-in-progress'),
            'review': document.getElementById('col-review'),
            'done': document.getElementById('col-done')
        };
        
        if (!columns.todo) return;
        populateTaskDropdowns();

        const priorityFilter = document.getElementById('filter-priority')?.value || '';

        try {
            let url = '/api/tasks';
            const params = new URLSearchParams();
            if (query) params.append('q', query);
            if (priorityFilter) params.append('priority', priorityFilter);
            
            if (params.toString()) url += `?${params.toString()}`;

            const resp = await fetch(url);
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
                    taskEl.style.cursor = 'pointer';
                    taskEl.onclick = (e) => {
                        if (e.target.closest('button')) return;
                        openTaskDetails(t.id);
                    };
                    taskEl.innerHTML = `
                        <button onclick="deleteTask('${t.id}')" class="btn btn-ghost" title="Delete Task" style="position: absolute; top: 8px; right: 8px; color: var(--danger); cursor: pointer; padding: 4px; font-size: 12px; opacity: 0.6; transition: opacity 0.2s">
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
                             <div title="Assigned To: ${t.assigned_to_name || 'Unassigned'}" style="width: 24px; height: 24px; border-radius: 50%; background: var(--primary); display: flex; align-items: center; justify-content: center; font-size: 10px; color: white; font-weight: bold; box-shadow: 0 0 10px rgba(79, 70, 229, 0.3)">
                                ${t.assigned_to_name ? t.assigned_to_name.charAt(0).toUpperCase() : '?'}
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
        const btnId = 'save-task-btn';

        if (!project_id) {
            showToast('Please select a project', 'error');
            return;
        }
        if (!title) {
            showToast('Task title is required', 'error');
            return;
        }

        setLoading(btnId, true);
        try {
            const resp = await fetch('/api/tasks', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ project_id, assigned_to, title, description, status, priority })
            });

            if (resp.ok) {
                showToast('Task created successfully');
                hideModal('task-modal');
                // Reset form
                document.getElementById('task-title').value = '';
                document.getElementById('task-desc').value = '';
                fetchTasks();
            } else {
                const errText = await resp.text();
                showToast(`Failed: ${errText}`, 'error');
            }
        } catch (err) {
            showToast('Server connection failed', 'error');
        } finally {
            setLoading(btnId, false);
        }
    };

    async function populateTaskDropdowns() {
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
    }

    // Profile Logic
    async function fetchProfile() {
        try {
            const resp = await fetch('/api/users/profile');
            const user = await resp.json();
            
            document.getElementById('profile-full-name').textContent = user.full_name || user.username;
            document.getElementById('profile-email').textContent = user.email;
            
            document.getElementById('edit-full-name').value = user.full_name || '';
            document.getElementById('edit-username').value = user.username;
            document.getElementById('edit-email').value = user.email;
            document.getElementById('edit-role').value = user.role;
        } catch (err) {
            showToast('Failed to fetch profile', 'error');
        }
    }

    window.updateProfile = async () => {
        const full_name = document.getElementById('edit-full-name').value;
        const email = document.getElementById('edit-email').value;
        const password = document.getElementById('new-password').value;
        const confirm = document.getElementById('confirm-password').value;

        if (password && password !== confirm) {
            showToast('Passwords do not match', 'error');
            return;
        }

        try {
            const resp = await fetch('/api/users/profile', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ full_name, email, password })
            });

            if (resp.ok) {
                showToast('Profile updated successfully');
                fetchProfile();
                // Clear password fields
                document.getElementById('new-password').value = '';
                document.getElementById('confirm-password').value = '';
            } else {
                showToast('Failed to update profile', 'error');
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
        const color = type === 'success' ? 'var(--success)' : (type === 'error' ? 'var(--danger)' : 'var(--warning)');
        toast.style.borderLeft = `4px solid ${color}`;
        
        document.body.appendChild(toast);
        
        setTimeout(() => {
            toast.style.opacity = '0';
            setTimeout(() => toast.remove(), 500);
        }, 3000);
    };

    window.formatDate = (dateStr) => {
        if (!dateStr) return 'N/A';
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-US', { 
            month: 'short', 
            day: 'numeric', 
            year: 'numeric'
        });
    };

    window.exportTasks = () => {
        window.location.href = '/api/tasks/export';
    };

    // Expose remaining local functions to global window
    window.deleteAccount = deleteAccount;

    // Loading State Helper
    window.setLoading = (elementId, isLoading) => {
        const el = document.getElementById(elementId);
        if (!el) return;
        if (isLoading) {
            el.dataset.oldContent = el.innerHTML;
            el.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Processing...';
            el.disabled = true;
        } else {
            el.innerHTML = el.dataset.oldContent || el.innerHTML;
            el.disabled = false;
        }
    };
});
