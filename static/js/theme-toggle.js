document.addEventListener("DOMContentLoaded", function() {
  var prefersDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;

  if (prefersDarkMode) {
    document.body.setAttribute("data-bs-theme", "dark");
  } else {
    document.body.setAttribute("data-bs-theme", "light");
  }
})

document.getElementById("theme-toggle").addEventListener("click", function () {
  var current = document.body.getAttribute("data-bs-theme");

  if (current === "dark" || current === null) {
    document.body.setAttribute("data-bs-theme", "light");
  } else {
    document.body.setAttribute("data-bs-theme", "dark");
  }
})
