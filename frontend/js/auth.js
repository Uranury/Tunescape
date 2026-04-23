import { api, toast, setLoading } from './api.js'
import { showDashboard, loadUserData } from './dashboard.js'

export function showAuth() {
  document.getElementById('auth-screen').style.display = 'flex'
  document.getElementById('app-shell').classList.remove('active')
  showLoginForm()
}

export function showLoginForm() {
  document.getElementById('login-form').style.display = 'block'
  document.getElementById('signup-form').style.display = 'none'
  document.getElementById('auth-title').textContent = 'Sign in'
}

export function showSignupForm() {
  document.getElementById('login-form').style.display = 'none'
  document.getElementById('signup-form').style.display = 'block'
  document.getElementById('auth-title').textContent = 'Create account'
}

export async function login() {
  const email = document.getElementById('login-email').value.trim()
  const password = document.getElementById('login-password').value
  if (!email || !password) { toast('Fill in all fields', 'error'); return }
  setLoading(true)
  const res = await api('/auth/login', 'POST', { email, password })
  setLoading(false)
  if (res.ok && res.data.access_token) {
    window.currentAccessToken = res.data.access_token
    await loadUserData()
    showDashboard()
    toast('Signed in', 'success')
  } else {
    toast(res.data?.error || 'Login failed', 'error')
  }
}

export async function signup() {
  const email = document.getElementById('signup-email').value.trim()
  const password = document.getElementById('signup-password').value
  const displayName = document.getElementById('signup-name').value.trim()
  if (!email || !password || !displayName) { toast('Fill in all fields', 'error'); return }
  if (password.length < 8) { toast('Password must be at least 8 characters', 'error'); return }
  setLoading(true)
  const res = await api('/auth/signup', 'POST', { email, password, display_name: displayName })
  setLoading(false)
  if (res.ok && res.data.access_token) {
    window.currentAccessToken = res.data.access_token
    await loadUserData()
    showDashboard()
    toast('Account created!', 'success')
  } else {
    toast(res.data?.error || 'Signup failed', 'error')
  }
}

export async function logout() {
  await api('/auth/logout', 'POST')
  window.currentAccessToken = null
  showAuth()
  toast('Signed out', 'success')
}
