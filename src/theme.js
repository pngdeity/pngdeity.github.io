(function () {
  var STORAGE_KEY = "pg-theme";

  function stored() {
    try {
      return localStorage.getItem(STORAGE_KEY);
    } catch (e) {
      return null;
    }
  }

  function preferred() {
    var value = stored();
    if (value === "light" || value === "dark") {
      return value;
    }
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  function apply(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    var buttons = document.querySelectorAll(".pg-theme-toggle");
    for (var i = 0; i < buttons.length; i++) {
      buttons[i].setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
      buttons[i].setAttribute(
        "aria-label",
        theme === "dark" ? "Switch to light theme" : "Switch to dark theme",
      );
    }
  }

  apply(preferred());

  document.addEventListener("click", function (event) {
    var button = event.target.closest ? event.target.closest(".pg-theme-toggle") : null;
    if (!button) {
      return;
    }
    var next = document.documentElement.getAttribute("data-theme") === "dark" ? "light" : "dark";
    try {
      localStorage.setItem(STORAGE_KEY, next);
    } catch (e) {}
    apply(next);
  });
})();
