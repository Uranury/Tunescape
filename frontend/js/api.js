// Global state
window.currentAccessToken = null
window.trendsChart = null

let _toastTimer = null

export function toast(message, type = 'info') {
  let el = document.getElementById('toast')
  if (!el) {
    el = document.createElement('div')
    el.id = 'toast'
    el.className = 'toast'
    document.body.appendChild(el)
  }
  el.textContent = message
  el.className = `toast ${type}`
  requestAnimationFrame(() => el.classList.add('show'))
  clearTimeout(_toastTimer)
  _toastTimer = setTimeout(() => el.classList.remove('show'), 3500)
}

export function escapeHtml(str) {
  const d = document.createElement('div')
  d.textContent = str
  return d.innerHTML
}

export async function api(url, method, body = null, auth = false) {
  const headers = { 'Content-Type': 'application/json' }
  if (auth && window.currentAccessToken)
    headers['Authorization'] = `Bearer ${window.currentAccessToken}`
  const opts = { method, credentials: 'same-origin', headers }
  if (body) opts.body = JSON.stringify(body)
  const res = await fetch(url, opts)
  const ct = res.headers.get('content-type') || ''
  const data = ct.includes('application/json') ? await res.json() : await res.text()
  return { ok: res.ok, status: res.status, data }
}

export function setLoading(show) {
  document.querySelectorAll('.loading-bar').forEach(el =>
    el.classList.toggle('active', show)
  )
}
