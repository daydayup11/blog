export const api = {
  getPosts: async ({ page = 1, tag = '', q = '' } = {}) => {
    let posts = [
      { id: 1, title: '置顶：Go 并发模式实践', slug: 'go-concurrency', summary: '深入理解 goroutine 和 channel 的设计模式，避免常见陷阱。', tags: '["Go","并发"]', is_pinned: true, word_count: 2400, created_at: '2026-04-10T10:00:00Z' },
      { id: 2, title: 'Docker 最佳实践总结', slug: 'docker-best-practices', summary: '从 Dockerfile 优化到多阶段构建，系统整理 Docker 使用心得。', tags: '["Docker","运维"]', is_pinned: false, word_count: 1800, created_at: '2026-04-01T10:00:00Z' },
      { id: 3, title: 'Java 虚拟线程初探', slug: 'java-virtual-threads', summary: 'Project Loom 带来的虚拟线程，彻底改变 Java 并发编程方式。', tags: '["Java"]', is_pinned: false, word_count: 1200, created_at: '2026-03-20T10:00:00Z' },
    ];
    if (tag) posts = posts.filter(p => JSON.parse(p.tags).includes(tag));
    if (q)   posts = posts.filter(p => p.title.includes(q) || (p.summary && p.summary.includes(q)));
    return { posts, total: posts.length, page };
  },
  getPost: async (slug) => ({
    post: { id: 1, title: 'Go 并发模式实践', slug, content: '', content_html: '<h2 id="intro">简介</h2><p>本文介绍 Go 并发模式...</p><h2 id="goroutine">Goroutine</h2><p>goroutine 是 Go 的轻量级线程...</p><h3 id="channel">Channel 用法</h3><p>channel 是 goroutine 间通信的桥梁。</p><pre><code class="language-go">ch := make(chan int)\ngo func() { ch &lt;- 42 }()\nfmt.Println(&lt;-ch)</code></pre>', tags: '["Go","并发"]', word_count: 2400, created_at: '2026-04-10T10:00:00Z' },
    views: 1234,
    reading_minutes: 8,
  }),
  getTags: async () => ({ tags: ['Go', 'Java', 'Python', 'Docker', '并发', '随笔'] }),
  getColumns: async () => ({
    columns: [
      { id: 1, slug: 'go-series', name: 'Go 从入门到实践', description: '系统学习 Go 语言，从语法基础到工程实践。', cover_url: '', is_published: true, created_at: '2026-01-01T00:00:00Z' },
      { id: 2, slug: 'system-design', name: '系统设计笔记', description: '分布式系统设计的思考与实践记录。', cover_url: '', is_published: true, created_at: '2026-02-01T00:00:00Z' },
    ],
  }),
  getColumn: async (slug) => ({
    column: { id: 1, slug, name: 'Go 从入门到实践', description: '系统学习 Go 语言，从语法基础到工程实践。共 5 篇。', cover_url: '' },
    posts: [
      { id: 10, title: '第一篇：Go 环境搭建与基础语法', slug: 'go-01', word_count: 1500 },
      { id: 11, title: '第二篇：Go 并发编程', slug: 'go-02', word_count: 2400 },
      { id: 12, title: '第三篇：Go 标准库精讲', slug: 'go-03', word_count: 2000 },
    ],
  }),
  getProjects: async () => ({
    projects: [
      { id: 1, name: 'blog', description: '本博客网站，Go + 原生前端，Docker 部署。', repo_url: 'https://github.com/daiyutong/blog', demo_url: '', tech_stack: '["Go","HTML","CSS","Docker"]', stars: 12, is_featured: true },
      { id: 2, name: 'go-tools', description: '日常开发工具集，Go 实现。', repo_url: 'https://github.com/daiyutong/go-tools', demo_url: '', tech_stack: '["Go"]', stars: 8, is_featured: true },
      { id: 3, name: 'data-pipeline', description: 'Python 数据处理管道，支持多数据源。', repo_url: 'https://github.com/daiyutong/data-pipeline', demo_url: '', tech_stack: '["Python"]', stars: 5, is_featured: false },
    ],
  }),
  getGardenSections: async () => ({
    sections: [
      { id: 1, slug: 'cs', name: '计算机知识', sort_order: 0 },
      { id: 2, slug: 'finance', name: '金融知识', sort_order: 1 },
      { id: 3, slug: 'psychology', name: '心理学知识', sort_order: 2 },
      { id: 4, slug: 'social', name: '我在哪里', sort_order: 3 },
    ],
  }),
  getGardenItems: async (sectionId, type = '') => {
    const all = {
      1: [
        { id: 1,  section_id: 1, type: 'article', title: 'Go 调度器原理深度解析', source: 'Go Blog', url: 'https://example.com', summary: '深入分析 GMP 模型和调度器工作原理。', embed_url: '' },
        { id: 2,  section_id: 1, type: 'video',   title: '如何做系统设计', source: 'YouTube', url: 'https://example.com', summary: '系统设计面试思路全解析。', embed_url: '' },
        { id: 3,  section_id: 1, type: 'book',    title: '《程序员的自我修养》', source: '豆瓣', url: 'https://example.com', summary: '链接、装载与库，深入理解程序运行原理。', embed_url: '' },
        { id: 4,  section_id: 1, type: 'website', title: 'roadmap.sh', source: 'roadmap.sh', url: 'https://roadmap.sh', summary: '开发者学习路径地图，覆盖前后端、DevOps 等方向。', embed_url: '' },
        { id: 5,  section_id: 1, type: 'book',    title: '《深入理解计算机系统》', source: '豆瓣', url: 'https://example.com', summary: 'CSAPP，计算机系统必读经典。', embed_url: '' },
      ],
      2: [
        { id: 6,  section_id: 2, type: 'book',    title: '《穷爸爸富爸爸》', source: '豆瓣', url: 'https://example.com', summary: '财务自由思维入门，改变金钱观的经典读物。', embed_url: '' },
        { id: 7,  section_id: 2, type: 'article', title: '指数基金投资入门', source: '雪球', url: 'https://example.com', summary: '定投指数基金的基本逻辑和操作方法。', embed_url: '' },
        { id: 8,  section_id: 2, type: 'video',   title: '巴菲特股东大会精选', source: 'YouTube', url: 'https://example.com', summary: '价值投资理念精华片段合集。', embed_url: '' },
      ],
      3: [
        { id: 9,  section_id: 3, type: 'book',    title: '《被讨厌的勇气》', source: '豆瓣', url: 'https://example.com', summary: '阿德勒心理学对话体入门，关于如何获得幸福的哲学。', embed_url: '' },
        { id: 10, section_id: 3, type: 'article', title: '认知偏见手册', source: 'Wikipedia', url: 'https://example.com', summary: '180+ 种认知偏见的系统整理，提升元认知能力。', embed_url: '' },
      ],
      4: [
        { id: 11, section_id: 4, type: 'account', title: 'GitHub', source: 'github.com/daiyutong', url: 'https://github.com', summary: '', embed_url: '' },
      ],
    };
    const items = all[sectionId] || [];
    return { items: type ? items.filter(i => i.type === type) : items };
  },
  trackPage: async () => null,
  trackPost: async () => null,
};
