let registered = false;

leftSketch = function (p) {
    leftP = p;
    p.setup = function () {
        canWidth = 256;
        canHeight = 280;
        can = p.createCanvas(canWidth, canHeight);
        can.mouseOut(unregister);
        can.mouseOver(register);
        bg_color = p.color(BG);
        dark = p.color(DARKC);
        light = p.color(LIGHTC);
        right = p.color(RIGHTC);
        other = p.color(OTHERLIGHTC);
        clickc = p.color(CLICKABLEC);
        hoverc = p.color(HOVERC);
        touch = is_touch_device4();
        leftobjects = [];
        current = undefined;
        console.log("left is getting");
        getgirllist();
    };
    p.draw = function () {
        p.background(bg_color);
        for (let obj of leftobjects) {
            if (obj.clickable && obj.clickTimer > 0) {
                obj.clickTimer--;
                if (obj.clickTimer === 0) {
                    obj.unclick();
                }
                obj.display();
            } else if (registered && obj.hoverable && obj.in()) { //found an "in"
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
    };

    p.mousePressed = function () {
        if (registered) {
            let x = p.mouseX;
            let y = p.mouseY;
            for (obj of leftobjects) {
                if (obj.clickable && obj.in(x, y)) {
                    obj.clicked();
                }
            }
        }
    };

};

function emptyLeft() {
    leftobjects = [];
}

function unregister() {
    registered = false;
}

function register() {
    registered = true;
}