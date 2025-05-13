var bcanvas;
var backP = 3;

backSketch = function (p) {
    bg_color = p.color(BG);
    dark = p.color(LEFTC);
    light = p.color(LIGHTC);
    backP = p;

    p.setup = function () {
        bcanvas = p.createCanvas(1280, 550);
        bcanvas.position(0, 0);
        bcanvas.style("z-index", "-1");
        myChar = new InterfaceImage(p, 0, 0, "", "myChar", "", 0, 0);
        oppChar = new InterfaceImage(p, 730, 0, "", "oppChar", "", 0, 0);
        plscreen = new LoadingScreen(p, 0.4 * 550, 0.4 * 550, 0.2 * 550, 0.2 * 550);
        plscreen.setColours(PlC, PlW, PlE, PlR);
        oppscreen = new LoadingScreen(p, 730 + 0.4 * 550, 0.4 * 550, 0.2 * 550, 0.2 * 550);
        oppscreen.setColours(OC, OW, OE, OR);
        if (myChar.name !== PlN && getResolution(PlNum)[0] !== 0) {
            myChar.open("/images/locked/" + PlN + "_left.png", PlN, getResolution(PlNum)[0], getResolution(PlNum)[1]);
        } else if (myChar.name !== PlNum) {
            myChar.open("/images/locked/Placeholder_left.png", PlNum, 350, 550);
        }
        if (oppChar.name !== ON && getResolution(ONum)[0] !== 0) {
            oppChar.open("/images/locked/" + ON + "_right.png", ON, getResolution(ONum)[0], getResolution(ONum)[1]);
        } else if (oppChar.name !== ON) {
            oppChar.open("/images/locked/Placeholder_right.png", ONum, 350, 550);
        }
    };
    p.draw = function () {
        let p = backP;
        p.background(bg_color);
        if (myChar.loaded() && oppChar.loaded() && plscreen.stopped > 1 && oppscreen.stopped > 1) {
            let k;
            if (isLight(PlC)) {
                k = 0.1
            } else {
                k = 0.35
            }
            let pColour = p.lerpColor(color(dark.toString()), PlC, k);
            if (isLight(OC)) {
                k = 0.10
            } else {
                k = 0.35
            }
            let oppColour = p.lerpColor(color(dark.toString()), OC, k);
            setGradient(p, 0, 230, 550, 320, 5, bg_color, pColour);
            setGradient(p, 730, 230, 550, 320, 5, bg_color, oppColour);
            setGradient(p, 550, 230, 180, 320, 5, bg_color, dark);
            myChar.display();
            oppChar.display();
            p.noLoop();
            console.log("drew girls");
        } else {
            if (myChar.loaded() && plscreen.stopped < 1) {
                //console.log("setting my to 1");
                plscreen.stop();
            }
            if (oppChar.loaded() && oppscreen.stopped < 1) {
                //console.log("setting other to 1");
                oppscreen.stop();
            }
            plscreen.display();
            oppscreen.display();

        }
    };

};

function setMyChar(PlayerName, PlayerNum, c1, c2, c3, c4) {
    PlN = PlayerName;
    PlNum = PlayerNum;
    PlC = c1;
    PlW = c2;
    PlE = c3;
    PlR = c4;
}

function setOppChar(OppName, OppNum, c1, c2, c3, c4) {
    ON = OppName;
    ONum = OppNum;
    OC = c1;
    OW = c2;
    OE = c3;
    OR = c4;
}


function setGradient(p, x, y, w, h, r, c1, c2) {
    p.noFill();
    p.strokeWeight(2);
    //circle at the beginning
    for (let i = y; (i < y + r) && (i < y + h); i++) {
        let inter = p.map(i, y, y + h, 0, 1);
        let c = p.lerpColor(c1, c2, inter);
        p.stroke(c);
        let top_x = x + r - p.sqrt(p.sq(r) - p.sq(i - y - r));
        let bot_x = x + w - r + p.sqrt(p.sq(r) - p.sq(i - y - r));
        p.line(top_x, i, bot_x, i);
    }

    for (let i = y + r; i < y + h - r; i++) {
        let inter = p.map(i, y, y + h, 0, 1);
        let c = p.lerpColor(c1, c2, inter);
        p.stroke(c);
        p.line(x, i, x + w, i);
    }
    //circle at the end
    for (let i = y + h; (i >= y + r) && (i >= y + h - r); i--) {
        let inter = p.map(i, y, y + h, 0, 1);
        let c = p.lerpColor(c1, c2, inter);
        p.stroke(c);
        let top_x = x + r - p.sqrt(p.sq(r) - p.sq(i - y - h + r));
        let bot_x = x + w - r + p.sqrt(p.sq(r) - p.sq(i - y - h + r));
        p.line(top_x, i, bot_x, i);
    }
}