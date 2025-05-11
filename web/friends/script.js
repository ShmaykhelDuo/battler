friends = [];
incoming = [];
pending = [];

function init(first) {
    fetchFriends();
    fetchIncoming();
    fetchPending();

    if (first) {
        //displaying current tab
        document.getElementById('Friends').style.display = "block";
        //clicking on the top button
        document.getElementById('defaultTab').className += " active";
    }
}

async function fetchFriends() {
    const response = await fetch("/friends");
    if (!response.ok) {
        if (response.status == 401) {
            window.location.href = "/web/auth/signin"
        }
        return;
    }
    setFriends(await response.json());
}

async function fetchIncoming() {
    const response = await fetch("/friends/incoming");
    if (!response.ok) {
        if (response.status == 401) {
            window.location.href = "/web/auth/signin"
        }
        return;
    }
    setIncoming(await response.json());
}

async function fetchPending() {
    const response = await fetch("/friends/outgoing");
    if (!response.ok) {
        if (response.status == 401) {
            window.location.href = "/web/auth/signin"
        }
        return;
    }
    setPending(await response.json());
}

async function removeFriend(name, fromFriendList) {
    console.log("remove: " + name);

    const response = await fetch(`/friends/${name}`, {
        method: "DELETE"
    });
    if (!response.ok) {
        if (response.status == 401) {
            window.location.href = "/web/auth/signin"
        }
        return;
    }

    if (fromFriendList) {
        init();
    }
}

function openTab(evt, tabName) {
    let i, tabcontent, standardbutton;
    //hiding all the tab contents
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    //deactivating all buttons
    standardbutton = document.getElementsByClassName("standardbutton");
    for (i = 0; i < standardbutton.length; i++) {
        standardbutton[i].className = standardbutton[i].className.replace(" active", "");
    }
    //displaying current tab
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}

function sortFriends(a, b) {
    return (a.username > b.username) ? 1 : ((b.username > a.username) ? -1 : 0)

    // if (a[1] === "Offline" && b[1] === "Offline" || a[1] !== "Offline" && b[1] !== "Offline") {
    //     return (a[0] > b[0]) ? 1 : ((b[0] > a[0]) ? -1 : 0)
    // } else if (a[1] === "Offline") {
    //     return 1
    // } else {
    //     return -1
    // }
}

function setFriends(friends) {
    let innerHTML = "";
    if (friends.length > 0) {
        friends = friends.sort(sortFriends);
        innerHTML += "<table width='100%'>";
        for (let friend of friends) {
            innerHTML += "<tr><td width=\"33%\">";
            console.log("you have friends");
            innerHTML += friend.username + "</td>";
            // let middle = "<td>" + friend[1];
            // if (friend.length === 3) {
            //     if (friend[1] === "Playing as") {
            //         middle += " " + friend[2];
            //     } else if (friend[1] === "Offline") {
            //         let time = parseSeconds(parseInt(friend[2]), true);
            //         if (time === "-1") {
            //             console.log("so long omg", friend[0])
            //         } else {
            //             middle += " for " + time;
            //         }
            //     }
            // }
            // middle += "</td>";
            // innerHTML += middle;
            let postfix = "<div align='right'><button class=\"standardbutton\" onclick='removeFriend(\"" + friend.id + "\", true)'>✖</button></div>";
            innerHTML += "<td width='30px'>" + postfix + "</td>";
            innerHTML += "</tr>";
        }
        innerHTML += "</table>";
    } else {
        console.log("you have 0 friends uwu");
        innerHTML = "<div align=\"center\">You have 0 friends uwu</div>";
    }
    document.getElementById("friends").innerHTML = innerHTML;
}

function setIncoming(incoming) {
    let innerHTML = "";
    if (incoming.length > 0) {
        console.log("you have incoming");
        incoming = incoming.sort(sortFriends);
        innerHTML += "<table width='100%'>";
        for (let friend of incoming) {
            innerHTML += "<tr><td width=\"33%\">" + friend.username + "</td><td>";
            innerHTML += "<div align='right'><button class=\"standardbutton\" onclick='addFriend(\"" + friend.id + "\", true)'>✔</button></div></td>";
            innerHTML += "</tr>"
        }
        innerHTML += "</table>";
    } else {
        console.log("you have 0 inc uwu");
        innerHTML = "<div align=\"center\">No incoming requests.</div>"
    }
    document.getElementById("incoming").innerHTML = innerHTML;

}

function setPending(pending) {
    let innerHTML = "";
    if (pending.length > 0) {
        console.log("you have pending");
        pending = pending.sort(sortFriends);
        innerHTML += "<table width='100%'>";
        for (let friend of pending) {
            innerHTML += "<tr><td width=\"33%\">" + friend.username + "</td><td>";
            innerHTML += "<div align='right'><button class=\"standardbutton\" onclick='removeFriend(\"" + friend.id + "\", true)'>✖</button></div></td></tr>";
        }
        innerHTML += "</table>";
    } else {
        console.log("you have 0 pending uwu");
        innerHTML = "<div align=\"center\">No pending requests. Add someone! :></div>";
    }
    document.getElementById("pending").innerHTML = innerHTML;
}

// init(true);
UpdateFreeData();
function handleVisibilityChange() {
    if (!document.hidden) {
        init(false);
    }
}
document.addEventListener("visibilitychange", handleVisibilityChange, false);
