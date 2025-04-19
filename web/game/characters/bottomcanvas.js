bottomSketch = function (p) {
    bottomP = p;
    p.setup = function () {
        bottomobjects = [];
        let textS = 30;
        can2 = p.createCanvas(1024, 155);
        can2.parent('bottomcanvas');
        bg_color = p.color(BG);
        let b1 = new InterfaceButton(p, 5, 5, "Set as Main", textS, 'setmain', "B", 255);
        b1.onClick = setasmain;
        let b2 = new InterfaceButton(p, 5, 80, "Set as Secondary", textS, 'setsec', "B", 255);
        b2.onClick = setassec;
        let b3 = new InterfaceButton(p, 265, 5, "Clear set", textS, 'clear', "B", 255);
        b3.onClick = clearsetgirls;
        b3.clickable = clearClickable;
        let b4 = new InterfaceButton(p, 265, 80, "Choose random", textS, 'random', "B", 255);
        b4.onClick = random;
        let b5 = new InterfaceButton(p, 844, 42.5, "Battle", textS + 25, 'battle', "B", 175, 70);
        b5.onClick = battle;
        let t1 = new InterfaceText(p, 525, 42.5, right, "Select two characters, then press \"Battle\".", textS, 'prompts', "B", 314);
        bottomobjects.push(b1);
        bottomobjects.push(b2);
        bottomobjects.push(b3);
        bottomobjects.push(b4);
        bottomobjects.push(b5);
        bottomobjects.push(t1);
        bottomReady = true;
    };
    p.draw = function () {
        p.background(bg_color);
        for (let obj of bottomobjects) {
            if (obj.clickTimer > 0) {
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
    };

    p.mousePressed = function () {
        let x = p.mouseX;
        let y = p.mouseY;
        for (obj of bottomobjects) {
            if (obj.clickable && obj.in(x, y)) {
                obj.clicked();
            }
        }
    };

};

function getElementBottom(id) {
    for (obj of bottomobjects) {
        if (obj.id === id) {
            return obj
        }
    }
}