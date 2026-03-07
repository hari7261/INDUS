const COMMANDS = [
  { "command": "ind about", "description": "show the INDUS runtime profile and release posture", "category": "core", "example": "ind about" },
  { "command": "ind dev bench", "description": "benchmark a native INDUS command with latency metrics", "category": "developer", "example": "ind dev bench --command \"ind sys stats\" --runs 5" },
  { "command": "ind dev cache", "description": "inspect or clear the INDUS execution cache", "category": "developer", "example": "ind dev cache" },
  { "command": "ind dev debug", "description": "print the current execution context and registry bindings", "category": "developer", "example": "ind dev debug" },
  { "command": "ind dev reload", "description": "reload registry, config, and terminal state into memory", "category": "developer", "example": "ind dev reload" },
  { "command": "ind dev report", "description": "write a diagnostics report for the current INDUS session", "category": "developer", "example": "ind dev report --output indus-report.json" },
  { "command": "ind dev watch", "description": "watch a path with polling and report INDUS-friendly change events", "category": "developer", "example": "ind dev watch --path . --seconds 5" },
  { "command": "ind docs", "description": "show documentation entry points and versioned docs", "category": "core", "example": "ind docs" },
  { "command": "ind doctor", "description": "run a full INDUS diagnostics sweep", "category": "core", "example": "ind doctor" },
  { "command": "ind env export", "description": "export managed INDUS environment values to a file", "category": "environment", "example": "ind env export --file indus-env.json" },
  { "command": "ind env import", "description": "import managed INDUS environment values from a file", "category": "environment", "example": "ind env import --file indus-env.json" },
  { "command": "ind env list", "description": "list environment variables managed by INDUS", "category": "environment", "example": "ind env list" },
  { "command": "ind env set", "description": "set a managed environment value without shell syntax", "category": "environment", "example": "ind env set INDUS_PROFILE production" },
  { "command": "ind env unset", "description": "remove a managed environment value from the INDUS state", "category": "environment", "example": "ind env unset INDUS_PROFILE" },
  { "command": "ind fs digest", "description": "compute a SHA-256 digest for a file using native INDUS tooling", "category": "filesystem", "example": "ind fs digest docs/index.html" },
  { "command": "ind fs find", "description": "find files or directories by substring without shell utilities", "category": "filesystem", "example": "ind fs find config --path ." },
  { "command": "ind fs inspect", "description": "inspect a file or directory with type, timestamps, and location data", "category": "filesystem", "example": "ind fs inspect README.md" },
  { "command": "ind fs size", "description": "measure total size for a file tree or a single file", "category": "filesystem", "example": "ind fs size ." },
  { "command": "ind fs sync", "description": "safely sync one directory tree into another", "category": "filesystem", "example": "ind fs sync docs docs-copy" },
  { "command": "ind fs tree", "description": "render a compact directory tree with a depth guard", "category": "filesystem", "example": "ind fs tree . --depth 2" },
  { "command": "ind net fetch", "description": "fetch a URL with native INDUS HTTP handling", "category": "network", "example": "ind net fetch https://example.com --method GET" },
  { "command": "ind net pingx", "description": "measure TCP reachability latency to a host and port", "category": "network", "example": "ind net pingx example.com --port 443" },
  { "command": "ind net ports", "description": "scan a local TCP port range and report responsive ports", "category": "network", "example": "ind net ports --from 8080 --to 8090" },
  { "command": "ind net scan", "description": "show local network interfaces or resolve a target host", "category": "network", "example": "ind net scan example.com" },
  { "command": "ind net status", "description": "report outbound network readiness and DNS availability", "category": "network", "example": "ind net status --url https://example.com" },
  { "command": "ind net trace", "description": "trace the soft network path with resolution and dial attempts", "category": "network", "example": "ind net trace example.com" },
  { "command": "ind pkg audit", "description": "validate installed INDUS packages against the local catalog", "category": "package", "example": "ind pkg audit" },
  { "command": "ind pkg install", "description": "install an INDUS-native package from the internal catalog", "category": "package", "example": "ind pkg install aurora-kit" },
  { "command": "ind pkg list", "description": "list installed INDUS-native packages", "category": "package", "example": "ind pkg list" },
  { "command": "ind pkg remove", "description": "remove an installed INDUS-native package", "category": "package", "example": "ind pkg remove aurora-kit" },
  { "command": "ind pkg search", "description": "search the internal INDUS package catalog", "category": "package", "example": "ind pkg search aurora" },
  { "command": "ind pkg update", "description": "update an installed INDUS-native package to the catalog version", "category": "package", "example": "ind pkg update aurora-kit" },
  { "command": "ind proj build", "description": "build the active project into an INDUS artifact bundle", "category": "project", "example": "ind proj build ." },
  { "command": "ind proj clean", "description": "clean build artifacts and local INDUS project cache data", "category": "project", "example": "ind proj clean ." },
  { "command": "ind proj create", "description": "create a new INDUS project scaffold", "category": "project", "example": "ind proj create orbit-app --dir ." },
  { "command": "ind proj init", "description": "initialize the current directory as an INDUS project", "category": "project", "example": "ind proj init --name orbit-app" },
  { "command": "ind proj list", "description": "list discovered INDUS projects beneath a path", "category": "project", "example": "ind proj list --path ." },
  { "command": "ind proj run", "description": "run the active INDUS project manifest in simulation mode", "category": "project", "example": "ind proj run ." },
  { "command": "ind scan", "description": "scan the local INDUS environment and summarize readiness", "category": "core", "example": "ind scan" },
  { "command": "ind status", "description": "show the active terminal, workspace, cache, and package status", "category": "core", "example": "ind status" },
  { "command": "ind sys clean", "description": "clean INDUS cache, reports, and stale runtime artifacts", "category": "system", "example": "ind sys clean" },
  { "command": "ind sys doctor", "description": "validate the core runtime, registry, docs, and writable paths", "category": "system", "example": "ind sys doctor" },
  { "command": "ind sys info", "description": "show platform, process, and directory information", "category": "system", "example": "ind sys info" },
  { "command": "ind sys stats", "description": "show fast local runtime and memory statistics", "category": "system", "example": "ind sys stats" },
  { "command": "ind sys watch", "description": "sample runtime stats on a short interval", "category": "system", "example": "ind sys watch --interval 500ms --count 4" },
  { "command": "ind term clearx", "description": "clear the INDUS terminal viewport without using shell commands", "category": "terminal", "example": "ind term clearx" },
  { "command": "ind term doctor", "description": "check terminal theme, history, and interactive readiness", "category": "terminal", "example": "ind term doctor" },
  { "command": "ind term history", "description": "show recent native INDUS command history", "category": "terminal", "example": "ind term history --limit 10" },
  { "command": "ind term reset", "description": "reset theme, cache, and transient terminal state", "category": "terminal", "example": "ind term reset" },
  { "command": "ind term speed", "description": "show recent command latency and cache hit information", "category": "terminal", "example": "ind term speed" },
  { "command": "ind term theme", "description": "set or inspect the active INDUS prompt theme", "category": "terminal", "example": "ind term theme saffron" },
  { "command": "ind version", "description": "print the INDUS Terminal release, build, and registry version", "category": "core", "example": "ind version" },
  { "command": "ind work archive", "description": "archive the active workspace into a zip bundle", "category": "workspace", "example": "ind work archive ." },
  { "command": "ind work clean", "description": "clean transient workspace artifacts and inactive state", "category": "workspace", "example": "ind work clean ." },
  { "command": "ind work init", "description": "register a directory as an INDUS workspace", "category": "workspace", "example": "ind work init orbit-space" },
  { "command": "ind work list", "description": "list known INDUS workspaces and their status", "category": "workspace", "example": "ind work list" },
  { "command": "ind work pin", "description": "pin the active workspace for quick reuse", "category": "workspace", "example": "ind work pin orbit-space" },
  { "command": "ind work switch", "description": "switch the active workspace without shell cd semantics", "category": "workspace", "example": "ind work switch orbit-space" }
];

const CATEGORY_ORDER = [
  "core",
  "system",
  "project",
  "environment",
  "filesystem",
  "network",
  "developer",
  "package",
  "terminal",
  "workspace"
];

const CATEGORY_META = {
  core: "Platform-level lifecycle and diagnostics commands.",
  system: "Runtime and host health checks.",
  project: "Project scaffold, build, run, and maintenance commands.",
  environment: "Managed environment state commands.",
  filesystem: "Native file and directory operations.",
  network: "Network diagnostics and connection tooling.",
  developer: "Benchmarking, reporting, and debug operations.",
  package: "Internal package lifecycle commands.",
  terminal: "Terminal state, theme, and history commands.",
  workspace: "Workspace setup and movement commands."
};

const state = {
  query: "",
  category: "all"
};

const refs = {
  search: document.querySelector("[data-command-search]"),
  categorySidebar: document.querySelector("[data-category-sidebar]"),
  sectionContainer: document.querySelector("[data-command-sections]"),
  pageToc: document.querySelector("[data-page-toc]"),
  totalCount: document.querySelector("[data-total-count]"),
  visibleCount: document.querySelector("[data-visible-count]"),
  categoryCount: document.querySelector("[data-category-count]")
};

function escapeHtml(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

function prettyCategory(category) {
  return category.charAt(0).toUpperCase() + category.slice(1);
}

function getCategoryCounts(commands) {
  const counts = new Map();
  commands.forEach((item) => {
    counts.set(item.category, (counts.get(item.category) || 0) + 1);
  });
  return counts;
}

function filterCommands() {
  const needle = state.query.trim().toLowerCase();
  return COMMANDS.filter((item) => {
    if (state.category !== "all" && item.category !== state.category) {
      return false;
    }
    if (!needle) {
      return true;
    }
    return [item.command, item.description, item.example, item.category]
      .join(" ")
      .toLowerCase()
      .includes(needle);
  });
}

function groupCommands(commands) {
  const grouped = new Map();
  commands.forEach((item) => {
    if (!grouped.has(item.category)) {
      grouped.set(item.category, []);
    }
    grouped.get(item.category).push(item);
  });
  return CATEGORY_ORDER
    .filter((category) => grouped.has(category))
    .map((category) => ({
      category,
      items: grouped.get(category).slice().sort((a, b) => a.command.localeCompare(b.command))
    }));
}

function renderSidebar(categoryCounts) {
  const parts = [];
  parts.push(
    `<button type="button" class="category-link ${state.category === "all" ? "active" : ""}" data-category="all"><span>All commands</span><span>${COMMANDS.length}</span></button>`
  );

  CATEGORY_ORDER.forEach((category) => {
    const count = categoryCounts.get(category) || 0;
    if (!count) {
      return;
    }
    parts.push(
      `<button type="button" class="category-link ${state.category === category ? "active" : ""}" data-category="${category}"><span>${prettyCategory(category)}</span><span>${count}</span></button>`
    );
  });

  refs.categorySidebar.innerHTML = parts.join("");
}

function renderSections(groupedCategories) {
  if (!groupedCategories.length) {
    refs.sectionContainer.innerHTML = `
      <div class="empty-state">
        No commands matched this filter. Try a shorter search term or switch back to "All commands".
      </div>
    `;
    return;
  }

  refs.sectionContainer.innerHTML = groupedCategories
    .map(({ category, items }) => {
      const itemsHtml = items
        .map((item) => `
          <article class="command-item" id="cmd-${item.command.replaceAll(" ", "-")}">
            <div class="command-item-top">
              <h4 class="command-name">${escapeHtml(item.command)}</h4>
              <button type="button" class="copy-btn" data-copy="${escapeHtml(item.example)}">Copy example</button>
            </div>
            <p class="command-desc">${escapeHtml(item.description)}</p>
            <div class="example-block">
              <p class="example-label">Example</p>
              <pre class="example-code"><code>${escapeHtml(item.example)}</code></pre>
            </div>
          </article>
        `)
        .join("");

      return `
        <section class="category-section" id="cat-${category}">
          <div class="category-section-header">
            <div>
              <h3>${prettyCategory(category)}</h3>
              <p class="category-meta">${CATEGORY_META[category]}</p>
            </div>
            <span class="category-count">${items.length} commands</span>
          </div>
          <div class="command-list">
            ${itemsHtml}
          </div>
        </section>
      `;
    })
    .join("");
}

function renderToc(groupedCategories) {
  if (!groupedCategories.length) {
    refs.pageToc.innerHTML = "<li>No sections visible</li>";
    return;
  }

  refs.pageToc.innerHTML = groupedCategories
    .map(({ category }) => `<li><a href="#cat-${category}">${prettyCategory(category)}</a></li>`)
    .join("");
}

function renderStats(filteredCommands, groupedCategories) {
  refs.totalCount.textContent = String(COMMANDS.length);
  refs.visibleCount.textContent = String(filteredCommands.length);
  refs.categoryCount.textContent = `${groupedCategories.length} categories visible`;
}

function render() {
  const categoryCounts = getCategoryCounts(COMMANDS);
  const filteredCommands = filterCommands();
  const groupedCategories = groupCommands(filteredCommands);

  renderSidebar(categoryCounts);
  renderSections(groupedCategories);
  renderToc(groupedCategories);
  renderStats(filteredCommands, groupedCategories);
}

function bindEvents() {
  refs.search.addEventListener("input", (event) => {
    state.query = event.target.value;
    render();
  });

  refs.categorySidebar.addEventListener("click", (event) => {
    const button = event.target.closest("[data-category]");
    if (!button) {
      return;
    }
    state.category = button.getAttribute("data-category");
    render();
  });

  refs.sectionContainer.addEventListener("click", (event) => {
    const button = event.target.closest("[data-copy]");
    if (!button) {
      return;
    }
    const text = button.getAttribute("data-copy");
    navigator.clipboard.writeText(text).then(() => {
      const original = button.textContent;
      button.textContent = "Copied";
      button.classList.add("copied");
      setTimeout(() => {
        button.textContent = original;
        button.classList.remove("copied");
      }, 1200);
    });
  });
}

bindEvents();
render();
