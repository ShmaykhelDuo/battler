function InterfaceButton(p, x, y, t, size, id, type, width, height) {
    this.p = p;
    this.id = id;
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = size;
    this.smooth = 5;
    this.type = type;
    if (this.type === "A") {
        this.height = size + 10;
        this.p.textSize(size);
        this.width = this.p.textWidth(t) + 10;
    } else if (this.type === "B") {
        this.width = width;
        while (this.p.textWidth(t) > (this.width - 10)) {
            this.textSize--;
            this.p.textSize(this.textSize);
        }

        if (height) {
            this.height = height;
        } else {
            this.height = 70;
        }

    }
    this.textColour = dark;
    this.clickable = true;
    this.hoverable = true;
    this.clickLinger = 4;
    this.clickTimer = 0;
    this.baseColour = this.p.color(light.toString());
    this.colour = this.baseColour;
    let hoverchange = 20;
    let clickchange = 35;
    this.maxframes = 7;
    this.frame = 0;
    this.destColour = this.baseColour;
    this.previousColour = this.baseColour;
    this.isTransitioning = false;
    this.hoverColour = this.p.color(this.colour.toString());
    this.clickedColour = this.p.color(this.colour.toString());
    this.hoverColour.setRed(this.p.red(this.colour) - hoverchange);
    this.hoverColour.setGreen(this.p.green(this.colour) - hoverchange);
    this.hoverColour.setBlue(this.p.blue(this.colour) - hoverchange);
    this.clickedColour.setRed(this.p.red(this.colour) - clickchange);
    this.clickedColour.setGreen(this.p.green(this.colour) - clickchange);
    this.clickedColour.setBlue(this.p.blue(this.colour) - clickchange);
    /*this.hoverColour.setAlpha(0.87 * 255);
    this.clickedColour.setAlpha(0.87 * 255);
    this.baseColour.setAlpha(0.87 * 255);*/

    this.setColour = function (colour) {
        this.baseColour = this.p.color(colour);
        let hoverchange = 17;
        let clickchange = 34;
        this.hoverColour = this.p.color(this.baseColour.toString());
        this.clickedColour = this.p.color(this.baseColour.toString());
        this.hoverColour.setRed(this.p.red(this.baseColour) - hoverchange);
        this.hoverColour.setGreen(this.p.green(this.baseColour) - hoverchange);
        this.hoverColour.setBlue(this.p.blue(this.baseColour) - hoverchange);
        this.clickedColour.setRed(this.p.red(this.baseColour) - clickchange);
        this.clickedColour.setGreen(this.p.green(this.baseColour) - clickchange);
        this.clickedColour.setBlue(this.p.blue(this.baseColour) - clickchange);
        /*this.hoverColour.setAlpha(0.87 * 255);
        this.clickedColour.setAlpha(0.87 * 255);
        this.baseColour.setAlpha(0.87 * 255);*/
        this.colour = this.baseColour;
    };

    this.display = function () {
        if (this.isTransitioning) {
            if (this.frame <= this.maxframes) {
                this.frame++;
                this.colour = this.p.lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
            } else {
                this.frame = 0;
                this.isTransitioning = false;
            }
        }
        this.p.noStroke();
        this.p.fill(this.colour);
        this.p.strokeWeight(1);
        this.p.textSize(this.textSize);
        this.p.rect(this.x, this.y, this.width, this.height, this.smooth);
        if (this.clickable) {
            this.textColour = dark;
        } else {
            this.textColour = this.p.color(110);
        }
        this.p.stroke(this.textColour);
        this.p.fill(this.textColour);
        if (this.type === "A") {
            this.p.textAlign(this.p.LEFT, this.p.CENTER);
            this.p.text(this.text, this.x + 5, this.y + this.height / 2);
        } else if (this.type === "B") {
            this.p.textAlign(this.p.CENTER, this.p.CENTER);
            this.p.text(this.text, this.x + this.width / 2, this.y + this.height / 2);

        }
        this.p.textAlign(this.p.LEFT, this.p.BASELINE);
    };

    this.in = function () {
        let x = this.p.mouseX;
        let y = this.p.mouseY;
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y))
    };

    this.clicked = function () {
        if (this.clickable) {
            this.colour = this.clickedColour;
            this.clickTimer = this.clickLinger;
            this.onClick();
        }
    };

    this.onClick = function () {

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
        this.p.textSize(this.textSize);
        if (this.type === "A") {
            this.width = this.p.textWidth(t);
        } else if (this.type === "B") {
            this.textSize = 45;
            while (this.p.textWidth(t) > (this.width - 10)) {
                this.textSize--;
                this.p.textSize(this.textSize);
            }
        }
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

    this.displayHover = function () {
    }
}

function InterfaceText(p, x, y, colour, t, size, id, type, width) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.textSize = size;
    this.text = t;
    this.type = type;
    if (type === "B" || type === "C") {
        this.width = width;
    } else if (type === "A" || type === "D") {
        p.textSize(size);
        this.width = p.textWidth(t);
    }

    this.height = size;
    this.p = p;
    this.textColour = colour;
    this.hoverable = false;

    this.display = function () {
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
        this.p.stroke(this.textColour);
        this.p.strokeWeight(1);
        this.p.fill(this.textColour);
        this.p.textSize(this.textSize);
        if (this.type === "B" || this.type === "C") {
            if (this.type === "B") {
                this.p.textAlign(this.p.CENTER);
            } else {
                this.p.textAlign(this.p.LEFT);
                this.p.textLeading(33);
            }
            this.p.text(this.text, this.x, this.y, this.width);
        } else if (this.type === "A" || this.type === "D") {
            if (this.type === "A") {
                this.p.textAlign(this.p.LEFT);
            } else if (this.type === "D") {
                this.p.textAlign(this.p.RIGHT);
            }
            this.p.text(this.text, this.x, this.y);
        }
    };

    this.in = function () {
        let x = this.p.mouseX;
        let y = this.p.mouseY;
        return (this.x <= x && x <= (this.width + this.x) && y <= this.y && (this.y - this.height) <= y);
    };

    this.setColour = function (c) {
        this.textColour = this.p.color(c);
    };

    this.setText = function (t) {
        this.text = t;
        this.height = this.textSize;
        this.p.textSize(this.textSize);
        if (type === "B" || type === "C") {
            this.width = width;
        } else if (type === "A") {
            p.textSize(size);
            this.width = p.textWidth(t);
        }
    };
}

function InterfaceImage(p, x, y, path, id, name, width, height, colour) {
    this.loadedObj = {};
    this.id = id;
    this.x = x;
    this.y = y;
    this.path = path;
    if (!!colour) {
        this.colour = p.lerpColor(leftP.color(colour), light, 0.45);
    } else {
        this.colour = p.color(DARKC);
    }
    if (!path) {
        this.width = width;
        this.height = height;
        this.image = undefined;
        this.name = "";
    } else {
        this.width = width;
        this.height = height;
        this.image = p.loadImage(path, img => {
            this.loadedObj.loaded = true;
            console.log(this.loadedObj);
        }, this.failed);
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

    this.loaded = function() {
        return this.loadedObj.hasOwnProperty("loaded")
    };

    this.copy = function () {
        let img = new InterfaceImage(p, 0, 0, "", "", "", this.width, this.height);
        img.image = this.image;
        img.path = this.path;
        img.colour = this.colour;
        return img;

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

    this.open = function (path, name, width, height) {
        this.x = this.x + (550 - width) / 2;
        this.width = width;
        this.height = height;
        this.image = p.loadImage(path, img => {
            this.loadedObj.loaded = true;
        }, this.failed);
        this.name = name;
    };

    this.failed = function (img) {
        console.log("FAILED TO LOAD", this.name);
    };

    this.display = function () {
        if (this.image) {
            p.image(this.image, this.x, this.y, this.width, this.height);
        }
    };
}

function InterfaceImageBox() {
    this.images = [];

    this.add = function (image) {
        this.images.push(image);
    };

    this.contains = function (name) {
        for (let image of this.images) {
            if (image.name === name) {
                return true;
            }
        }
        return false;
    };

    this.get = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                return this.images[i];
            }
        }
        return undefined;
    };


}

function LoadingScreen(p, x, y, w, h) {
    this.p = p;
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
            this.p.translate(center_x, center_y);
            this.p.rotate(-curr_phi * 2);
            //draw
            let prev = i !== 0 ? i - 1 : 3;
            let inter = this.p.map(this.phi, 0, -this.p.PI / 2, 0, 1);
            let c = this.p.lerpColor(s[6], this.squares[prev][6], inter);
            this.p.fill(c);
            this.p.noStroke();
            this.p.rect(-s[2] / 2, -s[3] / 2, s[2], s[3], s[2] / 9);
            //rotate back
            this.p.rotate(curr_phi * 2);
            this.p.translate(-center_x, -center_y);


            //move
            let x = s[0];
            let y = s[1];
            let r = this.p.sin(2 * curr_phi) * this.w;
            let new_x = r * this.p.cos(curr_phi) + s[4];
            let new_y = r * this.p.sin(curr_phi) + s[5];
            s[0] = new_x;
            s[1] = new_y;
            addition -= this.p.PI / 2;

        }
        if (this.stopped !== 2) {
            if (this.paused) {
                this.pause += 1;
                if (this.stopped > 0 && this.fade) {
                    for (let i = 0; i < this.squares.length; i++) {
                        let square = this.squares[i];
                        square[6].setAlpha(255 - 255*(0.45+this.pause/10))
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
                this.phi -= this.p.PI / (2 * (this.maxframes));
            }
        }

    };

    this.setColours = function (c1, c2, c3, c4) {
        this.squares[0][6] = this.p.color(c1.toString());
        this.squares[1][6] = this.p.color(c2.toString());
        this.squares[2][6] = this.p.color(c3.toString());
        this.squares[3][6] = this.p.color(c4.toString());

    }
}

function SkillButtonMiniP(p, x, y, t, id) {
    this.p = p;
    this.id = id;
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = 50;
    this.textColour = p.color(dark.toString());
    this.borderColour = p.color(dark.toString());
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
                this.colour = this.p.lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
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
        this.p.stroke(border);
        c.setAlpha(255 * 0.8);
        this.p.fill(c);
        this.p.strokeWeight(this.borderWidth);
        this.p.rect(this.x, this.y, this.width, this.height, 4, 4, 4, 4);
        this.p.noStroke();
        this.p.fill(this.textColour);
        this.p.textAlign(this.p.CENTER);
        this.p.textSize(this.textSize);
        /*stroke(color(255, 0, 0));
        line(x, y + h/2, x+w, y + h/2);*/
        for (let i = 0; i < len; i++) {
            this.p.text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
        }
    };

    this.in = function () {
        let x = this.p.mouseX;
        let y = this.p.mouseY;
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
            this.baseColour = this.p.color(stringColour);
            this.colour = this.baseColour;
            this.hoverColour = this.p.color(stringColour);
            this.clickedColour = this.p.color(stringColour);
            let hoverchange = 17;
            let clickchange = 34;
            if (isLightP(this.p, this.baseColour)) {
                this.textColour = this.p.color(dark.toString());
                this.hoverColour.setRed(this.p.red(this.colour) - hoverchange);
                this.hoverColour.setGreen(this.p.green(this.colour) - hoverchange);
                this.hoverColour.setBlue(this.p.blue(this.colour) - hoverchange);
                this.clickedColour.setRed(this.p.red(this.colour) - clickchange);
                this.clickedColour.setGreen(this.p.green(this.colour) - clickchange);
                this.clickedColour.setBlue(this.p.blue(this.colour) - clickchange);
            } else {
                this.hoverColour.setRed(this.p.red(this.colour) + hoverchange);
                this.hoverColour.setGreen(this.p.green(this.colour) + hoverchange);
                this.hoverColour.setBlue(this.p.blue(this.colour) + hoverchange);
                this.clickedColour.setRed(this.p.red(this.colour) + clickchange);
                this.clickedColour.setGreen(this.p.green(this.colour) + clickchange);
                this.clickedColour.setBlue(this.p.blue(this.colour) + clickchange);
                this.textColour = this.p.color(light.toString());
            }
        }
    };

    this.setText = function (t) {
        if (this.rawText !== t) {
            this.rawText = t;
            this.hoverText = SKILLDESCRIPTIONS.get(t);
            if (!!this.hoverText) {
                this.hoverLines = interfaceCalculateLines(this.p, this.hoverText);
            }
            this.text = this.p.split(t, " ");
            let w = this.width;
            this.textSize = 50;
            this.p.textSize(this.textSize);
            for (let word of this.text) {
                let width = this.p.textWidth(word);
                while (width >= w) {
                    this.textSize--;
                    this.p.textSize(this.textSize);
                    width = this.p.textWidth(word);
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
            displayStandardHoverBubbleP(this.p, this.hoverText, this.hoverLines);
        }
    };
}

function interfaceCalculateLines(p, hoverText, width, size) {
    if (!width) {
        width = 290;
    }
    if (!size) {
        size = 15;
    }
    let height = 0;
    p.textSize(size);
    for (let line of hoverText.split("\n")) {
        height += 1;
        let x_pos = 0;
        for (let word of line.split(" ")) {
            /*strokeWeight(5);
            stroke(red2);
            point(mouseX+10+x_pos + change, mouseY+height*size);*/
            if (x_pos + p.textWidth(word + " ") < width) { //we are still on that line
                x_pos += p.textWidth(word + " ")
            } else { //start a new line
                height += 1;
                x_pos = p.textWidth(word + " ");
            }
        }
    }
    return height;
}

function displayStandardHoverBubbleP(p, hoverText, lines) {
    let hoverSize = 15;
    p.textAlign(p.LEFT);
    p.textSize(hoverSize);
    p.noStroke();
    p.fill(hoverc);
    let changerForFlipping;
    let w = p.textWidth(hoverText);
    let amnt = 965;
    if (touch) amnt-= 40;
    if (p.mouseX >= amnt && w < 290) {
        changerForFlipping = -w - 20;
        if (touch) {
            changerForFlipping -= 40;
        }
    } else if (p.mouseX >= amnt) {
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
    if (p.mouseY + h + 10 >= 550) {
        changerForFlippingY = -h;
    } else {
        changerForFlippingY = 0;
    }
    if (h === hoverSize && w < 290) {
        p.rect(p.mouseX + changerForFlipping, p.mouseY + changerForFlippingY, w + 20, h + 10, 5);
    } else {
        p.rect(p.mouseX + changerForFlipping, p.mouseY + changerForFlippingY, 310, h + 10, 5);
    }
    p.strokeWeight(0.5);
    p.stroke(dark);
    p.fill(dark);
    p.textLeading(20);
    p.text(hoverText, p.mouseX + 10 + changerForFlipping, p.mouseY + changerForFlippingY + 5, 290);
}


function isLightP(p, colour) {
    return p.lightness(colour) > 50;
}

