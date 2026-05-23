import { api, toast, setLoading } from './api.js'

const MAX_PLAYLIST_TRACKS = 10

export async function loadPlaylistPreview() {
	const container = document.getElementById('playlist-track-list')
	if (!container) return

	container.innerHTML = '<p class="empty-state">Loading tracks…</p>'

	const listRes = await api('/me/snapshots', 'GET', null, true)
	if (!listRes.ok || !listRes.data?.length) {
		container.innerHTML =
			'<p class="empty-state">No snapshots yet — capture one first.</p>'
		return
	}

	const latestID = listRes.data[0].id
	const snapRes = await api(`/me/snapshots/${latestID}`, 'GET', null, true)
	if (!snapRes.ok) {
		container.innerHTML = '<p class="empty-state">Failed to load tracks.</p>'
		return
	}

	const tracks = (snapRes.data.tracks || []).slice(0, MAX_PLAYLIST_TRACKS)
	if (tracks.length === 0) {
		container.innerHTML = '<p class="empty-state">No tracks in snapshot.</p>'
		return
	}

	container.innerHTML = tracks
		.map((t, i) => {
			let imageHtml = ''
			if (t.image_url && t.image_url !== '') {
				imageHtml = `<div class="track-img"><img src="${escapeHtml(t.image_url)}" style="width:32px;height:32px;border-radius:4px;" onerror="this.style.display='none'"></div>`
			} else {
				imageHtml = `<div class="track-img"><div class="track-img-placeholder" style="width:32px;height:32px;font-size:1rem;">🎵</div></div>`
			}
			return `
            <div class="track-row" style="padding:4px 8px;">
                <span class="track-num" style="width:25px;">${i + 1}</span>
                ${imageHtml}
                <span class="track-name" style="font-size:0.82rem;">${escapeHtml(t.name)}</span>
            </div>
        `
		})
		.join('')
}

export async function createPlaylist() {
    setLoading(true)
    const res = await api('/me/playlists/top-tracks', 'POST', null, true)
    setLoading(false)

    if (res.ok) {
        toast('Playlist created! Opening in Spotify…', 'success')
        renderPlaylistResult(res.data)
        window.open(res.data.external_url, '_blank')
    } else if (res.status === 404) {
        toast('No snapshot found — capture one first', 'error')
    } else if (res.status === 422) {
        toast('Connect Spotify first', 'error')
    } else if (res.status === 502) {
        toast('Spotify is temporarily unavailable. Try again shortly.', 'error')
    } else if (res.status === 401) {
        const { logout } = await import('./auth.js')
        toast('Session expired', 'error')
        logout()
    } else {
        toast('Failed to create playlist', 'error')
    }
}

function renderPlaylistResult(data) {
    const container = document.getElementById('playlist-result')
    if (!container) return

    const embedId = 'playlist-embed-' + Date.now()

    container.innerHTML = `
    <div class="playlist-card">
      <div class="playlist-card-info">
        <div class="playlist-card-name">${escapeHtml(data.name)}</div>
        <a class="playlist-card-link" href="${escapeHtml(data.external_url)}" target="_blank" rel="noopener">
          Open in Spotify ↗
        </a>
      </div>
      <div id="${embedId}" style="height:352px;border-radius:8px;margin-top:0.75rem;background:#f0f4f8;display:flex;align-items:center;justify-content:center;color:#888;font-size:0.85rem;">
        Loading player…
      </div>
    </div>
  `

    setTimeout(() => {
        const slot = document.getElementById(embedId)
        if (!slot) return
        const iframe = document.createElement('iframe')
        iframe.src = data.embed_url
        iframe.width = '100%'
        iframe.height = '352'
        iframe.frameBorder = '0'
        iframe.allow = 'autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture'
        iframe.style.cssText = 'border-radius:8px;display:block;'
        slot.replaceWith(iframe)
    }, 3000)
}

function escapeHtml(str) {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
}