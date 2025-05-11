idCounter = 0;
transitionEvent = whichTransitionEvent();
notWS = null;

function whichTransitionEvent() {
    let t;
    let el = document.createElement('fakeelement');
    const animations = {
        'animation': 'animationend',
        'Oanimation': 'oanimationEnd',
        'Mozanimation': 'animationend',
        'Webkitanimation': 'webkitanimationEnd'
    };

    for (t in animations) {
        if (el.style[t] !== undefined) {
            return animations[t];
        }
    }
}

function showPopup(id) {
    let popup = document.getElementById(id);
    if (popup.classList.contains("remove")) {
        popup.classList.remove("remove");
    }
    if (!popup.classList.contains("show")) {
        popup.classList.add("show");
    }
}

function setText(id, t) {
    let regExpEmoji = /:([a-zA-Z0-9_]+):/;
    let regExpNextEmojiName = /(?<=:)[a-zA-Z0-9_]+(?=:)/;
    let final_text = t;
    let emoji = final_text.match(regExpNextEmojiName);
    while (emoji) {
        final_text = final_text.replace(regExpEmoji, "<img class='notifImage' src='/web/images/locked/emojis/" + emoji + ".png'>");
        emoji = final_text.match(regExpNextEmojiName);
    }
    document.getElementById(id).innerHTML = final_text;
}

function addPopup(text, where, ma) {
    id = "";
    if (!!ma) {
        id = ma;
    } else {
        id = "popup" + idCounter++;
    }
    let newNotif = document.createElement('span');
    let a1 = document.createAttribute("class");
    a1.value = "popup";
    let a2 = document.createAttribute("shouldredirect");
    a2.value = where;
    let a3 = document.createAttribute("id");
    a3.value = id;
    newNotif.setAttributeNode(a1);
    newNotif.setAttributeNode(a2);
    newNotif.setAttributeNode(a3);
    document.getElementById("popupWrapper").appendChild(newNotif);
    let elements = document.getElementsByClassName("popup");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("click", onClickForPopups, false);
        elements[i].WhatToDoWhenAPopupIsClicked = function () {
            let where = this.getAttribute("shouldredirect");
            if (where.length > 0) {
                console.log("where", where);
                location = where;
            }
        };
    }

    setText(id, text);
    showPopup(id);
    console.log("ID:", id);
    window.setTimeout("diePopup(\"" + id + "\")", 1000 * POPUPLIFETIME);
}

function diePopup(id) {
    let item = document.getElementById(id);
    console.log("DIE HAHA", id);
    if (!item.classList.contains("remove")) {
        item.classList.add("remove");
        item.addEventListener(transitionEvent, customFunction);
    }
}

function onClickForPopups(ev) {
    this.WhatToDoWhenAPopupIsClicked();
    if (!this.classList.contains("remove")) {
        this.classList.add("remove");
        this.addEventListener(transitionEvent, customFunction);
    }
}

function customFunction(ev) {
    this.removeEventListener(transitionEvent, customFunction);
    let elements = document.getElementById("popupWrapper");
    if (elements.contains(this)) {
        elements.removeChild(this);
    }
}

// function getNotifications() {
//     let xhr = new XMLHttpRequest();
//     xhr.open('GET', '/notifications', true);
//     xhr.send();
//     xhr.onreadystatechange = (e) => {
//         if (xhr.readyState === 4) {
//             if (xhr.status === 200) {
//                 let response = JSON.parse(xhr.responseText);
//                 for (let notification of response) {
//                     addPopup(notification[0], notification[1]);
//                 }
//             } else if (xhr.status === 400) {
//                 window.location = xhr.responseText;
//             } else {
//                 console.log(xhr.status, xhr.responseText);
//             }
//         }
//     };
// }

function connectNotifications() {
    notWS = new WebSocket("/notifications/");

    notWS.addEventListener("message", (event) => {
        const msg = JSON.parse(event.data);
        handleNotification(msg);

        notWS.send(JSON.stringify({notification_id: msg.id}));
    });
}

const conversionToEmoji = {
    1: "white_dust_small",
    2: "blue_dust_small",
    3: "yellow_dust_small",
    4: "purple_dust_small",
    5: "star_dust_small"
}

function handleNotification(n) {
    switch (n.type) {
        case 1: // currency conversion finished
            const text1 = `Your conversion for ${n.payload.amount} :${conversionToEmoji[n.payload.currency_id]}: is over!`;
            addPopup(text1, "/web/money/conversion", n.id);
            break;
        case 2: // new friend request
            const text2 = `<b>${n.payload.username}</b> has sent you a friend request.`;
            addPopup(text2, "/web/friends", n.id);
            break;
        case 3: // accepted friend request
            const text3 = `<b>${n.payload.username}</b> has accepted your friend request.`
            addPopup(text3, "/web/friends", n.id);
            break;
    }
}

//getNotifications();
// function handleVisibilityChange() {
//     if (!document.hidden) {
//         getNotifications();
//     }
// }
// document.addEventListener("visibilitychange", handleVisibilityChange, false);
window.onload = connectNotifications();
