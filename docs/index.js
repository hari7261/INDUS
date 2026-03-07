/* ── INDUS Terminal docs — script.js ── */

// ── NAV: scroll class + hamburger ──────────────────────────────────────────
const navbar    = document.getElementById('navbar');
const hamburger = document.getElementById('hamburger');
const navLinks  = document.getElementById('nav-links');

window.addEventListener('scroll', () => {
  navbar.classList.toggle('scrolled', window.scrollY > 20);
}, { passive: true });

hamburger.addEventListener('click', () => {
  hamburger.classList.toggle('open');
  navLinks.classList.toggle('open');
});

// Close mobile nav when a link is clicked
navLinks.querySelectorAll('a').forEach(a => {
  a.addEventListener('click', () => {
    hamburger.classList.remove('open');
    navLinks.classList.remove('open');
  });
});

// ── INSTALL TABS ────────────────────────────────────────────────────────────
document.querySelectorAll('.tab-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    const tab = btn.dataset.tab;
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
    btn.classList.add('active');
    document.getElementById(`tab-${tab}`).classList.add('active');
  });
});

// ── COPY BUTTONS ────────────────────────────────────────────────────────────
document.querySelectorAll('.copy-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    const text = btn.dataset.copy;
    if (!text) {
      // fallback: copy text from sibling code-block
      const block = btn.closest('.code-block-wrap')?.querySelector('.code-block');
      if (block) copyText(block.innerText, btn);
    } else {
      copyText(text, btn);
    }
  });
});

function copyText(text, btn) {
  navigator.clipboard.writeText(text).then(() => {
    const orig = btn.textContent;
    btn.textContent = 'Copied!';
    btn.style.background = 'var(--saffron)';
    btn.style.color = '#000';
    setTimeout(() => {
      btn.textContent = orig;
      btn.style.background = '';
      btn.style.color = '';
    }, 1800);
  }).catch(() => {
    // fallback for older browsers
    const el = document.createElement('textarea');
    el.value = text;
    el.style.position = 'fixed';
    el.style.opacity = '0';
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
  });
}

// ── HERO TERMINAL ANIMATION ─────────────────────────────────────────────────
const termLines = document.getElementById('term-lines');
const sequence = [
  { type: 'prompt', text: 'ind term theme saffron' },
  { type: 'out',    text: 'theme=saffron' },
  { type: 'prompt', text: 'ind sys stats' },
  { type: 'out',    text: 'runtime_go=go1.26.0  memory_alloc=288KB  cache_entries=1' },
  { type: 'prompt', text: 'ind proj create orbit-app --dir .' },
  { type: 'out',    text: 'project=orbit-app  root=.\\\\orbit-app' },
  { type: 'prompt', text: 'ind docs' },
  { type: 'ok',     text: '✓ commands.html • versions.html • v1.4.0 notes' },
];

let seqIdx = 0;
function typeNextLine() {
  if (seqIdx >= sequence.length) {
    // pause then restart
    setTimeout(() => {
      termLines.innerHTML = '';
      seqIdx = 0;
      setTimeout(typeNextLine, 800);
    }, 3500);
    return;
  }
  const item = sequence[seqIdx++];
  const div  = document.createElement('div');
  div.classList.add('term-line');
  if (item.type === 'prompt') {
    div.innerHTML = `<span class="term-line-prompt">INDUS ~/projects &gt;</span> `;
    div.classList.add('term-line');
  } else if (item.type === 'ok') {
    div.classList.add('term-line-ok');
  } else {
    div.classList.add('term-line-out');
  }
  termLines.appendChild(div);

  // type character by character for prompt lines
  if (item.type === 'prompt') {
    let ci = 0;
    const interval = setInterval(() => {
      if (ci < item.text.length) {
        div.innerHTML = `<span class="term-line-prompt">INDUS ~/projects &gt;</span> ${item.text.slice(0, ++ci)}`;
      } else {
        clearInterval(interval);
        setTimeout(typeNextLine, 300);
      }
    }, 28);
  } else {
    div.textContent = item.text;
    setTimeout(typeNextLine, 500);
  }
  termLines.scrollTop = termLines.scrollHeight;
}

// start animation after a short delay
setTimeout(typeNextLine, 1200);

// ── RELEASES FROM GITHUB API ─────────────────────────────────────────────────
const REPO = 'hari7261/INDUS';
const releasesContainer = document.getElementById('releases-list');

async function loadReleases() {
  try {
    const res = await fetch(`https://api.github.com/repos/${REPO}/releases`, {
      headers: { 'Accept': 'application/vnd.github+json' }
    });

    if (!res.ok) throw new Error(`GitHub API ${res.status}: ${res.statusText}`);
    const releases = await res.json();

    if (!releases.length) {
      // fall back to tags if no formal releases
      return loadFromTags();
    }

    renderReleases(releases);
  } catch (err) {
    // try tags as fallback
    loadFromTags().catch(() => showReleasesError(err.message));
  }
}

async function loadFromTags() {
  const res = await fetch(`https://api.github.com/repos/${REPO}/tags`, {
    headers: { 'Accept': 'application/vnd.github+json' }
  });
  if (!res.ok) throw new Error(`Tags API ${res.status}`);
  const tags = await res.json();
  renderTags(tags);
}

function renderReleases(releases) {
  releasesContainer.innerHTML = '';
  releases.forEach((rel, i) => {
    const card = document.createElement('div');
    card.className = 'release-card' + (i === 0 ? ' latest' : '');

    const date = new Date(rel.published_at || rel.created_at).toLocaleDateString('en-IN', {
      year: 'numeric', month: 'long', day: 'numeric'
    });

    const assetsHtml = rel.assets.map(a =>
      `<a class="asset-link" href="${a.browser_download_url}" target="_blank">
        ⬇ ${a.name} <span style="opacity:0.6;font-size:0.72em">${formatBytes(a.size)}</span>
      </a>`
    ).join('');

    const bodySnippet = rel.body
      ? rel.body.replace(/###/g,'').replace(/##/g,'').replace(/\*\*/g,'').slice(0, 280) + (rel.body.length > 280 ? '…' : '')
      : '';

    card.innerHTML = `
      <div class="release-header">
        <div class="release-title">
          <h3>${escHtml(rel.name || rel.tag_name)}</h3>
          ${i === 0 ? '<span class="release-badge latest">Latest</span>' : ''}
          ${rel.prerelease ? '<span class="release-badge pre">Pre-release</span>' : ''}
        </div>
        <span class="release-date">${date}</span>
      </div>
      ${assetsHtml ? `<div class="release-assets">${assetsHtml}</div>` : `
        <div class="release-assets">
          <a class="asset-link" href="${rel.html_url}" target="_blank">🔗 View on GitHub</a>
        </div>`}
      ${bodySnippet ? `<div class="release-body">${escHtml(bodySnippet)}</div>` : ''}
    `;
    releasesContainer.appendChild(card);
  });
}

function renderTags(tags) {
  releasesContainer.innerHTML = '';
  if (!tags.length) {
    showReleasesError('No releases found yet.');
    return;
  }
  tags.forEach((tag, i) => {
    const card = document.createElement('div');
    card.className = 'release-card' + (i === 0 ? ' latest' : '');
    card.innerHTML = `
      <div class="release-header">
        <div class="release-title">
          <h3>${escHtml(tag.name)}</h3>
          ${i === 0 ? '<span class="release-badge latest">Latest</span>' : ''}
        </div>
      </div>
      <div class="release-assets">
        <a class="asset-link" href="https://github.com/${REPO}/releases/tag/${encodeURIComponent(tag.name)}" target="_blank">🔗 View on GitHub</a>
        <a class="asset-link" href="https://github.com/${REPO}/archive/refs/tags/${encodeURIComponent(tag.name)}.zip" target="_blank">⬇ Source ZIP</a>
      </div>
    `;
    releasesContainer.appendChild(card);
  });
}

function showReleasesError(msg) {
  releasesContainer.innerHTML = `
    <div class="releases-error">
      <p>Could not load releases from GitHub. <a href="https://github.com/${REPO}/releases" target="_blank">View them directly on GitHub →</a></p>
      <p style="color:var(--text-dim);font-size:0.78rem;margin-top:0.5rem">${escHtml(msg)}</p>
    </div>`;
}

function formatBytes(bytes) {
  if (bytes > 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB';
  if (bytes > 1024) return (bytes / 1024).toFixed(0) + ' KB';
  return bytes + ' B';
}

function escHtml(str) {
  return String(str)
    .replace(/&/g,'&amp;')
    .replace(/</g,'&lt;')
    .replace(/>/g,'&gt;')
    .replace(/"/g,'&quot;');
}

loadReleases();

// ── COLOR PREVIEW MODAL ─────────────────────────────────────────────────────
const overlay = document.getElementById('modal-overlay');
const colorMap = {
  r: '#e74c3c', g: '#2ecc71', b: '#3498db', y: '#f1c40f',
  c: '#00bcd4', m: '#9b59b6', w: '#ecf0f1', o: '#FF9933',
  p: '#ff6b9d', d: '#00bcd4'
};

window.showColorDemo = function(key, name) {
  const hex = colorMap[key] || '#FF9933';
  document.getElementById('modal-prompt-text').textContent = `INDUS ~/projects >`;
  document.getElementById('modal-prompt-text').style.color = hex;
  document.getElementById('modal-desc').textContent =
    `${name} palette preview  •  use "ind term theme ..." inside INDUS to apply a runtime theme`;
  overlay.classList.add('open');

  // highlight selected swatch
  document.querySelectorAll('.color-swatch').forEach(s => s.classList.remove('active'));
  event.currentTarget.classList.add('active');
};

window.closeModal = function() {
  overlay.classList.remove('open');
};

document.addEventListener('keydown', e => {
  if (e.key === 'Escape') closeModal();
});

// ── SMOOTH ACTIVE NAV HIGHLIGHTING ──────────────────────────────────────────
const sections = document.querySelectorAll('section[id], header[id]');
const navAnchors = document.querySelectorAll('.nav-links a[href^="#"]');

const observer = new IntersectionObserver(entries => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      navAnchors.forEach(a => {
        a.style.color = a.getAttribute('href') === `#${entry.target.id}` ? 'var(--saffron)' : '';
      });
    }
  });
}, { rootMargin: '-40% 0px -55% 0px' });

sections.forEach(s => observer.observe(s));

// ── ANIMATE CARDS ON SCROLL ──────────────────────────────────────────────────
const cardObserver = new IntersectionObserver(entries => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      entry.target.style.opacity = '1';
      entry.target.style.transform = 'translateY(0)';
      cardObserver.unobserve(entry.target);
    }
  });
}, { threshold: 0.1 });

document.querySelectorAll(
  '.feature-card, .install-step, .release-card, .roadmap-col, .color-swatch'
).forEach(el => {
  el.style.opacity = '0';
  el.style.transform = 'translateY(20px)';
  el.style.transition = 'opacity 0.4s ease, transform 0.4s ease, border-color 0.2s ease, box-shadow 0.2s ease';
  cardObserver.observe(el);
});
