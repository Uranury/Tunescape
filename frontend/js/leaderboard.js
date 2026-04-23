import { toast, escapeHtml } from './api.js'

let activeFeature = null

export async function loadLeaderboard(feature) {
  activeFeature = feature
  document.querySelectorAll('.feature-tab').forEach(t =>
    t.classList.toggle('active', t.dataset.feature === feature)
  )

  const body = document.getElementById('lb-body')
  body.innerHTML = '<tr><td colspan="3" class="empty-state">Loading…</td></tr>'

  try {
    const res = await fetch(`/leaderboards/${feature}?limit=10`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()

    if (!data.entries?.length) {
      body.innerHTML = '<tr><td colspan="3" class="empty-state">No scores yet — analyze your taste first.</td></tr>'
      return
    }

    body.innerHTML = data.entries.map(e => `
      <tr>
        <td class="lb-rank">#${e.rank}</td>
        <td>${escapeHtml(e.display_name || e.user_id.slice(0, 8))}</td>
        <td class="lb-score">${e.score.toFixed(3)}</td>
      </tr>`).join('')
  } catch (err) {
    body.innerHTML = `<tr><td colspan="3" class="empty-state" style="color:var(--red)">Error: ${err.message}</td></tr>`
  }
}

export function getActiveFeature() { return activeFeature }
