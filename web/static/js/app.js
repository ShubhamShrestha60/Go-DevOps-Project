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

    // Fetch user info for dashboard
    if (path === '/' || path === '/projects' || path === '/tasks') {
        fetchUserInfo();
    }

    async function fetchUserInfo() {
        // In a real app, we'd have a /api/auth/me endpoint
        // For now, we'll try to get it from local storage or just show a fallback
        const userNameElement = document.getElementById('user-name');
        if (userNameElement) {
            userNameElement.textContent = "Production User";
        }
    }

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
