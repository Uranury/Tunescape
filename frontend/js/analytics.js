import { api, toast, setLoading, escapeHtml } from './api.js'

export async function analyzeTaste() {
  setLoading(true)
  const res = await api('/analytics/top-tracks', 'GET', null, true)
  setLoading(false)

  const grid = document.getElementById('taste-grid')
  if (!res.ok || !res.data) {
    toast('Analyze failed — make sure you have a snapshot first', 'error')
    grid.innerHTML = '<p class="empty-state">No data yet. Create a snapshot first.</p>'
    return
  }

  const a = res.data.averages
  grid.innerHTML = `
    <div class="taste-item">
      <div class="taste-item-label">Danceability</div>
      <div class="taste-item-value">${a.danceability.toFixed(3)}</div>
    </div>
    <div class="taste-item">
      <div class="taste-item-label">Valence</div>
      <div class="taste-item-value">${a.valence.toFixed(3)}</div>
    </div>
    <div class="taste-item">
      <div class="taste-item-label">Energy</div>
      <div class="taste-item-value">${a.energy.toFixed(3)}</div>
    </div>
    <div class="taste-item">
      <div class="taste-item-label">Acousticness</div>
      <div class="taste-item-value">${a.acousticness.toFixed(3)}</div>
    </div>`
  toast('Taste analyzed — scores pushed to leaderboards', 'success')
}

export async function loadTrends() {
  setLoading(true)
  const res = await api('/me/trends', 'GET', null, true)
  setLoading(false)

  const wrap = document.getElementById('trends-wrap')
  if (!res.ok || !res.data?.points?.length) {
    wrap.innerHTML = '<p class="empty-state">No trend data yet — create multiple snapshots over time.</p>'
    if (window.trendsChart) { window.trendsChart.destroy(); window.trendsChart = null }
    return
  }

  wrap.innerHTML = '<canvas id="trends-canvas"></canvas>'
  const pts = res.data.points
  const labels = pts.map(p => new Date(p.created_at).toLocaleDateString())
  const ctx = document.getElementById('trends-canvas').getContext('2d')
  if (window.trendsChart) window.trendsChart.destroy()
  window.trendsChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [
        { label: 'Danceability', data: pts.map(p => p.danceability), borderColor: '#F97316', tension: .3, fill: false },
        { label: 'Valence',      data: pts.map(p => p.valence),      borderColor: '#22C55E', tension: .3, fill: false },
        { label: 'Energy',       data: pts.map(p => p.energy),       borderColor: '#EF4444', tension: .3, fill: false },
        { label: 'Acousticness', data: pts.map(p => p.acousticness), borderColor: '#0078D7', tension: .3, fill: false },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      plugins: { legend: { labels: { font: { family: "'Segoe UI', sans-serif", size: 12 } } } },
      scales: {
        x: { grid: { color: '#E4EEF8' } },
        y: { grid: { color: '#E4EEF8' }, min: 0, max: 1 },
      },
    },
  })
}

export async function downloadReport() {
  setLoading(true)
  const res = await fetch('/me/report', {
    headers: { Authorization: `Bearer ${window.currentAccessToken}` },
    credentials: 'same-origin',
  })
  setLoading(false)
  if (res.ok) {
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'tunescape-report.pdf'
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
    toast('Report downloaded', 'success')
  } else {
    toast('Failed to generate report', 'error')
  }
}
