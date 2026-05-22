import { api, toast, setLoading, escapeHtml } from './api.js'

let compareChart = null

function avatarHtml(url, name) {
  if (url) {
    return `<div class="friends-user-avatar"><img src="${escapeHtml(url)}" alt="${escapeHtml(name)}"></div>`
  }
  return `<div class="friends-user-avatar">${escapeHtml(name.charAt(0).toUpperCase())}</div>`
}

export async function loadFriendsPanel() {
  await Promise.all([loadIncomingRequests(), loadFriendsList()])
}

export async function lookupUser() {
  const input = document.getElementById('friends-lookup-input')
  const email = input.value.trim()
  if (!email) {
    toast('Enter an email address to search', 'error')
    return
  }

  const resultEl = document.getElementById('friends-lookup-result')
  resultEl.innerHTML = '<p class="empty-state">Searching…</p>'

  const res = await api(`/users/lookup?email=${encodeURIComponent(email)}`, 'GET', null, true)
  if (!res.ok) {
    if (res.status === 404) {
      resultEl.innerHTML = '<p class="empty-state" style="color:var(--red)">No user found with that username.</p>'
    } else {
      resultEl.innerHTML = '<p class="empty-state" style="color:var(--red)">Search failed.</p>'
    }
    return
  }

  const u = res.data
  resultEl.innerHTML = `
    <div class="friends-lookup-found">
      ${avatarHtml(u.avatar_url ?? null, u.display_name)}
      <div class="friends-user-info">
        <div class="friends-user-name">${escapeHtml(u.display_name)}</div>
        <div class="friends-user-id">${escapeHtml(u.user_id)}</div>
      </div>
      <button class="btn btn-primary btn-sm" onclick="sendFriendRequest('${escapeHtml(u.user_id)}')">
        Add Friend
      </button>
    </div>`
}

export async function sendFriendRequest(receiverID) {
  const res = await api('/friends/requests', 'POST', { receiver_id: receiverID }, true)
  if (res.ok) {
    toast('Friend request sent!', 'success')
    document.getElementById('friends-lookup-result').innerHTML = ''
    document.getElementById('friends-lookup-input').value = ''
  } else if (res.status === 409) {
    toast(res.data?.message || 'Request already sent or already friends', 'error')
  } else if (res.status === 400) {
    toast(res.data?.message || 'Invalid request', 'error')
  } else {
    toast('Failed to send friend request', 'error')
  }
}

export async function loadIncomingRequests() {
  const res = await api('/friends/requests', 'GET', null, true)
  const el = document.getElementById('friends-requests-list')
  if (!res.ok) {
    el.innerHTML = '<p class="empty-state" style="color:var(--red)">Failed to load requests.</p>'
    return
  }

  const reqs = res.data || []
  if (reqs.length === 0) {
    el.innerHTML = '<p class="empty-state">No pending requests.</p>'
    return
  }

  el.innerHTML = reqs.map(r => `
    <div class="friends-request-row">
      ${avatarHtml(null, r.display_name)}
      <div class="friends-request-info">
        <div class="friends-user-name">${escapeHtml(r.display_name)}</div>
        <div class="friends-request-time">${new Date(r.created_at).toLocaleDateString()}</div>
      </div>
      <div class="friends-request-actions">
        <button class="btn btn-primary btn-sm" style="width:auto;flex-shrink:0" onclick="acceptRequest(${r.request_id})">Accept</button>
        <button class="btn btn-ghost btn-sm" onclick="rejectRequest(${r.request_id})">Reject</button>
      </div>
    </div>`).join('')
}

export async function acceptRequest(requestID) {
  const res = await api(`/friends/requests/${requestID}/accept`, 'POST', null, true)
  if (res.ok) {
    toast('Friend request accepted!', 'success')
    await Promise.all([loadIncomingRequests(), loadFriendsList()])
  } else {
    toast(res.data?.message || 'Failed to accept request', 'error')
  }
}

export async function rejectRequest(requestID) {
  const res = await api(`/friends/requests/${requestID}/reject`, 'POST', null, true)
  if (res.ok) {
    toast('Request rejected', 'info')
    await loadIncomingRequests()
  } else {
    toast(res.data?.message || 'Failed to reject request', 'error')
  }
}

export async function loadFriendsList() {
  const res = await api('/friends', 'GET', null, true)
  const el = document.getElementById('friends-list')
  if (!res.ok) {
    el.innerHTML = '<p class="empty-state" style="color:var(--red)">Failed to load friends.</p>'
    return
  }

  const friends = res.data || []
  if (friends.length === 0) {
    el.innerHTML = '<p class="empty-state">No friends yet — find someone above.</p>'
    return
  }

  el.innerHTML = friends.map(f => `
    <div class="friends-row">
      ${avatarHtml(f.avatar_url, f.display_name)}
      <div class="friends-user-info">
        <div class="friends-user-name">${escapeHtml(f.display_name)}</div>
        <div class="friends-user-sub">${f.spotify_connected ? '<span class="friends-spotify-badge">Spotify connected</span>' : '<span class="friends-spotify-badge disconnected">No Spotify</span>'}</div>
      </div>
      <button class="btn btn-ghost btn-sm" onclick="compareTastes('${escapeHtml(f.user_id)}', '${escapeHtml(f.display_name)}')">Compare Taste</button>
    </div>`).join('')
}

export async function compareTastes(friendID, friendName) {
  setLoading(true)
  const res = await api(`/friends/${friendID}/compare`, 'GET', null, true)
  setLoading(false)

  if (!res.ok) {
    if (res.status === 422) {
      toast(`${friendName} has no listening data yet`, 'error')
    } else if (res.status === 403) {
      toast('You are not friends with this user', 'error')
    } else {
      toast('Failed to load comparison', 'error')
    }
    return
  }

  const card = document.getElementById('friends-compare-card')
  const title = document.getElementById('friends-compare-title')
  const content = document.getElementById('friends-compare-content')

  title.textContent = `Your Taste vs ${friendName}`
  card.style.display = 'block'
  card.scrollIntoView({ behavior: 'smooth', block: 'nearest' })

  renderComparison(res.data, friendName)
}

function renderComparison(data, friendName) {
  const content = document.getElementById('friends-compare-content')
  const mine = data.mine
  const theirs = data.theirs
  const score = data.compatibility_score

  const scoreColor = score >= 70 ? 'var(--green)' : score >= 40 ? 'var(--skype-blue)' : 'var(--red)'

  content.innerHTML = `
    <div class="compare-score-row">
      <div class="compare-score-label">Compatibility Score</div>
      <div class="compare-score-value" style="color:${scoreColor}">${score.toFixed(1)}<span class="compare-score-unit">/100</span></div>
    </div>
    <div class="compare-grid">
      ${renderFeatureBar('Valence', mine.valence, theirs.valence, friendName)}
      ${renderFeatureBar('Energy', mine.energy, theirs.energy, friendName)}
      ${renderFeatureBar('Danceability', mine.danceability, theirs.danceability, friendName)}
      ${renderFeatureBar('Acousticness', mine.acousticness, theirs.acousticness, friendName)}
    </div>
    <div class="compare-chart-wrap">
      <canvas id="compare-radar-canvas"></canvas>
    </div>`

  renderRadarChart(mine, theirs, friendName)
}

function renderFeatureBar(label, myVal, theirVal, friendName) {
  const myPct = (myVal * 100).toFixed(1)
  const theirPct = (theirVal * 100).toFixed(1)
  return `
    <div class="compare-feature">
      <div class="compare-feature-label">${escapeHtml(label)}</div>
      <div class="compare-bars">
        <div class="compare-bar-row">
          <span class="compare-bar-who">You</span>
          <div class="compare-bar-track">
            <div class="compare-bar-fill mine" style="width:${myPct}%"></div>
          </div>
          <span class="compare-bar-val">${myVal.toFixed(3)}</span>
        </div>
        <div class="compare-bar-row">
          <span class="compare-bar-who">${escapeHtml(friendName)}</span>
          <div class="compare-bar-track">
            <div class="compare-bar-fill theirs" style="width:${theirPct}%"></div>
          </div>
          <span class="compare-bar-val">${theirVal.toFixed(3)}</span>
        </div>
      </div>
    </div>`
}

function renderRadarChart(mine, theirs, friendName) {
  if (compareChart) { compareChart.destroy(); compareChart = null }

  const ctx = document.getElementById('compare-radar-canvas').getContext('2d')
  compareChart = new Chart(ctx, {
    type: 'radar',
    data: {
      labels: ['Valence', 'Energy', 'Danceability', 'Acousticness'],
      datasets: [
        {
          label: 'You',
          data: [mine.valence, mine.energy, mine.danceability, mine.acousticness],
          borderColor: '#0078D7',
          backgroundColor: 'rgba(0,120,215,0.15)',
          pointBackgroundColor: '#0078D7',
          borderWidth: 2,
        },
        {
          label: friendName,
          data: [theirs.valence, theirs.energy, theirs.danceability, theirs.acousticness],
          borderColor: '#F97316',
          backgroundColor: 'rgba(249,115,22,0.15)',
          pointBackgroundColor: '#F97316',
          borderWidth: 2,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      scales: {
        r: {
          min: 0,
          max: 1,
          ticks: { stepSize: 0.25, font: { size: 10 } },
          grid: { color: '#D0E4F0' },
          pointLabels: { font: { family: "'Segoe UI', sans-serif", size: 12 } },
        },
      },
      plugins: {
        legend: { labels: { font: { family: "'Segoe UI', sans-serif", size: 12 } } },
      },
    },
  })
}

export function closeCompare() {
  const card = document.getElementById('friends-compare-card')
  card.style.display = 'none'
  if (compareChart) { compareChart.destroy(); compareChart = null }
}
