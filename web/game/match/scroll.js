function scrollToPointLeft(part, center, d) {
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0);
    let currentPos = document.scrollingElement.scrollLeft;
    let stuff = document.scrollingElement.scrollWidth;
    let point;
    if (center) {
        point = stuff * part - vw / 2;
    } else {
        point = part;
    }
    if (point < 0) {
        point = 0;
    }
    if (currentPos === point) return;
    let duration = 500;
    if (d) {
        duration = d;
    }
    let base = (currentPos + point) * 0.5;
    let difference = currentPos - base;
    let startTime = performance.now();

    function step() {
        let normalizedTime = (performance.now() - startTime) / duration;
        if (normalizedTime > 1) normalizedTime = 1;
        window.scrollTo(base + difference * Math.cos(normalizedTime * Math.PI), document.scrollingElement.scrollTop);
        if (normalizedTime < 1) window.requestAnimationFrame(step);
    }

    window.requestAnimationFrame(step);
}

function scrollToPointUp(part, center, d) {
    let currentPos = document.scrollingElement.scrollTop;
    let stuff = document.scrollingElement.scrollHeight;
    let point = part;
    if (point < 0) {
        point = 0;
    }
    if (currentPos === point) return;
    let duration = 350;
    if (d) {
        duration = d;
    }
    let base = (currentPos + point) * 0.5;
    let difference = currentPos - base;
    let startTime = performance.now();

    //console.log("duration", duration, "base", base, "difference", difference, "target", point, "current", currentPos);

    function step() {
        let normalizedTime = (performance.now() - startTime) / duration;
        if (normalizedTime > 1) normalizedTime = 1;

        window.scrollTo(document.scrollingElement.scrollLeft, base + difference * Math.cos(normalizedTime * Math.PI));
        if (normalizedTime < 1) window.requestAnimationFrame(step);
    }

    window.requestAnimationFrame(step);
}
