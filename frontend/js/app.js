import { api, toast } from './api.js'
import { showAuth, login, signup, logout, showLoginForm, showSignupForm } from './auth.js'
import {
  showDashboard, navigateTo, loadUserData,
  captureSnapshot, selectSnapshot, setTrackLimit, connectSpotify, disconnectSpotify,
} from './dashboard.js'
import { analyzeTaste, loadTrends, downloadReport } from './analytics.js'
import { loadLeaderboard } from './leaderboard.js'

// Expose to inline HTML handlers
window.login            = login
window.signup           = signup
window.logout           = logout
window.showLoginForm    = showLoginForm
window.showSignupForm   = showSignupForm
window.navigateTo        = navigateTo
window.captureSnapshot   = captureSnapshot
window.selectSnapshot    = selectSnapshot
window.setTrackLimit     = setTrackLimit
window.connectSpotify    = connectSpotify
window.disconnectSpotify = disconnectSpotify
window.analyzeTaste     = analyzeTaste
window.loadTrends       = loadTrends
window.downloadReport   = downloadReport
window.loadLeaderboard  = loadLeaderboard

async function checkOAuthCallback() {
  const params = new URLSearchParams(window.location.search)
  const connected = params.get('connected')
  const error = params.get('error')
  if (connected === '1') {
    toast('Spotify connected! Loading your tracks…', 'success')
  } else if (error) {
    const msgs = {
      not_logged_in:         'Please sign in first.',
      session_expired:       'Session expired — sign in again.',
      invalid_state:         'Security check failed.',
      spotify_already_linked:'This Spotify account is already linked to another user.',
      db_error:              'Server error — try again.',
    }
    toast(msgs[error] || `OAuth error: ${error}`, 'error')
  }
  if (connected || error)
    window.history.replaceState({}, '', window.location.pathname)
}

async function init() {
  const res = await api('/auth/refresh', 'POST')
  if (res.ok && res.data.access_token) {
    window.currentAccessToken = res.data.access_token
    await checkOAuthCallback()
    await loadUserData()
    showDashboard()
    navigateTo('tracks')
    loadLeaderboard('valence')
  } else {
    showAuth()
  }
}

init()
