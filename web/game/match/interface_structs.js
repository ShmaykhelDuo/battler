const hoverLinger = 300;
const taphoverLinger = 60;

function Panel(x, y, w, h, s) {
    this.x = x;
    this.y = y;
    this.width = w;
    this.height = h;
    this.smooth = s;
    this.objects = [];
    this.toplayerobjs = [];

    this.display = function () {
        for (obj of this.objects) {
            if (obj.clickable && obj.clickTimer > 0) {
                obj.clickTimer--;
                if (obj.clickTimer === 0) {
                    obj.unclick();
                }
                obj.display();
            } else if (!touch) {
                if (obj.hoverable && obj.in()) { //found an "in"
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
            } else {
                obj.display();
            }
        }
        //display the rect
        stroke(dark);
        noFill();
        strokeWeight(3);
        rect(this.x, this.y, this.width, this.height, this.smooth);
        for (let obj of this.toplayerobjs) {
            if (obj.clickable && obj.clickTimer > 0) {
                obj.clickTimer--;
                if (obj.clickTimer === 0) {
                    obj.unclick();
                }
                obj.display();
            } else if (!touch) {
                if (obj.hoverable && obj.in()) { //found an "in"
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
            } else {
                obj.display();
            }
        }
        //display hover
        if (!touch && current && current.in() || touch && current && current.isHovered && current.in()) {
            current.displayHover();
        }
    };

    this.in = function () {
        let x = fixCoordScale(mouseX);
        let y = fixCoordScale(mouseY);
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y))
    };

    this.add = function (obj) {
        this.objects.push(obj)
    };

    this.addTopLayer = function (obj) {
        this.toplayerobjs.push(obj)
    };

    this.discard = function (obj) {
        const index = this.objects.indexOf(obj);
        if (index > -1) {
            if (obj.clickable) {
                obj.unclick();
            }
            if (obj.hoverable) {
                obj.unhovered();
            }
            this.objects.splice(index, 1);
            return true;
        } else {
            const index = this.toplayerobjs.indexOf(obj);
            if (index > -1) {
                if (obj.clickable) {
                    obj.unclick();
                }
                if (obj.hoverable) {
                    obj.unhovered();
                }
                this.objects.splice(index, 1);
                return true;
            } else {
                return false;
            }
        }
    };


}

function CanvasImage(x, y, path, id, name, width, height, colour) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.path = path;
    if (!path) {
        this.width = width;
        this.height = height;
        if (colour) {
            this.loadedObj = {};
            this.colour = colour;
        } else {
            this.image = undefined;
        }
        this.name = "";
    } else {
        this.width = width;
        this.height = height;
        if (colour) {
            this.colour = colour;
            this.loadedObj = {};
            this.image = loadImage(path, img => {
                this.loadedObj.loaded = true;
                console.log(this.loadedObj);
            }, this.failed);
        } else {
            this.image = loadImage(path);
        }
        this.name = name;
    }
    this.hoverable = false;
    this.clickable = false;

    this.getImage = function () {
        return this.image;
    };

    this.getName = function () {
        return this.name;
    };

    this.copy = function () {
        let img = new CanvasImage(0, 0, "", "", "", this.width, this.height);
        img.image = this.image;
        img.path = this.path;
        return img;

    };

    this.open = function (path, name, width, height) {
        this.x = this.x + (550 - width) / 2;
        this.width = width;
        this.height = height;
        if (!!this.loadedObj) {
            this.image = loadImage(path, img => {
                this.loadedObj.loaded = true;
            }, this.failed);
        } else {
            this.image = loadImage(path);
        }
        this.name = name;
    };

    this.set = function (other) {
        this.id = other.id;
        this.x = other.x;
        this.y = other.y;
        this.path = other.path;
        this.colour = other.colour;
        if (!other.path) {
            this.loadedObj = {};
            this.width = other.width;
            this.height = other.height;
            this.image = undefined;
            this.name = other.name;
        } else {
            this.width = other.width;
            this.height = other.height;
            this.image = other.image;
            this.name = other.name;
            this.loadedObj = other.loadedObj;
        }
        this.hoverable = other.hoverable;
        this.clickable = other.clickable;
    };

    this.loaded = function () {
        if (!this.loadedObj) {
            return undefined;
        }
        return this.loadedObj.hasOwnProperty("loaded")
    };

    this.display = function () {
        if (this.image) {
            image(this.image, this.x, this.y, this.width, this.height);
        }
    }
}

function ImageBox() {
    this.images = [];
    this.isDisplayed = [];

    this.add = function (image) {
        this.images.push(image);
        this.isDisplayed.push(true);
    };

    this.addSimple = function (image) {
        this.images.push(image);
    };

    this.isTaken = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                return this.isDisplayed[i];
            }
        }
        return undefined;
    };

    this.take = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                this.isDisplayed[i] = true;
                return this.images[i];
            }
        }
        return undefined;
    };

    this.get = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                return this.images[i];
            }
        }
        return undefined;
    };

    this.clearDisplayed = function () {
        for (let i = 0; i < this.isDisplayed.length; i++) {
            this.isDisplayed[i] = false;
        }
    };

    this.contains = function (name) {
        for (let image of this.images) {
            if (image.name === name) {
                return true;
            }
        }
        return false;
    };

}

function TextInfo(x, y, colour, t, size, id, type, width, height, hoverable) {
    this.visible = true;
    this.id = id;
    this.x = x;
    this.y = y;
    this.textSize = size;
    this.type = type;

    if (this.type === "effects") {
        this.text = [];
        this.images = [];
        this.hoverTexts = [];
        this.hoverHeights = [];
        this.ids = [];
        this.width = width;
        this.height = 175;
    } else if (this.type === "info") {
        this.text = t;
        this.width = width;
        this.height = height;
    } else {
        if (this.id === "playerHP" || this.id === "oppHP") {
            this.HP = undefined;
            this.MaxHP = undefined;
            this.targetHP = undefined;
            this.speed = 0;
            this.framesLeft = 0;
            this.targetMaxHP = undefined;
            this.speedMax = 0;
            this.framesLeftMax = 0;
        }
        this.text = t;
        this.height = size;
        if (type === "B" || type === "C") {
            this.width = width;
        } else if (type === "A" || type === "D") {
            textSize(size);
            this.width = textWidth(t);
        } else {
            this.width = size * t.length;
        }
    }

    this.textColour = colour;
    this.clickable = false;
    this.hoverable = hoverable;
    if (this.hoverable) {
        this.isHovered = false;
        this.hoverTimer = 0;
        if (touch) {
            this.hoverLinger = taphoverLinger;

        } else {
            this.hoverLinger = hoverLinger;
        }
    }
    this.hoverText = "";


    this.hide = function () {
        this.visible = false;
        this.wasClickable = this.clickable;
        this.wasHoverable = this.hoverable;
        this.clickable = false;
        this.hoverable = false;
    };

    this.show = function () {
        this.visible = true;
        this.clickable = this.wasClickable;
        this.hoverable = this.wasHoverable;

    };

    this.display = function () {
        if (this.visible) {
            /*let temp_c = color(rightc.toString());
            temp_c.setAlpha(150);
            stroke(temp_c);
            strokeWeight(2);
            temp_c.setAlpha(50);
            fill(temp_c);
            if (this.type !== "info") {
                rect(this.x, this.y-this.height, this.width, this.height);
            } else {
                rect(this.x, this.y, this.width, this.height);
            }*/

            //do the rest
            if (this.type === "effects") {
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                let asc = textAscent();
                let desc = textDescent();
                if (this.id[this.id.length - 1] === "2") {
                    textAlign(RIGHT);
                    for (let i = this.text.length - 1; i >= 0; i--) {
                        this.images[i].display();
                        let y_pos = asc - desc + this.images[i].height / 4;
                        text(this.text[i], this.x + this.width - this.images[i].width - 5, this.y - this.height + this.images[i].height * (i) + y_pos);
                    }
                } else {
                    textAlign(LEFT);
                    for (let i = 0; i < this.text.length; i++) {
                        this.images[i].display();
                        let y_pos = asc - desc + this.images[i].height / 4;
                        text(this.text[i], this.x + this.images[i].width + 5, this.y - this.height + this.images[i].height * (i) + y_pos);
                    }
                }
            } else if (this.type === "info") {
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                textAlign(CENTER);
                text(this.text, this.x, this.y, this.width);
            } else if (this.type === "B" || this.type === "C") {
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                if (this.type === "B") {
                    textAlign(CENTER);
                } else {
                    textAlign(LEFT);
                    textLeading(33);
                }
                text(this.text, this.x, this.y, this.width);
            } else if (this.type === "A" || this.type === "D") {
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                if (this.type === "A") {
                    textAlign(LEFT);
                } else if (this.type === "D") {
                    textAlign(RIGHT);
                }
                text(this.text, this.x, this.y);
            } else {
                if ((this.id === "playerHP" || this.id === "oppHP") && this.framesLeft > 0) {
                    this.framesLeft--;
                    if (this.HP + this.speed > this.targetHP && this.speed < 0 ||
                        this.HP + this.speed < this.targetHP && this.speed > 0) { //if it's worth it yet
                        this.HP = this.HP + this.speed;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP)
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                    } else { //end
                        this.HP = this.targetHP;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP);
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                        this.framesLeft = 0;
                    }
                }
                if (((this.id === "playerHP" || this.id === "oppHP") && this.framesLeftMax > 0)) {
                    this.framesLeftMax -= 1;
                    if (this.MaxHP + this.speedMax > this.targetMaxHP && this.speedMax < 0 ||
                        this.MaxHP + this.speedMax < this.targetMaxHP && this.speedMax > 0) { //if it's worth it yet
                        this.MaxHP = this.MaxHP + this.speedMax;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP);
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                    } else { //end
                        this.MaxHP = this.targetMaxHP;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP);
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                        this.framesLeftMax = 0;
                    }
                }
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                textAlign(LEFT);
                text(this.text, this.x, this.y - 5);
            }
        }
    };

    this.in = function () {
        let x = fixCoordScale(mouseX);
        let y = fixCoordScale(mouseY);
        return (this.x <= x && x <= (this.width + this.x) && y <= this.y && (this.y - this.height) <= y);
    };

    this.hovered = function () {
        if (touch) {
            this.isHovered = true;
        }
    };

    this.unhovered = function () {
        this.hoverTimer = 0;

    };

    this.setColour = function (c) {
        this.textColour = color(c);
    };

    this.setText = function (t) {
        if (this.type === "info") {
            this.text = t;
        } else if (this.type === "effects") {
            this.text = [];
            this.images = [];
            this.hoverTexts = [];
            this.hoverHeights = [];
            let add = 0;
            let effIconSize = 35;
            for (let i = 0; i < t.length; i++) {
                if (this.id[this.id.length - 1] === "2") {
                    add = this.width - effIconSize;
                }

                if (t[i][1]) {
                    this.text.push(t[i][1]);
                }

                //preparing images!
                let name = t[i][0];
                //and btw also descriptions www
                this.hoverTexts.push(effectsDict[name].description);
                this.hoverHeights.push(calculateLines(effectsDict[name].description));
                //check if they are already in the box.
                if (!IMAGEBOX.contains(name)) {
                    //if they are not, add them (Adding a canvas image.).
                    let new_image = new CanvasImage(this.x + add, this.y - this.height + effIconSize * i,
                        "/images/locked/" + name + ".png", this.id + "_" + name, name, effIconSize,
                        effIconSize);
                    IMAGEBOX.add(new_image);
                    this.images.push(new_image);
                } else {
                    //if they are in the box, check if they are already displayed.
                    let image = IMAGEBOX.take(name);
                    if (IMAGEBOX.isTaken(name)) {
                        //if they are already displayed, make a copy with name + '_2' and add it.
                        //for now it's guaranteed there can't be more than 2 identical effect images.
                        let new_image = image.copy();
                        //now we gotta set x, y, id, name.
                        new_image.x = this.x + add;
                        new_image.y = this.y - this.height + effIconSize * i;
                        new_image.name = name + "_2";
                        new_image.id = this.id + "_" + new_image.name;
                        IMAGEBOX.add(new_image);
                        this.images.push(new_image);
                    } else {
                        //if they are not currently displayed, 'take' them.
                        //then change the image so that it would stay where you want it to.
                        image.x = this.x + add;
                        image.y = this.y - this.height + effIconSize * i;
                        this.images.push(image);
                    }
                }
            }
        } else {
            this.text = t;
            this.height = this.textSize;
            if (this.type && this.type !== "B" && this.type !== "C" || !this.type) {
                textSize(this.textSize);
                this.width = textWidth(this.text);
            }
        }
    };

    this.displayHover = function () {
        if (this.hoverTimer < this.hoverLinger) {
            this.hoverTimer += 1;
            return;
        }
        if (this.type === "effects" && this.images.length > 0) {
            this.hoverText = "";
            for (let i = 0; i < 5; i++) {
                if (i >= this.images.length) {
                    return;
                }
                if (Math.floor((fixCoordScale(mouseY) - this.y + this.height) / this.images[0].height) === i) {
                    this.hoverText = this.hoverTexts[i];
                    this.hoverLines = this.hoverHeights[i];
                    break;
                }
            }
        } else if (this.type === "effects") {
            return;
        }
        if ((this.id === "playerHP" || this.id === "oppHP") && this.hasOwnProperty("defenses")) {
            let hoverSize = 15;
            strokeWeight(0.5);
            textSize(hoverSize);
            noStroke();
            fill(hoverc);
            rect(fixCoordScale(mouseX), fixCoordScale(mouseY), 355, 70, 5);
            let y_pos = fixCoordScale(mouseY) + hoverSize + 5;
            let x_pos = fixCoordScale(mouseX) + 10;
            let map = new Map(COLOURIDS);
            for (let i = 0; i < COLOURIDS.length; i++) {
                //noStroke();
                textAlign(LEFT);
                let name = COLOURIDS[i][0];
                fill(map.get(name));
                stroke(map.get(name));
                text(name + ":", x_pos, y_pos); //Yellow:
                let amount = this.defenses[String(i + 1)]; //number
                let c;
                if (amount > 0) {
                    c = lerpColor(dark, rightc, amount / 5);
                } else {
                    c = lerpColor(dark, red2, amount / -5);
                }
                fill(c);
                stroke(c);
                textAlign(RIGHT);
                text(amount, x_pos + 75, y_pos);
                if ((i + 1) % 4 === 0) {
                    x_pos = fixCoordScale(mouseX) + 10;
                    y_pos += hoverSize + 5;
                } else {
                    x_pos += 85;
                }
            }
        } else if (this.hoverable) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines)
        }
    };

    this.startAnimation = function (speed, target) {
        this.speed = speed;
        this.targetHP = target;
        this.framesLeft = FRAMESFORANIMATIONS;
    };

    this.startAnimationMax = function (speed, target) {
        this.speedMax = speed;
        this.targetMaxHP = target;
        this.framesLeftMax = FRAMESFORANIMATIONS;
    }
}

function TurnLog(x, y, colour, size, id, width, height, hoverable) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.textSize = size;
    this.maxlen = 8;
    this.text = [];
    this.textColours = [];
    this.aligns = [];
    this.turns = [];
    this.len = 0;
    this.width = width;
    this.height = height;
    this.clickable = false;
    this.hoverable = hoverable;
    this.isHovered = false;
    this.hoverText = "";
    this.maxFramesDown = 45;
    this.maxFramesUp = 20;
    this.frame = 0;
    this.isTransitioning = false;
    this.transitioningUp = false;

    this.display = function () {
        strokeWeight(1);
        textSize(this.textSize);
        for (let i = this.len - 1; i >= 0; i--) {
            if (this.isTransitioning) {
                this.frame += 1;
                let c;
                let col;
                if (this.transitioningUp) {
                    col = lerpColor(bg_color, this.textColours[i], 1.05 - (this.len - i - 2) / (this.maxlen - 2));
                    c = lerpColor(col, this.textColours[i], this.frame / this.maxFramesUp);
                } else {
                    col = lerpColor(bg_color, this.textColours[i], 1.05 - (this.len - i - 2) / (this.maxlen - 2));
                    c = lerpColor(this.textColours[i], col, this.frame / this.maxFramesDown);
                }
                stroke(c);
                fill(c);
                if (this.frame > this.maxFramesUp && this.transitioningUp
                    || this.frame > this.maxFramesDown && !this.transitioningUp) {
                    this.isTransitioning = false;
                    this.frame = 0;
                }
            } else if (this.isHovered) {
                stroke(this.textColours[i]);
                fill(this.textColours[i]);
            } else {
                let col = lerpColor(bg_color, this.textColours[i], 1.05 - (this.len - i - 2) / (this.maxlen - 2));
                stroke(col);
                fill(col);
            }
            if (this.aligns[i]) {
                textAlign(LEFT);
                text(this.text[i], this.x + 5, this.y - (this.textSize) * (this.len - i - 1) - 5);
            } else {
                textAlign(RIGHT);
                text(this.text[i], this.x + this.width - 5, this.y - (this.textSize) * (this.len - i - 1) - 5);
            }
        }
    };

    this.push = function (text, colour, isMine) {
        if (this.len < this.maxlen) {
            this.text.push(text);
            //this.textColours.push(colour);
            this.textColours.push(lerpColor(color(dark.toString()), colour, 0.7));
            this.aligns.push(isMine);
            this.len++;
        } else {
            this.pop();
            this.push(text, colour, isMine);
        }

        console.log(this.text, this.aligns, this.textColours)
    };

    this.pop = function () {
        let len = this.len;
        if (len > 0) {
            this.text = this.text.slice(1, len);
            this.turns = this.turns.slice(1, len);
            this.textColours = this.textColours.slice(1, len);
            this.aligns = this.aligns.slice(1, len);
            this.len--;
        } else {
            console.log("popping when turn log is empty.")
        }
    };

    this.clear = function () {
        this.text.length = 0;
        this.turns.length = 0;
        this.textColours.length = 0;
        this.aligns.length = 0;
        this.len = 0;
    };

    this.in = function () {
        let x = fixCoordScale(mouseX);
        let y = fixCoordScale(mouseY);
        return (this.x <= x && x <= (this.width + this.x) && y <= this.y && (this.y - this.height) <= y);
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
            this.transitioningUp = true;
        } else {
            this.frame = Math.ceil((this.maxFramesDown - this.frame) / this.maxFramesDown * this.maxFramesUp);
            //this.frame = this.maxframes - this.frame;
        }
    };

    this.unhovered = function () {
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
            this.transitioningUp = false;
        } else {
            this.frame = Math.ceil((this.maxFramesUp - this.frame) / this.maxFramesUp * this.maxFramesDown);
            //this.frame = this.maxframes - this.frame;
        }
    };

    this.displayHover = function () {
        if (this.hoverText) {
            displayStandardHoverBubble(this.hoverText, calculateLines(this.hoverText));
        }
    }
}

function SkillButton(x, y, type, t, id, mine) {
    this.isMine = mine;
    this.id = id;
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = 50;
    this.textColour = color(dark.toString());
    this.borderColour = color(dark.toString());
    this.borderWidth = 5;
    this.type = type;
    if ((type === 1) || (type === 4)) {
        this.width = 82;
        this.height = 240;
    } else {
        this.width = 132;
        this.height = 210;
    }

    this.maxframes = 7;
    this.frame = 0;
    this.destColour = this.baseColour;
    this.previousColour = this.baseColour;
    this.isTransitioning = false;
    this.rawColour = "";
    this.rawText = "";
    this.colour = clickc;
    this.baseColour = clickc;
    this.hoverColour = clickc;
    this.clickedColour = clickc;
    this.isHovered = false;
    this.hoverText = "";
    this.clickable = mine;
    this.hoverable = true;
    this.clickLinger = 4;
    this.clickTimer = 0;
    if (this.hoverable) {
        this.hoverTimer = 0;
        if (touch) {
            this.hoverLinger = taphoverLinger;
        } else {
            this.hoverLinger = hoverLinger;
        }
    }

    this.display = function () {
        if (this.isTransitioning) {
            if (this.frame <= this.maxframes) {
                this.frame++;
                this.colour = lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
            } else {
                this.frame = 0;
                this.isTransitioning = false;
            }
        }

        let border = this.borderColour;
        let c = this.colour;
        let x = this.x;
        let y = this.y;
        let w = this.width;
        let h = this.height;
        let t = this.text;
        let deg = PI * 0.20872;
        let height = this.textSize;
        let len = this.text.length;
        switch (this.type) {
            case 5:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                rect(this.x, this.y, this.width, this.height, 15, 15, 15, 15);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                /*stroke(color(255, 0, 0));
                line(x, y + h/2, x+w, y + h/2);*/
                for (let i = 0; i < len; i++) {
                    text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
                }
                break;
            case 1:
                //shape
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                beginShape();
                vertex(x + 5, y - 25);
                vertex(x + w - 25, y + 15);
                vertex(x + w - 25, y + 15);
                quadraticVertex(x + w, y + 35, x + w, y + 70);
                vertex(x + w, y + h - 10);
                vertex(x + w, y + h - 10);
                quadraticVertex(x + w, y + h, x + w - 5, y + h - 5);
                vertex(x + 25, y + h - 45);
                vertex(x + 25, y + h - 45);
                quadraticVertex(x, y + h - 65, x, y + h - 100);
                vertex(x, y + h - 100);
                vertex(x, y - 20);
                vertex(x, y - 20);
                quadraticVertex(x, y - 25, x + 5, y - 25);
                endShape(CLOSE);

                //text
                if (this.text !== "") {
                    noStroke();
                    fill(this.textColour);
                    textAlign(CENTER);
                    textSize(this.textSize);
                    for (let i = 0; i < len; i++) {
                        translate(x + w / 2, y + h / 2 - height / 2 + height * (i + 1 - len / 2) / 0.75);
                        rotate(deg);
                        text(t[i], -textWidth("1") / 2, 0);
                        rotate(-deg);
                        translate(-x - w / 2, -y - h / 2 + height / 2 - height * (i + 1 - len / 2) / 0.75);
                    }
                }
                break;
            case 2:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                rect(this.x, this.y, this.width, this.height, 15, 15, 5, 40);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                /*stroke(color(255, 0, 0));
                line(x, y + h/2, x+w, y + h/2);*/
                for (let i = 0; i < len; i++) {
                    text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
                }
                break;
            case 3:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                rect(this.x, this.y, this.width, this.height, 15, 15, 40, 5);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                for (let i = 0; i < len; i++) {
                    text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
                }
                break;
            case 4:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                beginShape();
                vertex(x + w - 5, y - 25);
                vertex(x + 25, y + 15);
                vertex(x + 25, y + 15);
                quadraticVertex(x, y + 35, x, y + 70);
                vertex(x, y + h - 10);
                vertex(x, y + h - 10);
                quadraticVertex(x, y + h, x + 5, y + h - 5);
                vertex(x + w - 25, y + h - 45);
                vertex(x + w - 25, y + h - 45);
                quadraticVertex(x + w, y + h - 65, x + w, y + h - 100);
                vertex(x + w, y + h - 100);
                vertex(x + w, y - 20);
                vertex(x + w, y - 20);
                quadraticVertex(x + w, y - 25, x + w - 5, y - 25);
                endShape(CLOSE);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                deg = PI * 0.20872;
                for (let i = 0; i < len; i++) {
                    translate(x + w / 2, y + h / 2 - height / 2 + height * (i + 1 - len / 2) / 0.75);
                    rotate(-deg);
                    text(t[i], +textWidth("1") / 2, 0);
                    rotate(+deg);
                    translate(-x - w / 2, -y - h / 2 + height / 2 - height * (i + 1 - len / 2) / 0.75);
                }
                break;
            default:
                console.log("wrong skill button type: " + this.id);
                break;
        }
    };

    this.in = function () {
        let x = fixCoordScale(mouseX);
        let y = fixCoordScale(mouseY);
        if (this.type === 1) {
            let top_y = 0.77792 * x + (this.y - 25 - (this.x + 5) * 0.77792);

            return (this.x <= x && x <= (this.width + this.x) &&
                y - top_y <= this.height - 35 &&
                0 <= y - top_y);

        } else if (this.type === 4) {
            let top_y = -0.77792 * x + (this.y + 15 + (this.x + 25) * 0.77792);
            return (this.x <= x && x <= (this.width + this.x) &&
                y - top_y <= this.height - 35 &&
                0 <= y - top_y);
        } else {
            return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y));
        }
    };

    this.clicked = function () {
        console.log("CLICKED " + this.id);
        this.colour = this.clickedColour;
        this.clickTimer = this.clickLinger;
        sendSkill(this.id);
    };

    this.unclick = function () {
        if (this.isHovered) {
            this.colour = this.hoverColour;
        } else {
            this.colour = this.baseColour;
        }
        //console.log("UNCLICKED", this.id);
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;
        }
        this.previousColour = this.colour;
        this.destColour = this.hoverColour;
    };

    this.unhovered = function () {
        this.hoverTimer = 0;
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;

        }
        this.previousColour = this.colour;
        this.destColour = this.baseColour;
    };

    this.setColour = function (stringColour) {
        if (this.rawColour !== stringColour) {
            this.rawColour = stringColour;
            this.frame = 0;
            this.isTransitioning = false;
            this.baseColour = color(stringColour);
            this.baseColour.setAlpha(0.87 * 255);
            this.colour = this.baseColour;
            this.hoverColour = color(stringColour);
            this.clickedColour = color(stringColour);
            let hoverchange = 17;
            let clickchange = 34;
            if (isLight(this.baseColour)) {
                this.textColour = color(dark.toString());
                this.textColour.setAlpha(0.87 * 255);
                this.hoverColour.setRed(red(this.colour) - hoverchange);
                this.hoverColour.setGreen(green(this.colour) - hoverchange);
                this.hoverColour.setBlue(blue(this.colour) - hoverchange);
                this.clickedColour.setRed(red(this.colour) - clickchange);
                this.clickedColour.setGreen(green(this.colour) - clickchange);
                this.clickedColour.setBlue(blue(this.colour) - clickchange);
                this.hoverColour.setAlpha(0.87 * 255);
                this.clickedColour.setAlpha(0.87 * 255);
            } else {
                this.hoverColour.setRed(red(this.colour) + hoverchange);
                this.hoverColour.setGreen(green(this.colour) + hoverchange);
                this.hoverColour.setBlue(blue(this.colour) + hoverchange);
                this.clickedColour.setRed(red(this.colour) + clickchange);
                this.clickedColour.setGreen(green(this.colour) + clickchange);
                this.clickedColour.setBlue(blue(this.colour) + clickchange);
                this.hoverColour.setAlpha(0.87 * 255);
                this.clickedColour.setAlpha(0.87 * 255);
                /*this.hoverColour.setAlpha(0.87 * 255);
                this.clickedColour.setAlpha(0.91 * 255);*/
                this.textColour = color(light.toString());
                this.textColour.setAlpha(0.87 * 255);
            }
        }
    };

    this.setText = function (t) {
        if (this.rawText !== t) {
            this.rawText = t;
            this.hoverText = SKILLDESCRIPTIONS.get(t);
            if (!!this.hoverText) {
                this.hoverLines = calculateLines(this.hoverText);
            }
            this.text = split(t, " ");
            let w = this.width;
            this.textSize = 50;
            textSize(this.textSize);
            for (let word of this.text) {
                let width = textWidth(word);
                while (width >= w) {
                    this.textSize--;
                    textSize(this.textSize);
                    width = textWidth(word);
                }
            }
        }
    };

    this.getText = function () {
        if (this.text !== "") {
            return this.text.join(" ");
        } else {
            return "";
        }
    };

    this.setState = function (State) {
        this.state = State;
        switch (State) {
            //0 - active, -1 - on CD, -2 - dis by eff ??? -100 - disabled
            case 0:
                this.baseColour.setAlpha(0.87 * 255);
                this.hoverColour.setAlpha(0.87 * 255);
                this.clickedColour.setAlpha(0.87 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.87 * 255);
                this.borderColour = dark;
                // this.clickable = this.isMine;
                this.clickable = true;
                this.borderWidth = 4.5;
                break;
            case -1:
                this.baseColour.setAlpha(0.5 * 255);
                this.hoverColour.setAlpha(0.5 * 255);
                this.clickedColour.setAlpha(0.5 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.5 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 4.5;
                break;
            case -2:
                this.baseColour.setAlpha(0.5 * 255);
                this.hoverColour.setAlpha(0.5 * 255);
                this.clickedColour.setAlpha(0.5 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.5 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 2.5;
                break;
            case -100:
                this.baseColour.setAlpha(0.3 * 255);
                this.hoverColour.setAlpha(0.3 * 255);
                this.clickedColour.setAlpha(0.3 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.3 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 2.5;
                break;
            default:
                this.baseColour.setAlpha(0.5 * 255);
                this.hoverColour.setAlpha(0.5 * 255);
                this.clickedColour.setAlpha(0.5 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.5 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 4.5;
                break;
        }
        if (this.isHovered) {
            this.colour = this.hoverColour;
        } else {
            this.colour = this.baseColour;
        }
    };

    this.displayHover = function () {
        if (this.hoverTimer < this.hoverLinger) {
            this.hoverTimer += 1;
            return;
        }
        if (this.hoverText) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines);
        }
    };
}

function displayStandardHoverBubble(hoverText, lines) {
    let hoverSize = 15;
    textAlign(LEFT);
    textSize(hoverSize);
    noStroke();
    fill(hoverc);
    let changerForFlipping;
    let w = textWidth(hoverText);
    let amnt = 965;
    if (touch) amnt-= 40;
    if (fixCoordScale(mouseX) >= amnt && w < 290) {
        changerForFlipping = -w - 20;
        if (touch) {
            changerForFlipping -= 40;
        }
    } else if (fixCoordScale(mouseX) >= amnt) {
        changerForFlipping = -310;
        if (touch) {
            changerForFlipping -= 40;
        }
    } else {
        changerForFlipping = 0;
        if (touch) {
            changerForFlipping += 40;
        }
    }
    let changerForFlippingY;
    let h = lines * hoverSize + (lines - 1) * 5;
    if (fixCoordScale(mouseY) + h + 10 >= 550) {
        changerForFlippingY = -h;
    } else {
        changerForFlippingY = 0;
    }
    if (h === hoverSize && w < 290) {
        rect(fixCoordScale(mouseX) + changerForFlipping, fixCoordScale(mouseY) + changerForFlippingY, w + 20, h + 10, 5);
    } else {
        rect(fixCoordScale(mouseX) + changerForFlipping, fixCoordScale(mouseY) + changerForFlippingY, 310, h + 10, 5);
    }
    strokeWeight(0.5);
    stroke(dark);
    fill(dark);
    textLeading(20);

    text(hoverText, fixCoordScale(mouseX) + 10 + changerForFlipping, fixCoordScale(mouseY) + changerForFlippingY + 5, 290);
}

function StandardButton(x, y, s, t, size, id, col) {
    this.id = id;
    if (this.id === "GiveUp") {
        this.hoverText = "Click here to give up and end the match.";
        this.hoverLines = calculateLines(this.hoverText);
        this.warned = false;
    } else if (this.id === "back") {
        this.hoverText = "Click here to return to your rewards page.";
        this.hoverLines = calculateLines(this.hoverText);
    } else {
        this.hoverText = "";
    }
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = size;
    textSize(size);
    this.width = textWidth(t) + 10;
    this.smooth = s;
    if (!!col) {
        this.baseColour = color(col.toString());
    } else {
        this.baseColour = color(light.toString());
    }
    this.colour = this.baseColour;
    let hoverchange = 20;
    let clickchange = 35;
    this.maxframes = 7;
    this.frame = 0;
    this.destColour = this.baseColour;
    this.previousColour = this.baseColour;
    this.isTransitioning = false;
    this.hoverColour = color(this.colour.toString());
    this.clickedColour = color(this.colour.toString());
    this.hoverColour.setRed(red(this.colour) - hoverchange);
    this.hoverColour.setGreen(green(this.colour) - hoverchange);
    this.hoverColour.setBlue(blue(this.colour) - hoverchange);
    this.clickedColour.setRed(red(this.colour) - clickchange);
    this.clickedColour.setGreen(green(this.colour) - clickchange);
    this.clickedColour.setBlue(blue(this.colour) - clickchange);
    this.textColour = dark;
    this.height = size + 10;
    this.clickable = true;
    this.hoverable = true;
    this.clickLinger = 10;
    this.clickTimer = 0;
    if (this.hoverable) {
        this.hoverTimer = 0;
        if (touch) {
            this.hoverLinger = taphoverLinger;
        } else {
            this.hoverLinger = hoverLinger;
        }
    }
    this.visible = true;

    this.hide = function () {
        this.visible = false;
        this.wasClickable = this.clickable;
        this.wasHoverable = this.hoverable;
        this.clickable = false;
        this.hoverable = false;
    };

    this.show = function () {
        this.visible = true;
        this.clickable = this.wasClickable;
        this.hoverable = this.wasHoverable;

    };

    this.display = function () {
        if (this.visible) {
            if (this.isTransitioning) {
                if (this.frame <= this.maxframes) {
                    this.frame++;
                    this.colour = lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
                } else {
                    this.frame = 0;
                    this.isTransitioning = false;
                }
            }
            noStroke();
            fill(this.colour);
            strokeWeight(1);
            textSize(this.textSize);
            rect(this.x, this.y, this.width, this.height, this.smooth);
            stroke(this.textColour);
            fill(this.textColour);
            textAlign(LEFT, CENTER);
            text(this.text, this.x + 5, this.y + this.height / 2);
            textAlign(LEFT, BASELINE);
        }
    };

    this.in = function () {
        let x = fixCoordScale(mouseX);
        let y = fixCoordScale(mouseY);
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y))
    };

    this.clicked = function () {
        this.colour = this.clickedColour;
        this.clickTimer = this.clickLinger;
        if (this.id === "back") {
            console.log("RELOCATED!~", "/game/rewards/");
            window.location = '/game/rewards/';
            return
        } if (!this.warned && this.id === "GiveUp") {
            this.warned = true;
            let info = getElement("info");
            info.setColour(dark);
            info.setText("Are you sure you want to give up? If so, click this button again.");
        } else if (this.id === "GiveUp") {
            sendSkill(this.id);
        }
    };

    this.unclick = function () {
        if (this.isHovered) {
            this.colour = this.hoverColour;
        } else {
            this.colour = this.baseColour;
        }
    };

    this.setText = function (t) {
        this.text = t;
        textSize(this.textSize);
        this.width = textWidth(t) + 10;
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;
        }
        this.previousColour = this.colour;
        this.destColour = this.hoverColour;
    };

    this.unhovered = function () {
        this.hoverTimer = 0;
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;

        }
        this.previousColour = this.colour;
        this.destColour = this.baseColour;
    };

    this.displayHover = function () {
        if (this.hoverTimer < this.hoverLinger) {
            this.hoverTimer += 1;
        } else if (this.hoverText) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines);
        }
    }
}

function LoadingScreenNoP(x, y, w, h, id) {
    this.x = x;
    this.y = y;
    this.w = w;
    this.h = h;
    this.fade = false;
    this.frame = 0;
    this.maxframes = 20;
    this.pauseTime = 10;
    this.pause = 0;
    this.phi = 0.0;
    this.squares = [];
    this.paused = false;
    this.stopped = 0;
    this.hoverable = false;
    let width;
    let height;
    if (w > h) {
        width = h / 2;
        height = h / 2;
    } else {
        width = w / 2;
        height = w / 2;
    }
    if (!id) {
        this.id = 'loading_screen';
    }
    let distance = 0.05;
    this.squares.push([0 + x, height * (1 + distance) + y, width * (1 - distance), height * (1 - distance), 0 + x, height * (1 + distance) + y, undefined]); //1 bot left pi to 3pi/2
    this.squares.push([0 + x, 0 + y, width * (1 - distance), height * (1 - distance), 0 + x, 0 + y, undefined]); //4 top left 3pi/2 to 2pi
    this.squares.push([width * (1 + distance) + x, 0 + y, width * (1 - distance), height * (1 - distance), width * (1 + distance) + x, 0 + y, undefined]); //3 top right pi to 3pi/2
    this.squares.push([width * (1 + distance) + x, height * (1 + distance) + y, width * (1 - distance), height * (1 - distance), width * (1 + distance) + x, height * (1 + distance) + y, undefined]); //2 bot right pi/2 to pi

    this.stop = function () {
        if (this.stopped < 2) {
            this.stopped = 1;
        }
    };

    this.enableFade = function () {
        this.fade = true;
    };

    this.restart = function () {
        let width;
        let height;
        let distance = 0.05;
        if (this.w > this.h) {
            width = this.h / 2;
            height = this.h / 2;
        } else {
            width = this.w / 2;
            height = this.w / 2;
        }
        this.squares[0][0] = this.x;
        this.squares[0][1] = height * (1 + distance) + this.y;
        this.squares[3][0] = width * (1 + distance) + this.x;
        this.squares[3][1] = height * (1 + distance) + this.y;
        this.squares[2][0] = width * (1 + distance) + this.x;
        this.squares[2][1] = this.y;
        this.squares[1][0] = this.x;
        this.squares[1][1] = this.y;
        this.frame = 0;
        this.phi = 0.0;
        this.pause = 0;
        this.stopped = 0;
        for (let i = 0; i < this.squares.length; i++) {
            let square = this.squares[i];
            square[6].setAlpha(255)
        }
    };

    this.display = function () {
        let addition = 0.0;
        for (let i = 0; i < this.squares.length; i++) {
            let s = this.squares[i];
            let curr_phi = this.phi + addition;

            //rotate
            let center_x = s[0] + s[2] / 2;
            let center_y = s[1] + s[3] / 2;
            translate(center_x, center_y);
            rotate(-curr_phi * 2);
            //draw
            let prev = i !== 0 ? i - 1 : 3;
            let inter = map(this.phi, 0, -PI / 2, 0, 1);
            let c = lerpColor(s[6], this.squares[prev][6], inter);
            fill(c);
            noStroke();
            rect(-s[2] / 2, -s[3] / 2, s[2], s[3], s[2] / 9);
            //rotate back
            rotate(curr_phi * 2);
            translate(-center_x, -center_y);


            //move
            let x = s[0];
            let y = s[1];
            let r = sin(2 * curr_phi) * this.w;
            let new_x = r * cos(curr_phi) + s[4];
            let new_y = r * sin(curr_phi) + s[5];
            s[0] = new_x;
            s[1] = new_y;
            addition -= PI / 2;

        }
        if (this.stopped !== 2) {
            if (this.paused) {
                this.pause += 1;
                if (this.stopped > 0 && this.fade) {
                    for (let i = 0; i < this.squares.length; i++) {
                        let square = this.squares[i];
                        square[6].setAlpha(255 - 255 * (0.45 + this.pause / 10))
                    }
                }
                if (this.pause >= this.pauseTime) {
                    this.pause = 0;
                    this.paused = false;
                    if (this.stopped > 0) {
                        this.stopped = 2;
                    }
                }
            } else if (this.frame >= this.maxframes) {
                this.frame = 0;
                this.phi = 0.0;
                let lastCol = this.squares[3][6];
                for (let i = 3; i > 0; i--) {
                    let s = this.squares[i];
                    s[6] = this.squares[i - 1][6];
                }
                this.squares[0][6] = lastCol;
                this.paused = true;
            } else {
                this.frame += 1;
                this.phi -= PI / (2 * (this.maxframes));
            }
        }

    };

    this.setColours = function (c1, c2, c3, c4) {
        this.squares[0][6] = color(c1.toString());
        this.squares[1][6] = color(c2.toString());
        this.squares[2][6] = color(c3.toString());
        this.squares[3][6] = color(c4.toString());

    }
}

function SkillButtonMini(x, y, t, id) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = 50;
    this.textColour = color(dark.toString());
    this.borderColour = color(dark.toString());
    this.borderWidth = 2;
    this.width = 100;
    this.height = 100;
    this.maxframes = 7;
    this.frame = 0;
    this.destColour = this.baseColour;
    this.previousColour = this.baseColour;
    this.isTransitioning = false;
    this.rawColour = "";
    this.rawText = "";
    this.colour = clickc;
    this.baseColour = clickc;
    this.hoverColour = clickc;
    this.clickedColour = clickc;
    this.isHovered = false;
    this.hoverText = "";
    this.clickable = false;
    this.hoverable = true;

    this.display = function () {
        if (this.isTransitioning) {
            if (this.frame <= this.maxframes) {
                this.frame++;
                this.colour = lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
            } else {
                this.frame = 0;
                this.isTransitioning = false;
            }
        }
        let border = this.borderColour;
        let c = this.colour;
        let x = this.x;
        let y = this.y;
        let w = this.width;
        let h = this.height;
        let t = this.text;
        let height = this.textSize;
        let len = this.text.length;
        stroke(border);
        c.setAlpha(255 * 0.8);
        fill(c);
        strokeWeight(this.borderWidth);
        rect(this.x, this.y, this.width, this.height, 4, 4, 4, 4);
        noStroke();
        fill(this.textColour);
        textAlign(CENTER);
        textSize(this.textSize);
        /*stroke(color(255, 0, 0));
        line(x, y + h/2, x+w, y + h/2);*/
        for (let i = 0; i < len; i++) {
            text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
        }
    };

    this.in = function () {
        let x = fixCoordScale(mouseX);
        let y = fixCoordScale(mouseY);
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y));
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;
        }
        this.previousColour = this.colour;
        this.destColour = this.hoverColour;
    };

    this.unhovered = function () {
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;

        }
        this.previousColour = this.colour;
        this.destColour = this.baseColour;
    };

    this.setColour = function (stringColour) {
        if (this.rawColour !== stringColour) {
            this.rawColour = stringColour;
            this.frame = 0;
            this.isTransitioning = false;
            this.baseColour = color(stringColour);
            this.colour = this.baseColour;
            this.hoverColour = color(stringColour);
            this.clickedColour = color(stringColour);
            let hoverchange = 17;
            let clickchange = 34;
            if (isLight(this.baseColour)) {
                this.textColour = color(dark.toString());
                this.hoverColour.setRed(red(this.colour) - hoverchange);
                this.hoverColour.setGreen(green(this.colour) - hoverchange);
                this.hoverColour.setBlue(blue(this.colour) - hoverchange);
                this.clickedColour.setRed(red(this.colour) - clickchange);
                this.clickedColour.setGreen(green(this.colour) - clickchange);
                this.clickedColour.setBlue(blue(this.colour) - clickchange);
            } else {
                this.hoverColour.setRed(red(this.colour) + hoverchange);
                this.hoverColour.setGreen(green(this.colour) + hoverchange);
                this.hoverColour.setBlue(blue(this.colour) + hoverchange);
                this.clickedColour.setRed(red(this.colour) + clickchange);
                this.clickedColour.setGreen(green(this.colour) + clickchange);
                this.clickedColour.setBlue(blue(this.colour) + clickchange);
                this.textColour = color(light.toString());
            }
        }
    };

    this.setText = function (t) {
        if (this.rawText !== t) {
            this.rawText = t;
            this.hoverText = SKILLDESCRIPTIONS.get(t);
            if (!!this.hoverText) {
                this.hoverLines = calculateLines(this.hoverText);
            }
            this.text = split(t, " ");
            let w = this.width;
            this.textSize = 50;
            textSize(this.textSize);
            for (let word of this.text) {
                let width = textWidth(word);
                while (width >= w) {
                    this.textSize--;
                    textSize(this.textSize);
                    width = textWidth(word);
                }
            }
        }
    };

    this.getText = function () {
        if (this.text !== "") {
            return this.text.join(" ");
        } else {
            return "";
        }
    };

    this.displayHover = function () {
        if (this.hoverText) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines);
        }
    };
}