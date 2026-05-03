# 02 前端

## 页面结构

| 文件 | 路由 | 功能 |
|---|---|---|
| `index.html` | `/` | 主页：Hero、技能、时间线、精选项目、最新文章 |
| `blog.html` | `/blog.html` | 博客列表：文章 + 专栏两个 Tab |
| `post.html` | `/post.html?slug=xxx` | 单篇文章：正文、目录、上下篇 |
| `column.html` | `/column.html?slug=xxx` | 专栏详情：有序文章列表 |
| `works.html` | `/works.html` | 作品集：项目、演示、自定义板块 |
| `garden.html` | `/garden.html` | 知识花园：左侧板块导航 + 收藏卡片 |

## CSS 主题系统

主题切换通过 CSS 自定义变量实现，整个系统只需改 `<html>` 标签上的一个属性：

```html
<!-- 暗黑主题 -->
<html data-theme="dark">

<!-- 明亮主题 -->
<html data-theme="light">
```

每套主题在对应 CSS 文件里定义变量：

```css
/* theme-dark.css */
:root[data-theme="dark"] {
  --bg-primary:    #0a0a0a;
  --text-primary:  #e2e8f0;
  --accent:        #00ff88;
  /* ... */
}
```

所有组件用变量而不是硬编码颜色：

```css
.card {
  background: var(--bg-card);    /* 自动跟随主题 */
  border: 1px solid var(--border);
}
```

**主题切换逻辑**在 `js/theme.js`：
1. 页面加载时读 `localStorage` 里保存的偏好
2. 没有保存过则读系统偏好（`prefers-color-scheme`）
3. 用户手动切换时保存到 `localStorage`

## api.js — 统一的接口封装

所有 HTTP 请求都通过 `api.js` 发出，不在页面里直接写 `fetch`：

```js
// api.js 里的定义
export const api = {
  getPosts: (params = {}) => request('/posts?' + new URLSearchParams(params)),
  getPost:  (slug)        => request(`/posts/${slug}`),
  // ...
};

// 页面里这样用
const { posts } = await api.getPosts({ page: 1, tag: 'Go' });
```

好处：
- 接口地址集中管理，改一处全部生效
- 自动携带 JWT token（后台页面需要鉴权时）
- 统一的错误处理

## mock.js — 无需后端的开发模式

在 URL 后面加 `?mock=1` 参数，页面会加载 `mock.js` 而非调用真实 API：

```
http://localhost:3000/blog.html?mock=1
```

`mock.js` 返回和真实接口完全相同结构的假数据，让你在后端还没写好时就能开发和预览前端效果。

```js
// mock.js
export const api = {
  getPosts: async () => ({
    posts: [
      { id: 1, title: '测试文章', slug: 'test', ... },
    ],
    total: 1,
  }),
  // ...
};
```

## components.js — 共用渲染函数

避免每个页面重复写导航栏 HTML，统一放在 `components.js`：

```js
// 每个页面这样引入导航栏
import { renderNav, getApi } from '/js/components.js';
document.getElementById('nav-root').innerHTML = renderNav('blog'); // 传入当前页 key 高亮对应链接
```

其他共用函数：
- `renderTags(tagsJSON)` — 渲染标签列表
- `renderBadges(techStackJSON)` — 渲染技术栈徽章
- `renderPostMeta(...)` — 渲染字数/浏览/阅读时间
- `getApi()` — 自动判断是否使用 mock 模式

## 如何新增一个页面

1. 在 `frontend/` 创建 `newpage.html`
2. 复制任意现有页面的头部结构（`<head>` + FOUC 脚本 + CSS 引入）
3. 引入 `components.js` 和 `theme.js`
4. 在 `components.js` 的 `renderNav` 里加上新页面的导航项
5. 在后端 `router.go` 里加一行 `r.GET("/newpage.html", ...)`（或使用 NoRoute 兜底）

## 防止主题闪烁（FOUC）

每个 HTML 文件的 `<head>` 最顶部都有一段内联脚本，在任何 CSS 加载之前就设置主题：

```html
<script>
(function(){
  var t = localStorage.getItem('blog-theme') ||
    (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
  document.documentElement.setAttribute('data-theme', t);
})();
</script>
```

这段代码故意写成一行压缩形式，减少阻塞时间。
