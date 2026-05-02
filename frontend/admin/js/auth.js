export function requireAuth() {
  if (!localStorage.getItem('admin-token')) {
    location.href = '/admin/index.html';
  }
}

export function logout() {
  localStorage.removeItem('admin-token');
  location.href = '/admin/index.html';
}

export function renderAdminNav(activePage) {
  const pages = [
    { href: '/admin/dashboard.html', label: '统计',   key: 'dashboard' },
    { href: '/admin/posts.html',     label: '文章',   key: 'posts'     },
    { href: '/admin/columns.html',   label: '专栏',   key: 'columns'   },
    { href: '/admin/projects.html',  label: '项目',   key: 'projects'  },
    { href: '/admin/garden.html',    label: '花园',   key: 'garden'    },
  ];
  const links = pages.map(p =>
    `<a href="${p.href}" class="${p.key === activePage ? 'active' : ''}">${p.label}</a>`
  ).join('');
  return `
    <nav class="nav">
      <span class="nav-logo" style="font-size:14px">dyt <span style="color:var(--text-muted);font-weight:400">/ 后台</span></span>
      <div class="nav-links">${links}</div>
      <a href="/" style="font-size:13px;color:var(--text-muted);margin-right:8px">← 前台</a>
      <button class="nav-theme-btn" id="theme-btn">🌙</button>
      <button onclick="import('/admin/js/auth.js').then(m=>m.logout())"
        style="background:none;border:1px solid var(--border);border-radius:6px;padding:4px 10px;cursor:pointer;font-size:13px;color:var(--text-muted);margin-left:8px">
        退出
      </button>
    </nav>`;
}
