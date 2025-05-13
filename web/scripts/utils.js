async function UpdateFreeData(after) {
    const profResponse = await fetch("/api/v1/profile");
    if (!profResponse.ok) {
        if (profResponse.status === 401) {
            window.location.href = "/auth/signin";
        }
        return;
    }

    const prof = await profResponse.json();

    const response = await fetch("/api/v1/money/balance");
    if (!response.ok) {
        if (response.status === 401) {
            window.location.href = "/auth/signin";
        }
        return;
    }

    const res = await response.json();

    console.log(res);
    let welcome = "Welcome, " + prof.username;
    if (welcome.length > 19) {
        document.getElementById("username").innerText = "Hi, " + prof.username;
    } else {
        document.getElementById("username").innerText = welcome;
    }
    if (!!document.getElementById("moneytable")) {
        document.getElementById("wDust").innerText = res[1] ? res[1] : 0;
        document.getElementById("bDust").innerText = res[2] ? res[2] : 0;
        document.getElementById("yDust").innerText = res[3] ? res[3] : 0;
        document.getElementById("pDust").innerText = res[4] ? res[4] : 0;
        document.getElementById("sDust").innerText = res[5] ? res[5] : 0;
    }
    if (!!after) {
        after(res);
    }

    // let xhr = new XMLHttpRequest();
    // xhr.open('GET', '/freeinfo', true);
    // xhr.send();
    // xhr.onreadystatechange = (e) => {
    //     if (xhr.readyState === 4) {
    //         if (xhr.status === 200) {
    //             let response = JSON.parse(xhr.responseText);
    //             console.log(response);
    //             let welcome = "Welcome, " + response.Username;
    //             if (welcome.length > 19) {
    //                 document.getElementById("username").innerText = "Hi, " + response.Username;
    //             } else {
    //                 document.getElementById("username").innerText = welcome;
    //             }
    //             if (!!document.getElementById("moneytable")) {
    //                 document.getElementById("wDust").innerText = response.MoneyInfo["w"];
    //                 document.getElementById("bDust").innerText = response.MoneyInfo["b"];
    //                 document.getElementById("yDust").innerText = response.MoneyInfo["y"];
    //                 document.getElementById("pDust").innerText = response.MoneyInfo["p"];
    //                 document.getElementById("sDust").innerText = response.MoneyInfo["s"];
    //             }
    //             if (!!after) {
    //                 after(response);
    //             }
    //         } else {
    //             console.log(xhr.responseText);
    //         }
    //     }
    // };
}

async function UpdateProfileData(r) {
    const response = await fetch("/api/v1/profile");
    if (!response.ok) {
        if (response.status === 401) {
            window.location.href = "/auth/signin";
        }
        return;
    }

    const res = await response.json();
    document.getElementById("username2").innerText = res.username;

    // let xhr = new XMLHttpRequest();
    // xhr.open('GET', '/profileinfo', true);
    // xhr.send();
    // xhr.onreadystatechange = (e) => {
    //     if (xhr.readyState === 4) {
    //         if (xhr.status === 200) {
    //             let response = JSON.parse(xhr.responseText);
    //             document.getElementById("username2").innerText = r.Username;
    //             if (response.BattlesTotal > 0) {
    //                 document.getElementById("battles2").innerText = "Battle stats: " + response.BattlesWon + "/" + response.BattlesTotal + " (" + roundUp(response.BattlesWon / response.BattlesTotal * 100) + "% winrate)";
    //             } else {
    //                 document.getElementById("battles2").innerText = "Battle stats: " + 0 + "/" + 0 + " (" + 0 + "% winrate)";
    //             }
    //         } else {
    //             console.log(xhr.responseText);
    //         }
    //     }
    // };
}

function is_touch_device4() {

    let prefixes = ' -webkit- -moz- -o- -ms- '.split(' ');

    let mq = function (query) {
        return window.matchMedia(query).matches;
    };

    if (('ontouchstart' in window) || window.DocumentTouch && document instanceof DocumentTouch) {
        return true;
    }

    // include the 'heartz' as a way to have a non matching MQ to help terminate the join
    // https://git.io/vznFH
    let query = ['(', prefixes.join('touch-enabled),('), 'heartz', ')'].join('');
    return mq(query);
}

function getCookie(name) {
    let nameEQ = name + "=";
    let ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}

function setCookie(name, value, hrs) {
    let expires = "";
    if (hrs) {
        let date = new Date();
        date.setTime(date.getTime() + (hrs * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

function calculateLines(hoverText, width, size) {
    if (!width) {
        width = 290;
    }
    if (!size) {
        size = 15;
    }
    let height = 0;
    textSize(size);
    for (let line of hoverText.split("\n")) {
        height += 1;
        let x_pos = 0;
        for (let word of line.split(" ")) {
            /*strokeWeight(5);
            stroke(red2);
            point(mouseX+10+x_pos + change, mouseY+height*size);*/
            if (x_pos + textWidth(word + " ") < width) { //we are still on that line
                x_pos += textWidth(word + " ")
            } else { //start a new line
                height += 1;
                x_pos = textWidth(word + " ");
            }
        }
    }
    return height;
}

function isLight(colour) {
    return lightness(colour) > 50;
}

function roundUp(num) {
    return Math.round(num * 10) / 10
}

function setwhere(where) {
    redirectwhere = where;
}

function redirect(yesno) {
    doredirect = yesno;
}

function countdown(value, where, yesno) {
    timeleft = value - 1;
    if (timeleft < 0) {
        if (yesno) {
            console.log("RELOCATED!~", where);
            window.location = where;
            return;
        }
        isTicking = false;
        console.log("countdown times out");
        displayTimer(-1);
    } else {
        displayTimer(timeleft);
        console.log(timeleft);
        window.setTimeout("countdown(timeleft, redirectwhere, doredirect)", 1000);
    }
}

function parseSeconds(n, strip) {
    let mins = Math.floor(n / 60);
    let rem_secs = n - mins * 60;
    let hrs = Math.floor(mins / 60);
    let rem_mins = mins - hrs * 60;
    let days = Math.floor(hrs / 24);
    let rem_hrs = hrs - days * 24;
    let full = "";
    if (days > 0) {
        if (days === 1) {
            full += days + " day";
        } else {
            full += days + " days";
        }
        if (strip) {
            if (days > 364) {
                return "-1"
            } else {
                return full
            }
        } else {
            full += ", "
        }
    }
    if (rem_hrs > 0) {
        if (rem_hrs === 1) {
            full += rem_hrs + " hour";
        } else {
            full += rem_hrs + " hours";
        }
        if (strip) {
            return full
        } else {
            full += ", "
        }
    }
    if (rem_mins > 0) {
        if (rem_mins === 1) {
            full += rem_mins + " minute";
        } else {
            full += rem_mins + " minutes";
        }
        if (strip) {
            return full
        } else {
            full += ", "
        }
    }
    if (rem_secs >= 0) {
        if (rem_secs === 1) {
            full += rem_secs + " second";
        } else {
            full += rem_secs + " seconds";
        }
        if (strip) {
            return "less than a minute"
        }
    }

    return full;
}

async function addFriend(name, fromFriendList) {
    console.log("add: " + name);

    const response = await fetch(`/api/v1/friends/${name}`, {
        method: "POST"
    });
    if (!response.ok) {
        if (response.status == 401) {
            window.location.href = "/auth/signin"
        }
        return;
    }

    if (fromFriendList) {
        init();
    }
}

async function logout() {
    const response = await fetch("/api/v1/auth/signout", {
        method: "POST"
    });
    if (!response.ok) {
        error = await response.json();
        return;
    }

    window.location.href = "/auth/signin/";
}

var scaleFactor = 1.0;

window.addEventListener("load", () => {
    function updatebody() {
        const body = document.getElementsByTagName("body")[0]
        const bodyw = body.offsetWidth;
        const bodyh = body.offsetHeight;
        const ww = window.innerWidth;
        const wh = window.innerHeight;

        const scale = Math.min(ww / bodyw, wh / bodyh);
        body.style['transform'] = `translate(-50%, 0) scale(${scale})`;
        console.log("set body style", body.style['transform'], `scale(${scale})`, ww, wh, bodyw, bodyh);
        scaleFactor = scale;
    }

    window.onresize = updatebody;
    updatebody();

    new ResizeObserver(updatebody).observe(document.querySelector("body"));
})

function fixCoordScale(c) {
    return c / scaleFactor;
}

function testMatch() {
    if (window.location.pathname === "/game/match/") {
        return;
    }

    let loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    } else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.host + "/api/v1/game/match";
    const ws = new WebSocket(new_uri);

    ws.onopen = function (evt) {
        console.log("OPEN");
        ws.send(JSON.stringify({ type: 2, payload: {} }));
        connected = true;
    };
    
    ws.onmessage = function (evt) {
        let battleresponse = JSON.parse(evt.data);
        if (battleresponse.type !== 3) {
            window.location.href = "/game/match";
        }
        ws.close();
    };
}

// testMatch();
