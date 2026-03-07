document.addEventListener('DOMContentLoaded', () => {
    // === DOM Elements ===
    const timeDisplay = document.getElementById('timeDisplay');
    const startPauseBtn = document.getElementById('startPauseBtn');
    const resetBtn = document.getElementById('resetBtn');
    const cycleStatus = document.getElementById('cycleStatus');
    const currentCycleEl = document.getElementById('currentCycle');
    const totalCyclesEl = document.getElementById('totalCycles');
    const sessionTypeIndicator = document.getElementById('sessionTypeIndicator');
    const skipSessionBtn = document.getElementById('skipSessionBtn');

    const progressCircle = document.getElementById('timeProgress');
    const radius = progressCircle.r.baseVal.value;
    const circumference = radius * 2 * Math.PI;

    // Modals & Settings
    const settingsModal = document.getElementById('settingsModal');
    const openSettingsModalBtn = document.getElementById('openSettingsModalBtn');
    const closeSettingsBtn = document.getElementById('closeSettingsBtn');
    const cancelSettingsBtn = document.getElementById('cancelSettingsBtn');
    const settingsForm = document.getElementById('settingsForm');

    // Break Activity Card
    const breakActivityCard = document.getElementById('breakActivityCard');
    const activityTitle = document.getElementById('activityTitle');
    const activityDesc = document.getElementById('activityDesc');
    const newActivityBtn = document.getElementById('newActivityBtn');

    // Default Settings (in seconds)
    let settings = {
        focus: 25 * 60,
        shortBreak: 5 * 60,
        longBreak: 15 * 60,
        cycles: 4
    };

    // State Variables
    let timerId = null;
    let isRunning = false;
    let timeLeft = settings.focus;
    let currentSession = 'focus'; // 'focus', 'short-break', 'long-break'
    let currentCycle = 1;
    let totalTimeForSession = settings.focus;

    // === Break Activities ===
    const breakActivities = [
        { title: "10 Push-ups", desc: "Get up and do 10 push-ups to keep your blood flowing." },
        { title: "15 Sit-ups", desc: "Strengthen your core with 15 quick sit-ups." },
        { title: "20 Jumping Jacks", desc: "Get your heart rate up with 20 jumping jacks." },
        { title: "30-second Plank", desc: "Engage your core with a 30-second plank." },
        { title: "Stretch Arms & Shoulders", desc: "Reach for the ceiling, then stretch your arms across your chest." },
        { title: "Walk Around", desc: "Pace around your room or house for 1 minute." },
        { title: "Drink Water", desc: "Hydrate yourself! Drink a full glass of water." },
        { title: "Eye Relaxation", desc: "Look at something 20 feet away for 20 seconds." }
    ];

    // Initialize the Progress Circle
    progressCircle.style.strokeDasharray = `${circumference} ${circumference}`;
    progressCircle.style.strokeDashoffset = circumference;

    function setProgress(percent) {
        const offset = circumference - percent / 100 * circumference;
        progressCircle.style.strokeDashoffset = offset;
    }

    // === Core Timer Logic ===
    function updateDisplay() {
        const minutes = Math.floor(timeLeft / 60);
        const seconds = timeLeft % 60;
        timeDisplay.textContent = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;

        // Update Title tag to show time
        const sessionName = currentSession === 'focus' ? 'Focus' : 'Break';
        document.title = `${timeDisplay.textContent} - ${sessionName} - TaskTracker`;

        // Update Progress Circle
        const percentage = (timeLeft / totalTimeForSession) * 100;
        setProgress(percentage);
    }

    function switchSession(forced = false) {
        if (currentSession === 'focus') {
            if (currentCycle >= settings.cycles) {
                currentSession = 'long-break';
                timeLeft = settings.longBreak;
                totalTimeForSession = settings.longBreak;
                currentCycle = 1; // Reset cycles after long break
            } else {
                currentSession = 'short-break';
                timeLeft = settings.shortBreak;
                totalTimeForSession = settings.shortBreak;
            }
        } else {
            // Coming back from a break
            if (currentSession === 'short-break') {
                currentCycle++;
            }
            currentSession = 'focus';
            timeLeft = settings.focus;
            totalTimeForSession = settings.focus;
        }

        updateUIState();
        updateDisplay();
        saveState();

        if (!forced) notifySessionComplete();
    }

    function updateUIState() {
        // Update DOM attributes for CSS styling
        document.body.setAttribute('data-session', currentSession);

        // Update Labels
        if (currentSession === 'focus') {
            sessionTypeIndicator.textContent = 'Focus';
            breakActivityCard.style.display = 'none';
        } else if (currentSession === 'short-break') {
            sessionTypeIndicator.textContent = 'Short Break';
            showRandomActivity();
        } else if (currentSession === 'long-break') {
            sessionTypeIndicator.textContent = 'Long Break';
            showRandomActivity();
        }

        currentCycleEl.textContent = currentCycle;
        totalCyclesEl.textContent = settings.cycles;
    }

    function toggleTimer() {
        if (isRunning) {
            clearInterval(timerId);
            startPauseBtn.textContent = 'Resume';
        } else {
            // Request Notification Permission on first Start
            if (Notification.permission === 'default') {
                Notification.requestPermission();
            }

            // Start ticking
            timerId = setInterval(() => {
                if (timeLeft > 0) {
                    timeLeft--;
                    updateDisplay();
                    saveState();
                } else {
                    clearInterval(timerId);
                    switchSession();
                    if (isRunning) toggleTimer(); // Auto-start next session
                }
            }, 1000);
            startPauseBtn.textContent = 'Pause';
        }
        isRunning = !isRunning;
        startPauseBtn.classList.toggle('btn-primary', !isRunning);
        startPauseBtn.classList.toggle('btn-outline', isRunning);
    }

    function resetTimer() {
        clearInterval(timerId);
        isRunning = false;
        startPauseBtn.textContent = 'Start';
        startPauseBtn.classList.add('btn-primary');
        startPauseBtn.classList.remove('btn-outline');

        // Reset to initial focus state
        currentSession = 'focus';
        currentCycle = 1;
        timeLeft = settings.focus;
        totalTimeForSession = settings.focus;

        updateUIState();
        updateDisplay();
        saveState();
    }

    function skipSession() {
        clearInterval(timerId);
        isRunning = false;
        startPauseBtn.textContent = 'Start';
        startPauseBtn.classList.add('btn-primary');
        startPauseBtn.classList.remove('btn-outline');
        switchSession(true);
    }

    // === Notifications ===
    function notifySessionComplete() {
        // Play Sound
        playBeep();

        // Browser Notification
        if (Notification.permission === 'granted') {
            const title = currentSession === 'focus' ? 'Back to Work!' : 'Time for a break!';
            const options = {
                body: currentSession === 'focus' ? `Focus for ${settings.focus / 60} minutes.` : `Take a break for ${timeLeft / 60} minutes.`,
                icon: 'httpscdn-icons-png.flaticon.com/512/813/813473.png' // generic clock icon
            };
            new Notification(title, options);
        }
    }

    function playBeep() {
        const audioCtx = new (window.AudioContext || window.webkitAudioContext)();
        const oscillator = audioCtx.createOscillator();
        const gainNode = audioCtx.createGain();

        oscillator.type = 'sine';
        oscillator.frequency.setValueAtTime(880, audioCtx.currentTime); // A5

        gainNode.gain.setValueAtTime(0.1, audioCtx.currentTime);
        gainNode.gain.exponentialRampToValueAtTime(0.00001, audioCtx.currentTime + 1);

        oscillator.connect(gainNode);
        gainNode.connect(audioCtx.destination);

        oscillator.start(audioCtx.currentTime);
        oscillator.stop(audioCtx.currentTime + 1);
    }

    // === Break Activities ===
    function showRandomActivity() {
        breakActivityCard.style.display = 'flex';
        // Remove simple animation reflow hack to restart animation
        breakActivityCard.classList.remove('scale-in');
        void breakActivityCard.offsetWidth; // trigger reflow
        breakActivityCard.classList.add('scale-in');

        const randomIndex = Math.floor(Math.random() * breakActivities.length);
        const activity = breakActivities[randomIndex];

        activityTitle.textContent = activity.title;
        activityDesc.textContent = activity.desc;
    }

    newActivityBtn.addEventListener('click', showRandomActivity);

    // === Persistence (localStorage) ===
    function saveState() {
        const state = {
            timeLeft,
            currentSession,
            currentCycle,
            timestamp: Date.now()
        };
        localStorage.setItem('pomodoroState', JSON.stringify(state));
    }

    function loadState() {
        // Load Settings
        const savedSettings = localStorage.getItem('pomodoroSettings');
        if (savedSettings) {
            settings = JSON.parse(savedSettings);
        }

        // Load Timer State
        const savedStateStr = localStorage.getItem('pomodoroState');
        if (savedStateStr) {
            const savedState = JSON.parse(savedStateStr);
            const elapsedSeconds = Math.floor((Date.now() - savedState.timestamp) / 1000);

            currentSession = savedState.currentSession;
            currentCycle = savedState.currentCycle;

            // Determine correct total time
            if (currentSession === 'focus') totalTimeForSession = settings.focus;
            else if (currentSession === 'short-break') totalTimeForSession = settings.shortBreak;
            else totalTimeForSession = settings.longBreak;

            // Adjust time left as if timer was running in background?
            // Usually, Pomodoros don't automatically deduct missing time unless explicitly running.
            // For simplicity and exact resumption, we'll just restore the exact time left.
            timeLeft = savedState.timeLeft;
        } else {
            timeLeft = settings.focus;
            totalTimeForSession = settings.focus;
        }

        updateUIState();
        updateDisplay();
    }

    function saveSettings() {
        const f = parseInt(document.getElementById('focusDuration').value);
        const s = parseInt(document.getElementById('shortBreakDuration').value);
        const l = parseInt(document.getElementById('longBreakDuration').value);
        const c = parseInt(document.getElementById('cyclesCount').value);

        settings = {
            focus: f * 60,
            shortBreak: s * 60,
            longBreak: l * 60,
            cycles: c
        };

        localStorage.setItem('pomodoroSettings', JSON.stringify(settings));
        resetTimer();
    }

    function populateSettingsModal() {
        document.getElementById('focusDuration').value = settings.focus / 60;
        document.getElementById('shortBreakDuration').value = settings.shortBreak / 60;
        document.getElementById('longBreakDuration').value = settings.longBreak / 60;
        document.getElementById('cyclesCount').value = settings.cycles;
    }

    // === Event Listeners ===
    startPauseBtn.addEventListener('click', toggleTimer);
    resetBtn.addEventListener('click', resetTimer);
    skipSessionBtn.addEventListener('click', skipSession);

    // Modal Events
    openSettingsModalBtn.addEventListener('click', () => {
        populateSettingsModal();
        settingsModal.style.display = 'flex';
    });

    closeSettingsBtn.addEventListener('click', () => {
        settingsModal.style.display = 'none';
    });

    cancelSettingsBtn.addEventListener('click', () => {
        settingsModal.style.display = 'none';
    });

    settingsModal.addEventListener('click', (e) => {
        if (e.target === settingsModal) {
            settingsModal.style.display = 'none';
        }
    });

    settingsForm.addEventListener('submit', (e) => {
        e.preventDefault();
        saveSettings();
        settingsModal.style.display = 'none';
    });

    // Logout
    document.getElementById('logoutBtn').addEventListener('click', () => {
        localStorage.removeItem('token');
        window.location.href = 'login.html';
    });

    // Check auth
    if (!localStorage.getItem('token')) {
        window.location.href = 'login.html';
    }

    // Init
    loadState();
});
