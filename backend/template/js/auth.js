const API_URL = '/api';

// Utility to show errors
function showError(msg) {
    const errorEl = document.getElementById('errorMessage');
    if (errorEl) {
        errorEl.textContent = msg;
        errorEl.style.display = 'block';
    }
}

function hideError() {
    const errorEl = document.getElementById('errorMessage');
    if (errorEl) {
        errorEl.style.display = 'none';
    }
}

function toggleButtonLoading(btn, isLoading) {
    if (!btn) return;
    const span = btn.querySelector('span');
    const spinner = btn.querySelector('.spinner');

    if (isLoading) {
        btn.disabled = true;
        if (span) span.style.opacity = '0.7';
        if (spinner) spinner.style.display = 'block';
    } else {
        btn.disabled = false;
        if (span) span.style.opacity = '1';
        if (spinner) spinner.style.display = 'none';
    }
}

// Signup Flow
const signupForm = document.getElementById('signupForm');
if (signupForm) {
    signupForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        hideError();
        const btn = signupForm.querySelector('button[type="submit"]');
        toggleButtonLoading(btn, true);

        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            // 1. Register User
            let res = await fetch(`${API_URL}/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, email, password })
            });
            const data = await res.json();

            if (!res.ok) throw new Error(data.error || 'Registration failed');

            // 2. Send OTP
            res = await fetch(`${API_URL}/auth/send-otp`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email })
            });
            const otpData = await res.json();
            if (!res.ok) throw new Error(otpData.error || 'Failed to send OTP');

            // Save email temp for OTP page
            sessionStorage.setItem('tempEmail', email);
            window.location.replace('otp.html');

        } catch (err) {
            showError(err.message);
        } finally {
            toggleButtonLoading(btn, false);
        }
    });
}

// OTP Flow
const otpForm = document.getElementById('otpForm');
if (otpForm) {
    const tempEmail = sessionStorage.getItem('tempEmail');
    if (!tempEmail) {
        window.location.replace('signup.html');
    }

    document.getElementById('userEmailDisplay').textContent = tempEmail;

    // Auto-focus logic for OTP inputs
    const inputs = document.querySelectorAll('.otp-input');
    inputs.forEach((input, index) => {
        input.addEventListener('input', (e) => {
            if (e.target.value.length > 0) {
                if (index < inputs.length - 1) inputs[index + 1].focus();
            }
        });
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Backspace' && e.target.value === '' && index > 0) {
                inputs[index - 1].focus();
            }
        });
    });

    otpForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        hideError();
        const btn = otpForm.querySelector('button[type="submit"]');
        toggleButtonLoading(btn, true);

        let code = '';
        inputs.forEach(input => code += input.value);

        if (code.length !== 6) {
            showError('Please enter all 6 digits');
            toggleButtonLoading(btn, false);
            return;
        }

        try {
            const res = await fetch(`${API_URL}/auth/verify-otp`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: tempEmail, code })
            });
            const data = await res.json();

            if (!res.ok) throw new Error(data.error || 'OTP verification failed');

            localStorage.setItem('token', data.token);
            document.cookie = `token=${data.token}; max-age=${3600 * 24}; path=/`;
            sessionStorage.removeItem('tempEmail');
            window.location.replace('dashboard.html');

        } catch (err) {
            showError(err.message);
        } finally {
            toggleButtonLoading(btn, false);
        }
    });

    // Resend OTP logic
    document.getElementById('resendOtpBtn').addEventListener('click', async (e) => {
        e.preventDefault();
        try {
            const res = await fetch(`${API_URL}/auth/send-otp`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: tempEmail })
            });
            if (res.ok) {
                alert('A new OTP has been sent to your email.');
            } else {
                const data = await res.json();
                throw new Error(data.error);
            }
        } catch (err) {
            alert('Failed to resend OTP: ' + err.message);
        }
    });
}

// Login Flow
const loginForm = document.getElementById('loginForm');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        hideError();
        const btn = loginForm.querySelector('button[type="submit"]');
        toggleButtonLoading(btn, true);

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            const res = await fetch(`${API_URL}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });
            const data = await res.json();

            if (!res.ok) throw new Error(data.error || 'Login failed');

            localStorage.setItem('token', data.token);
            document.cookie = `token=${data.token}; max-age=${3600 * 24}; path=/`;
            window.location.replace('dashboard.html');

        } catch (err) {
            showError(err.message);
        } finally {
            toggleButtonLoading(btn, false);
        }
    });
}
