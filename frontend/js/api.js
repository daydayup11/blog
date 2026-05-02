const BASE = '/api';

async function request(path, options = {}) {
  const token = localStorage.getItem('admin-token');
  const headers = options.body instanceof FormData
    ? { ...options.headers }
    : { 'Content-Type': 'application/json', ...options.headers };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(BASE + path, { ...options, headers });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  getPosts:    (params = {}) => request('/posts?' + new URLSearchParams(params)),
  getPost:     (slug)        => request(`/posts/${slug}`),
  getTags:     ()            => request('/tags'),
  getColumns:  ()            => request('/columns'),
  getColumn:   (slug)        => request(`/columns/${slug}`),
  getProjects: ()            => request('/projects'),
  getGardenSections: ()                    => request('/garden/sections'),
  getWorksSections:  ()                    => request('/works/sections'),
  getGardenItems:    (sectionId, type = '', tag = '') => request(`/garden/items?section_id=${sectionId}${type ? '&type=' + type : ''}${tag ? '&tag=' + encodeURIComponent(tag) : ''}`),
  trackPage: (path, referrer = '') =>
    request('/track/page', { method: 'POST', body: JSON.stringify({ path, referrer }) }),
  trackPost: (id) =>
    request(`/track/post/${id}`, { method: 'POST' }),
  login: (username, password) =>
    request('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  adminStats: (days = 7) => request(`/admin/stats?days=${days}`),
  adminListPosts: (page = 1) => request(`/admin/posts?page=${page}`),
  adminGetPost: (slug) => request(`/admin/posts/${slug}`),
  adminCreatePost:  (data) => request('/admin/posts', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdatePost:  (id, data) => request(`/admin/posts/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeletePost:  (id) => request(`/admin/posts/${id}`, { method: 'DELETE' }),
  adminPinPost:     (id, pinned) => request(`/admin/posts/${id}/pin`, { method: 'PUT', body: JSON.stringify({ pinned }) }),
  adminPublishPost: (id, published) => request(`/admin/posts/${id}/publish`, { method: 'PUT', body: JSON.stringify({ published }) }),
  adminListColumns: () => request('/admin/columns'),
  adminListAllSections: () => request('/admin/garden/sections'),
  adminCreateWorksSection: (data) => request('/admin/works/sections', { method: 'POST', body: JSON.stringify(data) }),
  adminGetColumnPosts: (id) => request(`/admin/columns/${id}/posts`),
  adminCreateColumn: (data) => request('/admin/columns', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateColumn: (id, data) => request(`/admin/columns/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminPublishColumn: (id, published) => request(`/admin/columns/${id}`, { method: 'PUT', body: JSON.stringify({ published }) }),
  adminDeleteColumn: (id) => request(`/admin/columns/${id}`, { method: 'DELETE' }),
  adminSetColumnPosts: (id, postIds) => request(`/admin/columns/${id}/posts`, { method: 'PUT', body: JSON.stringify({ post_ids: postIds }) }),
  adminCreateProject: (data) => request('/admin/projects', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateProject: (id, data) => request(`/admin/projects/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteProject: (id) => request(`/admin/projects/${id}`, { method: 'DELETE' }),
  adminCreateSection: (data) => request('/admin/garden/sections', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateSection: (id, data) => request(`/admin/garden/sections/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteSection: (id) => request(`/admin/garden/sections/${id}`, { method: 'DELETE' }),
  adminCreateItem:  (data) => request('/admin/garden/items', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateItem:  (id, data) => request(`/admin/garden/items/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteItem:  (id) => request(`/admin/garden/items/${id}`, { method: 'DELETE' }),
  adminUpload: (file) => {
    const form = new FormData();
    form.append('file', file);
    return request('/admin/upload', { method: 'POST', headers: {}, body: form });
  },
};
