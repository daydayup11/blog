export function renderNav(activePage) {
  const pages = [
    { href: '/index.html', label: '主页',    key: 'index'  },
    { href: '/blog.html',  label: '博客',    key: 'blog'   },
    { href: '/works.html', label: '作品集',  key: 'works'  },
    { href: '/garden.html',label: '知识花园',key: 'garden' },
  ];
  const links = pages.map(p =>
    `<a href="${p.href}" class="${p.key === activePage ? 'active' : ''}">${p.label}</a>`
  ).join('');
  return `
    <nav class="nav">
      <span class="nav-logo">dyt</span>
      <div class="nav-links">${links}</div>
      <button class="nav-theme-btn" id="theme-btn">🌙</button>
    </nav>`;
}

export function renderPostMeta({ word_count, views, reading_minutes, created_at }) {
  const date = new Date(created_at).toLocaleDateString('zh-CN');
  return `
    <div class="post-meta">
      <span>📅 ${date}</span>
      <span>📝 ${word_count ?? 0} 字</span>
      <span>👁 ${views ?? 0} 次浏览</span>
      <span>📖 约 ${reading_minutes ?? 1} 分钟</span>
    </div>`;
}

export function renderTags(tagsJSON) {
  const tags = JSON.parse(tagsJSON || '[]');
  return tags.map(t => `<span class="tag">${t}</span>`).join(' ');
}

export function renderBadges(techStackJSON) {
  const stack = JSON.parse(techStackJSON || '[]');
  return stack.map(t => {
    const cls = ['Go','Java','Python','Docker'].includes(t)
      ? `badge-${t.toLowerCase()}` : 'badge-other';
    return `<span class="badge ${cls}">${t}</span>`;
  }).join(' ');
}

export function getApi() {
  const isMock = new URLSearchParams(location.search).has('mock');
  return isMock
    ? import('/js/mock.js').then(m => m.api)
    : import('/js/api.js').then(m => m.api);
}
