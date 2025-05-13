function GetInfo() {
    // let xhr = new XMLHttpRequest();
    // xhr.open('POST', '/afterbattle', true);
    // xhr.send();
    // xhr.onreadystatechange = (e) => {
    //     if (xhr.readyState === 4) {
    //         if (xhr.status === 200) {
    //             let stuff = JSON.parse(xhr.responseText);
    //             console.log(stuff);
    //             ParseMatches(stuff);
    //         } else {
    //             console.log(xhr.responseText);
    //         }
    //     }
    // };

    var res = window.localStorage.getItem("reward");
    if (res) {
        console.log(res);
        res = JSON.parse(res);
    }
    var state = window.localStorage.getItem("lastState");
    if (state) {
        console.log(state);
        state = JSON.parse(state);
    }
    ParseMatches(res, state);
}

function minusExtra(matches, to_add, level, curr_matches) {
    let a = to_add;
    let l = level;
    let c = curr_matches;
    while (a > 0) {
        a--;
        c--;
        if (c < 0) {
            l--;
            c = matches[l-1] - 1;
        }
    }
    return [l, c]
}

function ParseMatches(info, state) {
    switch (info.result) {
        case 1:
            result = "Victory! ★";
            resColour = win;
            break;
        case 2:
            result = "Defeat...";
            resColour = dark;
            break;
        case 3:
            result = "Draw!";
            resColour = light;
            break;
        // case 4:
        //     result = "Opponent gave up! ★";
        //     resColour = win;
        //     break;
        // case 5:
        //     result = "I gave up...";
        //     resColour = dark;
        //     break;
        default:
            result = "No new rewards...";
            break;
    }
    if (info.opponent && info.opponent.username && info.opponent.username.length > 0) {
        oppName = info.opponent.username;
        // if (!info.AreFriends) {
        //     let size = 20;
        //     let t = "Add friend";
        //     strokeWeight(1);
        //     textSize(size);
        //     let w = textWidth(t);
        //     objects.push(new StandardButton(x + initial_w - w - 10, y + 125 + size, 5, t, size, info.LastOpponentsName));
        //     getElement(info.LastOpponentsName).clicked = function () {
        //         this.colour = this.clickedColour;
        //         this.clickTimer = this.clickLinger;
        //         addFriend(this.id);
        //         removeElement(this.id);
        //     };
        // }
    }
    globalinfo = info;
    const rarity = characters[state.character.number].rarity;
    matches = MATCHES.get(rarity);
    if (info.level === info.prev_level && info.experience === info.prev_experience) {
        to_add = 0;
    } else {
        to_add = 1;
    }
    // to_add = info.ToAdd;
    // level = info.Level;
    // curr_matches = info.CurrentMatches;
    // if (to_add > 0) {
    //     let edited = minusExtra(matches, to_add, level, curr_matches);
    //     level = edited[0];
    //     curr_matches = edited[1];
    // }
    level = info.prev_level;
    curr_matches = info.prev_experience;
    percentage = 100 * curr_matches / matches[level - 1];
    getElement("bar").setPercentage(percentage);
    getElement("bar").setNewPercentage(percentage);
    // if ("w" in info.Dusts) {
    //     console.log("hasDust!");
    //     dust = info.Dusts["w"];
    // } else {
    //     console.log("no Dust!");
    // }
    dust = info.reward;
    gname = characters[state.character.number].name;
}

function setup() {
    can = createCanvas(600, 350);
    can.parent('rewards_sketch');
    touch = is_touch_device4();
    document.addEventListener("keydown", e => {
        if (e.code === 'Enter' || e.code === "Space") {
            e.preventDefault();
            getElement("girlList").clicked();
        }
    });
    bg_color = color(BG);
    dark = color(DARKC);
    light = color(LIGHTC);
    right = color(RIGHTC);
    win = color(WINC);
    rectangle = dark;
    matches = [];
    level = 1;
    curr_matches = 0;
    to_add = 0;
    dust = 0;
    added = 0;
    levelled_up = false;
    addingSpeed = 0.5;
    gname = "";
    oppName = "";
    result = "";
    resColour = dark;
    current = undefined;
    objects = [];
    x = 50;
    y = 100;
    initial_w = 500;
    objects.push(new LoadingBar(x, y, initial_w, 50, 10, "bar", color(dark.toString()), right));
    let size = 30;
    let t = "Battle again";
    strokeWeight(1);
    textSize(size);
    let w = textWidth(t);
    objects.push(new StandardButton(x + (initial_w - w - 10) / 2, y + 195, 5, t, size, "girlList"));
    getElement("girlList").clicked = function () {
        this.colour = this.clickedColour;
        this.clickTimer = this.clickLinger;
        window.location = "/game/characters";
    };
    UpdateFreeData();
    GetInfo();
}

function draw() {
    background(bg_color);
    let bar = getElement("bar");
    let percentage = bar.percentage;
    let new_percentage = bar.newPercentage;
    if (new_percentage > percentage) {
        percentage = percentage + addingSpeed;
        if (percentage > 100) {
            bar.setPercentage(100);
            bar.setNewPercentage(100);
        } else {
            bar.setPercentage(percentage);
        }

    } else if (percentage >= 100 && level < 20) {
        level = level + 1;
        levelled_up = true;
        if (level < 20) {
            bar.setPercentage(0.0);
            bar.setNewPercentage(0.0);
        }
    } else if (to_add > added) {
        add_match();
    }

    for (let obj of objects) {
        if (obj.clickable && obj.clickTimer > 0) {
            obj.clickTimer--;
            if (obj.clickTimer === 0) {
                obj.unclick();
            }
            obj.display();
        } else if (obj.hoverable && obj.in()) { //found an "in"
            if (!current) { //outside to something
                current = obj;
                obj.hovered();
                obj.display();
            } else if (current.id === obj.id) { //currently hovered
                obj.display();
            } else { //switched from another 2 this
                current.unhovered();
                current = obj;
                obj.hovered();
                obj.display();
            }
        } else if (obj.hoverable && current && obj.id === current.id) { //went outside
            obj.unhovered();
            current = undefined;
            obj.display();
        } else {
            obj.display();
        }
    }
    drawText(getElement("bar").stopColour);
}

function drawText(stop) {
    let bar = getElement("bar");
    strokeWeight(1);
    fill(stop);
    noStroke();
    textSize(30);
    textAlign(LEFT);
    text("Lv. " + level, x, y - 10);
    textAlign(CENTER);
    text(gname, x + 200, y - 10);
    textAlign(RIGHT);
    if (level === 20) {
        bar.setPercentage(100);
        bar.setNewPercentage(0);
        fill(255, 195, 13);
        text("★ Max ★", (initial_w - x) / 2 + 250, y - 10);
        right = color(130, 190, 255)
    } else {
        text(round(bar.percentage) + "%", (initial_w - x) / 2 + 310, y - 10);
        if (levelled_up) {
            fill(win);
            text("Lvl up! ★", (initial_w - x) / 2 + 250, y - 10);
        }
    }
    textAlign(LEFT);
    fill(dark);
    text("+ Dust: " + dust, x, y + 100);
    if (oppName) {
        textAlign(RIGHT);
        textSize(25);
        text("Opponent: " + oppName, x + initial_w, y + 130);
    }
    textAlign(CENTER);
    textSize(40);
    fill(resColour);
    stroke(resColour);
    text(result, 300, 40);
}

function add_match() {
    curr_matches = curr_matches + 1;
    added += 1;
    let to_level_up = matches[level - 1];
    if (to_add > to_level_up && to_add > 1) {
        addingSpeed = map(added, 0, to_add, to_add / 4, 0.5);
    } else {
        addingSpeed = 0.5
    }
    let new_percentage = 100 * curr_matches / to_level_up;
    getElement("bar").setNewPercentage(new_percentage);
    if (to_level_up === curr_matches) {
        curr_matches = 0
    }
}

function mousePressed() {
    let x = fixCoordScale(mouseX);
    let y = fixCoordScale(mouseY);
    for (obj of objects) {
        if (obj.clickable && obj.in(x, y)) {
            obj.clicked();
        }
    }
}

function getElement(id) {
    for (obj of objects) {
        if (obj.id === id) {
            return obj
        }
    }
}

function removeElement(id) {
    for (let i = 0; i < objects.length; i++) {
        let obj = objects[i];
        if (obj.id === id) {
            objects.splice(i, 1);
            return;
        }
    }
}
