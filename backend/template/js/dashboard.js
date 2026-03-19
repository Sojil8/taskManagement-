const API_URL = '/api';
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}
const token = getCookie('token') || localStorage.getItem('token');

// Toast Notification System
window.showToast = function(message, type = 'error') {
    let container = document.querySelector('.toast-container');
    if (!container) {
        container = document.createElement('div');
        container.className = 'toast-container';
        document.body.appendChild(container);
    }

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    let icon = '';
    if (type === 'error') {
        icon = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--danger)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>`;
    } else if (type === 'info') {
        icon = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--secondary)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`;
    } else if (type === 'success') {
        icon = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--success)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`;
    } else {
        icon = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`;
    }

    toast.innerHTML = `${icon} <span style="flex: 1;">${message}</span>`;
    container.appendChild(toast);

    setTimeout(() => {
        toast.classList.add('fade-out');
        toast.addEventListener('animationend', () => {
            toast.remove();
        });
    }, 4000);
}

// Logout Logic
document.getElementById('logoutBtn').addEventListener('click', () => {
    localStorage.removeItem('token');
    document.cookie = 'token=; Max-Age=0; path=/;';
    window.location.replace('login.html');
});

// Modal Logic
const modal = document.getElementById('newTaskModal');
const openModalBtn = document.getElementById('openNewTaskModalBtn');
const closeModalBtn = document.getElementById('closeModalBtn');
const cancelModalBtn = document.getElementById('cancelModalBtn');

let currentTasks = [];
let currentCategories = [];
let editingTaskId = null;
let taskToDeleteId = null;

function openModal() {
    modal.style.display = 'flex';
    document.querySelector('#newTaskModal h2').textContent = editingTaskId ? 'Edit Task' : 'Create New Task';
}
function closeModal() {
    modal.style.display = 'none';
    document.getElementById('newTaskForm').reset();
    editingTaskId = null;
}

openModalBtn.addEventListener('click', () => {
    if (currentCategories.length === 0) {
        window.showToast("Please create a category first before creating a task.", "info");
        openCategoryModal();
        return;
    }
    openModal();
});
closeModalBtn.addEventListener('click', closeModal);
cancelModalBtn.addEventListener('click', closeModal);
modal.addEventListener('click', (e) => {
    if (e.target === modal) closeModal();
});

// Category Modal Logic
const categoryModal = document.getElementById('newCategoryModal');
const openCategoryModalBtn = document.getElementById('openNewCategoryModalBtn');
const closeCategoryModalBtn = document.getElementById('closeCategoryModalBtn');
const cancelCategoryModalBtn = document.getElementById('cancelCategoryModalBtn');

function openCategoryModal() {
    categoryModal.style.display = 'flex';
}
function closeCategoryModal() {
    categoryModal.style.display = 'none';
    document.getElementById('newCategoryForm').reset();
}

openCategoryModalBtn.addEventListener('click', openCategoryModal);
closeCategoryModalBtn.addEventListener('click', closeCategoryModal);
cancelCategoryModalBtn.addEventListener('click', closeCategoryModal);
categoryModal.addEventListener('click', (e) => {
    if (e.target === categoryModal) closeCategoryModal();
});

// Fetch Categories
async function fetchCategories() {
    try {
        const res = await fetch(`${API_URL}/categories`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (res.ok) {
            currentCategories = await res.json() || [];
            updateCategoryDropdown();
        }
    } catch (err) {
        console.error('Failed to load categories', err);
    }
}

function updateCategoryDropdown() {
    const select = document.getElementById('taskCategory');
    if (currentCategories.length === 0) {
        select.innerHTML = '<option value="">Create a category first</option>';
        select.disabled = true;
    } else {
        select.innerHTML = currentCategories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
        select.disabled = false;
    }
}

// Fetch and Render Tasks
async function fetchTasks() {
    try {
        const res = await fetch(`${API_URL}/tasks`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (res.status === 401) {
            localStorage.removeItem('token');
            window.location.replace('login.html');
            return;
        }

        const tasks = await res.json();
        currentTasks = tasks || [];
        renderTasks(currentTasks);
    } catch (err) {
        console.error('Failed to load tasks', err);
        document.getElementById('tasksList').innerHTML = `<div class="error-message">Failed to load tasks.</div>`;
    }
}

function formatDate(dateString) {
    if (!dateString) return 'No deadline';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
}

function renderTasks(tasks) {
    const container = document.getElementById('tasksList');
    if (tasks.length === 0) {
        container.innerHTML = `
            <div class="loading-state" style="opacity: 0.7;">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 22C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zM12 8v4l3 3"/></svg>
                <p>No tasks found. Create one to get started!</p>
            </div>
        `;
        return;
    }

    container.innerHTML = tasks.map(task => {
        // Calculate progress percentage
        const total = task.checkpoints ? task.checkpoints.length : 4;
        const completed = task.checkpoints ? task.checkpoints.filter(cp => cp.completed).length : 0;

        // Let's assume exactly 4 checkpoints as per requirement.
        // The space between 4 points is 3 segments. So width = (completed / 3) * 100
        // Wait, if 1 completed -> no line filled entirely yet?
        // Let's say if it's node 1, progress = 0%. Node 2 = 33%. Node 3 = 66%. Node 4 = 100%.
        // But the tracker advances linearly. Actually, let's map completions:
        // 0 completed = 0%
        // 1 completed = 25% (Wait, space between 1&2 is segment1. If Node 1 is done, line goes till Node 1? Node 1 is at 12.5%. Let's just use CSS flex space-between).
        // Let's just calculate a simple percentage for the progress bar:
        let progressWidth = 0;
        if (completed === 1) progressWidth = 15;
        if (completed === 2) progressWidth = 50;
        if (completed === 3) progressWidth = 85;
        if (completed === 4) progressWidth = 100;
        if (completed === 0) progressWidth = 0;

        const isTaskCompleted = task.status === 'Completed';

        let checkpointsHtml = '';
        if (task.checkpoints) {
            // Find the index of the first incomplete checkpoint
            const firstIncompleteIdx = task.checkpoints.findIndex(c => !c.completed);

            checkpointsHtml = task.checkpoints.map((cp, idx) => {
                let interactHtml = '';
                if (isTaskCompleted) {
                    interactHtml = `<div class="cp-interact" title="Task completed - cannot edit"><div class="cp-circle"><span>${cp.order}</span></div></div>`;
                } else {
                    // Enforce sequential interaction visually
                    let canInteract = false;
                    let titleText = '';

                    if (cp.completed) {
                        // Can only untick if it's the LAST completed one (i.e. the one right before firstIncompleteIdx)
                        // If all are completed (isTaskCompleted = true), this block won't run.
                        const lastCompletedIdx = firstIncompleteIdx === -1 ? task.checkpoints.length - 1 : firstIncompleteIdx - 1;
                        if (idx === lastCompletedIdx) {
                            canInteract = true;
                            titleText = 'Click to untick';
                        } else {
                            titleText = 'Cannot untick while subsequent points are completed';
                        }
                    } else {
                        // Can only tick if it's the FIRST incomplete one
                        if (idx === firstIncompleteIdx) {
                            canInteract = true;
                            titleText = 'Click to mark as completed';
                        } else {
                            titleText = 'Complete previous points first';
                        }
                    }

                    if (canInteract) {
                        interactHtml = `<div class="cp-interact" title="${titleText}" onclick="completeCheckpoint(${task.id}, ${cp.id}, ${!cp.completed})">
                            <div class="cp-circle"><span>${cp.order}</span></div>
                        </div>`;
                    } else {
                        interactHtml = `<div class="cp-interact locked" title="${titleText}">
                            <div class="cp-circle"><span>${cp.order}</span></div>
                        </div>`;
                    }
                }

                return `
                <div class="checkpoint-node ${cp.completed ? 'completed' : ''}">
                    ${interactHtml}
                    <div class="cp-label" title="${cp.title}">${cp.title}</div>
                </div>
            `}).join('');
        }

        const percentage = Math.round((completed / total) * 100);
        const celebrateClass = (isTaskCompleted && completed === total) ? 'celebrate' : '';

        return `
            <div class="task-card glass-card scale-in ${celebrateClass}">
                <div class="task-header">
                    <div style="flex: 1; padding-right: 1rem; cursor: pointer;" onclick="viewTaskDetails(${task.id})">
                        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
                            <div class="task-title">${task.title}</div>
                            <div style="display: flex; gap: 0.5rem;" onclick="event.stopPropagation()">
                                <button class="btn-edit" title="Edit Task" onclick="editTask(${task.id})" style="background:transparent; border:none; cursor:pointer; color: var(--primary-color);">
                                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                                </button>
                                <button class="btn-delete" title="Delete Task" onclick="deleteTask(${task.id})">
                                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2M10 11v6M14 11v6"/></svg>
                                </button>
                            </div>
                        </div>
                        ${task.description ? `<div class="task-desc">${task.description}</div>` : ''}
                    </div>
                    <div class="badges">
                        ${task.category ? `<span class="badge" style="background: rgba(124, 58, 237, 0.1); color: var(--primary-color); border: 1px solid rgba(124, 58, 237, 0.2);">${task.category.name}</span>` : ''}
                        <span class="badge badge-deadline">${formatDate(task.deadline)}</span>
                        <span class="badge badge-progress">${completed}/${total} (${percentage}%)</span>
                        <span class="badge badge-status ${isTaskCompleted ? 'status-completed' : ''}">${task.status}</span>
                    </div>
                </div>

                <div class="checkpoint-tracker">
                    <div class="tracker-wrapper">
                        <div class="tracker-progress-bar" style="width: ${progressWidth}%"></div>
                        ${checkpointsHtml}
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

// Complete Checkpoint
window.completeCheckpoint = async function (taskId, cpId, targetStatus) {
    try {
        const res = await fetch(`${API_URL}/tasks/${taskId}/checkpoints/${cpId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ completed: targetStatus })
        });

        if (res.ok) {
            // Get original tasks to check if we just completed the 4th checkpoint
            const resTasks = await fetch(`${API_URL}/tasks`, { headers: { 'Authorization': `Bearer ${token}` } });
            const tasksData = await resTasks.json();
            const task = tasksData.find(t => t.id === taskId);

            // If the task just became fully completed
            if (targetStatus === true && task && task.status === 'Completed') {
                if (window.confetti) {
                    confetti({
                        particleCount: 150,
                        spread: 70,
                        origin: { y: 0.6 },
                        colors: ['#7C3AED', '#06B6D4', '#10B981', '#F59E0B']
                    });
                }
            }

            renderTasks(tasksData || []);
        } else {
            const data = await res.json();
            window.showToast('Error: ' + data.error, 'error');
        }
    } catch (err) {
        console.error(err);
        window.showToast('Failed to update checkpoint', 'error');
    }
}

// Delete Task Modal Logic
const deleteConfirmModal = document.getElementById('deleteConfirmModal');
const cancelDeleteBtn = document.getElementById('cancelDeleteBtn');
const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');

function openDeleteModal(taskId) {
    taskToDeleteId = taskId;
    deleteConfirmModal.style.display = 'flex';
}

function closeDeleteModal() {
    deleteConfirmModal.style.display = 'none';
    taskToDeleteId = null;
    
    // Reset button state
    const btn = confirmDeleteBtn;
    btn.disabled = false;
    btn.querySelector('.spinner').style.display = 'none';
    btn.querySelector('span').style.display = 'inline';
}

cancelDeleteBtn.addEventListener('click', closeDeleteModal);
deleteConfirmModal.addEventListener('click', (e) => {
    if (e.target === deleteConfirmModal) closeDeleteModal();
});

confirmDeleteBtn.addEventListener('click', async () => {
    if (!taskToDeleteId) return;

    const btn = confirmDeleteBtn;
    btn.disabled = true;
    btn.querySelector('.spinner').style.display = 'inline-block';
    btn.querySelector('span').style.display = 'none';

    try {
        const res = await fetch(`${API_URL}/tasks/${taskToDeleteId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (res.ok) {
            closeDeleteModal();
            fetchTasks();
            window.showToast('Task deleted successfully', 'success');
        } else {
            const data = await res.json();
            window.showToast('Error: ' + data.error, 'error');
            closeDeleteModal();
        }
    } catch (err) {
        console.error(err);
        window.showToast('Failed to delete task', 'error');
        closeDeleteModal();
    }
});

window.deleteTask = function (taskId) {
    openDeleteModal(taskId);
}

// Task Details Modal Logic
const taskDetailsModal = document.getElementById('taskDetailsModal');
const closeDetailsModalBtn = document.getElementById('closeDetailsModalBtn');

function closeDetailsModal() {
    taskDetailsModal.style.display = 'none';
}

closeDetailsModalBtn.addEventListener('click', closeDetailsModal);
taskDetailsModal.addEventListener('click', (e) => {
    if (e.target === taskDetailsModal) closeDetailsModal();
});

window.viewTaskDetails = function(taskId) {
    const task = currentTasks.find(t => t.id === taskId);
    if (!task) return;

    document.getElementById('detailsTitle').textContent = task.title;
    document.getElementById('detailsDescription').textContent = task.description || 'No description provided.';
    
    // Total and completed checkpoints
    const total = task.checkpoints ? task.checkpoints.length : 4;
    const completed = task.checkpoints ? task.checkpoints.filter(cp => cp.completed).length : 0;
    const percentage = Math.round((completed / total) * 100);
    const isTaskCompleted = task.status === 'Completed';

    let badgesHtml = '';
    if (task.category) {
        badgesHtml += `<span class="badge" style="background: rgba(124, 58, 237, 0.1); color: var(--primary-color); border: 1px solid rgba(124, 58, 237, 0.2);">${task.category.name}</span>`;
    }
    badgesHtml += `
        <span class="badge badge-deadline">${formatDate(task.deadline)}</span>
        <span class="badge badge-progress">${completed}/${total} (${percentage}%)</span>
        <span class="badge badge-status ${isTaskCompleted ? 'status-completed' : ''}">${task.status}</span>
    `;
    document.getElementById('detailsBadges').innerHTML = badgesHtml;

    let checkpointsHtml = '';
    if (task.checkpoints) {
        checkpointsHtml = task.checkpoints.map(cp => {
            return `
                <div style="display: flex; align-items: flex-start; gap: 1rem; padding-bottom: 0.5rem; ${!cp.completed ? 'opacity: 0.7;' : ''}">
                    <div style="flex-shrink: 0; width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; background: ${cp.completed ? 'var(--primary)' : 'rgba(255,255,255,0.1)'}; color: ${cp.completed ? '#fff' : 'var(--text-muted)'}; font-size: 0.8rem;">
                        ${cp.completed ? '✓' : cp.order}
                    </div>
                    <div style="color: ${cp.completed ? 'var(--text-main)' : 'var(--text-muted)'}; line-height: 1.4;">
                        ${cp.title}
                    </div>
                </div>
            `;
        }).join('');
    } else {
        checkpointsHtml = '<div style="color: var(--text-muted);">No checkpoints</div>';
    }
    document.getElementById('detailsCheckpoints').innerHTML = checkpointsHtml;

    taskDetailsModal.style.display = 'flex';
}

// Edit Task
window.editTask = function (taskId) {
    const task = currentTasks.find(t => t.id === taskId);
    if (!task) return;

    editingTaskId = taskId;

    document.getElementById('taskTitle').value = task.title;
    document.getElementById('taskDescription').value = task.description || '';

    if (task.deadline) {
        // Format to YYYY-MM-DD
        const date = new Date(task.deadline);
        const yyyy = date.getFullYear();
        const mm = String(date.getMonth() + 1).padStart(2, '0');
        const dd = String(date.getDate()).padStart(2, '0');
        document.getElementById('taskDeadline').value = `${yyyy}-${mm}-${dd}`;
    }

    if (task.checkpoints && task.checkpoints.length === 4) {
        document.getElementById('cp1').value = task.checkpoints[0].title;
        document.getElementById('cp2').value = task.checkpoints[1].title;
        document.getElementById('cp3').value = task.checkpoints[2].title;
        document.getElementById('cp4').value = task.checkpoints[3].title;
    }

    if (task.category_id) {
        document.getElementById('taskCategory').value = task.category_id;
    }

    openModal();
}

// Create Task
const newTaskForm = document.getElementById('newTaskForm');
newTaskForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    document.getElementById('modalError').style.display = 'none';
    const btn = newTaskForm.querySelector('button[type="submit"]');
    btn.disabled = true;
    btn.querySelector('.spinner').style.display = 'block';

    const title = document.getElementById('taskTitle').value;
    const description = document.getElementById('taskDescription').value;
    const deadlineStr = document.getElementById('taskDeadline').value;
    const cp1 = document.getElementById('cp1').value;
    const cp2 = document.getElementById('cp2').value;
    const cp3 = document.getElementById('cp3').value;
    const cp4 = document.getElementById('cp4').value;
    const categoryIdVal = document.getElementById('taskCategory').value;
    const category_id = categoryIdVal ? parseInt(categoryIdVal) : null;

    // Parse deadline
    const deadline = new Date(deadlineStr).toISOString();

    const checkpoints = [cp1, cp2, cp3, cp4];

    try {
        const url = editingTaskId ? `${API_URL}/tasks/${editingTaskId}` : `${API_URL}/tasks`;
        const method = editingTaskId ? 'PUT' : 'POST';

        const res = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                title, description, deadline, checkpoints, category_id
            })
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error || (editingTaskId ? 'Failed to update task' : 'Failed to create task'));

        closeModal();
        fetchTasks();
    } catch (err) {
        const errEl = document.getElementById('modalError');
        errEl.textContent = err.message;
        errEl.style.display = 'block';
    } finally {
        btn.disabled = false;
        btn.querySelector('.spinner').style.display = 'none';
    }
});

// Create Category
const newCategoryForm = document.getElementById('newCategoryForm');
newCategoryForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    document.getElementById('categoryModalError').style.display = 'none';
    const btn = newCategoryForm.querySelector('button[type="submit"]');
    btn.disabled = true;
    btn.querySelector('.spinner').style.display = 'block';

    const name = document.getElementById('categoryName').value;

    try {
        const res = await fetch(`${API_URL}/categories`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name })
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Failed to create category');

        closeCategoryModal();
        await fetchCategories(); // Refresh categories

        // Auto-select the newly created category if the task modal is open or immediately available
        const select = document.getElementById('taskCategory');
        if (select) {
            select.value = data.id;
        }
    } catch (err) {
        const errEl = document.getElementById('categoryModalError');
        errEl.textContent = err.message;
        errEl.style.display = 'block';
    } finally {
        btn.disabled = false;
        btn.querySelector('.spinner').style.display = 'none';
    }
});

// Init
async function init() {
    await fetchCategories();
    fetchTasks();
}
init();
