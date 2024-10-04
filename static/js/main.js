import "../css/main.css";

import "htmx.org";
import "alpinejs";

window.htmx = require("htmx.org");
window.Alpine = require("alpinejs").default;
window.dateFns = require("date-fns");

function htmxInit() {
  // Always swap content after HTMX request regardless of response status.
  document.body.addEventListener("htmx:beforeSwap", function (evt) {
    evt.detail.shouldSwap = true;
    evt.detail.isError = false;
  });
}

function alpineInit() {
  // local-time directive converts a UTC timestamp to a human-readable local time.
  window.Alpine.directive("local-time", (el) => {
    const dateString = el.getAttribute("data-time");
    const isToday = dateFns.isToday(dateString);
    const isTomorrow = dateFns.isTomorrow(dateString);
    const time = dateFns.format(dateString, "h:mm a");

    let day = "";
    if (isToday) {
      day = "Today";
    } else if (isTomorrow) {
      day = "Tomorrow";
    } else {
      day = dateFns.format(dateString, "MMMM d");
    }
    el.innerText = `${day} at ${time}`;
    el.classList.remove("opacity-0");
  });

  window.Alpine.start();
}

htmxInit();
alpineInit();
