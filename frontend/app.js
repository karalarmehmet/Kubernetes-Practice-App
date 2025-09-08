const API_BASE = (window.API_BASE || "/api"); // behind Ingress, /api routes to the API service

async function fetchJSON(url, opts) {
  const res = await fetch(url, { headers: { "Content-Type": "application/json" }, ...opts });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

const listEl = document.getElementById("list");
const statusEl = document.getElementById("status");

async function refresh() {
  statusEl.textContent = "Fetching…";
  try {
    const items = await fetchJSON(`${API_BASE}/messages`);
    listEl.innerHTML = "";
    if (items.length === 0) {
      listEl.innerHTML = '<li class="muted">No messages yet</li>';
    } else {
      for (const m of items) {
        const li = document.createElement("li");
        li.textContent = `#${m.id} ${m.text} — ${new Date(m.createdAt).toLocaleString()}`;
        listEl.appendChild(li);
      }
    }
    statusEl.textContent = "OK";
  } catch (e) {
    statusEl.textContent = `Error: ${e.message}`;
  }
}

const form = document.getElementById("msg-form");
form.addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const text = document.getElementById("text").value.trim();
  if (!text) return;
  try {
    await fetchJSON(`${API_BASE}/messages`, { method: "POST", body: JSON.stringify({ text }) });
    document.getElementById("text").value = "";
    await refresh();
  } catch (e) {
    alert(e.message);
  }
});

refresh();
