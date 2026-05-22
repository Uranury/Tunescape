import { api, toast, setLoading, escapeHtml } from './api.js'

let currentSnapshot = null
let trackLimit = 50

export function showDashboard() {
    document.getElementById('auth-screen').style.display = 'none'
    document.getElementById('app-shell').classList.add('active')
}

export async function navigateTo(panelId) {
    document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'))
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'))
    const panel = document.getElementById(`panel-${panelId}`)
    if (panel) panel.classList.add('active')
    const navItem = document.querySelector(`.nav-item[data-panel="${panelId}"]`)
    if (navItem) navItem.classList.add('active')
    document.getElementById('topbar-title').textContent =
        navItem?.querySelector('.nav-label')?.textContent || ''
    if (panelId === 'playlists') {
        const { loadPlaylistPreview } = await import('./playlist.js')
        await loadPlaylistPreview()
    }
    if (panelId === 'friends') {
        const { loadFriendsPanel } = await import('./friends.js')
        await loadFriendsPanel()
    }
}

export async function loadUserData() {
    const res = await api('/me/profile', 'GET', null, true)
    if (!res.ok) return

    const p = res.data
    const initials = (p.display_name || '?').charAt(0).toUpperCase()

    document.getElementById('sidebar-name').textContent = p.display_name
    document.getElementById('sidebar-email').textContent = p.email

    const avatarEl = document.getElementById('sidebar-avatar')
    if (p.avatar_url) {
        avatarEl.innerHTML = `<img src="${escapeHtml(p.avatar_url)}" alt="avatar">`
    } else {
        avatarEl.textContent = initials
    }

    const badge = document.getElementById('spotify-badge')
    const connectCard = document.getElementById('connect-card')
    const disconnectBtn = document.getElementById('disconnect-btn')

    if (p.spotify_connected) {
        badge.className = 'spotify-badge connected'
        badge.innerHTML = '<span class="status-dot online"></span> Spotify connected'
        if (connectCard) connectCard.style.display = 'none'
        if (disconnectBtn) disconnectBtn.style.display = 'inline-flex'
        await loadSnapshotHistory()
    } else {
        badge.className = 'spotify-badge disconnected'
        badge.innerHTML = '<span class="status-dot offline"></span> Not connected'
        if (connectCard) connectCard.style.display = 'flex'
        if (disconnectBtn) disconnectBtn.style.display = 'none'
        document.getElementById('tracks-panel-content').style.display = 'none'
    }
}

export async function loadSnapshotHistory() {
    setLoading(true)
    const res = await api('/me/snapshots', 'GET', null, true)
    setLoading(false)

    if (!res.ok) {
        toast('Failed to load snapshot history', 'error')
        return
    }

    const snapshots = res.data || []
    renderSnapshotPills(snapshots)
    document.getElementById('tracks-panel-content').style.display = 'block'

    if (snapshots.length > 0) {
        await selectSnapshot(snapshots[0].id)
    } else {
        clearTrackStats()
        document.getElementById('track-list').innerHTML =
            '<p class="empty-state">No snapshots yet — capture one to get started.</p>'
    }
}

function renderSnapshotPills(snapshots) {
    const container = document.getElementById('snapshot-pills')
    container.innerHTML = ''

    if (snapshots.length === 0) {
        container.innerHTML = '<span class="snapshot-empty">No snapshots captured yet</span>'
        return
    }

    snapshots.forEach((s, i) => {
        const btn = document.createElement('button')
        btn.className = 'snapshot-pill' + (i === 0 ? ' active' : '')
        btn.dataset.id = s.id
        const date = new Date(s.created_at)
        btn.textContent = date.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
        btn.title = date.toLocaleString()
        btn.onclick = () => selectSnapshot(s.id)
        container.appendChild(btn)
    })
}

export async function selectSnapshot(id) {
	document
		.querySelectorAll('.snapshot-pill')
		.forEach(p => p.classList.toggle('active', p.dataset.id === id))

	setLoading(true)
	const res = await api(`/me/snapshots/${id}`, 'GET', null, true)
	setLoading(false)

	if (!res.ok) {
		toast('Failed to load snapshot', 'error')
		return
	}

	currentSnapshot = res.data
	renderTracks()
}

export async function captureSnapshot() {
    setLoading(true)
    const timeRange = document.getElementById('time-range-select')?.value || 'medium_term'
    const res = await api(`/me/snapshots?time_range=${timeRange}`, 'POST', null, true)
    setLoading(false)

    if (res.ok && res.data.tracks) {
        toast('Snapshot captured!', 'success')
        // Reload the full history list so the new entry appears
        await loadSnapshotHistory()
    } else if (res.status === 502) {
        toast('Spotify is temporarily unavailable. Try again shortly.', 'error')
    } else if (res.status === 401) {
        const { logout } = await import('./auth.js')
        toast('Session expired', 'error')
        logout()
    } else {
        toast('Failed to capture snapshot', 'error')
    }
}

export function setTrackLimit(n) {
    trackLimit = n
    document.querySelectorAll('.track-limit-btn').forEach(b =>
        b.classList.toggle('active', Number(b.dataset.limit) === n)
    )
    if (currentSnapshot) renderTracks()
}

function renderTracks() {
	const snapshot = currentSnapshot
	const allTracks = snapshot.tracks ?? []
	const tracks = allTracks.slice(0, trackLimit)
	const total = allTracks.length
	const avgPop =
		total > 0 ? allTracks.reduce((s, t) => s + t.popularity, 0) / total : 0

	document.getElementById('stat-total').textContent = total
	document.getElementById('stat-avg-pop').textContent = avgPop.toFixed(1)
	document.getElementById('stat-updated').textContent = new Date(
		snapshot.created_at,
	).toLocaleDateString()

	const list = document.getElementById('track-list')
	list.innerHTML = ''

	if (tracks.length === 0) {
		list.innerHTML = '<p class="empty-state">No tracks in this snapshot.</p>'
		return
	}

	tracks.forEach((t, i) => {
		const row = document.createElement('div')
		row.className = 'track-row'

		let imageHtml = ''
		if (t.image_url && t.image_url !== '') {
			imageHtml = `<div class="track-img"><img src="${escapeHtml(t.image_url)}" alt="Album art" onerror="this.style.display='none'"></div>`
		} else {
			imageHtml = `<div class="track-img"><div class="track-img-placeholder">🎵</div></div>`
		}

		row.innerHTML = `
            <div class="track-num">${i + 1}</div>
            ${imageHtml}
            <div class="track-info">
                <div class="track-name">${escapeHtml(t.name)}</div>
                <div class="pop-bar-wrap">
                    <div class="pop-bar"><div class="pop-fill" style="width:${t.popularity}%"></div></div>
                    <span class="pop-label">${t.popularity}</span>
                </div>
            </div>`
		list.appendChild(row)
	})
}

function clearTrackStats() {
    document.getElementById('stat-total').textContent = '—'
    document.getElementById('stat-avg-pop').textContent = '—'
    document.getElementById('stat-updated').textContent = '—'
}

export function connectSpotify() {
    window.location.href = '/auth/spotify/login'
}

export async function disconnectSpotify() {
    if (!confirm('Disconnect your Spotify account?')) return
    const res = await api('/me/spotify', 'DELETE', null, true)
    if (res.ok) {
        const badge = document.getElementById('spotify-badge')
        badge.className = 'spotify-badge disconnected'
        badge.innerHTML = '<span class="status-dot offline"></span> Not connected'
        const connectCard = document.getElementById('connect-card')
        if (connectCard) connectCard.style.display = 'flex'
        document.getElementById('disconnect-btn').style.display = 'none'
        document.getElementById('tracks-panel-content').style.display = 'none'
        toast('Spotify disconnected', 'success')
    } else {
        toast('Failed to disconnect', 'error')
    }
}
