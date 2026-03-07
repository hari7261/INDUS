(() => {
  const root = document.documentElement;
  const toggle = document.querySelector("[data-theme-toggle]");
  if (!toggle) {
    return;
  }

  const stored = localStorage.getItem("indus-docs-theme");
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

  const applyTheme = (theme) => {
    root.setAttribute("data-theme", theme);
    toggle.textContent = theme === "dark" ? "Light" : "Dark";
  };

  applyTheme(stored || (prefersDark ? "dark" : "light"));

  toggle.addEventListener("click", () => {
    const current = root.getAttribute("data-theme") === "dark" ? "dark" : "light";
    const next = current === "dark" ? "light" : "dark";
    localStorage.setItem("indus-docs-theme", next);
    applyTheme(next);
  });
})();
